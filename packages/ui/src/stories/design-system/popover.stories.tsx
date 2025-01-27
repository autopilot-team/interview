import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { Button } from "../../components/button.js";
import { Input } from "../../components/input.js";
import { Label } from "../../components/label.js";
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from "../../components/popover.js";

const meta: Meta<typeof Popover> = {
	title: "Design System/Popover",
	component: Popover,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Popover>;

export const Default: Story = {
	render: () => (
		<Popover>
			<PopoverTrigger asChild>
				<Button variant="outline">Open popover</Button>
			</PopoverTrigger>
			<PopoverContent className="w-80">
				<div className="grid gap-4">
					<div className="space-y-2">
						<h4 className="font-medium leading-none">Dimensions</h4>
						<p className="text-sm text-muted-foreground">
							Set the dimensions for the layer.
						</p>
					</div>
					<div className="grid gap-2">
						<div className="grid grid-cols-3 items-center gap-4">
							<Label htmlFor="width">Width</Label>
							<Input
								id="width"
								defaultValue="100%"
								className="col-span-2 h-8"
							/>
						</div>
						<div className="grid grid-cols-3 items-center gap-4">
							<Label htmlFor="height">Height</Label>
							<Input
								id="height"
								defaultValue="25px"
								className="col-span-2 h-8"
							/>
						</div>
					</div>
				</div>
			</PopoverContent>
		</Popover>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const button = canvas.getByRole("button");
		await expect(button).toBeInTheDocument();
		await userEvent.click(button);

		// Wait for popover to be mounted and look for it in the document body
		const popover = await within(document.body).findByRole(
			"dialog",
			{},
			{ timeout: 1000 },
		);
		await new Promise((resolve) => setTimeout(resolve, 100)); // Add a small delay
		await expect(popover).toBeVisible();

		// Test popover content
		const popoverContent = within(popover);
		const widthInput = popoverContent.getByLabelText("Width");
		await expect(widthInput).toBeVisible();
	},
};

export const WithCustomStyling: Story = {
	render: () => (
		<Popover>
			<PopoverTrigger asChild>
				<Button>Custom Styled Popover</Button>
			</PopoverTrigger>
			<PopoverContent className="w-80 bg-secondary">
				<div className="grid gap-4">
					<div className="space-y-2">
						<h4 className="font-medium leading-none text-secondary-foreground">
							Custom Theme
						</h4>
						<p className="text-sm text-secondary-foreground/70">
							This popover uses custom background and text colors.
						</p>
					</div>
					<div className="grid gap-2">
						<Button variant="outline" className="w-full">
							Action Button
						</Button>
					</div>
				</div>
			</PopoverContent>
		</Popover>
	),
};

export const WithFooter: Story = {
	render: () => (
		<Popover>
			<PopoverTrigger asChild>
				<Button variant="outline">Show Details</Button>
			</PopoverTrigger>
			<PopoverContent className="w-80">
				<div className="grid gap-4">
					<div className="space-y-2">
						<h4 className="font-medium leading-none">Details</h4>
						<p className="text-sm text-muted-foreground">
							View additional information about this item.
						</p>
					</div>
					<div className="grid gap-2">
						<div className="text-sm">
							Lorem ipsum dolor sit amet, consectetur adipiscing elit.
						</div>
					</div>
				</div>
				<div className="mt-4 flex justify-end gap-2 border-t pt-4">
					<Button variant="outline" size="sm">
						Cancel
					</Button>
					<Button size="sm">Save</Button>
				</div>
			</PopoverContent>
		</Popover>
	),
};

export const WithForm: Story = {
	render: () => (
		<Popover>
			<PopoverTrigger asChild>
				<Button variant="outline">Edit Profile</Button>
			</PopoverTrigger>
			<PopoverContent className="w-80">
				<form className="grid gap-4">
					<div className="space-y-2">
						<h4 className="font-medium leading-none">Profile</h4>
						<p className="text-sm text-muted-foreground">
							Update your profile information.
						</p>
					</div>
					<div className="grid gap-2">
						<div className="grid gap-1">
							<Label htmlFor="name">Name</Label>
							<Input id="name" defaultValue="John Doe" className="h-8" />
						</div>
						<div className="grid gap-1">
							<Label htmlFor="email">Email</Label>
							<Input
								id="email"
								defaultValue="john@example.com"
								className="h-8"
							/>
						</div>
					</div>
					<div className="flex justify-end gap-2">
						<Button variant="outline" size="sm">
							Cancel
						</Button>
						<Button type="submit" size="sm">
							Save
						</Button>
					</div>
				</form>
			</PopoverContent>
		</Popover>
	),
};
