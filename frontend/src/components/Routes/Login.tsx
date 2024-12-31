import { FormEvent, useState } from 'react'

export default function Login() {
  const [password, setPassword] = useState("");
  const [username, setUsername] = useState("");

  async function login(e:FormEvent<HTMLFormElement>) {
    e.preventDefault();

    const res = await fetch('http://localhost:8081/api/auth/login', {
      method: "POST",
      body: JSON.stringify({
        "username": username,
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
    })

    if (res.ok) {
      console.log(res.status)
    }
  }


  return (
    <div>
      <h2>Login</h2>
      <form onSubmit={(e) => login(e)}>
        <label htmlFor='username'>Username</label>
        <input type='text' name='username' onChange={(e) => setUsername(e.target.value)} />
        <label htmlFor='password'>Password</label>
        <input type='text' name='password' onChange={(e) => setPassword(e.target.value)} />
        <button type='submit'>Submit</button>
      </form>

      <form onSubmit={(e) => logout(e)}>
        <button type='submit'>Logout</button>
      </form>
    </div>
  )
}
