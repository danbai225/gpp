package main

import (
	"embed"
	"github.com/getlantern/elevate"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"os"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if len(os.Args) == 1 {
		command := elevate.Command(os.Args[0], "sudo")
		command.Stderr = os.Stderr
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		_ = command.Run()
		os.Exit(0)
	}
	// Create an instance of the app structure
	app := NewApp()
	defer app.Stop()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "gpp",
		Width:  360,
		Height: 480,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
