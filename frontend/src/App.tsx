import React, { useContext, useEffect, useState } from 'react';
import {main} from "../wailsjs/go/models"
import { Playlist } from './components/Playlist';
import { SlideshowContext } from './context/SlideshowContext';
import { nextSlide, prevSlide } from './context/SlideshowContext/actions';
import { LoadSong, GetDataStore, GetSong } from "../wailsjs/go/main/App"
import { Song } from './global';
import { RiSlideshow4Line, RiStackLine, RiBookOpenLine, RiSettings5Line, RiPencilLine, RiAddLine, RiMusic2Line, RiVideoLine, RiMenuLine, RiCloseCircleFill, RiIndeterminateCircleFill, RiSkipForwardLine, RiSkipBackLine } from 'react-icons/ri'
import Display from './Display';

function App() {
    const [state, dispatch] = useContext(SlideshowContext)
    const [song, setSong] = useState<any>()
    const [isEditing, setIsEditing] = useState(false)
    const [lyrics, setLyrics] = useState("There is a Redeemer\nJesus, God's own son.")

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

    /* garbage test code */
    /*GetDataStore().then((dataStore: main.DataStore) => {
        console.log(dataStore)
        dataStore.GetSong("o-come-all-ye-faithful").then((song: Song) => {
            console.log(song)
        })
    })*/
    GetSong("o-come-all-ye-faithful").then((song: main.Song) => {
        console.log(song)
    })

    /* end garbage test code */

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

    useEffect(() => {
        setTimeout(() => {
            if (lyrics === "There is a Redeemer\nJesus, God's own son.") {
                setLyrics("Precious lamb of God, Messiah\nHoly one")
            } else {
                setLyrics("There is a Redeemer\nJesus, God's own son.")
            }
        }, 3000)
    })

    return (
        <div className='bg-gray-900 h-screen p-3 flex flex-col gap-2'>
            <div className='flex gap-2 overflow-hidden'>
                <div className='w-[350px] flex flex-col gap-2'>
                    <div>
                        <div className='bg-black rounded w-full aspect-video overflow-hidden'>
                            <Display lyrics={lyrics} />
                        </div>
                    </div>
                    <div className='flex gap-2'>
                        {["Blank", "Logo"].map(button => {
                            return <button className='bg-gray-600 px-3 py-1.5 text-sm text-gray-300 hover:bg-cyan-600 hover:text-cyan-100 rounded flex-1 font-semibold'>{button}</button>
                        })}
                    </div>
                    <div className='bg-gray-800 flex-1 rounded flex flex-col overflow-hidden group'>
                        <div className='py-5 px-6'>
                            <div className="flex justify-between items-center">
                                <h3 className='text-lg text-gray-300 flex items-center gap-2 font-semibold'><RiStackLine /> Sunday Morning</h3>
                                <button className='p-2 -m-2 text-gray-400 opacity-0 pointer-events-none group-hover:opacity-100 group-hover:pointer-events-auto' onClick={() => setIsEditing(!isEditing)}><RiPencilLine /></button>
                            </div>
                        </div>
                        <div className="flex-1 rounded-b flex flex-col overflow-scroll">
                            {[
                                {
                                    icon: <RiSlideshow4Line />,
                                    label: "Pre-Service Announcements",
                                    sublabel: "Announcement"
                                },
                                {
                                    icon: <RiMusic2Line />,
                                    label: "There Is A Redeemer",
                                    sublabel: "Song • Keith Green"
                                },
                                {
                                    icon: <RiVideoLine />,
                                    label: "Sermon Bumper",
                                    sublabel: "Video • 2m 43s"
                                },
                                {
                                    icon: <RiBookOpenLine />,
                                    label: "The Surprising Power of Grace",
                                    sublabel: "Sermon"
                                },
                            ].map(({ icon, label, sublabel }, i) => {
                                return (
                                    <div className={`px-6 py-2 relative ${i == 1 && 'bg-gray-700'}`}>
                                        <RiMenuLine className={`text-gray-400 absolute top-[50%] -mt-2 transition-all duration-75 ease-out ${isEditing ? 'left-6 opacity-100 pointer-events-auto' : 'left-0 opacity-0 pointer-events-none'}`} />
                                        <div className={`flex-1 transition-all duration-75 ease-out ${isEditing && 'pl-8'}`}>
                                            <p className={`${i == 1 ? 'text-cyan-400' : 'text-gray-400'} text-sm font-medium flex gap-1 items-center`}>{icon} {label}</p>
                                            <p className={`text-xs ${i == 1 ? 'text-gray-100' : 'text-gray-500'}`}>{sublabel}</p>
                                        </div>
                                        <RiIndeterminateCircleFill className={`text-rose-500 absolute top-[50%] -mt-2 transition-all duration-75 ease-out ${isEditing ? 'right-6 opacity-100 pointer-events-auto' : 'right-0 opacity-0 pointer-events-none'}`} />
                                    </div>
                                )
                            })}
                        </div>
                    </div>

                    {false && <div className='bg-gray-800 flex-1 rounded flex flex-col overflow-hidden'>
                        <div className='py-5 px-6'>
                            <div className="flex justify-between mb-4">
                                <h3 className='text-lg text-gray-300 font-bold flex items-center gap-2'><RiBookOpenLine /> Library</h3>
                                <button className='text-gray-400'><RiAddLine /></button>
                            </div>
                            <div className='gap-2 flex'>
                                <div className="text-xs bg-gray-600 text-gray-300 px-3 py-2 inline-block text-center rounded-full">Playlists</div>
                                <div className="text-xs bg-gray-600 text-gray-300 px-3 py-2 inline-block text-center rounded-full">Songs</div>
                                <div className="text-xs bg-gray-600 text-gray-300 px-3 py-2 inline-block text-center rounded-full">Videos</div>
                            </div>
                        </div>
                        <div className="flex-1 rounded-b flex flex-col overflow-scroll">
                            {Array(10).fill(null).map((_, i) => {
                                return <div className={`px-6 py-3 ${i == 0 && 'bg-gray-700'}`}>
                                    <p className={`text-gray-300 ${i == 0 && 'text-cyan-400'}`}>Sunday Morning</p>
                                    <p className={`text-xs ${i == 0 ? 'text-gray-400' : 'text-gray-500'}`}>Playlist • Sunday, September 10, 2023</p>
                                </div>
                            })}
                        </div>
                    </div>}
                </div>
                <div className='bg-gray-800 flex-1 rounded p-8 overflow-auto'>
                    <div className='flex justify-between items-center mb-4'>
                        <h1 className='text-3xl text-gray-300 font-bold top-0'>Sunday Morning</h1>
                        <button className='text-gray-400 hover:text-gray-300 p-2'><RiPencilLine /></button>
                    </div>
                    {Array(4).fill(null).map((_, i) => {
                        return (
                            <div className='mb-12'>
                                <h2 className={`text-xl mb-3 top-14 font-medium flex items-center gap-2 ${i == 0 ? 'text-cyan-300' : 'text-gray-300'}`}><RiMusic2Line /> There Is A Redeemer</h2>
                                <div className='grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-3'>
                                    {Array(7).fill("There is a Redeemer\nJesus, God's own son").map((lyrics, j) => {
                                        return <button className={`rounded text-left cursor-pointer ${i == 0 && j == 0 ? 'bg-cyan-600 hover:bg-cyan-500' : 'bg-gray-700 hover:bg-gray-600'}`}>
                                            <div className={`text-xs px-1.5 pt-1 rounded-t text-bold ${i == 0 && j == 0 ? 'text-gray-200' : 'text-gray-400'}`}>Verse 1</div>
                                            <div className='p-1'>
                                                <div className='bg-black rounded w-full aspect-video overflow-hidden'>
                                                    <Display lyrics={"There is a redeemer\nJesus, God's own son"} />
                                                </div>
                                            </div>
                                        </button>
                                    })}
                                </div>
                            </div>
                        )
                    })}
                </div>
            </div>
            <div className='flex'>
                <div className='flex-1'></div>
                <div className='flex-1 flex items-center justify-center gap-3'>
                    {[<RiSkipBackLine />,
                    <RiSkipForwardLine />].map(icon => <button className='text-3xl text-gray-300 rounded-full p-1 bg-gray-800 hover:bg-cyan-600 hover-text-gray-100'>{icon}</button>)}
                </div>
                <div className='flex-1 flex items-center justify-end'>
                    <button className="text-gray-400 text-lg hover:text-gray-300 p-2 -m-2"><RiSettings5Line /></button>
                </div>
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
