import { Brand } from "@autopilot/ui/components/brand";
import { Outlet } from "react-router";
import { useIdentity } from "@/components/identity-provider";

export default function Layout() {
	const { entity } = useIdentity();

	return (
		<div className="flex-1 grid overflow-y-auto overflow-x-hidden z-0 h-[100dvh]">
			<header className="sticky top-0 z-10 flex items-center p-5 border-b bg-background">
				<Brand href={`/${entity?.slug}`} spaNavigation={false} />
			</header>

			<Outlet />
		</div>
	);
}
