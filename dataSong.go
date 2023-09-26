package main

import (
	"encoding/json"
	"os"
	"strings"

	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Song represents a single song.
type Song struct {
	DataFile
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
	Lines []string `json:"lines"` // The actual lyrics
}

func (s *Song) ReloadFile() error {
	// open the file
	file, err := os.Open(s.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// decode the file
	err = json.NewDecoder(file).Decode(s)
	if err != nil {
		return err
	}

	return nil
}

// Save a Song file
// TODO: It would be nice is DataFile.SaveFile() could be used here with reflection somehow?
func (s *Song) SaveFile() error {
	// get the type of d
	rt.LogDebug(s.ctx, "[DATAFILE SAVE] Saving "+s.Filename)

	// open the file
	file, err := os.OpenFile(s.Filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	marshaled, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(marshaled)
	if err != nil {
		return err
	}

	dataMutationEvent(s.ctx, "update", "Song", s.Id)

	rt.LogDebug(s.ctx, "[DATAFILE SAVE] Saved "+s.Filename)
	return nil
}

// Generate the SongML for a song.
//
// [Verse 1]
// This is the first line of the first verse.
// This is the second line of the first verse.
//
// [Chorus]
// This is the first line of the chorus.
// This is the second line of the chorus.
//
// This is another line that would take another slide
// in the chorus.
//
// [Verse 2]
// This is the first line of the second verse.
// This is the second line of the second verse.
func (s *Song) ToSongML() string {
	songML := ""
	for _, part := range s.Parts {
		songML += "[" + part.Id + "]\n"
		for _, line := range part.Lines {
			songML += line + "\n\n"
		}
		songML += "\n"
	}
	return songML
}

// Update the song from SongML.
func (s *Song) UpdateFromSongML(songml string) {
	// parse line by line
	// if line starts with [ then it's a part
	// if line starts with anything else, it's a line

	parts := []SongPart{}
	part := SongPart{}
	l := ""
	lines := []string{}
	split := strings.Split(songml, "\n")
	for _, line := range split {
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// this is a part start
			if line != "" {
				// save the previous part
				lines = append(lines, l)
				part.Lines = lines
				parts = append(parts, part)
			}
			part = SongPart{
				Id:    line[1 : len(line)-1],
				Lines: []string{},
			}
		} else if strings.TrimSpace(line) == "" {
			// This is a break line, start a new part
			lines = append(lines, l)
		} else {
			l = strings.TrimSpace(l) + "\n"
		}
	}
	// save the last part
	part.Lines = lines
	parts = append(parts, part)
	s.Parts = parts
}
