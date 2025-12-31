// place files you want to import through the `$lib` alias in this folder.
import { readonly, writable } from 'svelte/store';

export type User = {
	email: string;
	role: string;
	profile_url?: string;
};

export const user = writable<User | null>(null);
export const readOnlyUser = readonly(user);