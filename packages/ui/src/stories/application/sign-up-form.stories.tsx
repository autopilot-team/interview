import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { createRoutesStub } from "react-router";
import { SignUpForm } from "../../components/sign-up-form.js";

const meta = {
	title: "Application/SignUpForm",
	component: SignUpForm,
	parameters: {
		layout: "centered",
	},
	decorators: [
		(Story) => {
			const Stub = createRoutesStub([
				{
					path: "/sign-up",
					Component: Story,
				},
				{
					path: "/sign-in",
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

			return <Stub initialEntries={["/sign-up"]} />;
		},
	],
	args: {
		cfTurnstileSiteKey: "",
		t: {
			title: "Create an account",
			name: "Full name",
			namePlaceholder: "Enter your full name",
			email: "Email",
			emailPlaceholder: "Enter your email",
			password: "Password",
			passwordPlaceholder: "Create a password",
			confirmPassword: "Confirm password",
			confirmPasswordPlaceholder: "Re-enter your password",
			signUp: "Sign up",
			haveAccount: "Already have an account?",
			signIn: "Sign in",
			termsText: "By clicking {button}, you agree to our {terms} and {privacy}",
			termsButton: "Sign up",
			termsOfService: "Terms of Service",
			privacyPolicy: "Privacy Policy",
			passwordStrength: {
				secure: "Password is secure",
				moderate: "Password is moderate",
				weak: "Password is weak",
				requirements: {
					minLength: "Use at least 8 characters",
					mixCase: "Mix uppercase & lowercase letters",
					number: "Add a number",
					special: "Add a special character",
				},
			},
			errors: {
				nameRequired: "Full name is required",
				nameMinLength: "Full name must be at least 2 characters",
				emailRequired: "Email is required",
				emailInvalid: "Please enter a valid email",
				passwordRequired: "Password is required",
				passwordMinLength: "Password must be at least 8 characters",
				passwordUppercase: "Password must contain an uppercase letter",
				passwordLowercase: "Password must contain a lowercase letter",
				passwordNumber: "Password must contain a number",
				passwordSpecial: "Password must contain a special character",
				confirmPasswordRequired: "Please confirm your password",
				confirmPasswordMatch: "Passwords do not match",
			},
		},
	},
	argTypes: {
		handleSignUp: { action: "handleSignUp" },
	},
} satisfies Meta<typeof SignUpForm>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const FillingForm: Story = {
	play: async ({ canvasElement, args }) => {
		const canvas = within(canvasElement);

		// Get form elements
		const nameInput = canvas.getByPlaceholderText(meta.args.t.namePlaceholder);
		const emailInput = canvas.getByPlaceholderText(
			meta.args.t.emailPlaceholder,
		);
		const passwordInput = canvas.getByPlaceholderText(
			meta.args.t.passwordPlaceholder,
		);
		const confirmPasswordInput = canvas.getByPlaceholderText(
			meta.args.t.confirmPasswordPlaceholder,
		);
		const submitButton = canvas.getByRole("button", {
			name: meta.args.t.signUp,
		});

		// Fill in name
		await userEvent.type(nameInput, "John Doe", { delay: 10 });
		await expect(nameInput).toHaveValue("John Doe");

		// Fill in email
		await userEvent.type(emailInput, "john@example.com", { delay: 10 });
		await expect(emailInput).toHaveValue("john@example.com");

		// Fill in password
		await userEvent.type(passwordInput, "StrongPass123!", { delay: 10 });
		await expect(passwordInput).toHaveValue("StrongPass123!");

		// Fill in confirm password
		await userEvent.type(confirmPasswordInput, "StrongPass123!", { delay: 10 });
		await expect(confirmPasswordInput).toHaveValue("StrongPass123!");

		// Submit form
		await userEvent.click(submitButton);
	},
};

export const ValidationErrors: Story = {
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Get form elements
		const nameInput = canvas.getByPlaceholderText(meta.args.t.namePlaceholder);
		const emailInput = canvas.getByPlaceholderText(
			meta.args.t.emailPlaceholder,
		);
		const passwordInput = canvas.getByPlaceholderText(
			meta.args.t.passwordPlaceholder,
		);
		const confirmPasswordInput = canvas.getByPlaceholderText(
			meta.args.t.confirmPasswordPlaceholder,
		);
		const submitButton = canvas.getByRole("button", {
			name: meta.args.t.signUp,
		});

		// Test empty form submission
		await userEvent.click(submitButton);
		await expect(
			await canvas.findByText(meta.args.t.errors.nameRequired),
		).toBeInTheDocument();
		await expect(
			await canvas.findByText(meta.args.t.errors.emailRequired),
		).toBeInTheDocument();
		await expect(
			await canvas.findByText(meta.args.t.errors.passwordRequired),
		).toBeInTheDocument();

		// Test invalid name
		await userEvent.type(nameInput, "a", { delay: 10 });
		await userEvent.click(submitButton);
		await expect(
			await canvas.findByText(meta.args.t.errors.nameMinLength),
		).toBeInTheDocument();

		// Test invalid email
		await userEvent.clear(nameInput);
		await userEvent.type(nameInput, "John Doe", { delay: 10 });
		await userEvent.type(emailInput, "invalid-email", { delay: 10 });
		await userEvent.click(submitButton);
		await expect(
			await canvas.findByText(meta.args.t.errors.emailInvalid),
		).toBeInTheDocument();

		// Test password requirements one by one
		await userEvent.clear(emailInput);
		await userEvent.type(emailInput, "john@example.com", { delay: 10 });

		// Test minimum length
		await userEvent.clear(passwordInput);
		await userEvent.type(passwordInput, "weak", { delay: 10 });
		await userEvent.click(submitButton);
		await expect(
			await canvas.findByText(meta.args.t.errors.passwordMinLength),
		).toBeInTheDocument();

		// Test uppercase requirement
		await userEvent.clear(passwordInput);
		await userEvent.type(passwordInput, "password123!", { delay: 10 });
		await userEvent.click(submitButton);
		await expect(
			await canvas.findByText(meta.args.t.errors.passwordUppercase),
		).toBeInTheDocument();

		// Test number requirement
		await userEvent.clear(passwordInput);
		await userEvent.type(passwordInput, "Password!", { delay: 10 });
		await userEvent.click(submitButton);
		await expect(
			await canvas.findByText(meta.args.t.errors.passwordNumber),
		).toBeInTheDocument();

		// Test special character requirement
		await userEvent.clear(passwordInput);
		await userEvent.type(passwordInput, "Password123", { delay: 10 });
		await userEvent.click(submitButton);
		await expect(
			await canvas.findByText(meta.args.t.errors.passwordSpecial),
		).toBeInTheDocument();

		// Test password mismatch
		await userEvent.clear(passwordInput);
		await userEvent.type(passwordInput, "StrongPass123!", { delay: 10 });
		await userEvent.type(confirmPasswordInput, "DifferentPass123!", {
			delay: 10,
		});
		await userEvent.click(submitButton);
		await expect(
			await canvas.findByText(meta.args.t.errors.confirmPasswordMatch),
		).toBeInTheDocument();
	},
};

export const PasswordStrengthIndicator: Story = {
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		const passwordInput = canvas.getByPlaceholderText(
			meta.args.t.passwordPlaceholder,
		);

		// Test weak password
		await userEvent.type(passwordInput, "weak", { delay: 10 });
		await expect(
			await canvas.findByText(
				meta.args.t.passwordStrength.requirements.minLength,
			),
		).toBeInTheDocument();

		// Test moderate password (meets length, case, and number requirements)
		await userEvent.clear(passwordInput);
		await userEvent.type(passwordInput, "Password123", { delay: 10 });
		await expect(
			await canvas.findByText(
				meta.args.t.passwordStrength.requirements.special,
			),
		).toBeInTheDocument();

		// Test secure password (meets all requirements)
		await userEvent.clear(passwordInput);
		await userEvent.type(passwordInput, "Password123!", { delay: 10 });

		// Wait for any requirement text to disappear since all requirements are met
		await new Promise((resolve) => setTimeout(resolve, 50));
		const hints = canvas.queryByText(
			meta.args.t.passwordStrength.requirements.special,
		);
		await expect(hints).not.toBeInTheDocument();
	},
};

export const PasswordVisibility: Story = {
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Get password inputs
		const passwordInput = canvas.getByPlaceholderText(
			meta.args.t.passwordPlaceholder,
		);
		const confirmPasswordInput = canvas.getByPlaceholderText(
			meta.args.t.confirmPasswordPlaceholder,
		);

		// Get toggle buttons by finding them within their respective input containers
		const passwordContainer = passwordInput.closest(".relative") as HTMLElement;
		const confirmContainer = confirmPasswordInput.closest(
			".relative",
		) as HTMLElement;

		if (!passwordContainer || !confirmContainer) {
			throw new Error("Password input containers not found");
		}

		const togglePasswordButton = within(passwordContainer).getByRole("button");
		const toggleConfirmButton = within(confirmContainer).getByRole("button");

		// Test password field visibility
		await expect(passwordInput).toHaveAttribute("type", "password");
		await userEvent.click(togglePasswordButton);
		await expect(passwordInput).toHaveAttribute("type", "text");
		await userEvent.click(togglePasswordButton);
		await expect(passwordInput).toHaveAttribute("type", "password");

		// Test confirm password field visibility
		await expect(confirmPasswordInput).toHaveAttribute("type", "password");
		await userEvent.click(toggleConfirmButton);
		await expect(confirmPasswordInput).toHaveAttribute("type", "text");
		await userEvent.click(toggleConfirmButton);
		await expect(confirmPasswordInput).toHaveAttribute("type", "password");
	},
};

