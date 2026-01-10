<script lang="ts">
    import { goto } from "$app/navigation";
	import { page } from "$app/state";
    import { user } from "$lib";
    import type { APIError } from "$lib/apiResponses";
    import { onMount } from "svelte";

	let email = '';
	let password = '';
	let error = '';
	let loading = false;
	let redirect_msg = page.url.searchParams.get("msg");

	onMount(() => {
		if (user) {
			goto('/app')
		}
	});

    async function login() {
		error = '';
		loading = true;

		// hp value
		let x = document.getElementById('username') as HTMLInputElement;
		let username = x.value ?? "";

		try {
			const res = await fetch(`http://localhost:8081/api/auth/login`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				credentials: 'include',
				body: JSON.stringify({ email, password, userAgent: navigator.userAgent, username })
			});

			if (!res.ok) {
				const data = await res.json() as APIError;

				error = data.error;
				return;
			}
		} catch (err) {
			console.error(err);
		} finally {
			loading = false;
		}

		// redirect on success
		goto('/app');
	}
</script>

<h1>Webapp Login</h1>

{#if redirect_msg}
	<p class="msg">{redirect_msg}</p>
{/if}

<form onsubmit={login}>
	<label for="email">Email</label>
	<input name="email" type="email" bind:value={email} required />
	
	<label for="password">Password</label>
	<input name="password" type="password" bind:value={password} required />

	<div class="hp-wrap">
		<label for="username">Username</label>
		<input id="username" name="username" type="text" tabindex="-1" autocomplete="off" />
	</div>

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
	.msg {
		color: #3d6185;
		text-align: center;
	}

	.hp-wrap {
		position: absolute;
		left: 0;
		bottom: 0;
		width: 0px;
		height: 0px;
		overflow: hidden;	
	}
</style>