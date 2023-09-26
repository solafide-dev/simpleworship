package main

import (
	"context"
)

// App struct
type App struct {
	ctx       context.Context
	DataStore *DataStore
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.DataStore = NewDataStore(ctx)
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

func (a *App) GetDataStore() *DataStore {
	return a.DataStore
}

// Get song from DataStore.
// For some reason wails doesn't generate functions for structs that are not
func (a *App) GetSong(id string) (Song, error) {
	return a.DataStore.GetSong(id)
}

type Slide struct {
	Section string `json:"section"`
	Text    string `json:"text"`
}

type Meta struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

type SongSlide struct {
	Meta   Meta    `json:"meta"`
	Slides []Slide `json:"slides"`
}

func (a *App) LoadSong() SongSlide {
	return SongSlide{
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
