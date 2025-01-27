import {
	AppSidebar,
	type Entity,
} from "@autopilot/ui/components/app-sidebar/index";
import {
	Breadcrumb,
	BreadcrumbItem,
	BreadcrumbLink,
	BreadcrumbList,
	BreadcrumbPage,
	BreadcrumbSeparator,
} from "@autopilot/ui/components/breadcrumb";
import {
	Activity,
	AlertOctagon,
	ArrowRightLeft,
	BookOpen,
	Brain,
	Building2,
	Calculator,
	CheckCircle2,
	CircleDollarSign,
	Code2,
	CreditCard,
	DollarSign,
	FileBarChart,
	FileCode,
	Gauge,
	Globe,
	Key,
	Landmark,
	LayoutDashboard,
	Lock,
	Network,
	Package,
	PieChart,
	Receipt,
	RefreshCw,
	Settings,
	ShieldAlert,
	ShieldCheck,
	Store,
	Users,
	Wallet,
	Wand2,
	Webhook,
} from "@autopilot/ui/components/icons";
import { Separator } from "@autopilot/ui/components/separator";
import {
	SidebarInset,
	SidebarProvider,
	SidebarTrigger,
} from "@autopilot/ui/components/sidebar";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Outlet } from "react-router";

// Example entities following the setup examples
const entities: Entity[] = [
	// Example 1: Single Merchant
	{
		id: "a1",
		name: "Boutique Store",
		type: "account",
		logo: Store,
	},

	// Example 2: Standalone Organizations
	{
		id: "o9",
		name: "Burger King",
		type: "organization",
		logo: Building2,
	},
	{
		id: "a19",
		name: "Times Square",
		type: "account",
		parentId: "o9",
		logo: Store,
	},
	{
		id: "a20",
		name: "Chicago Downtown",
		type: "account",
		parentId: "o9",
		logo: Store,
	},
	{
		id: "a21",
		name: "Miami Beach",
		type: "account",
		parentId: "o9",
		logo: Store,
	},

	{
		id: "o10",
		name: "Dunkin Donuts",
		type: "organization",
		logo: Building2,
	},
	{
		id: "a22",
		name: "Boston Central",
		type: "account",
		parentId: "o10",
		logo: Store,
	},
	{
		id: "a23",
		name: "Manhattan West",
		type: "account",
		parentId: "o10",
		logo: Store,
	},
	{
		id: "a24",
		name: "Brooklyn Heights",
		type: "account",
		parentId: "o10",
		logo: Store,
	},

	// Example 3: Franchise Model (McDonald's)
	{
		id: "p1",
		name: "McDonald's Global Payment System",
		type: "platform",
		logo: Building2,
	},
	{
		id: "o1",
		name: "McDonald's USA",
		type: "organization",
		parentId: "p1",
		logo: Globe,
	},
	{
		id: "a2",
		name: "NYC Franchise",
		type: "account",
		parentId: "o1",
		logo: Store,
	},
	{
		id: "a3",
		name: "LA Franchise",
		type: "account",
		parentId: "o1",
		logo: Store,
	},
	{
		id: "o2",
		name: "McDonald's Europe",
		type: "organization",
		parentId: "p1",
		logo: Globe,
	},
	{
		id: "a4",
		name: "Paris Franchise",
		type: "account",
		parentId: "o2",
		logo: Store,
	},
	{
		id: "a5",
		name: "Berlin Franchise",
		type: "account",
		parentId: "o2",
		logo: Store,
	},

	// Example 4: White-Label SaaS (PayRetailer)
	{
		id: "p2",
		name: "PayRetailer",
		type: "platform",
		logo: Building2,
	},
	{
		id: "o3",
		name: "ShopEasy",
		type: "organization",
		parentId: "p2",
		logo: Building2,
	},
	{
		id: "a6",
		name: "Boutique Clothing Store",
		type: "account",
		parentId: "o3",
		logo: Store,
	},
	{
		id: "a7",
		name: "Handmade Goods Store",
		type: "account",
		parentId: "o3",
		logo: Store,
	},
	{
		id: "o4",
		name: "RetailPro",
		type: "organization",
		parentId: "p2",
		logo: Building2,
	},
	{
		id: "a8",
		name: "Electronics Store",
		type: "account",
		parentId: "o4",
		logo: Store,
	},
	{
		id: "a9",
		name: "Furniture Outlet",
		type: "account",
		parentId: "o4",
		logo: Store,
	},
	{
		id: "a10",
		name: "OwnDays Store",
		type: "account",
		parentId: "p2",
		logo: Store,
	},

	// Example 5: Marketplace (Amazon)
	{
		id: "p3",
		name: "Amazon Payments",
		type: "platform",
		logo: Building2,
	},
	{
		id: "o5",
		name: "Third-Party Sellers",
		type: "organization",
		parentId: "p3",
		logo: Building2,
	},
	{
		id: "a11",
		name: "Electronics Seller",
		type: "account",
		parentId: "o5",
		logo: Store,
	},
	{
		id: "a12",
		name: "Clothing Seller",
		type: "account",
		parentId: "o5",
		logo: Store,
	},
	{
		id: "o6",
		name: "Amazon Basics",
		type: "organization",
		parentId: "p3",
		logo: Building2,
	},
	{
		id: "a13",
		name: "Home Appliances",
		type: "account",
		parentId: "o6",
		logo: Store,
	},
	{
		id: "a14",
		name: "Office Supplies",
		type: "account",
		parentId: "o6",
		logo: Store,
	},

	// Example 6: Multi-Regional Business (Starbucks)
	{
		id: "p4",
		name: "Starbucks Financial Platform",
		type: "platform",
		logo: Building2,
	},
	{
		id: "o7",
		name: "Starbucks USA",
		type: "organization",
		parentId: "p4",
		logo: Globe,
	},
	{
		id: "a15",
		name: "NYC Store",
		type: "account",
		parentId: "o7",
		logo: Store,
	},
	{
		id: "a16",
		name: "LA Store",
		type: "account",
		parentId: "o7",
		logo: Store,
	},
	{
		id: "o8",
		name: "Starbucks Asia-Pacific",
		type: "organization",
		parentId: "p4",
		logo: Globe,
	},
	{
		id: "a17",
		name: "Tokyo Store",
		type: "account",
		parentId: "o8",
		logo: Store,
	},
	{
		id: "a18",
		name: "Singapore Store",
		type: "account",
		parentId: "o8",
		logo: Store,
	},
];

