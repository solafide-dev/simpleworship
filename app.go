package main

import (
	"context"

	"github.com/solafide-dev/august"
	"github.com/solafide-dev/gobible"
)

// App struct
type App struct {
	ctx   context.Context  `json:"-"`
	Data  *august.August   `json:"-"`
	Bible *gobible.GoBible `json:"-"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.initAugust() // Initialize our data storage system
	a.initBible()  // Initialize our bible data system (will leverage august)
	a.Data.Run()   // Run Agust now that its all setup
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}
