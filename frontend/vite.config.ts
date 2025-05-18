import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { VitePWA } from 'vite-plugin-pwa'

// https://vite.dev/config/
export default defineConfig({
	plugins: [
		react(),
		VitePWA({
			registerType: 'autoUpdate',
			manifest: {
				name: 'TipaTwitter',
				short_name: 'TT',
				description: 'Социальная сеть для коротких постов',
				theme_color: '#ffffff',
				icons: [
					{
						src: '/vite.svg',
						sizes: '192x192',
						type: 'image/svg+xml',
					},
					{
						src: '/vite.svg',
						sizes: '512x512',
						type: 'image/svg+xml',
					},
				],
			},
		}),
	],
})
