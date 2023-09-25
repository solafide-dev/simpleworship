export type Slide = {
    id?: string
    text: string
}

export type Section = {
    id?: string
    name: string
    slides: Slide[]
}

export type Song = {
    name: string
    arrangement: string[]
    sections: Section[]
}

export type Playlist = {
    songs: Song[]
}