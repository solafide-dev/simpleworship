import React from 'react'
import { Song as SongType } from '../../global'
import SlideThumbnail from './SlideThumbnail'

const Song = ({ name, arrangement, sections }: SongType) => {
    // const arrangedSections = arrangement.map(a => sections.filter(b => b.id === a)[0])
    const arrangedSections = sections
    return (<div>
        <div style={{ background: '#1122ff', color: 'white', fontWeight: 'bold' }}>{name}</div>
        {arrangedSections.map(({ name, slides }) => {
            return (<div style={{ padding: '1em 0' }}>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '1em' }}>
                    {slides.map(slide => <div><SlideThumbnail {...slide} />{name}</div>)}
                </div>
            </div>)
        })}
    </div>)
}

export default Song