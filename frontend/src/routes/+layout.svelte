<script lang="ts">
    import { goto } from '$app/navigation';
    import { user } from '$lib';
	import favicon from '$lib/assets/favicon.svg';
	import '$lib/styles.css';

	let { children } = $props();

	async function logout() {
		await fetch(`http://localhost:8081/api/auth/logout`, {
			method: 'GET',
			credentials: 'include'
		});

		goto('/login');
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<nav>
	<ul class="navbar">
		{#if $user}
			<li><a href="/app">Home</a></li>
		{:else}
			<li><a href="/login">Login</a></li>
			<li><a href="/register">Register</a></li>
		{/if}
	</ul>

	{#if $user}
		<div class="profile-bubble">
			<span>{$user.email}</span>
			<svg xmlns="http://www.w3.org/2000/svg" fill="currentColor" class="profilePic" viewBox="0 0 640 640"><!--!Font Awesome Free v7.1.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2025 Fonticons, Inc.--><path d="M463 448.2C440.9 409.8 399.4 384 352 384L288 384C240.6 384 199.1 409.8 177 448.2C212.2 487.4 263.2 512 320 512C376.8 512 427.8 487.3 463 448.2zM64 320C64 178.6 178.6 64 320 64C461.4 64 576 178.6 576 320C576 461.4 461.4 576 320 576C178.6 576 64 461.4 64 320zM320 336C359.8 336 392 303.8 392 264C392 224.2 359.8 192 320 192C280.2 192 248 224.2 248 264C248 303.8 280.2 336 320 336z"/></svg>
			<button onclick={logout} class="logout">Logout</button>
		</div>
	{/if}
</nav>

<style>
	nav {
		display: flex;
		align-items: center;
		justify-content: space-between;
		width: 100%;
		background: #212931;
		color: white;
		padding: 1rem;
	}

	.navbar {
		display: flex;
		gap: 1rem;
		list-style: none;
	}

	a {
		color: white;
		text-decoration: none;
	}
	a:hover { text-decoration: underline; }

	.profile-bubble {
		display: flex;
		align-items: center;
	}

	.profilePic {
		width: 65px;
		height: 65px;
	}

	.logout {
		margin-left: auto;
		cursor: pointer;
		background-color: #2d4863;
		color: white;
		padding: .5rem .7rem;
		border-radius: 4px;
	}
	.logout:hover {
		background-color: #3d6185;
	}
</style>


{@render children()}
