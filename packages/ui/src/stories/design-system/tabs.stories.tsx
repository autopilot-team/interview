import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { Button } from "../../components/button.js";
import { Input } from "../../components/input.js";
import { Label } from "../../components/label.js";
import {
	Tabs,
	TabsContent,
	TabsList,
	TabsTrigger,
} from "../../components/tabs.js";

const meta: Meta<typeof Tabs> = {
	title: "Design System/Tabs",
	component: Tabs,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Tabs>;

export const Default: Story = {
	render: () => (
		<Tabs defaultValue="account" className="w-[400px]">
			<TabsList>
				<TabsTrigger value="account">Account</TabsTrigger>
				<TabsTrigger value="password">Password</TabsTrigger>
			</TabsList>
			<TabsContent value="account">
				<div className="space-y-4 p-4">
					<div className="space-y-2">
						<Label htmlFor="email">Email</Label>
						<Input id="email" type="email" placeholder="m@example.com" />
					</div>
					<div className="space-y-2">
						<Label htmlFor="username">Username</Label>
						<Input id="username" placeholder="@username" />
					</div>
					<Button>Save Changes</Button>
				</div>
			</TabsContent>
			<TabsContent value="password">
				<div className="space-y-4 p-4">
					<div className="space-y-2">
						<Label htmlFor="current">Current Password</Label>
						<Input id="current" type="password" />
					</div>
					<div className="space-y-2">
						<Label htmlFor="new">New Password</Label>
						<Input id="new" type="password" />
					</div>
					<Button>Change Password</Button>
				</div>
			</TabsContent>
		</Tabs>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const accountTab = canvas.getByRole("tab", { name: "Account" });
		await expect(accountTab).toBeInTheDocument();
		await expect(accountTab).toHaveAttribute("data-state", "active");
		await userEvent.click(canvas.getByRole("tab", { name: "Password" }));
		await expect(canvas.getByLabelText("Current Password")).toBeVisible();
	},
};

export const WithIcons: Story = {
	render: () => (
		<Tabs defaultValue="music" className="w-[400px]">
			<TabsList className="grid w-full grid-cols-3">
				<TabsTrigger value="music">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						strokeWidth="2"
						strokeLinecap="round"
						strokeLinejoin="round"
						className="mr-2 h-4 w-4"
					>
						<title>Music</title>
						<path d="M9 18V5l12-2v13" />
						<circle cx="6" cy="18" r="3" />
						<circle cx="18" cy="16" r="3" />
					</svg>
					Music
				</TabsTrigger>
				<TabsTrigger value="podcasts">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						strokeWidth="2"
						strokeLinecap="round"
						strokeLinejoin="round"
						className="mr-2 h-4 w-4"
					>
						<title>Podcasts</title>
						<circle cx="12" cy="12" r="4" />
						<path d="M16 8v5a3 3 0 0 0 6 0v-1a10 10 0 1 0-4 8" />
					</svg>
					Podcasts
				</TabsTrigger>
				<TabsTrigger value="live">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						strokeWidth="2"
						strokeLinecap="round"
						strokeLinejoin="round"
						className="mr-2 h-4 w-4"
					>
						<title>Live</title>
						<path d="M12 20v-6M6 20V10M18 20V4" />
					</svg>
					Live
				</TabsTrigger>
			</TabsList>
			<TabsContent value="music" className="p-4">
				Music content
			</TabsContent>
			<TabsContent value="podcasts" className="p-4">
				Podcasts content
			</TabsContent>
			<TabsContent value="live" className="p-4">
				Live content
			</TabsContent>
		</Tabs>
	),
};

export const WithCards: Story = {
	render: () => (
		<Tabs defaultValue="overview" className="w-[600px]">
			<TabsList>
				<TabsTrigger value="overview">Overview</TabsTrigger>
				<TabsTrigger value="analytics">Analytics</TabsTrigger>
				<TabsTrigger value="reports">Reports</TabsTrigger>
				<TabsTrigger value="notifications">Notifications</TabsTrigger>
			</TabsList>
			<TabsContent value="overview">
				<div className="grid gap-4 p-4 md:grid-cols-2">
					<div className="rounded-lg border p-4">
						<h3 className="font-semibold">Users</h3>
						<p className="text-2xl font-bold">1,234</p>
					</div>
					<div className="rounded-lg border p-4">
						<h3 className="font-semibold">Active Now</h3>
						<p className="text-2xl font-bold">567</p>
					</div>
				</div>
			</TabsContent>
			<TabsContent value="analytics" className="p-4">
				Analytics content
			</TabsContent>
			<TabsContent value="reports" className="p-4">
				Reports content
			</TabsContent>
			<TabsContent value="notifications" className="p-4">
				Notifications content
			</TabsContent>
		</Tabs>
	),
};

export const WithCustomStyling: Story = {
	render: () => (
		<Tabs defaultValue="tab1" className="w-[400px]">
			<TabsList className="bg-secondary">
				<TabsTrigger
					value="tab1"
					className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
				>
					Tab 1
				</TabsTrigger>
				<TabsTrigger
					value="tab2"
					className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
				>
					Tab 2
				</TabsTrigger>
			</TabsList>
			<TabsContent value="tab1" className="p-4 rounded-lg bg-secondary/10">
				Tab 1 content
			</TabsContent>
			<TabsContent value="tab2" className="p-4 rounded-lg bg-secondary/10">
				Tab 2 content
			</TabsContent>
		</Tabs>
	),
};
