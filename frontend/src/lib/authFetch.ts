import { goto } from '$app/navigation';

let refreshPromise: Promise<boolean> | null = null;

async function refreshAccessToken(): Promise<boolean> {
	if (!refreshPromise) {
		refreshPromise = (async () => {
			const res = await fetch('http://localhost:8081/api/refresh', {
				credentials: 'include'
			});
			refreshPromise = null;
			return res.ok;
		})();
	}

	return refreshPromise;
}

// Refresh token wrapper for fetch API calls
// handles automatic refresh ok access tokens if needed
export async function authFetch(
	input: RequestInfo,
	init: RequestInit = {}
): Promise<Response> {
	const res = await fetch(input, {
		...init,
		credentials: 'include'
	});

	// Skip attempting to refresh if auth is ok
	if (res.status === 200) {
		return res;
	}

	// Attempt refresh
	const refreshed = await refreshAccessToken();

	if (!refreshed) {
		console.log("Refresh Error: ", refreshed)
		goto('/login');
		throw new Error('Session expired');
	}

	// Retry original request
	return fetch(input, {
		...init,
		credentials: 'include'
	});
}
