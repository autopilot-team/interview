import { defineConfig } from "vite";

const host = "0.0.0.0";
const port = 2998;

export default defineConfig(async ({ mode }) => {
	return {
		build: {
			outDir: "build/client",
		},
		clearScreen: false,
		server: {
			host,
			port,
		},
	};
});
