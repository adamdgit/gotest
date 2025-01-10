import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
  const sessionCookie = cookies.get('session_id'); // Retrieve session cookie

  if (!sessionCookie) {
    throw redirect(302, '/login'); // Redirect if session doesn't exist
  }

  let userData;

  // // Optionally verify the session on the backend
  const response = await fetch('http://localhost:8081/api/auth/getUser', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', "Access-Control-Allow-Origin": "*" },
    credentials: "include",
  });

  console.log(response)

  if (!response.ok) {
    throw redirect(302, '/login'); // Redirect if session is invalid
  } else {
    userData = await response.json()
    console.log(userData)
  }

  return { userData };
}
