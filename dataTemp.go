package main

// THIS FILE CONTAINS STRUCTS AND METHODS THAT ARE USED BY THE FRONTEND
// FOR TESTING PURPOSES. THIS FILE SHOULD EVENTUALLY BE REMOVED.

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
