<script lang="ts">
	let email = '';
	let password = '';
	let error = '';
	let loading = false;

    async function login() {
		error = '';
		loading = true;

		const res = await fetch(`http://localhost:8081/api/auth/login`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ email, password })
		});

		loading = false;

		if (!res.ok) {
			error = 'Invalid email or password';
			return;
		}

		// redirect on success
		window.location.href = '/';
	}
</script>

<h1>Webapp Login</h1>

<form onsubmit={login}>
	<label for="email">Email</label>
	<input name="email" type="email" bind:value={email} required />
	
	<label for="password">Password</label>
	<input name="password" type="password" bind:value={password} required />

	<button disabled={loading}>
		{loading ? 'Logging in...' : 'Login'}
	</button>
</form>

{#if error}
	<p class="error">{error}</p>
{/if}

<style>
	h1 { 
		text-align: center;
		padding-top: 2rem;
	}

	form {
		max-width: 350px;
		margin: 2rem auto;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		background-color: #f8f8f8;
		padding: 2rem;
		border-radius: 4px;
		box-shadow: 0 4px 14px 0 #ddd;
	}

	input {
		padding: .4rem .7rem;
		border: 1px solid white;
		border-radius: 4px;
		border: 1px solid #ddd;
	}
	input:hover {
		border: 1px solid #bbb;
	}

	button {
		margin: 1rem auto 0 auto;
		width: 150px;
		padding: .7rem 1rem;
		border-radius: 4px;
		background-color: #2d4863;
		color: white;
	}
	button:hover {
		background-color: #3d6185;
	}

	.error {
		color: red;
		text-align: center;
	}
</style>