import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [tailwindcss()],
  publicDir: "public",
  server: {
    port: 3001,
    strictPort: true,
    cors: {
      origin: "*",
    },
    hmr: {
      host: "localhost",
      port: 3001,
    },
  },
  build: {
    manifest: true,
    rollupOptions: {
      input: "src/main.ts",
      output: {
        entryFileNames: "assets/[name]-[hash].js",
        chunkFileNames: "assets/[name]-[hash].js",
        assetFileNames: "assets/[name]-[hash][extname]",
      },
    },
    outDir: "../static",
  },
});

