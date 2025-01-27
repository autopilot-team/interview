/// <reference types="vite/client" />
import createFetchClient from "openapi-fetch";
import createClient from "openapi-react-query";
import type { paths as V1Paths } from "./contracts/v1.ts";

export { QueryClientProvider } from "./query-client-provider.tsx";

export const API_BASE_URL =
	import.meta.env.VITE_API_BASE_URL || "http://localhost:3001";

declare global {
	interface ImportMetaEnv {
		/**
		 * The base URL for the API.
		 *
		 * @default "http://localhost:3001"
		 */
		readonly VITE_API_BASE_URL: string;
	}

	interface ImportMeta {
		readonly env: ImportMetaEnv;
	}
}

export const api = {
	v1: createClient(
		createFetchClient<V1Paths>({
			baseUrl: API_BASE_URL,
		}),
	),
};
