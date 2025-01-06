import { FormEvent, useState } from 'react'

export default function Login() {
  const [password, setPassword] = useState("");
  const [email, setEmail] = useState("");

  async function login(e:FormEvent<HTMLFormElement>) {
    e.preventDefault();

    const res = await fetch('http://localhost:8081/api/auth/login', {
      method: "POST",
      body: JSON.stringify({
        "email": email,
        "password": password
      }),
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
    })

    if (res.ok) {
      const data = await res.json();
      console.log(data)
    }
  }

  async function logout(e:FormEvent<HTMLFormElement>) {
    e.preventDefault();

    const res = await fetch('http://localhost:8081/api/auth/logout', {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
    })

    if (res.ok) {
      console.log(res.status)
    }
  }


  return (
    <div>
      <h2>Login</h2>
      <form onSubmit={(e) => login(e)}>
        <label htmlFor='email'>Email</label>
        <input type='email' name='email' onChange={(e) => setEmail(e.target.value)} />
        <label htmlFor='password'>Password</label>
        <input type='password' name='password' onChange={(e) => setPassword(e.target.value)} />
        <button type='submit'>Submit</button>
      </form>

      <form onSubmit={(e) => logout(e)}>
        <button type='submit'>Logout</button>
      </form>
    </div>
  )
}
