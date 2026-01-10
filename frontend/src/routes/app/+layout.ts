import type { LayoutLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { authFetch } from '$lib/authFetch';
import type { UserSessionRes } from '$lib/apiResponses';

// CSR-only app boundary
export const ssr = false;
export const csr = true;
export const prerender = false;

export const load: LayoutLoad = async () => {
	const { response, error, shouldRedirect } = await authFetch('http://localhost:8081/api/user');

	if (shouldRedirect) throw redirect(302, '/login');
	if (error) throw redirect(302, '/login');

	const user = response as UserSessionRes;
	return { user };
};