const user = {
	name: "shadcn",
	email: "m@example.com",
	avatar: "/avatars/shadcn.jpg",
};

export function clientLoader() {
	return {};
}

export default function Component() {
	const { t } = useTranslation(["common"]);
	const [currentEntity, setCurrentEntity] = useState<Entity | undefined>(
		entities[0],
	);
	const navigation = [
		{
			title: t("common:nav.overview.title"),
			icon: LayoutDashboard,
			items: [
				{
					title: t("common:nav.overview.activity"),
					url: "/overview/activity",
					icon: Activity,
				},
				{
					title: t("common:nav.overview.balances"),
					url: "/overview/balances",
					icon: DollarSign,
				},
				{
					title: t("common:nav.overview.performance"),
					url: "/overview/performance",
					icon: Gauge,
				},
				{
					title: t("common:nav.overview.settings"),
					url: "/overview/settings",
					icon: Settings,
				},
			],
		},
		{
			title: t("common:nav.moneyIn.title"),
			icon: Wallet,
			items: [
				{
					title: t("common:nav.moneyIn.transactions"),
					url: "/money-in/transactions",
					icon: ArrowRightLeft,
				},
				{
					title: t("common:nav.moneyIn.paymentMethods"),
					url: "/money-in/methods",
					icon: CreditCard,
				},
				{
					title: t("common:nav.moneyIn.checkout"),
					url: "/money-in/checkout",
					icon: Store,
				},
				{
					title: t("common:nav.moneyIn.subscriptions"),
					url: "/money-in/subscriptions",
					icon: RefreshCw,
				},
				{
					title: t("common:nav.moneyIn.products"),
					url: "/money-in/products",
					icon: Package,
				},
			],
		},
		{
			title: t("common:nav.moneyFlow.title"),
			icon: CircleDollarSign,
			items: [
				{
					title: t("common:nav.moneyFlow.accounts"),
					url: "/money-flow/accounts",
					icon: BookOpen,
				},
				{
					title: t("common:nav.moneyFlow.transfers"),
					url: "/money-flow/transfers",
					icon: ArrowRightLeft,
				},
				{
					title: t("common:nav.moneyFlow.treasury"),
					url: "/money-flow/treasury",
					icon: Landmark,
				},
				{
					title: t("common:nav.moneyFlow.fx"),
					url: "/money-flow/fx",
					icon: Globe,
				},
			],
		},
		{
			title: t("common:nav.moneyOut.title"),
			icon: Building2,
			items: [
				{
					title: t("common:nav.moneyOut.payouts"),
					url: "/money-out/payouts",
					icon: CircleDollarSign,
				},
				{
					title: t("common:nav.moneyOut.expenses"),
					url: "/money-out/expenses",
					icon: Receipt,
				},
				{
					title: t("common:nav.moneyOut.vendors"),
					url: "/money-out/vendors",
					icon: Users,
				},
				{
					title: t("common:nav.moneyOut.tax"),
					url: "/money-out/tax",
					icon: Calculator,
				},
			],
		},
		{
			title: t("common:nav.operations.title"),
			icon: Settings,
			items: [
				{
					title: t("common:nav.operations.processors"),
					url: "/operations/processors",
					icon: Network,
				},
				{
					title: t("common:nav.operations.ledger"),
					url: "/operations/ledger",
					icon: BookOpen,
				},
				{
					title: t("common:nav.operations.reconciliation"),
					url: "/operations/reconciliation",
					icon: CheckCircle2,
				},
				{
					title: t("common:nav.operations.reports"),
					url: "/operations/reports",
					icon: FileBarChart,
				},
			],
		},
		{
			title: t("common:nav.risk.title"),
			icon: ShieldAlert,
			items: [
				{
					title: t("common:nav.risk.monitoring"),
					url: "/risk/monitoring",
					icon: Activity,
				},
				{
					title: t("common:nav.risk.prevention"),
					url: "/risk/prevention",
					icon: ShieldCheck,
				},
				{
					title: t("common:nav.risk.disputes"),
					url: "/risk/disputes",
					icon: AlertOctagon,
				},
				{
					title: t("common:nav.risk.compliance"),
					url: "/risk/compliance",
					icon: Lock,
				},
				{
					title: t("common:nav.risk.vault"),
					url: "/risk/vault",
					icon: Lock,
				},
			],
		},
		{
			title: t("common:nav.intelligence.title"),
			icon: Brain,
			items: [
				{
					title: t("common:nav.intelligence.analytics"),
					url: "/intelligence/analytics",
					icon: PieChart,
				},
				{
					title: t("common:nav.intelligence.insights"),
					url: "/intelligence/insights",
					icon: Wand2,
				},
				{
					title: t("common:nav.intelligence.reports"),
					url: "/intelligence/reports",
					icon: FileBarChart,
				},
				{
					title: t("common:nav.intelligence.models"),
					url: "/intelligence/models",
					icon: Brain,
				},
			],
		},
		{
			title: t("common:nav.developer.title"),
			icon: Code2,
			items: [
				{
					title: t("common:nav.developer.apiKeys"),
					url: "/developer/api-keys",
					icon: Key,
				},
				{
					title: t("common:nav.developer.webhooks"),
					url: "/developer/webhooks",
					icon: Webhook,
				},
				{
					title: t("common:nav.developer.documentation"),
					url: "/developer/docs",
					icon: FileCode,
				},
				{
					title: t("common:nav.developer.status"),
					url: "/developer/status",
					icon: Activity,
				},
			],
		},
	];

	const handleCreateEntity = () => {};

	return (
		<SidebarProvider>
			<AppSidebar
				t={{
					nav: t("common:nav", { returnObjects: true }),
					entitySwitcher: t("common:entitySwitcher", { returnObjects: true }),
					navUser: t("common:navUser", { returnObjects: true }),
				}}
				entities={entities}
				currentEntity={currentEntity}
				onEntityChange={setCurrentEntity}
				onCreateEntity={handleCreateEntity}
				user={user}
				navigation={navigation}
			/>

			<SidebarInset>
				<header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
					<div className="flex items-center gap-2 px-4">
						<SidebarTrigger className="-ml-1" />

						<Separator orientation="vertical" className="mr-2 h-4" />

						<Breadcrumb>
							<BreadcrumbList>
								<BreadcrumbItem className="hidden md:block">
									<BreadcrumbLink href="#">
										Building Your Application
									</BreadcrumbLink>
								</BreadcrumbItem>

								<BreadcrumbSeparator className="hidden md:block" />

								<BreadcrumbItem>
									<BreadcrumbPage>Data Fetching</BreadcrumbPage>
								</BreadcrumbItem>
							</BreadcrumbList>
						</Breadcrumb>
					</div>
				</header>

				<div className="flex flex-1 flex-col gap-4 p-4 pt-0">
					<Outlet />
				</div>
			</SidebarInset>
		</SidebarProvider>
	);
}
