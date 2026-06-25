import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "node:path";
import { fileURLToPath } from "node:url";
var __dirname = path.dirname(fileURLToPath(import.meta.url));
export default defineConfig({
    plugins: [react()],
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./src"),
        },
    },
    server: {
        port: 5173,
        proxy: {
            "/api": {
                target: "http://127.0.0.1:17878",
                changeOrigin: true,
            },
        },
    },
    build: {
        outDir: "dist",
        emptyOutDir: true,
        rollupOptions: {
            output: {
                manualChunks: function (id) {
                    // Keep chunking coarse-grained: isolate the heaviest markdown stack,
                    // then split framework/runtime and UI libraries into stable vendor bundles.
                    if (id.includes("react-markdown") ||
                        id.includes("remark-gfm") ||
                        id.includes("remark-breaks") ||
                        id.includes("rehype-raw")) {
                        return "markdown";
                    }
                    if (id.includes("/node_modules/react/") ||
                        id.includes("/node_modules/react-dom/") ||
                        id.includes("/node_modules/react-router-dom/") ||
                        id.includes("/node_modules/react-router/")) {
                        return "react-vendor";
                    }
                    if (id.includes("/node_modules/lucide-react/") ||
                        id.includes("/node_modules/@radix-ui/") ||
                        id.includes("/node_modules/embla-carousel-react/") ||
                        id.includes("/node_modules/sonner/") ||
                        id.includes("/node_modules/class-variance-authority/") ||
                        id.includes("/node_modules/clsx/") ||
                        id.includes("/node_modules/tailwind-merge/")) {
                        return "ui-vendor";
                    }
                },
            },
        },
    },
});
