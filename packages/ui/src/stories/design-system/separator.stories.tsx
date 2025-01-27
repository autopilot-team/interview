import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import { Separator } from "../../components/separator.js";

const meta: Meta<typeof Separator> = {
	title: "Design System/Separator",
	component: Separator,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Separator>;

export const Default: Story = {
	render: () => (
		<div className="w-[300px] space-y-1">
			<div className="text-sm font-medium">Radix UI</div>
			<div className="text-sm text-muted-foreground">
				An open-source UI component library.
			</div>
			<Separator className="my-4" />
			<div className="flex h-5 items-center space-x-4 text-sm">
				<div>Blog</div>
				<Separator orientation="vertical" />
				<div>Docs</div>
				<Separator orientation="vertical" />
				<div>Source</div>
			</div>
		</div>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Test horizontal separator
		const horizontalSeparators = canvas.getAllByRole("none", { hidden: true });
		const horizontalSeparator = horizontalSeparators.find(
			(sep) => sep.getAttribute("data-orientation") === "horizontal",
		);
		expect(horizontalSeparator).toBeDefined();
		await expect(horizontalSeparator).toBeVisible();
		await expect(horizontalSeparator).toHaveAttribute(
			"data-orientation",
			"horizontal",
		);

		// Test vertical separators
		const verticalSeparators = horizontalSeparators.filter(
			(sep) => sep.getAttribute("data-orientation") === "vertical",
		);
		expect(verticalSeparators).toHaveLength(2);
		for (const separator of verticalSeparators) {
			await expect(separator).toBeVisible();
			await expect(separator).toHaveAttribute("data-orientation", "vertical");
		}
	},
};

export const Horizontal: Story = {
	render: () => (
		<div className="w-[300px] space-y-4">
			<div className="space-y-1">
				<h4 className="text-sm font-medium">Section 1</h4>
				<p className="text-sm text-muted-foreground">Content for section 1</p>
			</div>
			<Separator />
			<div className="space-y-1">
				<h4 className="text-sm font-medium">Section 2</h4>
				<p className="text-sm text-muted-foreground">Content for section 2</p>
			</div>
		</div>
	),
};

export const Vertical: Story = {
	render: () => (
		<div className="flex h-[150px] items-center space-x-4">
			<div className="space-y-1">
				<h4 className="text-sm font-medium">Left Content</h4>
				<p className="text-sm text-muted-foreground">Description</p>
			</div>
			<Separator orientation="vertical" className="h-full" />
			<div className="space-y-1">
				<h4 className="text-sm font-medium">Right Content</h4>
				<p className="text-sm text-muted-foreground">Description</p>
			</div>
		</div>
	),
};

export const WithCustomStyling: Story = {
	render: () => (
		<div className="w-[300px] space-y-4">
			<Separator className="bg-primary/50" />
			<Separator className="bg-secondary" />
			<Separator className="bg-destructive/50" />
			<div className="flex h-5 items-center space-x-4 text-sm">
				<div>Item 1</div>
				<Separator orientation="vertical" className="bg-primary/50" />
				<div>Item 2</div>
				<Separator orientation="vertical" className="bg-secondary" />
				<div>Item 3</div>
			</div>
		</div>
	),
};

export const WithContent: Story = {
	render: () => (
		<div className="w-[300px] space-y-4">
			<div className="space-y-1">
				<h4 className="text-sm font-medium leading-none">Navigation</h4>
			</div>
			<Separator />
			<div className="flex flex-col space-y-2 text-sm">
				<div className="flex items-center">Home</div>
				<div className="flex items-center">About</div>
				<div className="flex items-center">Settings</div>
			</div>
			<Separator />
			<div className="flex flex-col space-y-2 text-sm text-muted-foreground">
				<div className="flex items-center">Privacy Policy</div>
				<div className="flex items-center">Terms of Service</div>
			</div>
		</div>
	),
};
