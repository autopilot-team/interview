import {
	Activity,
	BarChart3,
	CreditCard,
	Settings,
} from "@autopilot/ui/components/icons";
import { SidebarProvider } from "@autopilot/ui/components/sidebar";
import type { Meta, StoryObj } from "@storybook/react";
import { expect } from "@storybook/test";
import { userEvent, within } from "@storybook/test";
import { createRoutesStub } from "react-router";
import { AppSidebar } from "../../components/app-sidebar/app-sidebar.js";

const meta = {
	title: "Application/AppSidebar",
	component: AppSidebar,
	parameters: {
		layout: "centered",
		viewport: {
			defaultViewport: "desktop",
		},
	},
	decorators: [
		(Story) => {
			const Stub = createRoutesStub([
				{
					path: "/",
					Component: Story,
					children: [
						{ path: "activity", Component: () => null },
						{ path: "balances", Component: () => null },
						{ path: "transactions", Component: () => null },
						{ path: "payment-methods", Component: () => null },
						{ path: "checkout", Component: () => null },
						{ path: "subscriptions", Component: () => null },
						{ path: "products", Component: () => null },
					],
				},
			]);

			return (
				<SidebarProvider>
					<Stub initialEntries={["/"]} />
				</SidebarProvider>
			);
		},
	],
	args: {
		t: {
			nav: {
				overview: {
					title: "Overview",
					activity: "Activity",
					balances: "Balances",
					performance: "Performance",
					settings: "Settings",
				},
				moneyIn: {
					title: "Money In",
					transactions: "Transactions",
					paymentMethods: "Payment Methods",
					checkout: "Checkout",
					subscriptions: "Subscriptions",
					products: "Products",
				},
				moneyFlow: {
					title: "Money Flow",
					accounts: "Accounts",
					transfers: "Transfers",
					treasury: "Treasury",
					fx: "FX",
				},
				moneyOut: {
					title: "Money Out",
					payouts: "Payouts",
					expenses: "Expenses",
					vendors: "Vendors",
					tax: "Tax",
				},
				operations: {
					title: "Operations",
					processors: "Processors",
					ledger: "Ledger",
					reconciliation: "Reconciliation",
					reports: "Reports",
				},
				risk: {
					title: "Risk",
					monitoring: "Monitoring",
					prevention: "Prevention",
					disputes: "Disputes",
					compliance: "Compliance",
					vault: "Vault",
				},
				intelligence: {
					title: "Intelligence",
					analytics: "Analytics",
					insights: "Insights",
					reports: "Reports",
					models: "Models",
				},
				developer: {
					title: "Developer",
					apiKeys: "API Keys",
					webhooks: "Webhooks",
					documentation: "Documentation",
					status: "Status",
				},
			},
			entitySwitcher: {
				searchPlaceholder: "Search entities...",
				noEntities: "No entities found",
				noMatchingEntities: "No matching entities",
				selectEntity: "Select entity",
				platforms: "Platforms",
				organizations: "Organizations",
				accounts: "Accounts",
				addEntity: "Add Entity",
			},
			navUser: {
				upgradeToPro: "Upgrade to Pro",
				account: "Account",
				billing: "Billing",
				notifications: "Notifications",
				logOut: "Log Out",
			},
		},
		entities: [
			{
				id: "1",
				name: "Acme Inc",
				type: "organization",
				logo: Activity,
			},
			{
				id: "2",
				name: "Globex Corp",
				type: "organization",
				logo: Activity,
			},
		],
		currentEntity: {
			id: "1",
			name: "Acme Inc",
			type: "organization",
			logo: Activity,
		},
		user: {
			name: "John Doe",
			email: "john@example.com",
			avatar: "https://avatars.githubusercontent.com/u/1234567",
		},
		navigation: [
			{
				title: "Overview",
				icon: BarChart3,
				items: [
					{ title: "Activity", url: "/activity", icon: Activity },
					{ title: "Balances", url: "/balances", icon: CreditCard },
					{ title: "Performance", url: "/performance", icon: BarChart3 },
					{ title: "Settings", url: "/settings", icon: Settings },
				],
			},
			{
				title: "Money In",
				icon: CreditCard,
				items: [
					{ title: "Transactions", url: "/transactions", icon: Activity },
					{
						title: "Payment Methods",
						url: "/payment-methods",
						icon: CreditCard,
					},
					{ title: "Checkout", url: "/checkout", icon: CreditCard },
					{ title: "Subscriptions", url: "/subscriptions", icon: CreditCard },
					{ title: "Products", url: "/products", icon: CreditCard },
				],
			},
		],
	},
} satisfies Meta<typeof AppSidebar>;

