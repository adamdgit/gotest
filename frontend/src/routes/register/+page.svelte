<script lang="ts">
	let email = '';
	let password = '';
	let error = '';
	let loading = false;

	async function register() {
		error = '';
		loading = true;

		const res = await fetch(`http://localhost:8081/api/auth/register`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ email, password })
		});

		loading = false;

		if (!res.ok) {
			error = 'Registration failed';
			return;
		}

		window.location.href = '/login';
	}
</script>

<h1>Register</h1>

<form on:submit|preventDefault={register}>
	<input type="email" placeholder="Email" bind:value={email} required />
	<input type="password" placeholder="Password" bind:value={password} required />

	<button disabled={loading}>
		{loading ? 'Registering...' : 'Register'}
	</button>
</form>

{#if error}
	<p class="error">{error}</p>
{/if}

<style>
	form {
		max-width: 300px;
		margin: 2rem auto;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}
	.error {
		color: red;
		text-align: center;
	}
</style>
