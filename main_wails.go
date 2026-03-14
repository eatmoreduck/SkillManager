//go:build wails

package main

import (
	"embed"
	"io/fs"
	"log"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var frontendAssets embed.FS

func main() {
	container, err := NewApp(os.Getenv("SKILLMANAGER_CONFIG"))
	if err != nil {
		log.Fatal(err)
	}

	distFS, err := fs.Sub(frontendAssets, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}

	app := application.New(application.Options{
		Name:        "SkillManager",
		Description: "Cross-platform AI Agent Skills Manager",
		Services: []application.Service{
			application.NewService(container.ConfigBinding),
			application.NewService(container.AgentBinding),
			application.NewService(container.RegistryBinding),
			application.NewService(container.SkillBinding),
			application.NewService(container.InventoryBinding),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(distFS),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:              "SkillManager",
		URL:                "/",
		Width:              1280,
		Height:             860,
		MinWidth:           1024,
		MinHeight:          720,
		DevToolsEnabled:    true,
		EnableFileDrop:     true,
		UseApplicationMenu: true,
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
