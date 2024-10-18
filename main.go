package main

import (
	"embed"
	"github.com/danbai225/gpp/backend/config"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"net"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/windows/icon.ico
var logo []byte

func main() {
	dial, err := net.Dial("tcp", "127.0.0.1:54713")
	if err == nil {
		_, _ = dial.Write([]byte("SHOW_WINDOW"))
		_ = dial.Close()
		return
	}
	config.InitConfig()
	// Create an instance of the app structure
	app := NewApp()
	defer app.Stop()

	// Create application with options
	err = wails.Run(&options.App{
		Title:             "gpp",
		Width:             360,
		Height:            520,
		DisableResize:     true,
		HideWindowOnClose: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 0},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
