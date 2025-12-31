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
			<button onclick={logout} class="logout">Logout</button>
			<img src={$user.profile_url} alt="Default user circle" class="profilePic" />
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

	.profilePic {
		color: white;
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
