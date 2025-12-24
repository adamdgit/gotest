<script lang="ts">
    import { onMount } from 'svelte';
    import { user } from '$lib';
    import type { User } from '$lib';

    export let data: { user: User };

    // Set the user store on mount
    onMount(() => {
        user.set(data.user);
    });

    // Optional: refresh again on client if needed
    async function refreshAccessToken() {
        const res = await fetch('http://localhost:8081/api/refresh', {
            credentials: 'include'
        });

        if (res.ok) console.log('access token refreshed');
        else console.log(res.status, res.statusText);
    }

    onMount(() => {
        refreshAccessToken();
    });
</script>

<main>
    <h1>Welcome, {$user ? $user.email : ""}</h1>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit. Soluta excepturi animi enim magni, odio ipsam error aliquam natus nostrum, iste autem, nulla ipsa velit! Voluptatem ea iure officiis nesciunt officia.</p>
</main>
