import React, { useContext, useEffect } from 'react';
import './App.css';
import { Playlist } from './components/Playlist';
import { SlideshowContext } from './context/SlideshowContext';
import { nextSlide, prevSlide } from './context/SlideshowContext/actions';

function App() {
    const [state, dispatch] = useContext(SlideshowContext)

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
        window.addEventListener('keydown', handleKeyDown)

        return () => {
            window.removeEventListener('keydown', handleKeyDown);
        };
    })

    return (
        <div className="App" style={{ display: 'grid', gridTemplateColumns: '1fr 2fr' }}>
            <Playlist />
            <div>{state.prevSlide} | {state.activeSlide} | {state.nextSlide}<pre>{JSON.stringify(state, null, 2)}</pre></div>
        </div>
    );
}

export default App;
