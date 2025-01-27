import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { Button } from "../../components/button.js";
import { Input } from "../../components/input.js";
import { Label } from "../../components/label.js";
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetFooter,
	SheetHeader,
	SheetTitle,
	SheetTrigger,
} from "../../components/sheet.js";

const meta: Meta<typeof Sheet> = {
	title: "Design System/Sheet",
	component: Sheet,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Sheet>;

export const Default: Story = {
	render: () => (
		<Sheet>
			<SheetTrigger asChild>
				<Button variant="outline">Open Sheet</Button>
			</SheetTrigger>
			<SheetContent>
				<SheetHeader>
					<SheetTitle>Edit profile</SheetTitle>
					<SheetDescription>
						Make changes to your profile here. Click save when you're done.
					</SheetDescription>
				</SheetHeader>
				<div className="grid gap-4 py-4">
					<div className="grid grid-cols-4 items-center gap-4">
						<Label htmlFor="name" className="text-right">
							Name
						</Label>
						<Input id="name" value="Pedro Duarte" className="col-span-3" />
					</div>
					<div className="grid grid-cols-4 items-center gap-4">
						<Label htmlFor="username" className="text-right">
							Username
						</Label>
						<Input id="username" value="@peduarte" className="col-span-3" />
					</div>
				</div>
				<SheetFooter>
					<Button type="submit">Save changes</Button>
				</SheetFooter>
			</SheetContent>
		</Sheet>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const button = canvas.getByRole("button");
		await expect(button).toBeInTheDocument();
		await userEvent.click(button);
		// Wait for sheet to appear in the document body since it's rendered in a portal
		const sheet = await within(document.body).findByRole("dialog");
		await expect(sheet).toBeVisible();
		const nameInput = within(sheet).getByLabelText("Name");
		await expect(nameInput).toBeVisible();
	},
};

export const SideLeft: Story = {
	render: () => (
		<Sheet>
			<SheetTrigger asChild>
				<Button variant="outline">Open Left Sheet</Button>
			</SheetTrigger>
			<SheetContent side="left">
				<SheetHeader>
					<SheetTitle>Navigation</SheetTitle>
					<SheetDescription>
						Browse through different sections.
					</SheetDescription>
				</SheetHeader>
				<div className="grid gap-4 py-4">
					<Button variant="ghost" className="justify-start">
						Dashboard
					</Button>
					<Button variant="ghost" className="justify-start">
						Settings
					</Button>
					<Button variant="ghost" className="justify-start">
						Profile
					</Button>
				</div>
			</SheetContent>
		</Sheet>
	),
};

export const SideTop: Story = {
	render: () => (
		<Sheet>
			<SheetTrigger asChild>
				<Button variant="outline">Open Top Sheet</Button>
			</SheetTrigger>
			<SheetContent side="top">
				<SheetHeader>
					<SheetTitle>Quick Actions</SheetTitle>
					<SheetDescription>Access frequently used features.</SheetDescription>
				</SheetHeader>
				<div className="flex justify-center gap-4 py-4">
					<Button>New File</Button>
					<Button>Upload</Button>
					<Button>Share</Button>
				</div>
			</SheetContent>
		</Sheet>
	),
};

export const WithForm: Story = {
	render: () => (
		<Sheet>
			<SheetTrigger asChild>
				<Button variant="outline">Edit Settings</Button>
			</SheetTrigger>
			<SheetContent className="sm:max-w-[540px]">
				<SheetHeader>
					<SheetTitle>Account Settings</SheetTitle>
					<SheetDescription>
						Configure your account preferences and settings.
					</SheetDescription>
				</SheetHeader>
				<form className="grid gap-4 py-4">
					<div className="grid gap-2">
						<Label htmlFor="email">Email</Label>
						<Input id="email" type="email" placeholder="Enter your email" />
					</div>
					<div className="grid gap-2">
						<Label htmlFor="password">Password</Label>
						<Input id="password" type="password" />
					</div>
					<div className="grid gap-2">
						<Label htmlFor="bio">Bio</Label>
						<Input id="bio" placeholder="Tell us about yourself" />
					</div>
					<SheetFooter>
						<Button type="submit">Save Changes</Button>
					</SheetFooter>
				</form>
			</SheetContent>
		</Sheet>
	),
};

export const WithCustomStyling: Story = {
	render: () => (
		<Sheet>
			<SheetTrigger asChild>
				<Button>Custom Sheet</Button>
			</SheetTrigger>
			<SheetContent className="bg-secondary">
				<SheetHeader>
					<SheetTitle className="text-secondary-foreground">
						Custom Theme
					</SheetTitle>
					<SheetDescription className="text-secondary-foreground/70">
						This sheet uses custom background and text colors.
					</SheetDescription>
				</SheetHeader>
				<div className="py-4">
					<Button variant="outline" className="w-full">
						Action Button
					</Button>
				</div>
			</SheetContent>
		</Sheet>
	),
};
