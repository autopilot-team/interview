import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { createRoutesStub } from "react-router";
import { ForgotPasswordForm } from "../../components/forgot-password-form.js";

const meta = {
	title: "Application/ForgotPasswordForm",
	component: ForgotPasswordForm,
	parameters: {
		layout: "centered",
	},
	decorators: [
		(Story) => {
			const Stub = createRoutesStub([
				{
					path: "/forgot-password",
					Component: Story,
				},
				{
					path: "/sign-in",
					Component: () => null,
				},
			]);

			return <Stub initialEntries={["/forgot-password"]} />;
		},
	],
	args: {
		cfTurnstileSiteKey: "",
		t: {
			title: "Forgot password",
			description:
				"Enter your email address and we'll send you a link to reset your password.",
			email: "Email",
			emailPlaceholder: "Enter your email",
			resetPassword: "Send reset link",
			backToSignIn: "Back to sign in",
			errors: {
				emailRequired: "Email is required",
				emailInvalid: "Please enter a valid email",
			},
		},
	},
	argTypes: {
		handleForgotPassword: { action: "handleForgotPassword" },
	},
} satisfies Meta<typeof ForgotPasswordForm>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const FillingForm: Story = {
	play: async ({ canvasElement, args }) => {
		const canvas = within(canvasElement);

		// Get form elements
		const emailInput = canvas.getByPlaceholderText(
			meta.args.t.emailPlaceholder,
		);
		const submitButton = canvas.getByRole("button", {
			name: meta.args.t.resetPassword,
		});

		// Fill in email
		await userEvent.type(emailInput, "test@example.com", { delay: 10 });
		await expect(emailInput).toHaveValue("test@example.com");

		// Submit form
		await userEvent.click(submitButton);
	},
};

export const ValidationErrors: Story = {
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Get form elements
		const emailInput = canvas.getByPlaceholderText(
			meta.args.t.emailPlaceholder,
		);
		const submitButton = canvas.getByRole("button", {
			name: meta.args.t.resetPassword,
		});

		// Test empty form submission
		await userEvent.click(submitButton);
		await expect(
			await canvas.findByText(meta.args.t.errors.emailRequired),
		).toBeInTheDocument();

		// Test invalid email
		await userEvent.type(emailInput, "invalid-email", { delay: 10 });
		await userEvent.click(submitButton);
		await expect(
			await canvas.findByText(meta.args.t.errors.emailInvalid),
		).toBeInTheDocument();
	},
};

export const LoadingState: Story = {
	args: {
		isLoading: true,
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Get form elements
		const emailInput = canvas.getByPlaceholderText(
			meta.args.t.emailPlaceholder,
		);
		const submitButton = canvas.getByRole("button", {
			name: meta.args.t.resetPassword,
		});

		// Verify form elements are disabled
		await expect(emailInput).toBeDisabled();
		await expect(submitButton).toBeDisabled();

		// Verify back to sign in link is disabled
		const signInLink = canvas.getByRole("link", {
			name: meta.args.t.backToSignIn,
		});
		await expect(signInLink).toHaveAttribute("aria-disabled", "true");
		await expect(signInLink).toHaveAttribute("tabIndex", "-1");
	},
};

export const NavigationLinks: Story = {
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Verify back to sign in link
		const signInLink = canvas.getByRole("link", {
			name: meta.args.t.backToSignIn,
		});
		await expect(signInLink).toHaveAttribute("href", "/sign-in");
	},
};
