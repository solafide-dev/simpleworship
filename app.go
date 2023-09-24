package main

import (
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's YO time!", name)
}

type Slide struct {
	Section string `json:"section"`
	Text    string `json:"text"`
}

type Meta struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

type Song struct {
	Meta   Meta    `json:"meta"`
	Slides []Slide `json:"slides"`
}

func (a *App) LoadSong() Song {
	return Song{
		Meta: Meta{
			Title:  "There Is a Redeemer",
			Artist: "Keith Green",
		},
		Slides: []Slide{
			{
				Section: "verse 1",
				Text:    "There is a redeemer\nJesus, God's own Son",
			},
			{
				Section: "verse 1",
				Text:    "Precious Lamb of God, Messiah\nHoly One",
			},
		},
	}
}
