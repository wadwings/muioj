import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { version } from './package.json';
import { compilerOptions } from './tsconfig.json';
// https://vitejs.dev/config/
export default defineConfig({
  define: {
    __APP_VERSION__: `"${version}"`,
  },
  optimizeDeps: {
    entries: [
      "./src/main.tsx",
    ]
  },
  resolve: {
    alias: Object.fromEntries(
      Object.entries(compilerOptions.paths).map(([key, value]) => [
        key.replace('*', ''),
        `/${value[0].replace('*', '')}`,
      ])
    ),
  },
  build: {
    outDir: 'dist',
    minify: 'esbuild',
  },
  plugins: [
    react()
  ],
})
