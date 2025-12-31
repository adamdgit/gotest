<script lang="ts">
    import { onMount } from 'svelte';
    import { user } from '$lib';
    import type { PageProps } from '../$types';
    import { authFetch } from '$lib/authFetch';

	let { data }: PageProps = $props();

    onMount(() => {
        user.set(data.user);
    });
    
	let search = $state("")
	let users:any = $state(null)
	let error: string | null = $state(null)

	async function searchUser() {
		error = null;
		users = null;

		if (!search) {
			error = "Please enter a valid search address";
			return;
		}

		try {
			const res = await authFetch(
				`http://localhost:8081/api/admin/users?search=${encodeURIComponent(search)}`
			);

			if (!res.ok) {
				error = "Failed to fetch users";
				return;
			}

			users = await res.json();
		} catch (err) {
			console.error(err);
			error = "Network error";
		}
	}
</script>

<main>
    <h1>Welcome, {$user ? $user.email : ""}</h1>

    <label for="search">Search for users: </label>
	<input
		type="search"
		placeholder="Enter user email"
		bind:value={search}
		onkeydown={(e) => e.key === "Enter" && searchUser()}
	/>

	<button onclick={searchUser} class="searchtbtn">Search</button>

	{#if $error}
		<p style="color: red;">{$error}</p>
	{/if}

	{#if users}
		<h2>Users Found:</h2>
		<div class="result-list">
		{#each users as user}
			<a class="user" href="/app/user/{user.id}">
				<p><strong>Email</strong> {user.email}</p>
				<p><strong>Name</strong> {user.firstname} {user.lastname}</p>
				<p><strong>Address</strong> {user.address}</p>
				<p><strong>Role</strong> {user.role}</p>
			</a>
		{/each}
		</div>
	{/if}
</main>

<style>
	input {
		padding: 0.5rem;
		margin-right: 0.5rem;
	}

	.searchtbtn {
		padding: 0.5rem 1rem;
		color: #fff;
		background: #2d4863;
	}
	.searchtbtn:hover {
		background: #3d6185;
	}

	.result-list {
		display: grid;
		gap: 1rem;
	}

	.user {
		display: grid;
		grid-template-columns: 200px 170px 90px 230px 90px;
		gap: .4rem;
		background: #fff;
		padding: 1rem;
		color: black;
		text-decoration: none;
	}

	.user:hover {
		background: #d5e7ff;
	}

	.user p {
		display: grid;
	}
</style>