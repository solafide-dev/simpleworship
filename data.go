package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/solafide-dev/august"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	Registry_Song = "songs"
	DataType_Song = "song"

	Registry_OrderOfService = "services"
	DataType_OrderOfService = "service"
)

type DataMutationEvent struct {
	Type     string `json:"type"`     // update, delete, create
	DataType string `json:"dataType"` // OrderOfService, Song
	Id       string `json:"id"`
}

type SimpleWorshipDataType struct {
	DataType string `json:"swdt"`
}

// Storage directory -- maybe this is configurable eventually.
var StorageDir = "./storage"

func (a *App) initAugust() {
	// Before we init August, lets verify any sort of storage exists.
	// While august can create a new storage folder, we embed a default one
	// for first runs of the application, so lets verify storage manually first.
	checkStorage(a.ctx)

	a.Data = august.Init()
	a.Data.Config(august.Config_StorageDir, StorageDir)

	//a.Data.Verbose() // Remove in production

	a.Data.SetEventFunc(func(event, store, id string) {
		data := DataMutationEvent{
			Type:     event,
			DataType: store,
			Id:       id,
		}
		rt.EventsEmit(a.ctx, "data-mutate", data)
	})

	// Register our data types
	a.Data.Register(Registry_Song, Song{})
	a.Data.Register(Registry_OrderOfService, OrderOfService{})
}

func (a *App) importFile(filename string) error {
	// Lets get our data and see if it matches any of our stores that we support importing
	data, err := os.ReadFile(filename)
	if err != nil {
		rt.LogError(a.ctx, err.Error())
		return err
	}

	detect := SimpleWorshipDataType{}
	err = a.Data.Unmarshal(data, &detect)
	if err != nil {
		rt.LogError(a.ctx, err.Error())
		return err
	}

	rt.LogInfo(a.ctx, "Detected import type: "+detect.DataType)

	switch detect.DataType {
	case DataType_Song:
		return a.importSong(data)
	case DataType_OrderOfService:
		return a.importOrderOfService(data)
	default:
		rt.LogError(a.ctx, "File is not a supported import type")
		return fmt.Errorf("file is not a supported import type")
	}
}

func (a *App) importSong(data []byte) error {
	song := Song{}
	err := a.Data.Unmarshal(data, &song)
	if err != nil {
		rt.LogError(a.ctx, err.Error())
		return err
	}

	if song.Id != "" {
		_, err := a.GetSong(song.Id)
		if err == nil {
			rt.LogWarning(a.ctx, "Song already exists")
			resp, err := rt.MessageDialog(a.ctx, rt.MessageDialogOptions{
				Type:    rt.QuestionDialog,
				Title:   "Overwrite Song?",
				Message: "A song with this ID already exists. Do you want to overwrite it?\n\n(Selecting no will create a new song with a new ID.)",
			})
			if err != nil {
				rt.LogError(a.ctx, err.Error())
				return err
			}
			if resp == "No" {
				song.Id = "" // Clear the ID so we create a new one
			}
		}
	}

	// Lets enforce some rules
	if song.Parts == nil || len(song.Parts) == 0 {
		rt.LogError(a.ctx, "Song does not have any parts")
		return fmt.Errorf("song does not have any parts")
	}

	if song.Title == "" {
		rt.LogError(a.ctx, "Song does not have a title")
		return fmt.Errorf("song does not have a title")
	}

	// We have a song, lets save it
	savedSongId, err := a.SaveSong(song)
	if err != nil {
		rt.LogError(a.ctx, err.Error())
		return err
	}

	rt.MessageDialog(a.ctx, rt.MessageDialogOptions{
		Type:    rt.InfoDialog,
		Title:   "Song Imported Successfully",
		Message: "Song Imported Successfully\n\nID: " + savedSongId + "\nTitle: " + song.Title + "",
	})

	return nil
}

func (a *App) importOrderOfService(data []byte) error {
	service := OrderOfService{}
	err := a.Data.Unmarshal(data, &service)
	if err != nil {
		rt.LogError(a.ctx, err.Error())
		return err
	}

	// Assume a service has items
	if service.Items == nil || len(service.Items) == 0 {
		rt.LogError(a.ctx, "Service does not have any songs")
		return fmt.Errorf("service does not have any songs")
	}

	if service.Title == "" {
		rt.LogError(a.ctx, "Service does not have a title")
		return fmt.Errorf("service does not have a title")
	}

	savedServiceId, err := a.SaveOrderOfService(service)
	if err != nil {
		rt.LogError(a.ctx, err.Error())
		return err
	}

	rt.MessageDialog(a.ctx, rt.MessageDialogOptions{
		Type:    rt.InfoDialog,
		Title:   "Order of Service Imported Successfully",
		Message: "Order of Service Imported Successfully\n\nID: " + savedServiceId + "\nTitle: " + service.Title + "",
	})

	return nil
}

