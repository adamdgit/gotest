import { useNavigate } from "@solidjs/router";
import { useAuth } from "../AuthProvider"

export default function Header(props) {
    const navigate = useNavigate();
    const { userData } = useAuth();

    // Redirect non logged in users to login page
    if (!userData.email) {
      navigate("/");
    }
  
    return (
      <>
          <header>
            <h3>Appname</h3>
            {userData.email ? 
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
              : null
            }
          </header>
          {props.children}
      </>
    )
}
