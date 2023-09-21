import React, { useContext } from 'react'
import { SlideshowContext } from '../../context/SlideshowContext'
import { SlideshowActions } from '../../context/SlideshowContext/SlideshowContext.d'
import { Slide } from '../../global'

const SlideThumbnail = ({ text, id }: Slide) => {
    const [{ activeSlide }, dispatch] = useContext(SlideshowContext)

    return (<div style={{ border: '1px solid', borderColor: activeSlide === id ? 'red' : 'black', height: '0', paddingBottom: '56.25%' }}
        onClick={() => dispatch({ type: SlideshowActions.SET_ACTIVE_SLIDE, payload: id })}>
        {text.split('\n').map(line => <>{line}<br /></>)}
    </div>)
}

export default SlideThumbnail