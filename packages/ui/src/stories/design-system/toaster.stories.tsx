import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";

import { Button } from "../../components/button.js";
import { Toaster } from "../../components/toaster.js";
import { useToast } from "../../hooks/use-toast.js";

const meta: Meta<typeof Toaster> = {
	title: "Design System/Toaster",
	component: Toaster,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Toaster>;

const ToastDemo = () => {
	const { toast } = useToast();

	return (
		<div className="space-y-2">
			<Button
				variant="outline"
				onClick={() => {
					toast({
						title: "Default Toast",
						description: "This is a default toast notification",
					});
				}}
			>
				Show Toast
			</Button>
			<Toaster />
		</div>
	);
};

const ToastWithActionsDemo = () => {
	const { toast } = useToast();

	return (
		<div className="space-y-2">
			<Button
				variant="outline"
				onClick={() => {
					toast({
						title: "Scheduled: Catch up",
						description: "Friday, February 10, 2024 at 5:57 PM",
						action: (
							<Button variant="outline" size="sm">
								Undo
							</Button>
						),
					});
				}}
			>
				Show Toast with Action
			</Button>
			<Toaster />
		</div>
	);
};

const ToastVariantsDemo = () => {
	const { toast } = useToast();

	return (
		<div className="space-y-2">
			<div className="flex gap-2">
				<Button
					variant="default"
					onClick={() => {
						toast({
							title: "Success",
							description: "Your changes have been saved",
							variant: "default",
						});
					}}
				>
					Success
				</Button>
				<Button
					variant="destructive"
					onClick={() => {
						toast({
							title: "Error",
							description: "Something went wrong",
							variant: "destructive",
						});
					}}
				>
					Error
				</Button>
			</div>
			<Toaster />
		</div>
	);
};

export const Default: Story = {
	render: () => <ToastDemo />,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const button = canvas.getByRole("button", { name: "Show Toast" });
		await expect(button).toBeInTheDocument();

		await userEvent.click(button);
		await expect(canvas.getByText("Default Toast")).toBeVisible();
		await expect(
			canvas.getByText("This is a default toast notification"),
		).toBeVisible();
	},
};

export const WithActions: Story = {
	render: () => <ToastWithActionsDemo />,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const button = canvas.getByRole("button", {
			name: "Show Toast with Action",
		});

		await userEvent.click(button);
		await expect(canvas.getByText("Scheduled: Catch up")).toBeVisible();
		await expect(canvas.getByRole("button", { name: "Undo" })).toBeVisible();
	},
};

export const Variants: Story = {
	render: () => <ToastVariantsDemo />,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const successButton = canvas.getByRole("button", { name: "Success" });
		const errorButton = canvas.getByRole("button", { name: "Error" });

		// Test success toast
		await userEvent.click(successButton);
		const successToast = await within(document.body).findByRole("status");
		await expect(within(successToast).getByText("Success")).toBeVisible();
		await expect(
			within(successToast).getByText("Your changes have been saved"),
		).toBeVisible();

		// Test error toast
		await userEvent.click(errorButton);
		const errorToast = await within(document.body).findByRole("status");
		await expect(within(errorToast).getByText("Error")).toBeVisible();
		await expect(
			within(errorToast).getByText("Something went wrong"),
		).toBeVisible();
	},
};
