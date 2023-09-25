package main

import (
	"context"
	"embed"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
// var icon []byte

func main() {
	// setup display server (this is very temporary and just for proof of concept)
	// go startDisplayServer()

	// Create an instance of the app structure
	app := NewApp()
	displayServer := NewDisplayServer()

	AppMenu := menu.NewMenu()
	FileMenu := AppMenu.AddSubmenu("File")
	FileMenu.AddText("Open", keys.CmdOrCtrl("o"), nil)
	FileMenu.AddText("Settings", keys.CmdOrCtrl(","), nil)
	FileMenu.AddSeparator()
	FileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		rt.Quit(app.ctx)
	})
	PluginsMenu := AppMenu.AddSubmenu("Plugins")
	PluginsMenu.AddSeparator()
	PluginsMenu.AddText("Manage Plugins", nil, nil)
	HelpMenu := AppMenu.AddSubmenu("Help")
	HelpMenu.AddText("User's Guide", nil, nil)
	if runtime.GOOS == "darwin" {
		AppMenu.Append(menu.EditMenu()) // on macos platform, we should append EditMenu to enable Cmd+C,Cmd+V,Cmd+Z... shortcut
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "Simple Worship",
		WindowStartState: options.Maximised,
		MaxWidth:         99999,
		Assets:           assets,
		Menu:             AppMenu,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.Startup(ctx)
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
