// Refresh tokens rotate on every use, so concurrent requests that get a 401
// must share one refresh call, otherwise the later ones send a stale token
let refreshInFlight: Promise<Response> | null = null;

function refreshSession(): Promise<Response> {
	if (!refreshInFlight) {
		refreshInFlight = fetch('http://localhost:8081/api/auth/refresh', {
			credentials: 'include'
		}).finally(() => {
			refreshInFlight = null;
		});
	}
	return refreshInFlight;
}

// Automatically attempts to refresh access tokens if they expire
// and sends request again if refresh is successful. 
// Also a wrapper for fetch, returning errors
export async function authFetch(
	input: RequestInfo,
	init: RequestInit = {}
): Promise<{ 
	response: any | null, 
	error?: string | null, 
	shouldRedirect?: boolean 
}> {
	let response:  any | null = null;
	let error: string | null = null;
	let shouldRedirect = false;

	try {
		const res = await fetch(input, {
			...init,
			credentials: 'include'
		});

		// Skip attempting to refresh if auth is ok
		if (res.ok) {
			const data = await res.json();
			response = data;
			return { response, error, shouldRedirect }
		}
	} catch (err) {
		error = "Network Error";
		shouldRedirect = true;
		return { response, error, shouldRedirect }
	}

	// Attempt refresh
	try {
		const refreshRes = await refreshSession();

		// Refresh failed, return errors
		if (!refreshRes.ok) {
			error = refreshRes.statusText;
			shouldRedirect = true;
			return { response, error, shouldRedirect }
		}
	} catch (err) {
		error = "Network Error";
		shouldRedirect = true;
		return { response, error, shouldRedirect }
	}

	// retry after refreshing token is successful
	try {
		const res = await fetch(input, {
			...init,
			credentials: 'include'
		});

		if (res.ok) {
			const data = await res.json();
			response = data;
		} else {
			error = res.statusText;
			shouldRedirect = true;
		}
	} catch (err) {
		error = "Network Error";
		shouldRedirect = true;
	}

	// Retry request was successful, return results
	return { response, error, shouldRedirect }
}
