package main

import "strings"

// Song represents a single song.
type Song struct {
	Id          string     `json:"id"`          // Unique ID of the Song
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

// TODO: The two methods on Song are not properly moved to Wails typescript yet.
// Im not sure if this is a limitation of wails and we just need to arcetect differently,
// or if I am doing something wrong.

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
