import type { LayoutLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { authFetch } from '$lib/authFetch';

// CSR-only app boundary
export const ssr = false;
export const csr = true;
export const prerender = false;

export const load: LayoutLoad = async () => {
	try {
		const res = await authFetch('http://localhost:8081/api/user');

		if (!res.ok) {
			throw redirect(302, '/login');
		}

		const user = await res.json();
		return { user };
	} catch {
		throw redirect(302, '/login');
	}
};