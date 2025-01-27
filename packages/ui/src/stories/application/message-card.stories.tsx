import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import { createRoutesStub } from "react-router";
import { MessageCard } from "../../components/message-card.js";

const meta = {
	title: "Application/MessageCard",
	component: MessageCard,
	decorators: [
		(Story) => {
			const Stub = createRoutesStub([
				{
					path: "/",
					Component: Story,
				},
			]);

			return <Stub initialEntries={["/"]} />;
		},
	],
	parameters: {
		layout: "centered",
	},
} satisfies Meta<typeof MessageCard>;

export default meta;
type Story = StoryObj<typeof MessageCard>;

export const Default: Story = {
	args: {
		title: "Welcome",
		description: "Thank you for using our platform.",
		backButton: { show: true, label: "Go Back" },
		homeButton: { show: true, label: "Return Home", to: "/" },
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Verify title and description are rendered
		const title = canvas.getByText("Welcome");
		const description = canvas.getByText("Thank you for using our platform.");
		await expect(title).toBeVisible();
		await expect(description).toBeVisible();

		// Verify buttons are rendered
		const homeButton = canvas.getByText("Return Home");
		const backButton = canvas.getByText("Go Back");
		await expect(homeButton).toBeVisible();
		await expect(backButton).toBeVisible();
	},
};

export const Success: Story = {
	args: {
		title: "Payment Successful",
		description: "Your payment has been processed successfully.",
		variant: "success",
		footer: "Transaction ID: 123456789",
	},
};

export const ErrorState: Story = {
	args: {
		title: "Error Occurred",
		description: "We couldn't process your request at this time.",
		variant: "error",
		stack:
			"Error: Failed to process payment\n  at ProcessPayment (/src/payment.ts:42:5)\n  at async handlePayment (/src/handler.ts:15:3)",
	},
};

export const WithoutButtons: Story = {
	args: {
		title: "Information",
		description: "This is an informational message.",
		backButton: { show: false },
		homeButton: { show: false },
	},
};

export const CustomIcons: Story = {
	args: {
		title: "Custom Navigation",
		description: "This message card uses custom icons for navigation.",
		backButton: {
			show: true,
			label: "Previous",
			icon: <span>←</span>,
		},
		homeButton: {
			show: true,
			label: "Dashboard",
			icon: <span>🏠</span>,
			to: "/dashboard",
		},
	},
};
