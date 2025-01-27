import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { Input } from "../../components/input.js";

const meta = {
	title: "Design System/Input",
	component: Input,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
	argTypes: {
		type: {
			control: "select",
			options: ["text", "password", "email", "number", "search", "tel", "url"],
		},
		disabled: {
			control: "boolean",
		},
		placeholder: {
			control: "text",
		},
	},
} satisfies Meta<typeof Input>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {
		placeholder: "Enter text...",
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const input = canvas.getByPlaceholderText("Enter text...");
		await expect(input).toBeInTheDocument();
		await userEvent.type(input, "Hello, World!");
		await expect(input).toHaveValue("Hello, World!");
	},
};

export const Disabled: Story = {
	args: {
		disabled: true,
		placeholder: "Disabled input",
		value: "Can't edit this",
	},
};

export const WithLabel: Story = {
	args: {
		id: "email",
		type: "email",
		placeholder: "Enter your email",
	},
	decorators: [
		(Story) => (
			<div className="grid w-full max-w-sm gap-1.5">
				<label
					htmlFor="email"
					className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
				>
					Email
				</label>
				<Story />
			</div>
		),
	],
};

export const File: Story = {
	args: {
		type: "file",
		className: "cursor-pointer",
	},
};

export const WithIcon: Story = {
	args: {
		type: "search",
		placeholder: "Search...",
	},
	decorators: [
		(Story) => (
			<div className="relative w-full max-w-sm">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="16"
					height="16"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					strokeWidth="2"
					strokeLinecap="round"
					strokeLinejoin="round"
					className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground"
				>
					<title>Search</title>
					<circle cx="11" cy="11" r="8" />
					<path d="m21 21-4.35-4.35" />
				</svg>
				<Story />
			</div>
		),
	],
};
