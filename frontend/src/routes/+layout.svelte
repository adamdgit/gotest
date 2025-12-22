<script lang="ts">
    import { user } from '$lib';
	import favicon from '$lib/assets/favicon.svg';
	import '$lib/styles.css';

	let { children } = $props();

	async function logout() {
		await fetch(`http://localhost:8081/api/auth/logout`, {
			method: 'GET',
			credentials: 'include'
		});

		window.location.href = '/login';
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<nav>
	<ul class="navbar">
		{#if $user}
			<li><a href="/">Home</a></li>
		{/if}

		{#if !$user}
			<li><a href="/login">Login</a></li>
		{/if}

		<li><a href="/register">Register</a></li>
	</ul>

	{#if $user}
		<div class="profile-bubble">
			<span>{$user.email}</span>
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
