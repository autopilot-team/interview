import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { ChevronsUpDown, Plus } from "lucide-react";

import { Button } from "../../components/button.js";
import {
	Collapsible,
	CollapsibleContent,
	CollapsibleTrigger,
} from "../../components/collapsible.js";

const meta: Meta<typeof Collapsible> = {
	title: "Design System/Collapsible",
	component: Collapsible,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Collapsible>;

export const Default: Story = {
	render: () => (
		<Collapsible className="w-[350px] space-y-2">
			<div className="flex items-center justify-between space-x-4 px-4">
				<h4 className="text-sm font-semibold">Notifications</h4>
				<CollapsibleTrigger asChild>
					<Button variant="ghost" size="sm" className="w-9 p-0">
						<ChevronsUpDown className="h-4 w-4" />
						<span className="sr-only">Toggle</span>
					</Button>
				</CollapsibleTrigger>
			</div>
			<CollapsibleContent className="space-y-2">
				<div className="rounded-md border px-4 py-3 font-mono text-sm">
					You have a new message from @johndoe
				</div>
				<div className="rounded-md border px-4 py-3 font-mono text-sm">
					Your subscription will expire soon
				</div>
			</CollapsibleContent>
		</Collapsible>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByRole("button", { name: "Toggle" });
		await expect(trigger).toBeInTheDocument();

		// Test collapsible interaction
		await userEvent.click(trigger);
		const content = canvas.getByText("You have a new message from @johndoe");
		await expect(content).toBeVisible();
	},
};

export const WithCustomTrigger: Story = {
	render: () => (
		<Collapsible className="w-[350px] space-y-2">
			<div className="flex items-center justify-between space-x-4 px-4">
				<h4 className="text-sm font-semibold">Advanced Settings</h4>
				<CollapsibleTrigger asChild>
					<Button variant="outline" size="sm">
						<Plus className="mr-2 h-4 w-4" />
						Show Options
					</Button>
				</CollapsibleTrigger>
			</div>
			<CollapsibleContent className="space-y-2">
				<div className="rounded-md border px-4 py-2 shadow-sm">
					<div className="flex items-center justify-between">
						<span className="text-sm font-medium">Developer Mode</span>
						<Button variant="ghost" size="sm">
							Enable
						</Button>
					</div>
				</div>
				<div className="rounded-md border px-4 py-2 shadow-sm">
					<div className="flex items-center justify-between">
						<span className="text-sm font-medium">Experimental Features</span>
						<Button variant="ghost" size="sm">
							Configure
						</Button>
					</div>
				</div>
			</CollapsibleContent>
		</Collapsible>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByRole("button", { name: "Show Options" });

		await userEvent.click(trigger);
		await expect(canvas.getByText("Developer Mode")).toBeVisible();
		await expect(canvas.getByText("Experimental Features")).toBeVisible();
	},
};

export const WithAnimation: Story = {
	render: () => (
		<Collapsible className="w-[350px]">
			<div className="flex items-center justify-between space-x-4 rounded-lg border p-4">
				<div className="space-y-1">
					<h4 className="text-sm font-semibold">Filter Results</h4>
					<p className="text-sm text-muted-foreground">
						Customize your search filters
					</p>
				</div>
				<CollapsibleTrigger asChild>
					<Button variant="ghost" size="sm" className="w-9 p-0">
						<ChevronsUpDown className="h-4 w-4 transition-transform duration-200" />
						<span className="sr-only">Toggle</span>
					</Button>
				</CollapsibleTrigger>
			</div>
			<CollapsibleContent className="space-y-2 overflow-hidden transition-all data-[state=closed]:animate-collapse data-[state=open]:animate-expand">
				<div className="rounded-md border p-4 mt-2">
					<div className="space-y-2">
						<div className="flex items-center space-x-2">
							<input
								type="checkbox"
								id="price"
								className="rounded border-gray-300"
							/>
							<label htmlFor="price" className="text-sm">
								Price Range
							</label>
						</div>
						<div className="flex items-center space-x-2">
							<input
								type="checkbox"
								id="date"
								className="rounded border-gray-300"
							/>
							<label htmlFor="date" className="text-sm">
								Date Added
							</label>
						</div>
						<div className="flex items-center space-x-2">
							<input
								type="checkbox"
								id="rating"
								className="rounded border-gray-300"
							/>
							<label htmlFor="rating" className="text-sm">
								Rating
							</label>
						</div>
					</div>
				</div>
			</CollapsibleContent>
		</Collapsible>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByRole("button", { name: "Toggle" });

		await userEvent.click(trigger);
		await expect(canvas.getByText("Price Range")).toBeVisible();
		await expect(canvas.getByText("Date Added")).toBeVisible();
		await expect(canvas.getByText("Rating")).toBeVisible();
	},
};
