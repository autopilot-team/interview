import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { Button } from "../../components/button.js";
import {
	Drawer,
	DrawerClose,
	DrawerContent,
	DrawerDescription,
	DrawerFooter,
	DrawerHeader,
	DrawerTitle,
	DrawerTrigger,
} from "../../components/drawer.js";
import { Input } from "../../components/input.js";
import { Label } from "../../components/label.js";

const meta: Meta<typeof Drawer> = {
	title: "Design System/Drawer",
	component: Drawer,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Drawer>;

export const Default: Story = {
	render: () => (
		<Drawer>
			<DrawerTrigger asChild>
				<Button variant="outline">Open Drawer</Button>
			</DrawerTrigger>
			<DrawerContent>
				<DrawerHeader>
					<DrawerTitle>Edit Profile</DrawerTitle>
					<DrawerDescription>
						Make changes to your profile here. Click save when you're done.
					</DrawerDescription>
				</DrawerHeader>
				<div className="p-4 pb-0">
					<div className="grid gap-4">
						<div className="grid gap-2">
							<Label htmlFor="name">Name</Label>
							<Input id="name" placeholder="Enter your name" />
						</div>
						<div className="grid gap-2">
							<Label htmlFor="username">Username</Label>
							<Input id="username" placeholder="Enter your username" />
						</div>
					</div>
				</div>
				<DrawerFooter>
					<Button>Save changes</Button>
					<DrawerClose asChild>
						<Button variant="outline">Cancel</Button>
					</DrawerClose>
				</DrawerFooter>
			</DrawerContent>
		</Drawer>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const openButton = canvas.getByRole("button", { name: "Open Drawer" });
		await expect(openButton).toBeInTheDocument();
		await userEvent.click(openButton);

		// Wait for drawer to be mounted and look for it in the document body
		const drawer = await within(document.body).findByRole("dialog");
		await expect(drawer).toBeVisible();

		// Test drawer content
		const drawerContent = within(drawer);
		const nameInput = drawerContent.getByLabelText("Name");
		const usernameInput = drawerContent.getByLabelText("Username");
		await expect(nameInput).toBeVisible();
		await expect(usernameInput).toBeVisible();
	},
};

export const WithDestructiveAction: Story = {
	render: () => (
		<Drawer>
			<DrawerTrigger asChild>
				<Button variant="outline">Delete Account</Button>
			</DrawerTrigger>
			<DrawerContent>
				<DrawerHeader>
					<DrawerTitle>Are you absolutely sure?</DrawerTitle>
					<DrawerDescription>
						This action cannot be undone. This will permanently delete your
						account and remove your data from our servers.
					</DrawerDescription>
				</DrawerHeader>
				<DrawerFooter>
					<DrawerClose asChild>
						<Button variant="outline">Cancel</Button>
					</DrawerClose>
					<Button variant="destructive">Delete Account</Button>
				</DrawerFooter>
			</DrawerContent>
		</Drawer>
	),
};

export const WithCustomContent: Story = {
	render: () => (
		<Drawer>
			<DrawerTrigger asChild>
				<Button variant="outline">View Details</Button>
			</DrawerTrigger>
			<DrawerContent className="h-[75vh]">
				<DrawerHeader>
					<DrawerTitle>Product Details</DrawerTitle>
					<DrawerDescription>
						View detailed information about this product.
					</DrawerDescription>
				</DrawerHeader>
				<div className="p-4 space-y-4">
					<div className="aspect-video rounded-lg bg-muted" />
					<h3 className="text-lg font-semibold">Product Features</h3>
					<ul className="list-disc list-inside space-y-2 text-muted-foreground">
						<li>High-quality materials</li>
						<li>Durable construction</li>
						<li>Modern design</li>
						<li>Easy to maintain</li>
					</ul>
					<p className="text-muted-foreground">
						Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do
						eiusmod tempor incididunt ut labore et dolore magna aliqua.
					</p>
				</div>
				<DrawerFooter>
					<Button>Add to Cart</Button>
					<DrawerClose asChild>
						<Button variant="outline">Close</Button>
					</DrawerClose>
				</DrawerFooter>
			</DrawerContent>
		</Drawer>
	),
};
