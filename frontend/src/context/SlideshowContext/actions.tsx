import { Playlist } from "../../global";
import { SlideshowAction, SlideshowActions } from "./SlideshowContext.d";

export const setActiveSlide = (id: string): SlideshowAction => {
    return {
        type: SlideshowActions.SET_ACTIVE_SLIDE,
        payload: id
    }
}

export const nextSlide = (): SlideshowAction => {
    return {
        type: SlideshowActions.NEXT_SLIDE
    }
}

export const prevSlide = (): SlideshowAction => {
    return {
        type: SlideshowActions.PREV_SLIDE
    }
}

export const loadPlaylist = (playlist: Playlist): SlideshowAction => {
    return {
        type: SlideshowActions.LOAD_PLAYLIST,
        payload: playlist
    }
}