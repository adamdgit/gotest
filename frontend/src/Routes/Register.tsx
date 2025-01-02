import { FormEvent, useState } from 'react'

export default function Register() {
  const [password, setPassword] = useState("");
  const [username, setUsername] = useState("");

  async function register(e:FormEvent<HTMLFormElement>) {
    e.preventDefault();

    const formData = new FormData();
    formData.append("username", username)
    formData.append("password", password)

    const res = await fetch('http://localhost:8081/api/auth/register', {
      method: "POST",
      body: formData
    })

    if (res.ok) {
      const data = await res.json();
      console.log(data)
    }
  }

  return (
    <div>
      <h2>Register an Account</h2>
      <form onSubmit={(e) => register(e)}>
        <label htmlFor='username'>Username</label>
        <input type='text' name='username' onChange={(e) => setUsername(e.target.value)} />
        <label htmlFor='password'>Password</label>
        <input type='text' name='password' onChange={(e) => setPassword(e.target.value)} />
        <button type='submit'>Submit</button>
      </form>
    </div>
  )
}