export default meta;
type Story = StoryObj<typeof meta>;

// Default story to show the component
export const Default: Story = {};

// Test navigation and interaction
export const Navigation: Story = {
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Test navigation sections are rendered
		const overviewButton = canvas.getByRole("button", {
			name: /overview/i,
		});
		await expect(overviewButton).toBeInTheDocument();

		const moneyInButton = canvas.getByRole("button", {
			name: /money in/i,
		});
		await expect(moneyInButton).toBeInTheDocument();

		// Test navigation items
		await userEvent.click(overviewButton);
		await new Promise((resolve) => setTimeout(resolve, 300)); // Wait for animation

		// Find links within the Overview section
		let links = canvas.getAllByRole("link");
		const activityLink = links.find((link) =>
			link.textContent?.includes("Activity"),
		);
		const balancesLink = links.find((link) =>
			link.textContent?.includes("Balances"),
		);

		await expect(activityLink).toBeInTheDocument();
		await expect(activityLink).toHaveAttribute("href", "/activity");
		await expect(balancesLink).toBeInTheDocument();
		await expect(balancesLink).toHaveAttribute("href", "/balances");

		// Test Money In section
		await userEvent.click(moneyInButton);
		await new Promise((resolve) => setTimeout(resolve, 300));

		links = canvas.getAllByRole("link");
		const transactionsLink = links.find((link) => {
			return link.textContent?.includes("Transactions");
		});
		await expect(transactionsLink).toBeInTheDocument();
		await expect(transactionsLink).toHaveAttribute("href", "/transactions");
	},
};

// Test entity switcher functionality
export const EntitySwitcher: Story = {
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Find and click entity switcher button
		const entitySwitcherButton = canvas.getByRole("button", {
			name: /select entity/i,
		});
		await expect(entitySwitcherButton).toBeInTheDocument();
		await userEvent.click(entitySwitcherButton);
		await new Promise((resolve) => setTimeout(resolve, 300));

		// Get the portal content
		const portal = within(document.body);
		const entityList = portal.getByRole("listbox");
		await expect(entityList).toBeInTheDocument();

		// Find and verify entities
		const entities = portal.getAllByRole("option");
		const entityNames = entities.map((entity) => entity.textContent);
		expect(entityNames).toContain("Acme Inc");
		expect(entityNames).toContain("Globex Corp");
	},
};

// Test keyboard navigation
export const KeyboardNavigation: Story = {
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Focus first navigation button
		await userEvent.tab();
		const entitySwitcher = canvas.getByRole("button", {
			name: /select entity/i,
		});
		await expect(entitySwitcher).toHaveFocus();

		// Navigate to Overview section
		await userEvent.tab();
		const overviewButton = canvas.getByRole("button", {
			name: /overview/i,
		});
		await expect(overviewButton).toHaveFocus();

		// Open Overview section
		await userEvent.keyboard("{Enter}");
		await new Promise((resolve) => setTimeout(resolve, 300));

		// Navigate to first link
		await userEvent.tab();
		const links = canvas.getAllByRole("link");
		const activityLink = links.find((link) =>
			link.textContent?.includes("Activity"),
		);
		await expect(activityLink).toHaveFocus();
	},
};

// Test responsive variants
export const Mobile: Story = {
	parameters: {
		viewport: {
			defaultViewport: "mobile1",
		},
	},
};

export const Tablet: Story = {
	parameters: {
		viewport: {
			defaultViewport: "tablet",
		},
	},
};

export const Desktop: Story = {
	parameters: {
		viewport: {
			defaultViewport: "desktop",
		},
	},
};
