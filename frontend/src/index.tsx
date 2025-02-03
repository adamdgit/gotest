import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";
import "./styles/index.css"
import Home from "./pages/Home";
import Login from "./pages/Login";
import NotFound from "./pages/NotFound";
import { AuthProvider } from "./AuthProvider";
import Register from "./pages/Register";

const root = document.getElementById("root");
if (!root) throw new Error("Somehow root doesn't exist?")

render(
  () => (
    <AuthProvider>
      <Router>
        <Route path="/" component={Login} />
        <Route path="/register" component={Register} />
        <Route path="/home" component={Home} />
        <Route path="*paramName" component={NotFound} />
      </Router>
    </AuthProvider>
  ),
  root
);