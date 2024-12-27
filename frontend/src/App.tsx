import { BrowserRouter, Routes, Route } from "react-router-dom";
import Home from './components/Routes/Home';
import Navbar from './components/Navbar';
import React from "react";
import Posts from "./components/Routes/Posts";

function App() {

  return (
    <React.Fragment>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Navbar />}>
            <Route index element={<Home />} />
            <Route path="Posts" element={<Posts />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </React.Fragment>
  )
}

export default App
