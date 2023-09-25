import { Reducer } from 'react'
import { Section, Slide, Song } from '../../global.d'
import { SlideshowAction, SlideshowActions, SlideshowInitialState } from './SlideshowContext.d'

export const SlideshowReducer: Reducer<SlideshowInitialState, SlideshowAction> = (state, action): SlideshowInitialState => {
    switch (action.type) {
        case SlideshowActions.SET_ACTIVE_SLIDE: {
            if (!state.activeSong) return state
            const { activeSlide, nextSlide, prevSlide } = setActiveSlide(state.activeSong, action.payload)
            return { ...state, activeSlide, nextSlide, prevSlide }
        }

        case SlideshowActions.NEXT_SLIDE: {
            if (!state.activeSong || !state.nextSlide) return state
            const { activeSlide, nextSlide, prevSlide } = setActiveSlide(state.activeSong, state.nextSlide)
            return { ...state, activeSlide, nextSlide, prevSlide }
        }

        case SlideshowActions.PREV_SLIDE: {
            if (!state.activeSong || !state.prevSlide) return state
            const { activeSlide, nextSlide, prevSlide } = setActiveSlide(state.activeSong, state.prevSlide)
            return { ...state, activeSlide, nextSlide, prevSlide }
        }

        case SlideshowActions.LOAD_PLAYLIST: {
            const songs = action.payload.songs.map((song: Song, i: number) => {
                const sections = song.arrangement.map((sectionId, j) => {
                    const section: any = song.sections.find(section => section.id === sectionId)
                    const slides = section.slides.map((slide: Slide, k: number) => {
                        return { ...slide, id: `${i}_${j}_${k}` }
                    })
                    return { ...section, slides, id: `${i}_${j}` }
                })
                return { ...song, sections, id: `${i}` }
            })

            const { activeSlide, nextSlide, prevSlide } = setActiveSlide(songs[0], '0_0_0')

            return {
                ...state,
                activeSong: songs[0],
                activeSlide,
                nextSlide,
                prevSlide,
                songs
            }
        }
    }
}

const setActiveSlide = (song: Song, activeSlide: string) => {
    let section = 0
    let slide = 0
    let nextSlide
    let prevSlide

    do {
        const currentSlide = song.sections[section].slides[slide]

        if (slide + 1 < song.sections[section].slides.length) {
            slide++
        } else {
            section++
            slide = 0
        }

        if (currentSlide.id === activeSlide) {
            const slides = song.sections[section].slides

            if (slide < slides.length) {
                nextSlide = slides[slide].id
            } else {
                nextSlide = song.sections[section + 1].slides[0].id
            }

            if (section > 0) {
                const slides = song.sections[section - 1].slides
                prevSlide = slides[slides.length - 1].id
            } else {
                if (slide - 2 > 0) {
                    prevSlide = slides[slide - 2].id
                    console.log('here')
                } else {
                    prevSlide = '0_0_0'
                }
            }
        }
    } while (!nextSlide)

    return { activeSlide, nextSlide, prevSlide }
}
