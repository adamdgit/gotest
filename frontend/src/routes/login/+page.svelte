<script lang="ts">
  let { data } = $props();
  let email = $state("");
  let password = $state("");

  console.log(data)

  async function login() {
    const res = await fetch("http://localhost:8081/api/auth/login", {
      method: "post",
      headers: {"content-type": "application/json"},
      body: JSON.stringify({"email": email, "password": password}),
      credentials: "include"
    })

    // Handle errors
    if (!res.ok) {
      const err = await res.json();
      return console.log(err)
    }

    const data = await res.json();
    console.log(data)
  };

  async function logout() {
    const res = await fetch("http://localhost:8081/api/auth/logout", {
      method: "post",
      headers: {"content-type": "application/json"},
      credentials: "include"
    })

    // Handle errors
    if (!res.ok) {
      return console.log(res.status, res.statusText)
    }

    const data = await res.json();
    console.log(data)
  };
</script>

<h1>Login</h1>

<p>{data.title} {data.content}</p>

<form>
  <label for="email">Email</label>
  <input type="email" name="email" bind:value={email}>

  <label for="password">Password</label>
  <input type="password" name="password" bind:value={password}>

  <button onclick={() => login()}>Login</button>
</form>

<button onclick={() => logout()}>Logout</button>

<style>
  form {
    width: fit-content;
    display: grid;
    padding: 2rem;
    border-radius: 4px;
    background-color: #fafafa;
    box-shadow: 0 4px 15px 0 #ddd;
  }

  form input, form label, form button {
    padding: .5rem;
    max-width: 400px;
    min-width: 250px;
  }

  form button {
    margin-top: 1.4rem;
    display: flex;
    justify-self: center;
    justify-content: center;
    width: 200px;
  }
</style>