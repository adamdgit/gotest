import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";
import { AuthProvider } from "./AuthProvider";
import Home from "./pages/Home";
import Login from "./pages/Login";
import NotFound from "./pages/NotFound";
import Register from "./pages/Register";
import Inventory from "./pages/inventory/Inventory";
import "./styles/index.css"
import Header from "./pages/Header";

const root = document.getElementById("root");
if (!root) throw new Error("Somehow root doesn't exist?")

render(
  () => (
    <AuthProvider>
      <Router root={Header}>
        <Route path="/" component={Login} />
        <Route path="/register" component={Register} />
        <Route path="/home" component={Home} />
        <Route path="/inventory" component={Inventory} />
        <Route path="*paramName" component={NotFound} />
      </Router>
    </AuthProvider>
  ),
  root
);