import { useNavigate } from "@solidjs/router";
import { useAuth } from "../AuthProvider"
import { onMount, Show } from "solid-js";

export default function Header(props) {
    const { userData, validateUser, signOut } = useAuth();
    const navigate = useNavigate()

    onMount(async () => {
      const isLoggedIn = await validateUser();

      if (!isLoggedIn) {
        navigate("/")
      } 
    })

    async function handleLogout() {  
      const res = await signOut();
  
      if (res) {
        navigate("/");
      } else {
        console.log("Error logging out")
      }
    }
  
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
              <button class='btn' onclick={() => handleLogout()}>Logout</button>
            </div>
          </>
        </Show>
      </header>
      {props.children}
      </>
    )
}
