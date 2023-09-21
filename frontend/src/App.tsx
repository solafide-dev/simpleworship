import React, { useContext, useEffect, useState } from 'react';
import { Playlist } from './components/Playlist';
import { SlideshowContext } from './context/SlideshowContext';
import { nextSlide, prevSlide } from './context/SlideshowContext/actions';
import { Greet } from "../wailsjs/go/main/App"

function App() {
    const [state, dispatch] = useContext(SlideshowContext)
    const [message, setMessage] = useState("")

    const handleKeyDown = (e: KeyboardEvent) => {
        if (e.code === 'Space' || e.code === 'ArrowRight') {
            e.preventDefault()
            dispatch(nextSlide())
        }
        if (e.code === 'Space' || e.code === 'ArrowLeft') {
            e.preventDefault()
            dispatch(prevSlide())
        }
    }

    useEffect(() => {
        (async () => {
            const message = await Greet('Michael')
            setMessage(message)
        })()

        window.addEventListener('keydown', handleKeyDown)

        return () => {
            window.removeEventListener('keydown', handleKeyDown);
        };
    })

    return (
        <div className="App" style={{ display: 'grid', gridTemplateColumns: '1fr 2fr' }}>
            {message}
            {/* <Playlist />
            <div>{state.prevSlide} | {state.activeSlide} | {state.nextSlide}<pre>{JSON.stringify(state, null, 2)}</pre></div> */}
        </div>
    );
}

export default App;
