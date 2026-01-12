<script lang="ts">
    import { onMount } from 'svelte';
    import { user } from '$lib';
    import { authFetch } from '$lib/authFetch';
    import { goto } from '$app/navigation';
    import type { AdminUserDataRes } from '$lib/apiResponses';

	let { data } = $props();

    onMount(() => {
        user.set(data.user);
    });
    
	let search = $state("");
	let users:any = $state(null);
	let errors: string | null = $state(null);
	let selectedUser = $state<AdminUserDataRes | null>(null);

	async function searchUser() {
		errors = null;
		users = null;

		if (!search) {
			errors = "Please enter a search query";
			return;
		}

		const { response, error, shouldRedirect } = await authFetch(
			`http://localhost:8081/api/admin/users?search=${encodeURIComponent(search)}`
		);

		if (response) users = response as AdminUserDataRes;
		if (error) errors = error;
		if (shouldRedirect) goto('/login');
	}
</script>

<main>
    <h1>Users Details</h1>

	<label for="search" style="display: block;">Search by name, address, email..</label>
	<div class="search-wrap">
		<input
			class="searchinput"
			type="search"
			placeholder="search..."
			bind:value={search}
			onkeydown={(e) => e.key === "Enter" && searchUser()}
		/>

		<button onclick={searchUser} class="searchtbtn">
			<svg xmlns="http://www.w3.org/2000/svg" fill="currentColor" width="20" height="20" viewBox="0 0 640 640"><!--!Font Awesome Free v7.1.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2026 Fonticons, Inc.--><path d="M480 272C480 317.9 465.1 360.3 440 394.7L566.6 521.4C579.1 533.9 579.1 554.2 566.6 566.7C554.1 579.2 533.8 579.2 521.3 566.7L394.7 440C360.3 465.1 317.9 480 272 480C157.1 480 64 386.9 64 272C64 157.1 157.1 64 272 64C386.9 64 480 157.1 480 272zM272 416C351.5 416 416 351.5 416 272C416 192.5 351.5 128 272 128C192.5 128 128 192.5 128 272C128 351.5 192.5 416 272 416z"/></svg>
		</button>
	</div>

	{#if errors}
		<p style="color: red;">{errors}</p>
	{/if}

	<div class="results-wrap">
		<h2>Results</h2>
		{#if users}
			<div class="result-list">
				{#each users as user}
					<button class="user" onclick={() => selectedUser = user}>
						<div><strong>Email</strong> {user.email}</div>
						<div><strong>Name</strong> {user.firstname} {user.lastname}</div>
						<div><strong>Address</strong> {user.address}</div>
						<div><strong>Role</strong> {user.role}</div>
					</button>
				{/each}
			</div>
		{/if}
	</div>

	{#if selectedUser}
		<div class="user-detailed-modal">
			<button class="close" onclick={() => selectedUser = null}>
				<svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" fill="currentColor" viewBox="0 0 640 640"><!--!Font Awesome Free v7.1.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2026 Fonticons, Inc.--><path d="M320 576C461.4 576 576 461.4 576 320C576 178.6 461.4 64 320 64C178.6 64 64 178.6 64 320C64 461.4 178.6 576 320 576zM231 231C240.4 221.6 255.6 221.6 264.9 231L319.9 286L374.9 231C384.3 221.6 399.5 221.6 408.8 231C418.1 240.4 418.2 255.6 408.8 264.9L353.8 319.9L408.8 374.9C418.2 384.3 418.2 399.5 408.8 408.8C399.4 418.1 384.2 418.2 374.9 408.8L319.9 353.8L264.9 408.8C255.5 418.2 240.3 418.2 231 408.8C221.7 399.4 221.6 384.2 231 374.9L286 319.9L231 264.9C221.6 255.5 221.6 240.3 231 231z"/></svg>
			</button>
			<div class="label-wrap">
				<div>
					<img src={selectedUser.profile_url} alt="user profile" class="profilepic" />
				</div>
				<div class="label"><strong>ID</strong> {selectedUser.id}</div>
				<div class="label"><strong>Email</strong> {selectedUser.email}</div>
				<div class="label"><strong>Name</strong> {selectedUser.firstname} {selectedUser.lastname}</div>
				<div class="label"><strong>Address</strong> {selectedUser.address}</div>
				<div class="label"><strong>Phone</strong> {selectedUser.phone}</div>
				<div class="label"><strong>Role</strong> {selectedUser.role}</div>
				<div class="label"><strong>Created At</strong>
					{new Date(selectedUser.created_at).toLocaleDateString("en", {day: "2-digit", month: "2-digit", year: "numeric", hour: "numeric", minute: "numeric"})}
				</div>
				<div class="label"><strong>Updated At</strong>
					{new Date(selectedUser.updated_at).toLocaleDateString("en", {day: "2-digit", month: "2-digit", year: "numeric", hour: "numeric", minute: "numeric"})}
				</div>
			</div>
		</div>
	{/if}
</main>

<style>
	input {
		padding: 0.5rem;
		margin-right: 0.5rem;
	}

	.search-wrap {
		display: flex;
		align-items: center;
	}

	.searchinput {
		width: 300px;
		border: 1px solid #ddd;
		height: 40px;
	}

	.results-wrap {
		margin-top: 1rem;
	}

	.user-detailed-modal {
		position: absolute;
		left: 50%;
		top: 50%;
		transform: translateX(-50%) translateY(-50%);
		border-radius: 4px;
		background-color: white;
		padding: 2rem;
		box-shadow: 0 0 0 6000px rgba(0, 0, 0, .7);
	}

	.label-wrap {
		display: grid;
		justify-items: center;
	}

	.label {
		display: grid;
		grid-template-columns: 1fr 3fr;
		align-items: center;
		gap: 1rem;
		width: 400px;
	}

	.close {
		position: absolute;
		top: -12px;
		left: -12px;
		background: white;
		border-radius: 50%;
		width: 40px;
		height: 40px;
	}
	.close svg:hover {
		fill: #3d6185;
	}

	.profilepic {
		width: 200px;
	}

	.searchtbtn {
		padding: 0.5rem 1rem;
		color: #fff;
		background: #2d4863;
		border-radius: 5px;
		height: 40px;
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
		grid-template-columns: 2fr 1fr 2fr .5fr;
		gap: .4rem;
		background: #fff;
		padding: 1rem;
		color: black;
		text-decoration: none;
	}

	.user:hover {
		background: #d5e7ff;
	}

	.user div {
		display: grid;
		text-align: left;
	}
</style>