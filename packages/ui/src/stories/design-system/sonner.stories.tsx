import type { Meta, StoryObj } from "@storybook/react";
import { expect, waitFor, within } from "@storybook/test";
import { useEffect, useRef } from "react";
import { toast } from "sonner";

import { Button } from "../../components/button.js";
import { Toaster } from "../../components/sonner.js";

const meta: Meta<typeof Toaster> = {
	title: "Design System/Sonner",
	component: Toaster,
	parameters: {
		layout: "centered",
	},
	decorators: [
		(Story) => (
			<div>
				<Story />
				<Toaster />
			</div>
		),
	],
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Toaster>;

const ToastDemo = () => {
	const initialized = useRef(false);

	useEffect(() => {
		if (!initialized.current) {
			initialized.current = true;
			// Clear any existing toasts
			toast.dismiss();
		}
	}, []);

	return (
		<div className="space-y-2">
			<div className="flex gap-2">
				<Button
					variant="outline"
					onClick={() =>
						toast("Event has been created", {
							description: "Sunday, December 03, 2023 at 9:00 AM",
						})
					}
				>
					Show Toast
				</Button>
				<Button
					variant="outline"
					onClick={() =>
						toast.success("Profile updated", {
							description: "Your changes have been saved successfully.",
						})
					}
				>
					Success
				</Button>
				<Button
					variant="outline"
					onClick={() =>
						toast.error("Error occurred", {
							description: "There was a problem with your request.",
						})
					}
				>
					Error
				</Button>
			</div>
		</div>
	);
};

const ToastWithActionsDemo = () => {
	const initialized = useRef(false);

	useEffect(() => {
		if (!initialized.current) {
			initialized.current = true;
			// Clear any existing toasts
			toast.dismiss();
		}
	}, []);

	return (
		<div className="space-y-2">
			<div className="flex gap-2">
				<Button
					variant="outline"
					onClick={() => {
						toast("Scheduled: Catch up", {
							description: "Friday, February 10, 2024 at 5:57 PM",
							action: {
								label: "Undo",
								onClick: () => console.log("Undo"),
							},
						});
					}}
				>
					With Action
				</Button>
				<Button
					variant="outline"
					onClick={() => {
						toast("Multiple actions", {
							description: "Choose what you want to do",
							action: {
								label: "Accept",
								onClick: () => console.log("Accepted"),
							},
							cancel: {
								label: "Cancel",
								onClick: () => console.log("Cancelled"),
							},
						});
					}}
				>
					Multiple Actions
				</Button>
			</div>
		</div>
	);
};

const ToastCustomDemo = () => {
	const initialized = useRef(false);

	useEffect(() => {
		if (!initialized.current) {
			initialized.current = true;
			// Clear any existing toasts
			toast.dismiss();
		}
	}, []);

	return (
		<div className="space-y-2">
			<div className="flex gap-2">
				<Button
					variant="outline"
					onClick={() => {
						toast.custom((t) => (
							<div className="flex items-center gap-4 rounded-lg bg-primary p-4 text-primary-foreground">
								<div className="flex-1">
									<div className="font-semibold">Custom Toast</div>
									<div className="text-sm opacity-90">With custom styling</div>
								</div>
								<Button
									variant="secondary"
									size="sm"
									onClick={() => toast.dismiss(t)}
								>
									Dismiss
								</Button>
							</div>
						));
					}}
				>
					Custom Style
				</Button>
				<Button
					variant="outline"
					onClick={() => {
						toast.promise(
							() => new Promise((resolve) => setTimeout(resolve, 2000)),
							{
								loading: "Loading...",
								success: "Data loaded successfully",
								error: "Error loading data",
							},
						);
					}}
				>
					Promise
				</Button>
			</div>
		</div>
	);
};

export const Default: Story = {
	render: () => <ToastDemo />,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Wait for component to be fully mounted
		await new Promise((resolve) => setTimeout(resolve, 500));

		const button = canvas.getByRole("button", { name: "Show Toast" });
		await expect(button).toBeInTheDocument();

		// Click the button
		await button.click();

		// Wait for toast to appear and be visible
		await waitFor(
			async () => {
				const toastText = await within(document.body).findByText(
					/Event has been created/i,
				);
				expect(toastText).toBeVisible();
			},
			{ timeout: 3000 },
		);

		await waitFor(
			async () => {
				const descriptionText = await within(document.body).findByText(
					/Sunday, December 03, 2023 at 9:00 AM/i,
				);
				expect(descriptionText).toBeVisible();
			},
			{ timeout: 3000 },
		);
	},
};

export const WithActions: Story = {
	render: () => <ToastWithActionsDemo />,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Wait for component to be fully mounted
		await new Promise((resolve) => setTimeout(resolve, 500));

		const button = canvas.getByRole("button", { name: "With Action" });
		await expect(button).toBeInTheDocument();

		// Click the button
		await button.click();

		// Wait for toast to appear and be visible
		await waitFor(
			async () => {
				const headingText = await within(document.body).findByText(
					/Scheduled: Catch up/i,
				);
				expect(headingText).toBeVisible();
			},
			{ timeout: 3000 },
		);

		await waitFor(
			async () => {
				const descriptionText = await within(document.body).findByText(
					/Friday, February 10, 2024 at 5:57 PM/i,
				);
				expect(descriptionText).toBeVisible();
			},
			{ timeout: 3000 },
		);

		await waitFor(
			async () => {
				const actionButton = await within(document.body).findByRole("button", {
					name: "Undo",
				});
				expect(actionButton).toBeVisible();
			},
			{ timeout: 3000 },
		);
	},
};

export const CustomStyles: Story = {
	render: () => <ToastCustomDemo />,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Wait for component to be fully mounted
		await new Promise((resolve) => setTimeout(resolve, 500));

		const button = canvas.getByRole("button", { name: "Custom Style" });
		await expect(button).toBeInTheDocument();

		// Click the button
		await button.click();

		// Wait for toast to appear and be visible
		await waitFor(
			async () => {
				const headingText = await within(document.body).findByText(
					/Custom Toast/i,
				);
				expect(headingText).toBeVisible();
			},
			{ timeout: 3000 },
		);

		await waitFor(
			async () => {
				const descriptionText = await within(document.body).findByText(
					/With custom styling/i,
				);
				expect(descriptionText).toBeVisible();
			},
			{ timeout: 3000 },
		);

		await waitFor(
			async () => {
				const dismissButton = await within(document.body).findByRole("button", {
					name: "Dismiss",
				});
				expect(dismissButton).toBeVisible();
			},
			{ timeout: 3000 },
		);
	},
};
