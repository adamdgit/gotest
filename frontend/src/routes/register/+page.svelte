<script lang="ts">
  let email = $state("")
  let password = $state("")

  async function register() {
    const res = await fetch("http://localhost:8081/api/auth/register", {
      method: "post",
      headers: {"content-type": "application/json"},
      body: JSON.stringify({"email": email, "password": password}),
      credentials: "include"
    })

    if (res.ok) {
      console.log(res.statusText, res.status)
      const data = await res.json();
      console.log(data)
    } else {
      const err = await res.json();
      console.log(err)
    }
  }
</script>

<h1>Register an Account</h1>

<form>
  <label for="email">Email</label>
  <input type="email" name="email" bind:value={email}>

  <label for="password">Password</label>
  <input type="password" name="password" bind:value={password}>

  <button onclick={() => register()}>Register</button>
</form>

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