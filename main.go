package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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

	AppMenu := menu.NewMenu()
	FileMenu := AppMenu.AddSubmenu("File")
	FileMenu.AddText("Open", keys.CmdOrCtrl("o"), nil)
	FileMenu.AddText("Settings", keys.CmdOrCtrl(","), nil)
	FileMenu.AddSeparator()
	FileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		runtime.Quit(app.ctx)
	})
	PluginsMenu := AppMenu.AddSubmenu("Plugins")
	PluginsMenu.AddSeparator()
	PluginsMenu.AddText("Manage Plugins", nil, nil)
	HelpMenu := AppMenu.AddSubmenu("Help")
	HelpMenu.AddText("User's Guide", nil, nil)

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "SimpleWorship",
		WindowStartState: options.Maximised,
		MaxWidth:         99999,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Menu:             AppMenu,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
