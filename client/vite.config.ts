import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  
  plugins: [react()],
  define: {
    'import.meta.env.SERVER_PORT': JSON.stringify('http://127.0.0.1:4000/api')
  },
})
