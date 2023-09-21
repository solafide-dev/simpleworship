import React, { useContext, useEffect } from 'react'
import { SlideshowContext } from '../../context/SlideshowContext'
import { loadPlaylist, setActiveSlide } from '../../context/SlideshowContext/actions'
import sundayMorning from '../../data/playlists/sundayMorning'
import { Playlist as PlaylistType, Song } from '../../global'
import SongComponent from './Song'

const Playlist = () => {
    const [{ songs }, dispatch] = useContext(SlideshowContext)

    useEffect(() => {
        dispatch(loadPlaylist(sundayMorning))
    }, [])

    return (<div>
        {/* <pre>{JSON.stringify(songs, null, 2)}</pre> */}
        {songs && songs.map((song: Song) => <SongComponent {...song} />)}
    </div>)
}

export default Playlist