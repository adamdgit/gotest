import { useAuth } from "../AuthProvider"
import { Show } from "solid-js";

export default function Header(props) {
    const { userData } = useAuth();
  
    return (
      <>
      <header>
        <h3>Appname</h3>
        <Show when={userData.email}>
          <>
            <nav class="nabvar">
                <ul class="navlist">
                    <li><a href="/home">Home</a></li>
                    <li><a href="#">Admin</a></li>
                    <li><a href="#">Forms</a></li>
                    <li><a href="#">Staff</a></li>
                    <li><a href="/inventory">Inventory</a></li>
                </ul>
            </nav>
            <div class="profilewrapper">
              {userData?.email}
              <button class='btn'>Logout</button>
            </div>
          </>
        </Show>
      </header>
      {props.children}
      </>
    )
}
