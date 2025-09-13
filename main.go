package main

import (
	"embed"
	"log"

	"todo-lits-DMARK/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	appInstance := app.NewApp()

	err := wails.Run(&options.App{
		Title:  "TodoApp - Task Management",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        appInstance.OnStartup,
		OnShutdown:       appInstance.OnShutdown,
		Bind: []interface{}{
			appInstance,
		},
	})

	if err != nil {
		log.Printf("Error starting application: %v", err)
	}
}
