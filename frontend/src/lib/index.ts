// place files you want to import through the `$lib` alias in this folder.
import { readonly, writable } from 'svelte/store';
import type { UserSessionRes } from './apiResponses';

export const user = writable<UserSessionRes | null>(null);
export const readOnlyUser = readonly(user);