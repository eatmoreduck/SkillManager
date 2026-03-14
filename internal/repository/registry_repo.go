package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"skillmanager/internal/model"
)

// RegistryRepository reads skills from remote registries.
type RegistryRepository interface {
	FetchRegistry(url string, client *http.Client) ([]model.RegistrySkill, error)
	FetchSkillMeta(sourceURL string, client *http.Client) (*model.RegistrySkill, error)
}

type HTTPRegistryRepository struct{}

type legacyRegistryEnvelope struct {
	Skills []model.RegistrySkill `json:"skills"`
	Data   []model.RegistrySkill `json:"data"`
	Items  []model.RegistrySkill `json:"items"`
}

type skillsSHSearchEnvelope struct {
	Skills []skillsSHSearchItem `json:"skills"`
}

type skillsSHSearchItem struct {
	ID       string `json:"id"`
	SkillID  string `json:"skillId"`
	Name     string `json:"name"`
	Source   string `json:"source"`
	Installs int    `json:"installs"`
}

var skillsSHBrowseCardPattern = regexp.MustCompile(`(?s)<a[^>]+href="/([^"]+/[^"]+/[^"]+)"[^>]*>.*?<h3[^>]*>\s*([^<]+?)\s*</h3>\s*<p[^>]*>\s*([^<]+?)\s*</p>.*?<div[^>]*text-right[^>]*>.*?<span[^>]*>\s*([^<]+?)\s*</span>`)

func NewHTTPRegistryRepository() *HTTPRegistryRepository {
	return &HTTPRegistryRepository{}
}

