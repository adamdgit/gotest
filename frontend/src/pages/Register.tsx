import { createSignal } from "solid-js";
import { useNavigate } from "@solidjs/router";

export default function Register() {
    const navigate = useNavigate();
    const [email, setEmail] = createSignal("");
    const [password, setPassword] = createSignal("");
    const [honeypot, setHoneypot] = createSignal("");
    const [error, setError] = createSignal("");


    async function handleRegister(e: Event) {
      e.preventDefault();

      // if honeypot is filled out, don't send request
      if (honeypot() !== "") {
        return
      }

      const res = await fetch("http://localhost:8081/api/auth/register", {
        method: 'POST',
        credentials: 'include',
        headers: {
            "content-type": "application/json",
        }, 
        body: JSON.stringify({
          "email": email(),
          "password": password()
        })
      })

      if (res.ok) {
        navigate("/?message=account%20created");
      } else {
        const { message } = await res.json();
        setError(`Error: ${message}`);
      }

    }

  return (
    <div>
      <h1>Register</h1>

      <p style={{color: 'red'}}>{error()}</p>
      <form class='loginform'>
        <div class='inputwrap'>
          <label for="email">Email</label>
          <input type="email" name='email' 
            placeholder='email@example.com' 
            onchange={(e) => setEmail(e.target.value)} 
          />
        </div>

        <div style={{ display: "none" }} aria-disabled="true">
          <label for="hp">Hp</label>
          <input type="text" name='hp' 
            placeholder='something' 
            onchange={(e) => setHoneypot(e.target.value)}
          />
        </div>

        <div class='inputwrap'>
          <label for='password'>Password</label>
          <input type='password' name='password' 
            onchange={(e) => setPassword(e.target.value)}
          />
        </div>
        
        <button onclick={(e) => handleRegister(e)}>Register</button>
      </form>
    </div>
  )
}
