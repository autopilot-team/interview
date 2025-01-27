import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import { Label } from "../../components/label.js";
import { Textarea } from "../../components/textarea.js";

const meta: Meta<typeof Textarea> = {
	title: "Design System/Textarea",
	component: Textarea,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Textarea>;

export const Default: Story = {
	render: () => <Textarea placeholder="Type your message here." />,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const textarea = canvas.getByPlaceholderText("Type your message here.");
		await expect(textarea).toBeInTheDocument();
		await expect(textarea).toHaveClass("min-h-[60px]", "w-full", "rounded-md");
	},
};

export const WithLabel: Story = {
	render: () => (
		<div className="grid w-full gap-1.5">
			<Label htmlFor="message">Your message</Label>
			<Textarea id="message" placeholder="Type your message here." />
		</div>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const label = canvas.getByText("Your message");
		const textarea = canvas.getByPlaceholderText("Type your message here.");
		await expect(label).toBeInTheDocument();
		await expect(textarea).toBeInTheDocument();
		await expect(label).toHaveAttribute("for", "message");
	},
};

export const Disabled: Story = {
	render: () => <Textarea disabled placeholder="You cannot type here." />,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const textarea = canvas.getByPlaceholderText("You cannot type here.");
		await expect(textarea).toBeInTheDocument();
		await expect(textarea).toBeDisabled();
		await expect(textarea).toHaveClass(
			"disabled:cursor-not-allowed",
			"disabled:opacity-50",
		);
	},
};

export const WithDefaultValue: Story = {
	render: () => (
		<Textarea defaultValue="This is some default text that can be edited." />
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const textarea = canvas.getByDisplayValue(
			"This is some default text that can be edited.",
		);
		await expect(textarea).toBeInTheDocument();
		await expect(textarea).toHaveValue(
			"This is some default text that can be edited.",
		);
	},
};

export const WithRows: Story = {
	render: () => <Textarea placeholder="This textarea has 10 rows" rows={10} />,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const textarea = canvas.getByPlaceholderText("This textarea has 10 rows");
		await expect(textarea).toBeInTheDocument();
		await expect(textarea).toHaveAttribute("rows", "10");
	},
};
