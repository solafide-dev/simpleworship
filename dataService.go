package main

import (
	"encoding/json"
	"os"

	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

type ServiceItemType string

const (
	SongServiceItemType      ServiceItemType = "song"
	ScriptureServiceItemType ServiceItemType = "scripture"
	CategoryServiceItemType  ServiceItemType = "category"
)

// ServiceItem represents a single item in an Order of Service.
type ServiceItem struct {
	Title string          `json:"title"` // Title of the item
	Type  ServiceItemType `json:"type"`  // Song, Scripture, etc.
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

// OrderOfService represents a single Order of Service.
// This of this as a single "slide show" that contains everything for a single day's service.
type OrderOfService struct {
	DataFile
	Title string        `json:"title"` // Title of the Order of Service
	Date  string        `json:"date"`  // Date of the Order of Service
	Items []ServiceItem `json:"items"` // Items in the Order of Service
}

func (o *OrderOfService) ReloadFile() error {
	// open the file
	file, err := os.Open(o.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// decode the file
	err = json.NewDecoder(file).Decode(o)
	if err != nil {
		return err
	}

	return nil
}

// Save an OrderOfService file
// TODO: It would be nice is DataFile.SaveFile() could be used here with reflection somehow?
func (o *OrderOfService) SaveFile() error {
	// get the type of d
	rt.LogDebug(o.ctx, "[DATAFILE SAVE] Saving "+o.Filename)

	// open the file
	file, err := os.OpenFile(o.Filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	marshaled, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(marshaled)
	if err != nil {
		return err
	}

	rt.LogDebug(o.ctx, "[DATAFILE SAVE] Saved "+o.Filename)
	return nil
}
