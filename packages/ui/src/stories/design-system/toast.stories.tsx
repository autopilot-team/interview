import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import { Button } from "../../components/button.js";
import {
	Toast,
	ToastAction,
	ToastClose,
	ToastDescription,
	ToastProvider,
	ToastTitle,
	ToastViewport,
} from "../../components/toast.js";

const meta: Meta<typeof Toast> = {
	title: "Design System/Toast",
	component: Toast,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
	decorators: [
		(Story) => (
			<ToastProvider>
				<Story />
				<ToastViewport />
			</ToastProvider>
		),
	],
};

export default meta;
type Story = StoryObj<typeof Toast>;

export const Default: Story = {
	render: () => (
		<Toast>
			<div className="grid gap-1">
				<ToastTitle>Scheduled: Catch up</ToastTitle>
				<ToastDescription>
					Friday, February 10, 2024 at 5:57 PM
				</ToastDescription>
			</div>
		</Toast>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const toast = canvas.getByRole("status");
		await expect(toast).toBeInTheDocument();
		const title = canvas.getByText("Scheduled: Catch up");
		await expect(title).toBeVisible();
	},
};

export const WithAction: Story = {
	render: () => (
		<Toast>
			<div className="grid gap-1">
				<ToastTitle>Success!</ToastTitle>
				<ToastDescription>Your message has been sent.</ToastDescription>
			</div>
			<ToastAction altText="Undo send message" asChild>
				<Button variant="outline" size="sm">
					Undo
				</Button>
			</ToastAction>
		</Toast>
	),
};

export const Destructive: Story = {
	render: () => (
		<Toast variant="destructive">
			<div className="grid gap-1">
				<ToastTitle>Error</ToastTitle>
				<ToastDescription>
					There was a problem with your request.
				</ToastDescription>
			</div>
			<ToastClose />
		</Toast>
	),
};

export const WithCustomStyling: Story = {
	render: () => (
		<Toast className="bg-primary text-primary-foreground">
			<div className="grid gap-1">
				<ToastTitle className="text-primary-foreground">
					Custom Theme
				</ToastTitle>
				<ToastDescription className="text-primary-foreground/90">
					This toast uses custom background and text colors.
				</ToastDescription>
			</div>
			<ToastAction
				altText="Try again"
				className="bg-primary-foreground text-primary hover:bg-primary-foreground/90"
				asChild
			>
				<Button variant="outline" size="sm">
					Try again
				</Button>
			</ToastAction>
		</Toast>
	),
};

export const WithLongContent: Story = {
	render: () => (
		<Toast>
			<div className="grid gap-1">
				<ToastTitle>Message</ToastTitle>
				<ToastDescription>
					This is a longer message that demonstrates how the toast handles
					multiple lines of text. It should wrap nicely and maintain readability
					while staying within the bounds of the toast container.
				</ToastDescription>
			</div>
			<ToastClose />
		</Toast>
	),
};
