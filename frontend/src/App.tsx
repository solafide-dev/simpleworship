import wailsLogo from './assets/wails.png'
import './App.css'
import * as Display from '../wailsjs/go/main/DisplayServer'
import * as rt from '../wailsjs/runtime/runtime'

function App() {

    // Set Slide Data
    var verse = 1

    setInterval(() => {
        if (verse == 1) {
            verse = 2
            Display.SetData({
                type: "song", 
                meta: {section: "Verse 1", text: "This is amazing grace, this is unfailing love\n\nNext line I can't remember yeahhhh"},
                data: "This is ignored because its not part of the struct"
            })
        } else {
            verse = 1
            Display.SetData({
                type: "song", 
                meta: {section: "Verse 2", text: "Who brings our chaos back into order\n\nWho makes the orphan a son and daughter"},
                data: "This is ignored because its not part of the struct"
            })
        }
    }, 5000)

    rt.LogInfo("Hello from React!")
    
    return (
        <div className="min-h-screen bg-white grid grid-cols-1 place-items-center justify-items-center mx-auto py-8">
            <div className="text-blue-900 text-2xl font-bold font-mono">
                <h1 className="content-center">Vite + React + TS + Tailwind</h1>
            </div>
            <div className="w-fit max-w-md">
                <a href="https://wails.io" target="_blank">
                    <img src={wailsLogo} className="logo wails" alt="Wails logo" />
                </a>
            </div>
        </div>
    )
}

export default App
