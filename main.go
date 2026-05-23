package main

import (
	"embed"
	"io/fs"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app, err := NewApp()
	if err != nil {
		panic(err)
	}
	distAssets, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		panic(err)
	}

	err = wails.Run(&options.App{
		Title:  "apex-custom-match-director",
		Width:  1180,
		Height: 760,
		AssetServer: &assetserver.Options{
			Assets: distAssets,
		},
		BackgroundColour:   &options.RGBA{R: 242, G: 245, B: 248, A: 1},
		LogLevel:           logger.DEBUG,
		LogLevelProduction: logger.DEBUG,
		Debug: options.Debug{
			OpenInspectorOnStartup: true,
		},
		OnStartup:  app.startup,
		OnDomReady: app.domReady,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		panic(err)
	}
}
