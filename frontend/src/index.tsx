import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import { SlideshowProvider } from './context/SlideshowContext';
import '@fontsource-variable/inter';
import Display from './Display';
import './offscreenCanvasPolyfill.js'

import { EventsOn } from '../wailsjs/runtime';

(() => {
  // This is probably not remotely ok in react but I'm not a frontend dev so ¯\_(ツ)_/¯
  EventsOn("data-mutate", (data: any) => {
    console.log("[DATA MUTATED]",data)
  })
})();

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <SlideshowProvider>
      <App />
    </SlideshowProvider>
  </React.StrictMode>
);
