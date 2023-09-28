package main

import (
	"github.com/solafide-dev/august"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	Registry_Song           = "songs"
	Registry_OrderOfService = "services"
)

func (a *App) initAugust() {
	a.Data = august.Init()

	a.Data.Verbose() // Remove in production

	// Register our data types
	a.Data.Register(Registry_Song, Song{})
	a.Data.Register(Registry_OrderOfService, OrderOfService{})

	a.Data.Run()
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
		s.Id = id
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
	s.Id = id

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
		return store.New(song)
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
		s.Id = id
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
	s.Id = id

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
		return store.New(service)
	}

	return service.Id, store.Set(service.Id, service)
}
