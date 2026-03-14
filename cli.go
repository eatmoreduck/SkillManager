package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"skillmanager/internal/model"
)

func runCLI(args []string, stdout io.Writer, stderr io.Writer) error {
	app, err := NewApp(os.Getenv("SKILLMANAGER_CONFIG"))
	if err != nil {
		return err
	}

	if len(args) == 0 {
		printUsage(stdout)
		return nil
	}

	switch args[0] {
	case "doctor":
		report, err := app.InventoryBinding.BuildReport()
		if err != nil {
			return err
		}
		printDoctor(stdout, report)
		return nil
	case "agents":
		report, err := app.InventoryBinding.BuildReport()
		if err != nil {
			return err
		}
		printAgents(stdout, report.Agents)
		return nil
	case "skills":
		report, err := app.InventoryBinding.BuildReport()
		if err != nil {
			return err
		}
		printSkills(stdout, report)
		return nil
	case "config":
		cfg, err := app.ConfigBinding.GetConfig()
		if err != nil {
			return err
		}
		printConfig(stdout, app, cfg)
		return nil
	case "help", "-h", "--help":
		printUsage(stdout)
		return nil
	default:
		printUsage(stderr)
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "SkillManager local diagnostics")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  go run . doctor   # full local report")
	fmt.Fprintln(w, "  go run . agents   # configured agent summary")
	fmt.Fprintln(w, "  go run . skills   # managed/external skill inventory")
	fmt.Fprintln(w, "  go run . config   # config summary")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Optional:")
	fmt.Fprintln(w, "  SKILLMANAGER_CONFIG=/path/to/config.yaml go run . doctor")
}

func printDoctor(w io.Writer, report *model.InventoryReport) {
	fmt.Fprintf(w, "Platform: %s\n", report.Platform)
	fmt.Fprintf(w, "Config:   %s\n", report.ConfigPath)
	fmt.Fprintf(w, "Store:    %s\n", report.ManagerSkillsDir)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Agents")
	printAgents(w, report.Agents)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Skills")
	printSkills(w, report)
}

func printAgents(w io.Writer, agents []model.Agent) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tNAME\tENABLED\tDETECTED\tCUSTOM\tSKILLS DIR")
	for _, agent := range agents {
		fmt.Fprintf(
			tw,
			"%s\t%s\t%s\t%s\t%s\t%s\n",
			agent.ID,
			agent.Name,
			boolLabel(agent.IsEnabled),
			boolLabel(agent.IsInstalled),
			boolLabel(agent.IsCustom),
			agent.SkillsDir,
		)
	}
	_ = tw.Flush()
}

func printSkills(w io.Writer, report *model.InventoryReport) {
	fmt.Fprintln(w, "SkillManager managed storage")
	if len(report.ManagedSkills) == 0 {
		fmt.Fprintln(w, "  (empty)")
	} else {
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "NAME\tASSIGNED AGENTS\tPATH")
		for _, skill := range report.ManagedSkills {
			agents := "-"
			if len(skill.Agents) > 0 {
				agents = strings.Join(skill.Agents, ", ")
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\n", skill.Name, agents, skill.LocalPath)
		}
		_ = tw.Flush()
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Agent-visible skills")
	for _, inventory := range report.AgentInventories {
		fmt.Fprintf(w, "\n[%s] %s\n", inventory.Agent.ID, inventory.Agent.Name)
		fmt.Fprintf(w, "  dir: %s\n", inventory.Agent.SkillsDir)
		fmt.Fprintf(w, "  enabled: %s | detected: %s\n", boolLabel(inventory.Agent.IsEnabled), boolLabel(inventory.Agent.IsInstalled))
		printInventoryGroup(w, "managed", inventory.Managed)
		printInventoryGroup(w, "external", inventory.External)
		printInventoryGroup(w, "broken", inventory.Broken)
	}
}

func printInventoryGroup(w io.Writer, title string, items []model.SkillInventoryItem) {
	fmt.Fprintf(w, "  %s:\n", title)
	if len(items) == 0 {
		fmt.Fprintln(w, "    (none)")
		return
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	for _, item := range items {
		line := fmt.Sprintf("    - %s", item.Name)
		if item.Description != "" {
			line += " | " + item.Description
		}
		line += " | " + item.Path
		if item.ResolvedPath != "" && item.ResolvedPath != item.Path {
			line += " -> " + item.ResolvedPath
		}
		fmt.Fprintln(w, line)
	}
}

func printConfig(w io.Writer, app *App, cfg *model.Config) {
	fmt.Fprintf(w, "Config path: %s\n", app.ConfigPath)
	fmt.Fprintf(w, "Version:     %s\n", cfg.Version)
	fmt.Fprintf(w, "Proxy:       %s\n", proxySummary(cfg.Proxy))
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Registries")
	if len(cfg.Registries) == 0 {
		fmt.Fprintln(w, "  (none)")
	} else {
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "ID\tNAME\tDEFAULT\tURL")
		for _, registry := range cfg.Registries {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", registry.ID, registry.Name, boolLabel(registry.IsDefault), registry.URL)
		}
		_ = tw.Flush()
	}
}

func boolLabel(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func proxySummary(proxy model.ProxyConfig) string {
	if !proxy.Enabled {
		return "off"
	}
	return proxy.URL()
}
