package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// setup display server (this is very temporary and just for proof of concept)
	// go startDisplayServer()

	// Create an instance of the app structure
	app := NewApp()
	displayServer := NewDisplayServer()
	AppMenu := appMenu(app)

	// Create application with options
	err := wails.Run(&options.App{
		Title: "Simple Worship",
		//WindowStartState: options.Maximised,
		MaxWidth:         99999,
		Assets:           assets,
		Menu:             AppMenu,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			displayServer.startup(ctx)

		},
		OnShutdown: func(ctx context.Context) {
			app.shutdown(ctx)
			displayServer.shutdown(ctx)
		},
		OnDomReady: func(ctx context.Context) {
			app.domReady(ctx)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			return app.beforeClose(ctx)
		},
		Bind: []interface{}{
			app,
			displayServer,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
