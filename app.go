package main

import (
	"context"
	"errors"
)

// App struct
type App struct {
	ctx       context.Context `json:"-"`
	DataStore *DataStore      `json:"dataStore"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.DataStore = &DataStore{}
	a.DataStore.init(ctx)
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

func (a *App) GetSongs() []Song {
	return a.DataStore.Songs
}

func (a *App) GetOrderOfServices() []OrderOfService {
	return a.DataStore.OrderOfServices
}

// Get song from DataStore.
func (a *App) GetSong(id string) (Song, error) {
	for _, song := range a.DataStore.Songs {
		if song.Id == id {
			return song, nil
		}
	}
	return Song{}, errors.New("song not found")
}

// Get order of service from DataStore.
func (a *App) GetOrderOfService(id string) (OrderOfService, error) {
	for _, service := range a.DataStore.OrderOfServices {
		if service.Id == id {
			return service, nil
		}
	}
	return OrderOfService{}, errors.New("service not found")
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
