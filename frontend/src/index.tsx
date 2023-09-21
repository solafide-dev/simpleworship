import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import { SlideshowProvider } from './context/SlideshowContext';
import '@fontsource-variable/inter';

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
