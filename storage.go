package main

import (
	"context"
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/simpleworship/gobible"
	"github.com/simpleworship/gobible/bible"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Storage handles the storage of songs and services, and provides loading of default data.

// OrderOfService represents a single Order of Service.
// This of this as a single "slide show" that contains everything for a single day's service.
type OrderOfService struct {
	Id    string        `json:"id"`    // UUID
	Title string        `json:"title"` // Title of the Order of Service
	Date  string        `json:"date"`  // Date of the Order of Service
	Items []ServiceItem `json:"items"` // Items in the Order of Service
}

type ServiceItemType string

const (
	SongServiceItemType      ServiceItemType = "song"
	ScriptureServiceItemType ServiceItemType = "scripture"
	CategoryServiceItemType  ServiceItemType = "category"
)

// ServiceItem represents a single item in an Order of Service.
type ServiceItem struct {
	Id    string          `json:"id"`    // UUID
	Title string          `json:"title"` // Title of the item
	Type  ServiceItemType `json:"type"`  // Song, Scripture, etc.
}

// Song represents a single song.
type Song struct {
	Id          string     `json:"id"`          // UUID
	Title       string     `json:"title"`       // Title of the Song
	Attribution string     `json:"attribution"` // Who wrote the song / who owns the song
	License     string     `json:"license"`     // License of the song
	Notes       string     `json:"notes"`       // Notes about the song
	Parts       []SongPart `json:"parts"`       // Parts of the song
	Order       []string   `json:"order"`       // Default song order
}

// SongPart represents a single part of a song (verse, chorus, etc.)
type SongPart struct {
	Id    string   `json:"id"`    // Verse 1, Chorus, etc.
	Title string   `json:"title"` // Title of the part
	Lines []string `json:"lines"` // The actual lyrics
}

// SongServiceItem represents a single song in an Order of Service.
type SongServiceItem struct {
	ServiceItem
	SongId string `json:"songId"` // UUID
}

// Scripture represents a single scripture.
// TODO: Integrate with GoBible
type ScriptureServiceItem struct {
	ServiceItem
	Reference string `json:"reference"` // Reference to the scripture
}

type DataStore struct {
	OrderOfServices []OrderOfService `json:"services"`
	Songs           []Song           `json:"songs"`
	Bibles          []*bible.Bible   `json:"bibles"`
}

//go:embed all:default_storage
var defaultStorage embed.FS

func InitDataStore(ctx context.Context) *DataStore {
	checkStorage(ctx)

	// load bibles from storage/bibles
	bibles := []*bible.Bible{}
	bibleFiles, err := os.ReadDir("./storage/bibles")
	if err != nil {
		rt.LogError(ctx, err.Error())
	}
	for _, bibleFile := range bibleFiles {
		bibleName := bibleFile.Name()
		// TODO: Update Gobible to add a loader function that does this
		if strings.HasSuffix(bibleName, ".json") {
			bible := gobible.New("./storage/bibles/" + bibleName)
			bibles = append(bibles, bible)
		}
		if strings.HasSuffix(bibleName, ".osis") {
			bible := gobible.NewOSIS("./storage/bibles/" + bibleName)
			bibles = append(bibles, bible)
		}
	}

	songFiles, err := os.ReadDir("./storage/songs")
	if err != nil {
		rt.LogError(ctx, err.Error())
	}
	songs := []Song{}
	for _, songFile := range songFiles {
		// parse each file as a song
		song := Song{}
		// open the file
		file, err := os.Open("./storage/songs/" + songFile.Name())
		if err != nil {
			rt.LogError(ctx, err.Error())
			continue
		}
		defer file.Close()

		// decode the file
		err = json.NewDecoder(file).Decode(&song)
		if err != nil {
			rt.LogError(ctx, err.Error())
			continue
		}

		// add the song to the list of songs
		songs = append(songs, song)
	}

	serviceFiles, err := os.ReadDir("./storage/services")
	if err != nil {
		rt.LogError(ctx, err.Error())
	}
	services := []OrderOfService{}
	for _, serviceFile := range serviceFiles {
		// parse each file as a service
		service := OrderOfService{}
		// open the file
		file, err := os.Open("./storage/services/" + serviceFile.Name())
		if err != nil {
			rt.LogError(ctx, err.Error())
			continue
		}
		defer file.Close()

		// decode the file
		err = json.NewDecoder(file).Decode(&service)
		if err != nil {
			rt.LogError(ctx, err.Error())
			continue
		}

		// add the service to the list of services
		services = append(services, service)
	}

	return &DataStore{
		OrderOfServices: services,
		Songs:           songs,
		Bibles:          bibles,
	}
}

// checkStorage checks if the storage folder exists, and if not, creates it and copies the default data to it
// currently default data can only be 1 folder deep with this checking logic
func checkStorage(ctx context.Context) {
	// check if a storage folder exists next  to the executable
	if _, err := os.Stat("./storage"); os.IsNotExist(err) {
		// if not, create it
		err := os.Mkdir("./storage", 0755)
		if err != nil {
			rt.LogError(ctx, err.Error())
		}

		entries, err := fs.ReadDir(defaultStorage, "default_storage")
		if err != nil {
			rt.LogError(ctx, err.Error())
		}

		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() {
				err := os.Mkdir("./storage/"+name, 0755)
				if err != nil {
					rt.LogError(ctx, err.Error())
				}
				subEntries, err := fs.ReadDir(defaultStorage, "default_storage/"+name)
				if err != nil {
					rt.LogError(ctx, err.Error())
				}
				for _, subEntry := range subEntries {
					subName := subEntry.Name()
					copyEmbededToStorage(ctx, name+"/"+subName)
				}
			} else {
				copyEmbededToStorage(ctx, name)
			}
		}
	}
}

func copyEmbededToStorage(ctx context.Context, name string) {
	rt.LogDebug(ctx, "Copying "+name+" to storage")
	// create the file
	file, err := os.Create("./storage/" + name)
	if err != nil {
		rt.LogError(ctx, err.Error())
	}
	defer file.Close()

	// open the embedded file
	embeddedFile, err := defaultStorage.Open("default_storage/" + name)
	if err != nil {
		rt.LogError(ctx, err.Error())
	}
	defer embeddedFile.Close()

	// copy the embedded file to the new file
	_, err = io.Copy(file, embeddedFile)
	if err != nil {
		rt.LogError(ctx, err.Error())
	}
}