func (a *App) getData(store, id string) (interface{}, error) {
	s, err := a.Data.GetStore(store)
	if err != nil {
		rt.LogError(a.ctx, err.Error())
		return nil, err
	}

	return s.Get(id)
}

// Get All Songs from DataStore.
func (a *App) GetSongs() []Song {
	store, err := a.Data.GetStore(Registry_Song)
	if err != nil {
		rt.LogError(a.ctx, err.Error())
	}

	songIds := store.GetIds()

	songs := []Song{}
	for _, id := range songIds {
		song, err := store.Get(id)
		if err != nil {
			rt.LogError(a.ctx, err.Error())
		}
		s := song.(Song)
		s.Id = id // make sure the ID always matches our august managed ID
		songs = append(songs, s)
	}
	return songs
}

// Get Song from DataStore.
func (a *App) GetSong(id string) (Song, error) {
	song, err := a.getData(Registry_Song, id)
	if err != nil {
		return Song{}, err
	}

	s := song.(Song)

	if s.Id != id {
		rt.LogWarning(a.ctx, "Song ID does not match. Updating ID in local storage")
		s.Id = id
		a.SaveSong(s)
	}

	s.Id = id // make sure the ID always matches our august managed ID

	return s, nil
}

// Save a song to the datastore.
func (a *App) SaveSong(song Song) (songId string, err error) {
	store, err := a.Data.GetStore(Registry_Song)
	if err != nil {
		rt.LogError(a.ctx, err.Error())
		return "", err
	}

	if song.Id == "" {
		rt.LogInfo(a.ctx, "Creating new song")
		newId, err := store.New(song)
		if err != nil {
			rt.LogError(a.ctx, err.Error())
			return "", err
		}
		song.Id = newId
	}

	return song.Id, store.Set(song.Id, song)
}

// Get All Order of Services from DataStore.
func (a *App) GetOrderOfServices() []OrderOfService {
	store, err := a.Data.GetStore(Registry_OrderOfService)
	if err != nil {
		rt.LogError(a.ctx, err.Error())
	}

	serviceIds := store.GetIds()

	services := []OrderOfService{}
	for _, id := range serviceIds {
		service, err := store.Get(id)
		if err != nil {
			rt.LogError(a.ctx, err.Error())
		}
		s := service.(OrderOfService)
		s.Id = id // make sure the ID always matches our august managed ID
		services = append(services, s)
	}
	return services
}

// Get order of service from DataStore.
func (a *App) GetOrderOfService(id string) (OrderOfService, error) {
	service, err := a.getData(Registry_OrderOfService, id)
	if err != nil {
		return OrderOfService{}, err
	}

	s := service.(OrderOfService)

	if s.Id != id {
		rt.LogWarning(a.ctx, "Order of Service ID does not match. Updating ID in local storage")
		s.Id = id
		a.SaveOrderOfService(s)
	}

	s.Id = id // make sure the ID always matches our august managed ID

	return s, nil
}

// Save an order of service to the datastore.
func (a *App) SaveOrderOfService(service OrderOfService) (serviceId string, err error) {
	store, err := a.Data.GetStore(Registry_OrderOfService)
	if err != nil {
		rt.LogError(a.ctx, err.Error())
		return "", err
	}

	if service.Id == "" {
		newId, err := store.New(service)
		if err != nil {
			rt.LogError(a.ctx, err.Error())
			return "", err
		}
		service.Id = newId
	}

	return service.Id, store.Set(service.Id, service)
}

//go:embed all:default_storage
var defaultStorage embed.FS

// checkStorage checks if the storage folder exists, and if not, creates it and copies the default data to it
// currently default data can only be 1 folder deep with this checking logic
func checkStorage(ctx context.Context) {
	// check if a storage folder exists next  to the executable
	if _, err := os.Stat(StorageDir); os.IsNotExist(err) {
		// if not, create it
		err := os.Mkdir(StorageDir, 0755)
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
				err := os.Mkdir(StorageDir+"/"+name, 0755)
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
	file, err := os.Create(StorageDir + "/" + name)
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
