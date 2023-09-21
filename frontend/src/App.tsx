import React, { useContext, useEffect, useState } from 'react';
import { Playlist } from './components/Playlist';
import { SlideshowContext } from './context/SlideshowContext';
import { nextSlide, prevSlide } from './context/SlideshowContext/actions';
import { LoadSong } from "../wailsjs/go/main/App"
import { Song } from './global';
import { FaBook, FaBookOpen, FaGear, FaPencil, FaPlus } from 'react-icons/fa6'

function App() {
    const [state, dispatch] = useContext(SlideshowContext)
    const [song, setSong] = useState<any>()

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

    const handleClick = async () => {
        const song = await LoadSong()
        setSong(song)
    }

    useEffect(() => {
        window.addEventListener('keydown', handleKeyDown)

        return () => {
            window.removeEventListener('keydown', handleKeyDown);
        };
    }, [])

    return (
        <div className='bg-slate-900 h-screen p-3 flex flex-col gap-2'>
            <div className='flex gap-2 overflow-hidden'>
                <div className='w-[350px] flex flex-col gap-2'>
                    <div>
                        <div className='bg-black rounded w-full aspect-video'></div>
                    </div>
                    <div className='bg-slate-800 flex-1 rounded flex flex-col overflow-hidden'>
                        <div className='py-5 px-6'>
                            <div className="flex justify-between mb-4">
                                <h3 className='text-lg text-slate-300 font-bold flex items-center gap-2'><FaBookOpen /> Library</h3>
                                <button className='text-slate-400'><FaPlus /></button>
                            </div>
                            <div className='gap-2 flex'>
                                <div className="text-xs bg-slate-600 text-slate-300 px-3 py-2 inline-block text-center rounded-full">Playlists</div>
                                <div className="text-xs bg-slate-600 text-slate-300 px-3 py-2 inline-block text-center rounded-full">Songs</div>
                                <div className="text-xs bg-slate-600 text-slate-300 px-3 py-2 inline-block text-center rounded-full">Videos</div>
                            </div>
                        </div>
                        <div className="flex-1 rounded-b flex flex-col overflow-scroll">
                            {Array(10).fill(null).map((_, i) => {
                                return <div className={`px-6 py-3 ${i == 0 && 'bg-slate-700'}`}>
                                    <p className={`text-slate-300 ${i == 0 && 'text-indigo-400'}`}>Sunday Morning</p>
                                    <p className={`text-xs ${i == 0 ? 'text-slate-400' : 'text-slate-500'}`}>Playlist • Sunday, September 10, 2023</p>
                                </div>
                            })}
                        </div>
                    </div>
                </div>
                <div className='bg-slate-800 flex-1 rounded p-8 overflow-auto'>
                    <div className='flex justify-between items-center mb-4'>
                        <h1 className='text-3xl text-slate-300 font-bold top-0'>Sunday Morning</h1>
                        <button className='text-slate-400'><FaPencil /></button>
                    </div>
                    {Array(4).fill(null).map((_, i) => {
                        return (
                            <div className='mb-12'>
                                <h2 className={`text-xl text-slate-300 mb-3  top-14 ${i == 0 && 'text-indigo-400'}`}>There Is A Redeemer</h2>
                                <div className='grid grid-cols-3 gap-2'>
                                    {Array(7).fill(null).map((_, j) => {
                                        return <div className={`bg-slate-700 rounded ${i == 0 && j == 0 ? 'bg-indigo-600' : ''}`}>
                                            <div className={`text-xs px-1.5 pt-1 rounded-t text-bold ${i == 0 ? 'text-indigo-200' : 'text-indigo-400'}`}>Verse 1</div>
                                            <div className='p-1'>
                                                <div className='bg-black rounded w-full aspect-video'></div>
                                            </div>
                                        </div>
                                    })}
                                </div>
                            </div>
                        )
                    })}
                </div>
            </div>
            <div className='py-3 flex flex-1 items-center justify-end'>
                <button className="text-slate-400"><FaGear /></button>
            </div>
            {/* <button onClick={handleClick} className='bg-blue-600 text-white px-4 py-2 rounded'>Load a Song</button>
            {song && <div>
                <h1>{song.meta.title}</h1>
                <p>{song.meta.artist}</p>
                {song.slides.map(slide => <p>
                    <small>{slide.section}<br /></small>
                    {slide.text.split('\n').map(line => <>{line}<br /></>)}
                </p>)}
            </div>}
            {/* <Playlist />
            <div>{state.prevSlide} | {state.activeSlide} | {state.nextSlide}<pre>{JSON.stringify(state, null, 2)}</pre></div> */}
        </div>
    );
}

export default App;
