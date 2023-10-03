package main

type ServiceItemType string

const (
	SongServiceItemType      ServiceItemType = "song"
	ScriptureServiceItemType ServiceItemType = "scripture"
	CategoryServiceItemType  ServiceItemType = "category"
)

// ServiceItem represents a single item in an Order of Service.
type ServiceItem struct {
	Title string          `json:"title"`          // Title of the item
	Type  ServiceItemType `json:"type"`           // Song, Scripture, etc.
	Meta  ServiceItemMeta `json:"meta,omitempty"` // Meta data for the item
}

type ServiceItemMeta struct {
	// Song Meta
	SongId string `json:"songId,omitempty"`
	// Scripture Meta
	VerseReference string `json:"verseReference,omitempty"`
}

// SongServiceItem represents a single song in an Order of Service
// OrderOfService represents a single Order of Service.
// This of this as a single "slide show" that contains everything for a single day's service.
type OrderOfService struct {
	SimpleWorshipDataType
	Id    string        `json:"id"`    // Unique ID of the item
	Title string        `json:"title"` // Title of the Order of Service
	Date  string        `json:"date"`  // Date of the Order of Service
	Items []ServiceItem `json:"items"` // Items in the Order of Service
}
