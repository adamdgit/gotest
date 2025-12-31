import type { User } from '$lib';
import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ fetch }) => {
    // Try to fetch the current user
    const res = await fetch('http://localhost:8081/api/user', {
        credentials: 'include',
    });

    // If user not found or access token expired
    if (res.status === 401) {
        // Try refreshing
        const refreshRes = await fetch('http://localhost:8081/api/refresh', {
            credentials: 'include',
        });

        if (!refreshRes.ok) {
            // Refresh failed → redirect to login
            throw redirect(302, '/login');
        }

        // Refresh succeeded → fetch user again
        const newUserRes = await fetch('http://localhost:8081/api/user', {
            credentials: 'include',
        });

        if (!newUserRes.ok) throw redirect(302, '/login');

        const user = await newUserRes.json();
        return { user };
    }

    if (!res.ok) {
        throw redirect(302, '/login');
    }

    const user: User = await res.json();
    return { user };
};
