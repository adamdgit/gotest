import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import { VitePWA } from 'vite-plugin-pwa';

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    svelte(),
    VitePWA({
      devOptions: {
        enabled: true
      },
      injectRegister: 'auto',
      workbox: {
        globPatterns: ['**/*.{js,css,html,ico,png,svg,webp}']
      },
      includeAssets: ['vite.svg'],
      manifest: {
        name: 'My Svelte PWA',
        short_name: 'SveltePWA',
        description: 'A simple Progressive Web App with offline support.',
        theme_color: '#ffffff',
        background_color: '#ffffff',
        display: 'standalone',
        start_url: '/',
        icons: [
          {
            src: 'pwa-icon-192.webp',
            sizes: '192x192',
            type: 'image/webp',
          },
          {
            src: 'pwa-icon-512.webp',
            sizes: '512x512',
            type: 'image/webp',
          },
        ],
      },
      registerType: 'autoUpdate', // Automatically update the service worker
    }),
  ],
});