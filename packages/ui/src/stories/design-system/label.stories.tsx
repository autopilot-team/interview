import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import { Checkbox } from "../../components/checkbox.js";
import { Input } from "../../components/input.js";
import { Label } from "../../components/label.js";

const meta: Meta<typeof Label> = {
	title: "Design System/Label",
	component: Label,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Label>;

export const Default: Story = {
	render: () => <Label>Default Label</Label>,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const label = canvas.getByText("Default Label");
		await expect(label).toBeInTheDocument();
		await expect(label).toHaveClass("text-sm", "font-medium");
	},
};

export const WithInput: Story = {
	render: () => (
		<div className="grid w-full max-w-sm gap-1.5">
			<Label htmlFor="email">Email</Label>
			<Input type="email" id="email" placeholder="Email" />
		</div>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const label = canvas.getByText("Email");
		const input = canvas.getByPlaceholderText("Email");
		await expect(label).toBeInTheDocument();
		await expect(input).toBeInTheDocument();
		await expect(label).toHaveAttribute("for", "email");
	},
};

export const WithCheckbox: Story = {
	render: () => (
		<div className="flex items-center space-x-2">
			<Checkbox id="terms" />
			<Label htmlFor="terms">Accept terms and conditions</Label>
		</div>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const label = canvas.getByText("Accept terms and conditions");
		await expect(label).toBeInTheDocument();
		await expect(label).toHaveAttribute("for", "terms");
	},
};

export const Required: Story = {
	render: () => (
		<div className="grid w-full max-w-sm gap-1.5">
			<Label
				htmlFor="username"
				className="after:text-red-500 after:content-['*']"
			>
				Username
			</Label>
			<Input id="username" required />
		</div>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const label = canvas.getByText("Username");
		await expect(label).toBeInTheDocument();
		await expect(label).toHaveClass(
			"after:text-red-500",
			"after:content-['*']",
		);
	},
};

export const Disabled: Story = {
	render: () => (
		<div className="grid w-full max-w-sm gap-1.5">
			<Label htmlFor="disabled">Disabled Field</Label>
			<Input id="disabled" disabled />
		</div>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const label = canvas.getByText("Disabled Field");
		await expect(label).toBeInTheDocument();
		await expect(label.parentElement?.querySelector("input")).toBeDisabled();
	},
};
