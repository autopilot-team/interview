import {
	Links,
	Meta,
	Outlet,
	Scripts,
	ScrollRestoration,
	isRouteErrorResponse,
} from "react-router";
import "@autopilot/ui/globals.css";
import { IdentityProvider, useIdentity } from "@/components/identity-provider";
import { QueryClientProvider } from "@autopilot/api";
import {
	MessageCard,
	type MessageCardVariant,
} from "@autopilot/ui/components/message-card";
import { Toaster as SonnerToaster } from "@autopilot/ui/components/sonner";
import { Toaster } from "@autopilot/ui/components/toaster";
import { useTranslation } from "react-i18next";
import type { Route } from "./+types/root";

export const links: Route.LinksFunction = () => [
	{
		rel: "preconnect",
		href: "https://fonts.googleapis.com",
	},
	{
		rel: "preconnect",
		href: "https://fonts.gstatic.com",
		crossOrigin: "anonymous",
	},
	{
		rel: "stylesheet",
		href: "https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&display=swap",
	},
];

export function meta({}: Route.MetaArgs) {
	return [
		{ title: "Merchants Dashboard | Autopilot" },
		{
			name: "description",
			content: "A dashboard for merchants to manage their entities.",
		},
	];
}

export function Layout({ children }: { children: React.ReactNode }) {
	return (
		<html lang="en">
			<head>
				<meta charSet="utf-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1" />
				<Meta />
				<Links />
			</head>

			<body>
				<QueryClientProvider>
					<IdentityProvider>
						{children}
						<SonnerToaster />
						<Toaster />
					</IdentityProvider>
				</QueryClientProvider>
				<ScrollRestoration />
				<Scripts />
			</body>
		</html>
	);
}

export default function App() {
	const { isInitializing, user } = useIdentity();
	const { t, ready } = useTranslation(["common"]);

	if (isInitializing || typeof user === "undefined") {
		return (
			<LoadingScreen
				message={ready ? t("common:loading.initializing") : "Initializing..."}
			/>
		);
	}

	return <Outlet />;
}

export function ErrorBoundary({ error }: Route.ErrorBoundaryProps) {
	const { t, ready } = useTranslation();

	if (typeof window === "undefined" || !ready) {
		return (
			<MessageCard
				title="Not found"
				description="This page doesn't exist or you don't have access."
				variant="default"
			/>
		);
	}

	let title = t("error.unexpected.title");
	let description = t("error.unexpected.description");
	let stack: string | undefined;
	let variant: MessageCardVariant = "default";

	if (isRouteErrorResponse(error)) {
		if (error.status === 404) {
			title = t("error.notFound.title");
			description = t("error.notFound.description");
			variant = "default";
		} else {
			title = `${error.status} Error`;
			description = error.statusText || description;
			variant = "error";
		}
	} else if (import.meta.env.DEV && error instanceof Error) {
		title = error.name;
		description = error.message;
		stack = error.stack;
		variant = "error";
	}

	return (
		<MessageCard
			backButton={{ label: t("error.action.back", "Back"), show: true }}
			description={description}
			footer={t(
				"error.support.text",
				"If this error persists, please contact support.",
			)}
			homeButton={{
				label: t("error.action.home", "Home"),
				show: true,
				to: "/",
			}}
			stack={stack}
			title={title}
			variant={variant}
		/>
	);
}

function LoadingSpinner() {
	return (
		<div className="relative flex items-center justify-center">
			<div className="absolute size-16 animate-[spin_2s_linear_infinite] rounded-full border-[3px] border-primary/5 border-t-primary/40" />
			<div className="absolute size-12 animate-[spin_1.5s_linear_infinite] rounded-full border-[3px] border-primary/10 border-t-primary/60" />
			<div className="absolute size-8 animate-[spin_1s_linear_infinite] rounded-full border-[3px] border-primary/20 border-t-primary" />
			<div className="size-2 rounded-full bg-primary animate-pulse" />
		</div>
	);
}

function LoadingScreen({ message }: { message: string }) {
	return (
		<div className="fixed inset-0 flex flex-col items-center justify-center gap-2 bg-background">
			<div className="flex flex-col items-center gap-6">
				<LoadingSpinner />
				<p className="text-sm font-medium text-muted-foreground mt-3">
					{message}
				</p>
			</div>

			<div className="h-1 w-48 overflow-hidden rounded-full bg-primary/10">
				<div className="h-full w-3/4 rounded-full bg-primary/50" />
			</div>
		</div>
	);
}

export function HydrateFallback() {
	if (typeof window === "undefined") {
		return <LoadingScreen message="Initializing..." />;
	}

	const { t, ready } = useTranslation(["common"]);

	return (
		<LoadingScreen
			message={ready ? t("common:loading.initializing") : "Initializing..."}
		/>
	);
}
