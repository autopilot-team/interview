import { type ClientLoaderFunctionArgs, Outlet, redirect } from "react-router";
import type { RouteHandle } from "@/routes/_app+/_layout";

export const handle = {
	breadcrumb: "common:breadcrumbs.settings.title",
} satisfies RouteHandle;

export function clientLoader({ request }: ClientLoaderFunctionArgs) {
	const url = new URL(request.url);
	const pathname = url.pathname;

	if (pathname === "/settings") {
		throw redirect("/settings/profile");
	}
}

export default function Component() {
	return <Outlet />;
}
