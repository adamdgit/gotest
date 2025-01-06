import { FormEvent, useState } from 'react'

export default function Register() {
  const [password, setPassword] = useState("");
  const [email, setEmail] = useState("");

  async function register(e:FormEvent<HTMLFormElement>) {
    e.preventDefault();

    const formData = new FormData();
    formData.append("email", email)
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
        <label htmlFor='email'>Email</label>
        <input type='email' name='email' onChange={(e) => setEmail(e.target.value)} />
        <label htmlFor='password'>Password</label>
        <input type='password' name='password' onChange={(e) => setPassword(e.target.value)} />
        <button type='submit'>Submit</button>
      </form>
    </div>
  )
}
