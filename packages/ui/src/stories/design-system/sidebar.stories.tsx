import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import {
	BarChart3,
	Boxes,
	CircleUser,
	FileText,
	Home,
	LayoutDashboard,
	LifeBuoy,
	LogOut,
	Search,
	Settings,
} from "lucide-react";

import { Button } from "../../components/button.js";
import { Input } from "../../components/input.js";
import { Sidebar, SidebarProvider } from "../../components/sidebar.js";

const meta: Meta<typeof Sidebar> = {
	title: "Design System/Sidebar",
	component: Sidebar,
	parameters: {
		layout: "fullscreen",
	},
	tags: ["autodocs"],
	decorators: [
		(Story) => (
			<SidebarProvider>
				<Story />
			</SidebarProvider>
		),
	],
};

export default meta;
type Story = StoryObj<typeof Sidebar>;

const SidebarContent = () => (
	<>
		<div className="flex h-14 items-center border-b px-4">
			<Button variant="ghost" className="h-auto p-0">
				<LayoutDashboard className="h-6 w-6" />
				<span className="ml-3 text-lg font-semibold">Dashboard</span>
			</Button>
		</div>
		<div className="flex-1 overflow-auto">
			<div className="space-y-4 py-4">
				<div className="px-3">
					<div className="relative">
						<Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
						<Input placeholder="Search" className="pl-8" />
					</div>
				</div>
				<div className="px-3">
					<h2 className="mb-2 px-4 text-lg font-semibold tracking-tight">
						Overview
					</h2>
					<div className="space-y-1">
						<Button variant="ghost" className="w-full justify-start">
							<Home className="mr-2 h-4 w-4" />
							Home
						</Button>
						<Button variant="ghost" className="w-full justify-start">
							<BarChart3 className="mr-2 h-4 w-4" />
							Analytics
						</Button>
						<Button variant="ghost" className="w-full justify-start">
							<FileText className="mr-2 h-4 w-4" />
							Reports
						</Button>
					</div>
				</div>
				<div className="px-3">
					<h2 className="mb-2 px-4 text-lg font-semibold tracking-tight">
						Resources
					</h2>
					<div className="space-y-1">
						<Button variant="ghost" className="w-full justify-start">
							<Boxes className="mr-2 h-4 w-4" />
							Products
						</Button>
						<Button variant="ghost" className="w-full justify-start">
							<CircleUser className="mr-2 h-4 w-4" />
							Customers
						</Button>
					</div>
				</div>
			</div>
		</div>
		<div className="mt-auto border-t p-4">
			<div className="flex items-center justify-between">
				<Button variant="ghost" size="icon">
					<Settings className="h-4 w-4" />
				</Button>
				<Button variant="ghost" size="icon">
					<LifeBuoy className="h-4 w-4" />
				</Button>
				<Button variant="ghost" size="icon">
					<LogOut className="h-4 w-4" />
				</Button>
			</div>
		</div>
	</>
);

export const Default: Story = {
	render: () => (
		<div className="flex h-screen">
			<Sidebar>
				<SidebarContent />
			</Sidebar>
			<div className="flex-1 p-8">
				<h1 className="text-2xl font-bold">Main Content</h1>
			</div>
		</div>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		await expect(canvas.getByText("Dashboard")).toBeVisible();
		await expect(canvas.getByPlaceholderText("Search")).toBeVisible();
	},
};

export const Floating: Story = {
	render: () => (
		<div className="flex h-screen">
			<Sidebar variant="floating">
				<SidebarContent />
			</Sidebar>
			<div className="flex-1 p-8">
				<h1 className="text-2xl font-bold">Main Content</h1>
			</div>
		</div>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		await expect(canvas.getByText("Dashboard")).toBeVisible();
	},
};

export const Inset: Story = {
	render: () => (
		<div className="flex h-screen">
			<Sidebar variant="inset">
				<SidebarContent />
			</Sidebar>
			<div className="flex-1 p-8">
				<h1 className="text-2xl font-bold">Main Content</h1>
			</div>
		</div>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		await expect(canvas.getByText("Dashboard")).toBeVisible();
	},
};

export const RightSide: Story = {
	render: () => (
		<div className="flex h-screen">
			<div className="flex-1 p-8">
				<h1 className="text-2xl font-bold">Main Content</h1>
			</div>
			<Sidebar side="right">
				<SidebarContent />
			</Sidebar>
		</div>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		await expect(canvas.getByText("Dashboard")).toBeVisible();
	},
};
