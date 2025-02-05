import { createEffect, createSignal } from 'solid-js'
import "../styles/Login.css"
import { useNavigate } from '@solidjs/router';
import { useAuth } from '../AuthProvider';

export default function Login() {
  const navigate = useNavigate();
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [error, setError] = createSignal("");
  const { userData, signIn } = useAuth();

  // Redirect non logged in users to login page
  createEffect(() => {
    if (userData.email) {
      navigate("/home");
    }
  })

  async function handleLogin(e: Event) {
    e.preventDefault();
    setError(""); // Reset error state

    const res = await signIn(email(), password());

    if (res) {
      navigate("/home");
    } else {
      setError("Error logging in");
    }
  }

  return (
    <main>
      <h1>Login</h1>
      <p style={{color: 'red'}}>{error()}</p>
      <form class='loginform'>
        <div class='inputwrap'>
          <label for="email">Email</label>
          <input type="email" name='email' 
            placeholder='email@example.com' 
            onchange={(e) => setEmail(e.target.value)} 
          />
        </div>

        <div class='inputwrap'>
          <label for='password'>Password</label>
          <input type='password' name='password' 
            onchange={(e) => setPassword(e.target.value)}
          />
        </div>
        <div>
          Don't have an account? 
          <a href='/register'>Register account</a>
        </div>
        <button onclick={(e) => handleLogin(e)}>Login</button>
      </form>
    </main>
  )
}
