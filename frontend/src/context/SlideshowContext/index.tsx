import React, { Context, createContext, useReducer } from 'react'
import { SlideshowReducer } from './reducer'


export const SlideshowContext: Context<[any, any]> = createContext([{}, () => { }])

export const SlideshowProvider = ({ children }: any) => {
    const [state, dispatch] = useReducer(SlideshowReducer, {})

    return <SlideshowContext.Provider value={[state, dispatch]}>{children}</SlideshowContext.Provider>
}