export const NavigationLinks: Story = {
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		const signInLink = canvas.getByRole("link", { name: meta.args.t.signIn });
		const termsLink = canvas.getByRole("link", {
			name: meta.args.t.termsOfService,
		});
		const privacyLink = canvas.getByRole("link", {
			name: meta.args.t.privacyPolicy,
		});

		await expect(signInLink).toHaveAttribute("href", "/sign-in");
		await expect(termsLink).toHaveAttribute("href", "/terms-of-service");
		await expect(privacyLink).toHaveAttribute("href", "/privacy-policy");
	},
};

export const LoadingState: Story = {
	args: {
		isLoading: true,
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Get form elements
		const nameInput = canvas.getByPlaceholderText(meta.args.t.namePlaceholder);
		const emailInput = canvas.getByPlaceholderText(
			meta.args.t.emailPlaceholder,
		);
		const passwordInput = canvas.getByPlaceholderText(
			meta.args.t.passwordPlaceholder,
		);
		const confirmPasswordInput = canvas.getByPlaceholderText(
			meta.args.t.confirmPasswordPlaceholder,
		);
		const submitButton = canvas.getByRole("button", {
			name: meta.args.t.signUp,
		});

		// Get password toggle buttons
		const passwordContainer = passwordInput.closest(".relative") as HTMLElement;
		const confirmContainer = confirmPasswordInput.closest(
			".relative",
		) as HTMLElement;

		if (!passwordContainer || !confirmContainer) {
			throw new Error("Password input containers not found");
		}

		const togglePasswordButton = within(passwordContainer).getByRole("button");
		const toggleConfirmButton = within(confirmContainer).getByRole("button");

		// Verify form elements are disabled
		await expect(nameInput).toBeDisabled();
		await expect(emailInput).toBeDisabled();
		await expect(passwordInput).toBeDisabled();
		await expect(confirmPasswordInput).toBeDisabled();
		await expect(submitButton).toBeDisabled();
		await expect(togglePasswordButton).toBeDisabled();
		await expect(toggleConfirmButton).toBeDisabled();

		// Verify navigation links are disabled
		const signInLink = canvas.getByRole("link", { name: meta.args.t.signIn });
		const termsLink = canvas.getByRole("link", {
			name: meta.args.t.termsOfService,
		});
		const privacyLink = canvas.getByRole("link", {
			name: meta.args.t.privacyPolicy,
		});

		await expect(signInLink).toHaveAttribute("aria-disabled", "true");
		await expect(signInLink).toHaveAttribute("tabIndex", "-1");
		await expect(termsLink).toHaveAttribute("aria-disabled", "true");
		await expect(termsLink).toHaveAttribute("tabIndex", "-1");
		await expect(privacyLink).toHaveAttribute("aria-disabled", "true");
		await expect(privacyLink).toHaveAttribute("tabIndex", "-1");
	},
};
