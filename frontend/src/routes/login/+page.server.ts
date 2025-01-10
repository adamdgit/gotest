import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
  // Simulate fetching data from an API or database
  return { title: 'SSR Example', content: 'This page is server-side rendered.' };
};