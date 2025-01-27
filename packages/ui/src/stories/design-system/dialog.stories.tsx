import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { Button } from "../../components/button.js";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "../../components/dialog.js";

const meta = {
	title: "Design System/Dialog",
	component: Dialog,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
} satisfies Meta<typeof Dialog>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<Dialog>
			<DialogTrigger asChild>
				<Button variant="outline">Open Dialog</Button>
			</DialogTrigger>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Edit Profile</DialogTitle>
					<DialogDescription>
						Make changes to your profile here. Click save when you're done.
					</DialogDescription>
				</DialogHeader>
				<div className="grid gap-4 py-4">
					<div className="grid grid-cols-4 items-center gap-4">
						<label htmlFor="name" className="text-right text-sm font-medium">
							Name
						</label>
						<input
							id="name"
							className="col-span-3 h-9 rounded-md border border-input px-3"
							placeholder="Enter your name"
						/>
					</div>
					<div className="grid grid-cols-4 items-center gap-4">
						<label
							htmlFor="username"
							className="text-right text-sm font-medium"
						>
							Username
						</label>
						<input
							id="username"
							className="col-span-3 h-9 rounded-md border border-input px-3"
							placeholder="Enter username"
						/>
					</div>
				</div>
				<DialogFooter>
					<Button type="submit">Save changes</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const button = canvas.getByRole("button", { name: "Open Dialog" });
		await expect(button).toBeInTheDocument();
		await userEvent.click(button);

		// Wait for dialog to be mounted and look for it in the document body
		const dialog = await within(document.body).findByRole(
			"dialog",
			{},
			{ timeout: 1000 },
		);
		await new Promise((resolve) => setTimeout(resolve, 100)); // Add a small delay for animation
		await expect(dialog).toBeVisible();

		// Test dialog content
		const dialogContent = within(dialog);
		const nameInput = dialogContent.getByLabelText("Name");
		await expect(nameInput).toBeVisible();
	},
};

export const WithDestructiveAction: Story = {
	render: () => (
		<Dialog>
			<DialogTrigger asChild>
				<Button variant="outline">Delete Account</Button>
			</DialogTrigger>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Are you absolutely sure?</DialogTitle>
					<DialogDescription>
						This action cannot be undone. This will permanently delete your
						account and remove your data from our servers.
					</DialogDescription>
				</DialogHeader>
				<DialogFooter>
					<Button variant="outline">Cancel</Button>
					<Button variant="destructive">Delete Account</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	),
};

export const WithForm: Story = {
	render: () => (
		<Dialog>
			<DialogTrigger asChild>
				<Button variant="outline">Create Project</Button>
			</DialogTrigger>
			<DialogContent className="sm:max-w-[425px]">
				<DialogHeader>
					<DialogTitle>Create project</DialogTitle>
					<DialogDescription>
						Add a new project to your workspace.
					</DialogDescription>
				</DialogHeader>
				<form className="grid gap-4 py-4">
					<div className="grid grid-cols-4 items-center gap-4">
						<label
							htmlFor="project-name"
							className="text-right text-sm font-medium"
						>
							Name
						</label>
						<input
							id="project-name"
							className="col-span-3 h-9 rounded-md border border-input px-3"
							placeholder="Project name"
						/>
					</div>
					<div className="grid grid-cols-4 items-center gap-4">
						<label
							htmlFor="description"
							className="text-right text-sm font-medium"
						>
							Description
						</label>
						<textarea
							id="description"
							className="col-span-3 h-24 rounded-md border border-input p-3"
							placeholder="Project description"
						/>
					</div>
					<DialogFooter>
						<Button type="submit">Create Project</Button>
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	),
};
