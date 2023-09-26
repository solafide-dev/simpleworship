package main

import (
	"context"
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

// This file handles all the main data storage and loading for the application.
// Some of the types may be defined in other files.

// Datafile is the base for the different types of data files
// It contains the ID and the filename of the file

type DataFile struct {
	ctx      context.Context `json:"-"`
	typeName string          `json:"-"`
	Id       string          `json:"id"` // UUID
	Filename string          `json:"-"`  // Filename of the file
}

type DataStore struct {
	ctx             context.Context
	OrderOfServices []OrderOfService `json:"services"`
	Songs           []Song           `json:"songs"`
	//Bibles          []*bible.Bible   `json:"bibles"`
}

type DataMutationEvent struct {
	Type     string `json:"type"`     // update, delete, create
	DataType string `json:"dataType"` // OrderOfService, Song
	Id       string `json:"id"`
}

func dataMutationEvent(ctx context.Context, t, dataType, id string) {
	rt.LogDebug(ctx, "[DATA MUTATION EVENT] "+t+" "+dataType+" "+id)
	data := DataMutationEvent{
		Type:     t,
		DataType: dataType,
		Id:       id,
	}
	rt.EventsEmit(ctx, "data-mutate", data)
}

func (d *DataStore) HandleFileDelete(t string, filename string) {
	rt.LogDebug(d.ctx, "[DATAFILE DELETED] "+t+" file "+filename+" was deleted. Removing its data.")

	switch t {
	case `OrderOfService`:
		for i, service := range d.OrderOfServices {
			if service.Filename == filename {
				d.OrderOfServices = append(d.OrderOfServices[:i], d.OrderOfServices[i+1:]...)
				dataMutationEvent(d.ctx, "delete", t, service.Id)
				return
			}
		}
	case `Song`:
		for i, song := range d.Songs {
			if song.Filename == filename {
				d.Songs = append(d.Songs[:i], d.Songs[i+1:]...)
				dataMutationEvent(d.ctx, "delete", t, song.Id)
				return
			}
		}
	}
}

func (d *DataStore) LoadDataFile(t string, filename string) error {

	// Open the file
	rt.LogDebug(d.ctx, "[DATAFILE LOAD] Loading "+t+" "+filename)
	file, err := os.ReadFile(filename)
	if err != nil {
		rt.LogError(d.ctx, "[DATAFILE LOAD ERROR] "+err.Error())
		return err
	}

	// add the data to the correct slice
	// There is probably a fancy way of doing this with relfection
	// where we don't need a switch statement, but I am a tired idiot
	switch t {
	case `OrderOfService`:
		newData := OrderOfService{}
		newData.ctx = d.ctx
		newData.typeName = t
		err := json.Unmarshal(file, &newData)
		if err != nil {
			rt.LogError(d.ctx, "[DATAFILE UNMARSHAL ERROR] "+err.Error())
			return err
		}
		newData.Filename = filename
		if newData.Id == "" {
			newData.Id = uuid.New().String()
			err = newData.SaveFile()
			if err != nil {
				rt.LogError(d.ctx, "[DATAFILE SAVE ERROR] "+err.Error())
				return err
			}
		}

		for i, service := range d.OrderOfServices {
			if service.Id == newData.Id {
				d.OrderOfServices[i] = newData
				dataMutationEvent(d.ctx, "update", "OrderOfService", newData.Id)
				return nil
			}
		}
		dataMutationEvent(d.ctx, "create", "OrderOfService", newData.Id)
		d.OrderOfServices = append(d.OrderOfServices, newData)
	case `Song`:
		newData := Song{}
		newData.ctx = d.ctx
		newData.typeName = t
		err := json.Unmarshal(file, &newData)
		if err != nil {
			rt.LogError(d.ctx, "[DATAFILE UNMARSHAL ERROR] "+err.Error())
			return err
		}
		newData.Filename = filename
		if newData.Id == "" {
			newData.Id = uuid.New().String()
			err = newData.SaveFile()
			if err != nil {
				rt.LogError(d.ctx, "[DATAFILE SAVE ERROR] "+err.Error())
				return err
			}
		}
		for i, song := range d.Songs {
			if song.Id == newData.Id {
				d.Songs[i] = newData
				dataMutationEvent(d.ctx, "update", "Song", newData.Id)
				return nil
			}
		}
		dataMutationEvent(d.ctx, "create", "Song", newData.Id)
		d.Songs = append(d.Songs, newData)
	}

	return nil
}

func (d *DataStore) loadOrderOfServiceData() {
	storage := "./storage/services"
	serviceFiles, err := os.ReadDir(storage)
	if err != nil {
		rt.LogError(d.ctx, err.Error())
	}
	d.OrderOfServices = []OrderOfService{}
	for _, serviceFile := range serviceFiles {
		d.LoadDataFile(`OrderOfService`, storage+"/"+serviceFile.Name())
	}
}

func (d *DataStore) loadSongData() {
	storage := "./storage/songs"
	songFiles, err := os.ReadDir(storage)
	if err != nil {
		rt.LogError(d.ctx, err.Error())
	}
	d.Songs = []Song{}
	for _, songFile := range songFiles {
		d.LoadDataFile(`Song`, storage+"/"+songFile.Name())
	}
}

/*func (d *DataStore) loadBibleData() {
	// load bibles from storage/bibles
	bibles := []*bible.Bible{}
	bibleFiles, err := os.ReadDir("./storage/bibles")
	if err != nil {
		rt.LogError(d.ctx, err.Error())
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
	d.Bibles = bibles
}*/

var WatchDirectories = map[string]string{
	"OrderOfService": "./storage/services",
	"Song":           "./storage/songs",
	//"./storage/bibles",
}

// monitorFiles monitors the data files for changes and reloads them
// TODO: Make this actually work. Currently it just logs when a file is changed.
func (d *DataStore) monitorFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		rt.LogError(d.ctx, err.Error())
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//rt.LogDebugf(d.ctx, "event: %v", event)

				filename := strings.Replace(event.Name, "\\", "/", -1) // windows support

				dataType := getFileDataType(filename)
				if dataType == "" {
					continue
				}

				if event.Has(fsnotify.Rename) || event.Has(fsnotify.Remove) {
					// rename events are followed up by create events, so find the
					d.HandleFileDelete(dataType, filename)
				}
				if event.Has(fsnotify.Write) {
					//log.Println("modified file:", filename)
					rt.LogInfof(d.ctx, "[DATASTORE] Reloading data from %s", filename)
					d.LoadDataFile(dataType, filename)
				}
				if event.Has(fsnotify.Create) {
					//log.Println("created file:", filename)
					rt.LogInfof(d.ctx, "[DATASTORE] Loading new data from %s", filename)
					d.LoadDataFile(dataType, filename)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				rt.LogError(d.ctx, err.Error())
			}
		}
	}()

	for _, directory := range WatchDirectories {
		// Add a path.
		err = watcher.Add(directory)
		if err != nil {
			rt.LogError(d.ctx, err.Error())
		}
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

func (d *DataStore) init(ctx context.Context) {
	d.ctx = ctx
	checkStorage(ctx)
	d.loadOrderOfServiceData()
	d.loadSongData()
	//d.loadBibleData(ctx)

	go d.monitorFiles()
}

func getFileDataType(filename string) string {
	// if filename contains "services" return "OrderOfService"
	if strings.Contains(filename, "/services/") {
		return "OrderOfService"
	}
	// if filename contains "songs" return "Song"
	if strings.Contains(filename, "/songs/") {
		return "Song"
	}
	// if filename contains "bibles" return "Bible"
	if strings.Contains(filename, "/bibles/") {
		return "Bible"
	}
	return ""
}

//go:embed all:default_storage
var defaultStorage embed.FS

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
