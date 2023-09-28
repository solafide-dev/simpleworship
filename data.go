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
		songs = append(songs, song.(Song))
	}
	return songs
}

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
		services = append(services, service.(OrderOfService))
	}
	return services
}

// Get song from DataStore.
func (a *App) GetSong(id string) (Song, error) {
	song, err := a.getData(Registry_Song, id)
	if err != nil {
		return Song{}, err
	}

	return song.(Song), nil
}

// Get order of service from DataStore.
func (a *App) GetOrderOfService(id string) (OrderOfService, error) {
	service, err := a.getData(Registry_OrderOfService, id)
	if err != nil {
		return OrderOfService{}, err
	}

	return service.(OrderOfService), nil
}
