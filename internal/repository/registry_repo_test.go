package repository

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPRegistryRepositoryFetchRegistrySupportsEnvelope(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"skills":[{"id":"demo/skill","name":"Demo Skill","installUrl":"https://example.com/demo.git"}]}`))
	}))
	defer server.Close()

	repo := NewHTTPRegistryRepository()
	skills, err := repo.FetchRegistry(server.URL, server.Client())
	if err != nil {
		t.Fatalf("FetchRegistry() error = %v", err)
	}

	if len(skills) != 1 {
		t.Fatalf("FetchRegistry() len = %d, want 1", len(skills))
	}

	if skills[0].ID != "demo/skill" {
		t.Fatalf("FetchRegistry() id = %q, want %q", skills[0].ID, "demo/skill")
	}
}

func TestHTTPRegistryRepositoryFetchRegistrySupportsSkillsSHSearchEnvelope(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"query":"react","skills":[{"id":"openai/skills/linear","skillId":"linear","name":"linear","installs":929,"source":"openai/skills"}],"count":1}`))
	}))
	defer server.Close()

	repo := NewHTTPRegistryRepository()
	skills, err := repo.FetchRegistry(server.URL, server.Client())
	if err != nil {
		t.Fatalf("FetchRegistry() error = %v", err)
	}

	if len(skills) != 1 {
		t.Fatalf("FetchRegistry() len = %d, want 1", len(skills))
	}

	if skills[0].InstallURL != "https://github.com/openai/skills.git" {
		t.Fatalf("FetchRegistry() installUrl = %q, want %q", skills[0].InstallURL, "https://github.com/openai/skills.git")
	}
	if skills[0].Author != "openai/skills" {
		t.Fatalf("FetchRegistry() author = %q, want %q", skills[0].Author, "openai/skills")
	}
	if skills[0].Stars != 929 {
		t.Fatalf("FetchRegistry() stars = %d, want %d", skills[0].Stars, 929)
	}
}

func TestHTTPRegistryRepositoryFetchRegistrySupportsSkillsSHBrowseHTML(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`
<div class="h-[72px] lg:h-[52px]">
  <a class="group" href="/vercel-labs/agent-browser/agent-browser">
    <div class="lg:col-span-1 text-left"><span>25</span></div>
    <div class="lg:col-span-13 min-w-1">
      <h3>agent-browser</h3>
      <p>vercel-labs/agent-browser</p>
    </div>
    <div class="lg:col-span-2 text-right"><span>94.8K</span></div>
  </a>
</div>`))
	}))
	defer server.Close()

	repo := NewHTTPRegistryRepository()
	skills, err := repo.FetchRegistry(server.URL, server.Client())
	if err != nil {
		t.Fatalf("FetchRegistry() error = %v", err)
	}

	if len(skills) != 1 {
		t.Fatalf("FetchRegistry() len = %d, want 1", len(skills))
	}

	if skills[0].ID != "vercel-labs/agent-browser/agent-browser" {
		t.Fatalf("FetchRegistry() id = %q, want %q", skills[0].ID, "vercel-labs/agent-browser/agent-browser")
	}
	if skills[0].Stars != 94800 {
		t.Fatalf("FetchRegistry() stars = %d, want %d", skills[0].Stars, 94800)
	}
}

func TestBuildRegistrySearchURLMigratesSkillsSHHost(t *testing.T) {
	t.Parallel()

	got, err := BuildRegistrySearchURL("https://api.skills.sh", "react")
	if err != nil {
		t.Fatalf("BuildRegistrySearchURL() error = %v", err)
	}

	want := "https://skills.sh/api/search?limit=100&q=react"
	if got != want {
		t.Fatalf("BuildRegistrySearchURL() = %q, want %q", got, want)
	}
}

func TestBuildRegistryBrowseURLMigratesSkillsSHHost(t *testing.T) {
	t.Parallel()

	got, err := BuildRegistryBrowseURL("https://api.skills.sh", "")
	if err != nil {
		t.Fatalf("BuildRegistryBrowseURL() error = %v", err)
	}

	want := "https://skills.sh"
	if got != want {
		t.Fatalf("BuildRegistryBrowseURL() = %q, want %q", got, want)
	}
}
