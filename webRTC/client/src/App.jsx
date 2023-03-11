import React from 'react'
import {BrowserRouter, Routes, Route, createBrowserRouter} from 'react-router-dom';
import { createBrowserHistory } from 'history';

import CreateRoom from './components/CreateRoom';
import Room from './components/Room';

function App() {
  return (
    <div className="App">
      <BrowserRouter history={history}>
        <Routes>
          <Route path="/" element={<CreateRoom />} />
          <Route path="/room/:roomID" element={<Room />} />
        </Routes>
      </BrowserRouter>
    </div>
  )
}

export default App;