package main

import (
	"runtime"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Application Menu Definition, included in the main.go file

func appMenu(app *App) *menu.Menu {

	AppMenu := menu.NewMenu()

	// File Menu

	FileMenu := AppMenu.AddSubmenu("File")
	FileMenu.AddText("Open", keys.CmdOrCtrl("o"), nil)
	FileMenu.AddText("Settings", keys.CmdOrCtrl(","), nil)
	FileMenu.AddSeparator()
	FileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		rt.Quit(app.ctx)
	})

	// Plugin Menu

	PluginsMenu := AppMenu.AddSubmenu("Plugins")
	PluginsMenu.AddSeparator()
	PluginsMenu.AddText("Manage Plugins", nil, func(cd *menu.CallbackData) {
		rt.MessageDialog(app.ctx, rt.MessageDialogOptions{
			Type:    rt.InfoDialog,
			Title:   "Plugins",
			Message: "Plugins are not currently supported in this version of Simple Worship.",
		})
	})

	// Library

	LibraryMenu := AppMenu.AddSubmenu("Library")
	LibraryMenu.AddText("Songs", nil, func(cd *menu.CallbackData) {
		rt.MessageDialog(app.ctx, rt.MessageDialogOptions{
			Type:    rt.InfoDialog,
			Title:   "Library",
			Message: "This will open the song library for management.... WHEN WE HAVE ONE.",
		})
	})
	LibraryMenu.AddText("Bibles", nil, func(cd *menu.CallbackData) {
		rt.MessageDialog(app.ctx, rt.MessageDialogOptions{
			Type:    rt.InfoDialog,
			Title:   "Library",
			Message: "This will open the bible library for management.... WHEN WE HAVE ONE.",
		})
	})
	LibraryMenu.AddSeparator()
	LibraryMenu.AddText("Import...", nil, func(cd *menu.CallbackData) {
		files, err := rt.OpenMultipleFilesDialog(app.ctx, rt.OpenDialogOptions{
			Title: "Import",
			Filters: []rt.FileFilter{
				{
					DisplayName: "Simple Worship Song / GoBible File (*.json)",
					Pattern:     "*.json",
				},
			},
		})
		if err != nil {
			rt.LogInfo(app.ctx, "Import cancelled")
			rt.LogError(app.ctx, err.Error())
			return
		}

		for _, file := range files {
			if err := app.importFile(file); err != nil {
				rt.LogError(app.ctx, err.Error())
				rt.MessageDialog(app.ctx, rt.MessageDialogOptions{
					Type:    rt.ErrorDialog,
					Title:   "Import Error",
					Message: err.Error(),
				})
			}
		}
	})

	// Display

	DisplayMenu := AppMenu.AddSubmenu("Display")
	DisplayMenu.AddText("Open in Browser", nil, func(cd *menu.CallbackData) {
		rt.BrowserOpenURL(app.ctx, "http://localhost:7777")
	})

	// Help Menu

	HelpMenu := AppMenu.AddSubmenu("Help")
	HelpMenu.AddText("User's Guide", nil, func(cd *menu.CallbackData) {
		rt.BrowserOpenURL(app.ctx, "https://solafide.dev/simpleworship/")
	})

	// Extras

	if runtime.GOOS == "darwin" {
		AppMenu.Append(menu.EditMenu()) // on macos platform, we should append EditMenu to enable Cmd+C,Cmd+V,Cmd+Z... shortcut
	}

	return AppMenu
}
