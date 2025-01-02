import { BrowserRouter, Routes, Route } from "react-router-dom";
import Login from './Routes/Login';
import Navbar from './components/Navbar';
import React from "react";
import Posts from "./Routes/Posts";
import Register from "./Routes/Register";
import Home from "./Routes/Home";

function App() {

  return (
    <React.Fragment>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Navbar />}>
            <Route index element={<Home />} />
            <Route path="login" element={<Login />} />
            <Route path="register" element={<Register />} />
            <Route path="posts" element={<Posts />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </React.Fragment>
  )
}

export default App
