import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import {
	InputOTP,
	InputOTPGroup,
	InputOTPSeparator,
	InputOTPSlot,
} from "../../components/input-otp.js";

const meta: Meta<typeof InputOTP> = {
	title: "Design System/InputOTP",
	component: InputOTP,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof InputOTP>;

const getInput = (inputs: HTMLElement[], index: number): HTMLInputElement => {
	const input = inputs[index];
	if (!input) throw new Error(`Input at index ${index} not found`);
	return input as HTMLInputElement;
};

export const Default: Story = {
	render: () => (
		<InputOTP maxLength={6}>
			<InputOTPGroup>
				<InputOTPSlot index={0} />
				<InputOTPSlot index={1} />
				<InputOTPSlot index={2} />
				<InputOTPSlot index={3} />
				<InputOTPSlot index={4} />
				<InputOTPSlot index={5} />
			</InputOTPGroup>
		</InputOTP>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Wait for input to be mounted
		await new Promise((resolve) => setTimeout(resolve, 100));

		// Find the OTP input
		const input = await canvas.findByRole("textbox", {}, { timeout: 1000 });
		await expect(input).toBeInTheDocument();
		await expect(input).toHaveAttribute("maxlength", "6");
		await expect(input).toHaveAttribute("data-input-otp", "true");

		// Test initial state
		await expect(input).not.toBeDisabled();
		await expect(input).toHaveValue("");

		// Test input sequence
		await userEvent.type(input, "1");
		await expect(input).toHaveValue("1");

		await userEvent.type(input, "23456");
		await expect(input).toHaveValue("123456");
	},
};

export const WithSeparator: Story = {
	render: () => (
		<InputOTP maxLength={6}>
			<InputOTPGroup>
				<InputOTPSlot index={0} />
				<InputOTPSlot index={1} />
				<InputOTPSlot index={2} />
				<InputOTPSeparator />
				<InputOTPSlot index={3} />
				<InputOTPSlot index={4} />
				<InputOTPSlot index={5} />
			</InputOTPGroup>
		</InputOTP>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Wait for input to be mounted
		await new Promise((resolve) => setTimeout(resolve, 100));

		// Find the OTP input
		const input = await canvas.findByRole("textbox", {}, { timeout: 1000 });
		await expect(input).toBeInTheDocument();
		await expect(input).toHaveAttribute("maxlength", "6");

		const separator = await canvas.findByRole(
			"separator",
			{},
			{ timeout: 1000 },
		);
		await expect(separator).toBeInTheDocument();

		// Test input sequence
		await userEvent.type(input, "123456");
		await expect(input).toHaveValue("123456");
	},
};

export const Disabled: Story = {
	render: () => (
		<InputOTP maxLength={4} disabled>
			<InputOTPGroup>
				<InputOTPSlot index={0} />
				<InputOTPSlot index={1} />
				<InputOTPSlot index={2} />
				<InputOTPSlot index={3} />
			</InputOTPGroup>
		</InputOTP>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Wait for input to be mounted
		await new Promise((resolve) => setTimeout(resolve, 100));

		// Find the OTP input
		const input = await canvas.findByRole("textbox", {}, { timeout: 1000 });
		await expect(input).toBeInTheDocument();
		await expect(input).toHaveAttribute("maxlength", "4");
		await expect(input).toBeDisabled();
		await expect(input).toHaveValue("");

		// Attempt to type should not work
		await userEvent.type(input, "1");
		await expect(input).toHaveValue("");
	},
};

export const WithValue: Story = {
	render: () => (
		<InputOTP maxLength={4} value="1234">
			<InputOTPGroup>
				<InputOTPSlot index={0} />
				<InputOTPSlot index={1} />
				<InputOTPSlot index={2} />
				<InputOTPSlot index={3} />
			</InputOTPGroup>
		</InputOTP>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Wait for input to be mounted
		await new Promise((resolve) => setTimeout(resolve, 100));

		// Find the OTP input
		const input = await canvas.findByRole("textbox", {}, { timeout: 1000 });
		await expect(input).toBeInTheDocument();
		await expect(input).toHaveAttribute("maxlength", "4");
		await expect(input).toHaveAttribute("data-input-otp", "true");
		await expect(input).toHaveValue("1234");

		// Since this is a controlled input, we can only verify the initial value
		// We cannot test value modifications as they need to be handled by the parent component
	},
};

export const CustomStyle: Story = {
	render: () => (
		<InputOTP maxLength={4}>
			<InputOTPGroup className="gap-4">
				<InputOTPSlot
					index={0}
					className="rounded-xl border-2 border-primary"
				/>
				<InputOTPSlot
					index={1}
					className="rounded-xl border-2 border-primary"
				/>
				<InputOTPSlot
					index={2}
					className="rounded-xl border-2 border-primary"
				/>
				<InputOTPSlot
					index={3}
					className="rounded-xl border-2 border-primary"
				/>
			</InputOTPGroup>
		</InputOTP>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Wait for input to be mounted
		await new Promise((resolve) => setTimeout(resolve, 100));

		// Find the OTP input
		const input = await canvas.findByRole("textbox", {}, { timeout: 1000 });
		await expect(input).toBeInTheDocument();
		await expect(input).toHaveAttribute("maxlength", "4");
		await expect(input).toHaveClass("disabled:cursor-not-allowed");

		// Test input sequence
		await userEvent.type(input, "1234");
		await expect(input).toHaveValue("1234");
	},
};
