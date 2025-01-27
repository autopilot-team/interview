import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { createRoutesStub } from "react-router";
import { SignInForm } from "../../components/sign-in-form.js";

const meta = {
	title: "Application/SignInForm",
	component: SignInForm,
	parameters: {
		layout: "centered",
	},
	decorators: [
		(Story) => {
			const Stub = createRoutesStub([
				{
					path: "/sign-in",
					Component: Story,
				},
				{
					path: "/sign-up",
					Component: () => null,
				},
				{
					path: "/forgot-password",
					Component: () => null,
				},
				{
					path: "/terms-of-service",
					Component: () => null,
				},
				{
					path: "/privacy-policy",
					Component: () => null,
				},
			]);

			return <Stub initialEntries={["/sign-in"]} />;
		},
	],
	args: {
		cfTurnstileSiteKey: "",
		t: {
			title: "Welcome back",
			email: "Email",
			emailPlaceholder: "Enter your email",
			password: "Password",
			passwordPlaceholder: "Enter your password",
			forgotPassword: "Forgot password?",
			signIn: "Sign in",
			noAccount: "Don't have an account?",
			signUp: "Sign up",
			termsText: "By clicking {button}, you agree to our {terms} and {privacy}",
			termsButton: "Sign in",
			termsOfService: "Terms of Service",
			privacyPolicy: "Privacy Policy",
			errors: {
				emailRequired: "Email is required",
				emailInvalid: "Please enter a valid email",
				passwordRequired: "Password is required",
				passwordMinLength: "Password must be at least 8 characters",
			},
		},
	},
	argTypes: {
		handleSignIn: { action: "handleSignIn" },
	},
} satisfies Meta<typeof SignInForm>;

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
		const passwordInput = canvas.getByPlaceholderText(
			meta.args.t.passwordPlaceholder,
		);
		const submitButton = canvas.getByRole("button", {
			name: meta.args.t.signIn,
		});

		// Fill in email
		await userEvent.type(emailInput, "test@example.com", { delay: 10 });
		await expect(emailInput).toHaveValue("test@example.com");

		// Fill in password
		await userEvent.type(passwordInput, "Password123!", { delay: 10 });
		await expect(passwordInput).toHaveValue("Password123!");

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
		const passwordInput = canvas.getByPlaceholderText(
			meta.args.t.passwordPlaceholder,
		);
		const submitButton = canvas.getByRole("button", {
			name: meta.args.t.signIn,
		});

		// Test empty form submission
		await userEvent.click(submitButton);

		// Verify error messages
		await expect(
			await canvas.findByText(meta.args.t.errors.emailRequired),
		).toBeInTheDocument();
		await expect(
			await canvas.findByText(meta.args.t.errors.passwordRequired),
		).toBeInTheDocument();

		// Test invalid email
		await userEvent.type(emailInput, "invalid-email", { delay: 10 });
		await userEvent.click(submitButton);
		await expect(
			await canvas.findByText(meta.args.t.errors.emailInvalid),
		).toBeInTheDocument();

		// Test short password
		await userEvent.clear(emailInput);
		await userEvent.type(emailInput, "test@example.com", { delay: 10 });
		await userEvent.type(passwordInput, "short", { delay: 10 });
		await userEvent.click(submitButton);
		await expect(
			await canvas.findByText(meta.args.t.errors.passwordMinLength),
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
		const passwordInput = canvas.getByPlaceholderText(
			meta.args.t.passwordPlaceholder,
		);
		const submitButton = canvas.getByRole("button", {
			name: meta.args.t.signIn,
		});

		// Verify form elements are disabled
		await expect(emailInput).toBeDisabled();
		await expect(passwordInput).toBeDisabled();
		await expect(submitButton).toBeDisabled();

		// Verify navigation links are disabled
		const forgotPasswordLink = canvas.getByRole("link", {
			name: meta.args.t.forgotPassword,
		});
		const signUpLink = canvas.getByRole("link", { name: meta.args.t.signUp });

		await expect(forgotPasswordLink).toHaveAttribute("aria-disabled", "true");
		await expect(forgotPasswordLink).toHaveAttribute("tabIndex", "-1");
		await expect(signUpLink).toHaveAttribute("aria-disabled", "true");
		await expect(signUpLink).toHaveAttribute("tabIndex", "-1");
	},
};

export const PasswordVisibility: Story = {
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Get password input and toggle button
		const passwordInput = canvas.getByPlaceholderText(
			meta.args.t.passwordPlaceholder,
		);
		const toggleButton = canvas.getByRole("button", {
			name: "Show password",
		});

		// Verify password is initially hidden
		await expect(passwordInput).toHaveAttribute("type", "password");

		// Toggle password visibility
		await userEvent.click(toggleButton);
		await expect(passwordInput).toHaveAttribute("type", "text");

		// Toggle back to hidden
		await userEvent.click(toggleButton);
		await expect(passwordInput).toHaveAttribute("type", "password");
	},
};

export const NavigationLinks: Story = {
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Verify all navigation links are present
		const forgotPasswordLink = canvas.getByRole("link", {
			name: meta.args.t.forgotPassword,
		});
		const signUpLink = canvas.getByRole("link", { name: meta.args.t.signUp });
		const termsLink = canvas.getByRole("link", {
			name: meta.args.t.termsOfService,
		});
		const privacyLink = canvas.getByRole("link", {
			name: meta.args.t.privacyPolicy,
		});

		await expect(forgotPasswordLink).toHaveAttribute(
			"href",
			"/forgot-password",
		);
		await expect(signUpLink).toHaveAttribute("href", "/sign-up");
		await expect(termsLink).toHaveAttribute("href", "/terms-of-service");
		await expect(privacyLink).toHaveAttribute("href", "/privacy-policy");
	},
};
