/// <reference types="vite/client" />
import createFetchClient from "openapi-fetch";
import createClient from "openapi-react-query";
import type { paths as V1Paths } from "./contracts/v1.ts";

export * as v1 from "./contracts/v1.ts";
export { QueryClientProvider } from "./query-client-provider.tsx";

export const API_BASE_URL =
	import.meta.env.VITE_API_URL || "http://localhost:3001";
export const ASSETS_BASE_URL =
	import.meta.env.VITE_ASSETS_URL || "http://localhost:2998";

/**
 * Returns the full URL for an asset path based on the environment.
 * Uses Vite's environment variables to determine the assets URL.
 *
 * @param path - The path to the asset, without leading slash
 * @returns The full URL to the asset
 */
export function asset(path: string): string {
	const cleanPath = path.startsWith("/") ? path.slice(1) : path;
	return `${ASSETS_BASE_URL}/${cleanPath}`;
}

// Track if we're currently refreshing to prevent multiple refresh calls
let isRefreshing = false;
// Queue of requests to retry after refresh
const refreshQueue: Array<{
	resolve: (value: Response | PromiseLike<Response>) => void;
	reject: (reason?: unknown) => void;
	request: Request;
}> = [];

// Process the queue of failed requests
const processQueue = (error?: unknown) => {
	if (error) {
		// If there was an error refreshing, reject all queued requests
		for (const { reject } of refreshQueue) {
			reject(error);
		}
	} else {
		// If refresh was successful, retry all queued requests
		for (const { resolve, reject, request } of refreshQueue) {
			// We catch each retry individually to prevent one failure from affecting others
			void retry(request).then(resolve).catch(reject);
		}
	}
	refreshQueue.splice(0, refreshQueue.length);
};

// Retry a request with updated session
const retry = async (request: Request): Promise<Response> => {
	const retryResponse = await fetch(request);
	if (!retryResponse.ok) {
		throw new Error(`Request failed with status ${retryResponse.status}`);
	}
	return retryResponse;
};

// Create a custom fetch function that handles token refresh
const createCustomFetch = () => {
	return async (input: Request): Promise<Response> => {
		const response = await fetch(input);

		// Return early for successful responses or identity endpoints
		if (response.ok || input.url.includes("/v1/identity")) {
			return response;
		}

		// Handle 401 Unauthorized errors
		if (response.status === 401) {
			if (!isRefreshing) {
				isRefreshing = true;

				try {
					// Attempt to refresh the token
					const refreshResponse = await fetch(
						`${API_BASE_URL}/v1/identity/refresh-session`,
						{
							method: "POST",
							credentials: "include",
						},
					);

					if (!refreshResponse.ok) {
						throw new Error(
							`Token refresh failed with status ${refreshResponse.status}`,
						);
					}

					isRefreshing = false;
					processQueue();

					// Retry the original request
					return retry(input);
				} catch (error) {
					isRefreshing = false;
					processQueue(error);
					throw error;
				}
			}

			// If we're already refreshing, queue this request
			return new Promise((resolve, reject) => {
				refreshQueue.push({ resolve, reject, request: input });
			});
		}

		// If it's not a 401 but still an error, throw it
		if (!response.ok) {
			throw new Error(`Request failed with status ${response.status}`);
		}

		return response;
	};
};

export const api = {
	v1: createClient(
		createFetchClient<V1Paths>({
			baseUrl: API_BASE_URL,
			credentials: "include",
			fetch: createCustomFetch(),
		}),
	),
};

declare global {
	interface ImportMetaEnv {
		/**
		 * The base URL for the API.
		 *
		 * @default "http://localhost:3001"
		 */
		readonly VITE_API_URL: string;

		/**
		 * The base URL for static assets.
		 *
		 * @default "http://localhost:2998"
		 */
		readonly VITE_ASSETS_URL: string;
	}

	interface ImportMeta {
		readonly env: ImportMetaEnv;
	}
}
