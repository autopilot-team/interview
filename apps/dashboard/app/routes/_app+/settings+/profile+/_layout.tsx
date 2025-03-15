import {
	KeyIcon,
	ShieldCheckIcon,
	UserIcon,
} from "@autopilot/ui/components/icons";
import {
	Tabs,
	TabsContent,
	TabsList,
	TabsTrigger,
} from "@autopilot/ui/components/tabs";
import { useTranslation } from "react-i18next";
import { Link, Outlet, useLocation } from "react-router";

export default function Layout() {
	const { t } = useTranslation(["common"]);
	const location = useLocation();
	const currentPath = location.pathname.split("/").pop() || "";
	const activeTab = currentPath === "profile" ? "index" : currentPath;
	const tabs = [
		{
			value: "index",
			icon: <UserIcon className="h-4 w-4" />,
			label: t("settings.tabs.profile"),
			to: ".",
		},
		{
			value: "security",
			icon: <ShieldCheckIcon className="h-4 w-4" />,
			label: t("settings.tabs.security"),
			to: "security",
		},
		{
			value: "sessions",
			icon: <KeyIcon className="h-4 w-4" />,
			label: t("settings.tabs.sessions"),
			to: "sessions",
		},
	];

	return (
		<div className="p-5">
			<div className="mb-8">
				<h1 className="text-3xl font-bold">{t("settings.title")}</h1>
				<p className="text-muted-foreground">{t("settings.description")}</p>
			</div>

			<Tabs value={activeTab} className="space-y-6">
				<TabsList className="h-auto w-full justify-start rounded-none border-b bg-transparent p-0">
					{tabs.map((tab) => (
						<TabsTrigger
							key={tab.value}
							value={tab.value}
							className="relative h-auto w-auto rounded-none border-b-2 border-transparent bg-transparent px-4 pb-3 pt-2 font-medium ring-offset-0 hover:bg-transparent hover:text-foreground focus-visible:ring-0 data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:shadow-none"
							asChild
						>
							<Link to={tab.to}>
								<div className="flex items-center gap-2">
									{tab.icon}
									{tab.label}
								</div>
							</Link>
						</TabsTrigger>
					))}
				</TabsList>

				<TabsContent value={activeTab}>
					<Outlet />
				</TabsContent>
			</Tabs>
		</div>
	);
}
