import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { Trash2 } from "lucide-react";
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
	AlertDialogTrigger,
} from "../../components/alert-dialog.js";
import { Button } from "../../components/button.js";

const meta: Meta<typeof AlertDialog> = {
	title: "Design System/AlertDialog",
	component: AlertDialog,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof AlertDialog>;

export const Default: Story = {
	render: () => (
		<AlertDialog>
			<AlertDialogTrigger asChild>
				<Button variant="outline">Show Dialog</Button>
			</AlertDialogTrigger>
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
					<AlertDialogDescription>
						This action cannot be undone. This will permanently delete your
						account and remove your data from our servers.
					</AlertDialogDescription>
				</AlertDialogHeader>
				<AlertDialogFooter>
					<AlertDialogCancel>Cancel</AlertDialogCancel>
					<AlertDialogAction>Continue</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const button = canvas.getByRole("button");
		await expect(button).toBeInTheDocument();
		await button.click();

		// Wait for dialog to be mounted in document body
		const dialog = await within(document.body).findByRole(
			"alertdialog",
			{},
			{ timeout: 2000 },
		);
		await expect(dialog).toBeInTheDocument();

		// Wait for animation to complete
		await new Promise((resolve) => setTimeout(resolve, 300));
		await expect(dialog).toBeVisible();

		// Check dialog content
		const title = within(dialog).getByRole("heading");
		await expect(title).toBeVisible();
		const description = within(dialog).getByText(
			/This action cannot be undone/i,
		);
		await expect(description).toBeVisible();
	},
};

export const DestructiveAction: Story = {
	render: () => (
		<AlertDialog>
			<AlertDialogTrigger asChild>
				<Button variant="destructive">
					<Trash2 className="mr-2 h-4 w-4" />
					Delete Account
				</Button>
			</AlertDialogTrigger>
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle className="text-destructive">
						Delete Account
					</AlertDialogTitle>
					<AlertDialogDescription>
						This will permanently delete your account and all associated data.
						This action is irreversible and your data cannot be recovered.
					</AlertDialogDescription>
				</AlertDialogHeader>
				<AlertDialogFooter>
					<AlertDialogCancel>Keep Account</AlertDialogCancel>
					<AlertDialogAction className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
						Delete Account
					</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByRole("button", { name: "Delete Account" });
		await expect(trigger).toHaveClass("bg-destructive");
		await trigger.click();

		// Wait for dialog to be mounted in document body
		const dialog = await within(document.body).findByRole(
			"alertdialog",
			{},
			{ timeout: 2000 },
		);
		await expect(dialog).toBeInTheDocument();

		// Wait for animation to complete
		await new Promise((resolve) => setTimeout(resolve, 300));
		await expect(dialog).toBeVisible();

		// Check dialog content
		const deleteButton = within(dialog).getByRole("button", {
			name: "Delete Account",
		});
		await expect(deleteButton).toHaveClass("bg-destructive");
	},
};

export const WithCustomContent: Story = {
	render: () => (
		<AlertDialog>
			<AlertDialogTrigger asChild>
				<Button>Upgrade Plan</Button>
			</AlertDialogTrigger>
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>Upgrade to Pro Plan</AlertDialogTitle>
					<AlertDialogDescription className="space-y-2">
						<p>You're about to upgrade to our Pro plan. This includes:</p>
						<ul className="list-disc pl-4">
							<li>Unlimited projects</li>
							<li>Priority support</li>
							<li>Custom domain</li>
							<li>Analytics dashboard</li>
						</ul>
						<p className="font-medium">Price: $29/month</p>
					</AlertDialogDescription>
				</AlertDialogHeader>
				<AlertDialogFooter>
					<AlertDialogCancel>Maybe Later</AlertDialogCancel>
					<AlertDialogAction>Upgrade Now</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByRole("button", { name: "Upgrade Plan" });
		await expect(trigger).toBeInTheDocument();
		await trigger.click();

		// Wait for dialog to be mounted in document body
		const dialog = await within(document.body).findByRole(
			"alertdialog",
			{},
			{ timeout: 2000 },
		);
		await expect(dialog).toBeInTheDocument();

		// Wait for animation to complete
		await new Promise((resolve) => setTimeout(resolve, 300));
		await expect(dialog).toBeVisible();

		// Test custom content
		const content = within(dialog);
		await expect(content.getByText("Unlimited projects")).toBeVisible();
		await expect(content.getByText("Price: $29/month")).toBeVisible();
	},
};

export const WithForm: Story = {
	render: () => (
		<AlertDialog>
			<AlertDialogTrigger asChild>
				<Button variant="outline">Leave Team</Button>
			</AlertDialogTrigger>
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>Leave Team</AlertDialogTitle>
					<AlertDialogDescription>
						<p className="mb-4">
							Are you sure you want to leave this team? You will lose access to
							all team projects and resources.
						</p>
						<div className="space-y-2">
							<label
								htmlFor="confirm"
								className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
							>
								Type "confirm" to continue
							</label>
							<input
								id="confirm"
								className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
								placeholder="confirm"
							/>
						</div>
					</AlertDialogDescription>
				</AlertDialogHeader>
				<AlertDialogFooter>
					<AlertDialogCancel>Cancel</AlertDialogCancel>
					<AlertDialogAction>Leave Team</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByRole("button", { name: "Leave Team" });
		await expect(trigger).toBeInTheDocument();
		await trigger.click();

		// Wait for dialog to be mounted in document body
		const dialog = await within(document.body).findByRole(
			"alertdialog",
			{},
			{ timeout: 2000 },
		);
		await expect(dialog).toBeInTheDocument();

		// Wait for animation to complete
		await new Promise((resolve) => setTimeout(resolve, 300));
		await expect(dialog).toBeVisible();

		// Test form input
		const input = within(dialog).getByPlaceholderText("confirm");
		await expect(input).toBeVisible();
		await userEvent.type(input, "confirm");
		await expect(input).toHaveValue("confirm");
	},
};
