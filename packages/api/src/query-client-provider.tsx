import {
	QueryClient,
	QueryClientProvider as TanstackQueryClientProvider,
} from "@tanstack/react-query";
import { useState } from "react";

export function QueryClientProvider({
	children,
}: { children: React.ReactNode }) {
	const [queryClient] = useState(
		() =>
			new QueryClient({
				defaultOptions: {
					queries: {
						staleTime: 0, // Consider data stale immediately
						gcTime: 60 * 1000, // Keep inactive queries for 1 minute only
						retry: 2, // Retry failed requests twice
						retryDelay: (attemptIndex) =>
							Math.min(1000 * 2 ** attemptIndex, 30000),
						refetchOnWindowFocus: true, // Refetch when user returns to the app
						refetchOnReconnect: true, // Refetch on network reconnection
						refetchOnMount: true, // Always refetch on component mount
					},
					mutations: {
						retry: 1, // Retry failed mutations once
						retryDelay: 1000, // Wait 1 second before retrying
					},
				},
			}),
	);

	return (
		<TanstackQueryClientProvider client={queryClient}>
			{children}
		</TanstackQueryClientProvider>
	);
}
