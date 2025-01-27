import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import { Folder, House, Settings } from "lucide-react";
import {
	Breadcrumb,
	BreadcrumbEllipsis,
	BreadcrumbItem,
	BreadcrumbLink,
	BreadcrumbList,
	BreadcrumbPage,
	BreadcrumbSeparator,
} from "../../components/breadcrumb.js";

const meta: Meta<typeof Breadcrumb> = {
	title: "Design System/Breadcrumb",
	component: Breadcrumb,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Breadcrumb>;

export const Default: Story = {
	render: () => (
		<Breadcrumb>
			<BreadcrumbList>
				<BreadcrumbItem>
					<BreadcrumbLink href="/">Home</BreadcrumbLink>
				</BreadcrumbItem>
				<BreadcrumbSeparator />
				<BreadcrumbItem>
					<BreadcrumbLink href="/docs">Documentation</BreadcrumbLink>
				</BreadcrumbItem>
				<BreadcrumbSeparator />
				<BreadcrumbItem>
					<BreadcrumbPage>Components</BreadcrumbPage>
				</BreadcrumbItem>
			</BreadcrumbList>
		</Breadcrumb>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const home = canvas.getByText("Home");
		const docs = canvas.getByText("Documentation");
		const components = canvas.getByText("Components");

		await expect(home).toBeInTheDocument();
		await expect(docs).toBeInTheDocument();
		await expect(components).toBeInTheDocument();
		await expect(components).toHaveAttribute("aria-current", "page");
	},
};

export const WithIcons: Story = {
	render: () => (
		<Breadcrumb>
			<BreadcrumbList>
				<BreadcrumbItem>
					<BreadcrumbLink href="/" className="flex items-center gap-2">
						<House className="h-4 w-4" />
						Home
					</BreadcrumbLink>
				</BreadcrumbItem>
				<BreadcrumbSeparator />
				<BreadcrumbItem>
					<BreadcrumbLink href="/files" className="flex items-center gap-2">
						<Folder className="h-4 w-4" />
						Files
					</BreadcrumbLink>
				</BreadcrumbItem>
				<BreadcrumbSeparator />
				<BreadcrumbItem>
					<BreadcrumbLink className="flex items-center gap-2">
						<Settings className="h-4 w-4" />
						Settings
					</BreadcrumbLink>
				</BreadcrumbItem>
			</BreadcrumbList>
		</Breadcrumb>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Check navigation
		const nav = canvas.getByRole("navigation", { name: "breadcrumb" });
		await expect(nav).toBeInTheDocument();

		// Check links with icons
		const homeLink = canvas.getByRole("link", { name: /home/i });
		await expect(homeLink).toBeInTheDocument();
		await expect(
			homeLink.querySelector("svg.lucide-house"),
		).toBeInTheDocument();

		const filesLink = canvas.getByRole("link", { name: /files/i });
		await expect(filesLink).toBeInTheDocument();
		await expect(
			filesLink.querySelector("svg.lucide-folder"),
		).toBeInTheDocument();

		const settingsLink = canvas.getByText("Settings").closest("a");
		await expect(settingsLink).toBeInTheDocument();
		await expect(
			settingsLink?.querySelector("svg.lucide-settings"),
		).toBeInTheDocument();

		// Check separators
		const list = canvas.getByRole("list");
		const separators = Array.from(list.children).filter((item) =>
			item.querySelector("svg.lucide-chevron-right"),
		);
		await expect(separators).toHaveLength(2);
	},
};

export const WithEllipsis: Story = {
	render: () => (
		<Breadcrumb>
			<BreadcrumbList>
				<BreadcrumbItem>
					<BreadcrumbLink href="/">Home</BreadcrumbLink>
				</BreadcrumbItem>
				<BreadcrumbSeparator />
				<BreadcrumbItem>
					<BreadcrumbEllipsis />
				</BreadcrumbItem>
				<BreadcrumbSeparator />
				<BreadcrumbItem>
					<BreadcrumbLink href="/docs/components">Components</BreadcrumbLink>
				</BreadcrumbItem>
				<BreadcrumbSeparator />
				<BreadcrumbItem>
					<BreadcrumbPage>Button</BreadcrumbPage>
				</BreadcrumbItem>
			</BreadcrumbList>
		</Breadcrumb>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Check navigation structure
		const nav = canvas.getByRole("navigation", { name: "breadcrumb" });
		await expect(nav).toBeInTheDocument();

		// Check list structure
		const list = canvas.getByRole("list");
		const items = canvas.getAllByRole("listitem");
		await expect(items).toHaveLength(4); // Home, Ellipsis, Components, Button

		// Check links and current page
		const homeLink = canvas.getByRole("link", { name: "Home" });
		const componentsLink = canvas.getByRole("link", { name: "Components" });
		await expect(homeLink).toBeInTheDocument();
		await expect(componentsLink).toBeInTheDocument();
		await expect(canvas.getByText("Button")).toHaveAttribute(
			"aria-current",
			"page",
		);

		// Check ellipsis is between Home and Components
		const itemsArray = Array.from(list.children);
		const ellipsisIndex = itemsArray.findIndex((item) =>
			item.textContent?.includes("More"),
		);
		await expect(ellipsisIndex).toBeGreaterThan(0);
		await expect(ellipsisIndex).toBeLessThan(itemsArray.length - 1);
	},
};

export const LongPath: Story = {
	render: () => (
		<Breadcrumb>
			<BreadcrumbList>
				<BreadcrumbItem>
					<BreadcrumbLink href="/">Home</BreadcrumbLink>
				</BreadcrumbItem>
				<BreadcrumbSeparator />
				<BreadcrumbItem>
					<BreadcrumbLink href="/users">Users</BreadcrumbLink>
				</BreadcrumbItem>
				<BreadcrumbSeparator />
				<BreadcrumbItem>
					<BreadcrumbLink href="/users/settings">Settings</BreadcrumbLink>
				</BreadcrumbItem>
				<BreadcrumbSeparator />
				<BreadcrumbItem>
					<BreadcrumbLink href="/users/settings/profile">
						Profile
					</BreadcrumbLink>
				</BreadcrumbItem>
				<BreadcrumbSeparator />
				<BreadcrumbItem>
					<BreadcrumbPage>Edit</BreadcrumbPage>
				</BreadcrumbItem>
			</BreadcrumbList>
		</Breadcrumb>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Check navigation structure
		const nav = canvas.getByRole("navigation", { name: "breadcrumb" });
		await expect(nav).toBeInTheDocument();

		// Check list items (actual navigation items, not separators)
		const items = canvas.getAllByRole("listitem");
		await expect(items).toHaveLength(5); // Home, Users, Settings, Profile, Edit

		// Verify all links are present
		const links = [
			{ text: "Home", href: "/" },
			{ text: "Users", href: "/users" },
			{ text: "Settings", href: "/users/settings" },
			{ text: "Profile", href: "/users/settings/profile" },
		];

		for (const link of links) {
			const element = canvas.getByRole("link", { name: link.text });
			await expect(element).toBeInTheDocument();
			await expect(element).toHaveAttribute("href", link.href);
		}

		// Check current page
		const currentPage = canvas.getByText("Edit");
		await expect(currentPage).toHaveAttribute("aria-current", "page");
	},
};

export const CustomSeparator: Story = {
	render: () => (
		<Breadcrumb>
			<BreadcrumbList>
				<BreadcrumbItem>
					<BreadcrumbLink href="/">Home</BreadcrumbLink>
				</BreadcrumbItem>
				<BreadcrumbSeparator className="text-primary">→</BreadcrumbSeparator>
				<BreadcrumbItem>
					<BreadcrumbLink href="/docs">Documentation</BreadcrumbLink>
				</BreadcrumbItem>
				<BreadcrumbSeparator className="text-primary">→</BreadcrumbSeparator>
				<BreadcrumbItem>
					<BreadcrumbPage>API</BreadcrumbPage>
				</BreadcrumbItem>
			</BreadcrumbList>
		</Breadcrumb>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const separators = canvas.getAllByText("→");
		await expect(separators).toHaveLength(2);
		await expect(separators[0]).toHaveClass("text-primary");
	},
};