func (r *HTTPRegistryRepository) FetchRegistry(rawURL string, client *http.Client) ([]model.RegistrySkill, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json, text/html;q=0.9, */*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("registry request failed with status %d", resp.StatusCode)
	}

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if skills, ok, err := parseRegistryJSON(payload); err != nil {
		return nil, err
	} else if ok {
		return skills, nil
	}

	if skills, ok := parseSkillsSHBrowseHTML(payload); ok {
		return skills, nil
	}

	return nil, fmt.Errorf("unsupported registry response format")
}

func (r *HTTPRegistryRepository) FetchSkillMeta(sourceURL string, client *http.Client) (*model.RegistrySkill, error) {
	req, err := http.NewRequest(http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("registry meta request failed with status %d", resp.StatusCode)
	}

	var skill model.RegistrySkill
	if err := json.NewDecoder(resp.Body).Decode(&skill); err != nil {
		return nil, err
	}

	return &skill, nil
}

func BuildRegistryBrowseURL(baseURL, category string) (string, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	if isSkillsSHRegistry(parsed) {
		if category != "" {
			return BuildRegistrySearchURL(baseURL, category)
		}
		parsed = normalizeSkillsSHURL(parsed)
		parsed.Path = "/"
		parsed.RawQuery = ""
		return strings.TrimRight(parsed.String(), "/"), nil
	}

	if !strings.HasSuffix(parsed.Path, "/") {
		parsed.Path += "/"
	}
	parsed.Path += "skills"

	q := parsed.Query()
	if category != "" {
		q.Set("category", category)
	}
	parsed.RawQuery = q.Encode()

	return parsed.String(), nil
}

func BuildRegistrySearchURL(baseURL, query string) (string, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	if isSkillsSHRegistry(parsed) {
		parsed = normalizeSkillsSHURL(parsed)
		parsed.Path = "/api/search"
		q := parsed.Query()
		q.Set("q", query)
		q.Set("limit", "100")
		parsed.RawQuery = q.Encode()
		return parsed.String(), nil
	}

	if !strings.HasSuffix(parsed.Path, "/") {
		parsed.Path += "/"
	}
	parsed.Path += "skills/search"

	q := parsed.Query()
	q.Set("q", query)
	parsed.RawQuery = q.Encode()

	return parsed.String(), nil
}

func parseRegistryJSON(payload []byte) ([]model.RegistrySkill, bool, error) {
	trimmed := bytes.TrimSpace(payload)
	if len(trimmed) == 0 {
		return nil, false, nil
	}

	switch trimmed[0] {
	case '[':
		var skills []model.RegistrySkill
		if err := json.Unmarshal(trimmed, &skills); err != nil {
			return nil, false, err
		}
		return normalizeRegistrySkills(skills), true, nil
	case '{':
	default:
		return nil, false, nil
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(trimmed, &raw); err != nil {
		return nil, false, err
	}

	if _, ok := raw["skills"]; ok && (raw["query"] != nil || raw["searchType"] != nil || raw["count"] != nil) {
		var current skillsSHSearchEnvelope
		if err := json.Unmarshal(trimmed, &current); err == nil && current.Skills != nil {
			return mapSkillsSHSearchResults(current.Skills), true, nil
		}
	}

	if _, ok := raw["skills"]; ok || raw["data"] != nil || raw["items"] != nil {
		var legacy legacyRegistryEnvelope
		if err := json.Unmarshal(trimmed, &legacy); err != nil {
			return nil, false, err
		}
		switch {
		case legacy.Skills != nil:
			return normalizeRegistrySkills(legacy.Skills), true, nil
		case legacy.Data != nil:
			return normalizeRegistrySkills(legacy.Data), true, nil
		case legacy.Items != nil:
			return normalizeRegistrySkills(legacy.Items), true, nil
		}
		return []model.RegistrySkill{}, true, nil
	}

	return nil, false, nil
}

func mapSkillsSHSearchResults(items []skillsSHSearchItem) []model.RegistrySkill {
	skills := make([]model.RegistrySkill, 0, len(items))
	for _, item := range items {
		id := strings.Trim(item.ID, "/")
		if id == "" && item.Source != "" && item.SkillID != "" {
			id = strings.Trim(item.Source, "/") + "/" + strings.Trim(item.SkillID, "/")
		}

		name := strings.TrimSpace(item.Name)
		if name == "" {
			name = strings.TrimSpace(item.SkillID)
		}

		source := strings.TrimSpace(item.Source)
		description := ""
		if source != "" {
			description = "Source: " + source
		}

		skills = append(skills, model.RegistrySkill{
			ID:          id,
			Name:        name,
			Description: description,
			Author:      source,
			Stars:       item.Installs,
			Tags:        []string{},
			InstallURL:  githubCloneURL(source),
		})
	}
	return skills
}

func parseSkillsSHBrowseHTML(payload []byte) ([]model.RegistrySkill, bool) {
	matches := skillsSHBrowseCardPattern.FindAllSubmatch(payload, -1)
	if len(matches) == 0 {
		return nil, false
	}

	skills := make([]model.RegistrySkill, 0, len(matches))
	seen := make(map[string]struct{}, len(matches))
	for _, match := range matches {
		if len(match) != 5 {
			continue
		}

		id := html.UnescapeString(string(match[1]))
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}

		name := html.UnescapeString(strings.TrimSpace(string(match[2])))
		source := html.UnescapeString(strings.TrimSpace(string(match[3])))
		installs := parseCompactCount(html.UnescapeString(strings.TrimSpace(string(match[4]))))

		description := ""
		if source != "" {
			description = "Source: " + source
		}

		skills = append(skills, model.RegistrySkill{
			ID:          id,
			Name:        name,
			Description: description,
			Author:      source,
			Stars:       installs,
			Tags:        []string{},
			InstallURL:  githubCloneURL(source),
		})
	}

	return skills, len(skills) > 0
}

func normalizeRegistrySkills(skills []model.RegistrySkill) []model.RegistrySkill {
	result := make([]model.RegistrySkill, 0, len(skills))
	for _, skill := range skills {
		if skill.Tags == nil {
			skill.Tags = []string{}
		}
		result = append(result, skill)
	}
	return result
}

func isSkillsSHRegistry(parsed *url.URL) bool {
	host := strings.ToLower(parsed.Hostname())
	return host == "skills.sh" || host == "www.skills.sh" || host == "api.skills.sh"
}

func normalizeSkillsSHURL(parsed *url.URL) *url.URL {
	copy := *parsed
	copy.Scheme = "https"
	copy.Host = "skills.sh"
	return &copy
}

func githubCloneURL(source string) string {
	source = strings.Trim(source, "/")
	if source == "" {
		return ""
	}
	return "https://github.com/" + source + ".git"
}

func parseCompactCount(raw string) int {
	raw = strings.TrimSpace(strings.ReplaceAll(raw, ",", ""))
	if raw == "" {
		return 0
	}

	multiplier := 1.0
	switch suffix := strings.ToUpper(raw[len(raw)-1:]); suffix {
	case "K":
		multiplier = 1_000
		raw = raw[:len(raw)-1]
	case "M":
		multiplier = 1_000_000
		raw = raw[:len(raw)-1]
	case "B":
		multiplier = 1_000_000_000
		raw = raw[:len(raw)-1]
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0
	}

	return int(value * multiplier)
}
