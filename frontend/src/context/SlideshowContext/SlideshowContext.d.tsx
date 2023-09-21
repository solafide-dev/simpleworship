import { Song } from "../../global.d"

export type SlideshowInitialState = {
    activeSong?: Song
    activeSlide?: string
    nextSlide?: string
    prevSlide?: string
    songs?: Song[]
}

export enum SlideshowActions {
    SET_ACTIVE_SLIDE = 'SET_ACTIVE_SLIDE',
    LOAD_PLAYLIST = 'LOAD_PLAYLIST',
    NEXT_SLIDE = 'NEXT_SLIDE',
    PREV_SLIDE = 'PREV_SLIDE'
}

export type SlideshowAction = {
    type: SlideshowActions,
    payload?: any
}