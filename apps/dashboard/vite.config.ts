import { reactRouter } from "@react-router/dev/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig(async ({ mode }) => {
	const devPlugins = [];

	if (mode === "development") {
		const { i18nextHMRPlugin } = await import("i18next-hmr/vite");

		devPlugins.push(
			i18nextHMRPlugin({
				localesDir: "./public/locales",
			}),
		);
	}

	return {
		clearScreen: false,
		plugins: [reactRouter(), tailwindcss(), tsconfigPaths()].concat(devPlugins),
	};
});
