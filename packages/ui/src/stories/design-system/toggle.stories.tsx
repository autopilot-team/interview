import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { Toggle } from "../../components/toggle.js";

const meta: Meta<typeof Toggle> = {
	title: "Design System/Toggle",
	component: Toggle,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Toggle>;

export const Default: Story = {
	args: {
		"aria-label": "Toggle italic",
		children: (
			<svg
				xmlns="http://www.w3.org/2000/svg"
				width="24"
				height="24"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				strokeWidth="2"
				strokeLinecap="round"
				strokeLinejoin="round"
				className="h-4 w-4"
				aria-hidden="true"
			>
				<line x1="19" y1="4" x2="10" y2="4" />
				<line x1="14" y1="20" x2="5" y2="20" />
				<line x1="15" y1="4" x2="9" y2="20" />
			</svg>
		),
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const toggle = canvas.getByRole("button");
		await expect(toggle).toBeInTheDocument();
		await userEvent.click(toggle);
		await expect(toggle).toHaveAttribute("data-state", "on");
		await userEvent.click(toggle);
		await expect(toggle).toHaveAttribute("data-state", "off");
	},
};

export const WithText: Story = {
	args: {
		"aria-label": "Toggle italic",
		children: (
			<>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					strokeWidth="2"
					strokeLinecap="round"
					strokeLinejoin="round"
					className="h-4 w-4"
					aria-hidden="true"
				>
					<line x1="19" y1="4" x2="10" y2="4" />
					<line x1="14" y1="20" x2="5" y2="20" />
					<line x1="15" y1="4" x2="9" y2="20" />
				</svg>
				Italic
			</>
		),
	},
};

export const Outline: Story = {
	args: {
		variant: "outline",
		"aria-label": "Toggle bold",
		children: (
			<>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					strokeWidth="2"
					strokeLinecap="round"
					strokeLinejoin="round"
					className="h-4 w-4"
					aria-hidden="true"
				>
					<path d="M6 4h8a4 4 0 0 1 4 4 4 4 0 0 1-4 4H6z" />
					<path d="M6 12h9a4 4 0 0 1 4 4 4 4 0 0 1-4 4H6z" />
				</svg>
				Bold
			</>
		),
	},
};

export const Small: Story = {
	args: {
		size: "sm",
		"aria-label": "Toggle underline",
		children: (
			<svg
				xmlns="http://www.w3.org/2000/svg"
				width="24"
				height="24"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				strokeWidth="2"
				strokeLinecap="round"
				strokeLinejoin="round"
				className="h-4 w-4"
				aria-hidden="true"
			>
				<path d="M6 3v7a6 6 0 0 0 6 6 6 6 0 0 0 6-6V3" />
				<line x1="4" y1="21" x2="20" y2="21" />
			</svg>
		),
	},
};

export const Large: Story = {
	args: {
		size: "lg",
		"aria-label": "Toggle strikethrough",
		children: (
			<>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					strokeWidth="2"
					strokeLinecap="round"
					strokeLinejoin="round"
					className="h-4 w-4"
					aria-hidden="true"
				>
					<line x1="4" y1="12" x2="20" y2="12" />
					<path
						d="M16 6c0-2.21-3.58-4-8-4S0 3.79 0 6"
						transform="translate(4 0)"
					/>
					<path d="M8 22c4.42 0 8-1.79 8-4" transform="translate(4 0)" />
				</svg>
				Strikethrough
			</>
		),
	},
};

export const Disabled: Story = {
	args: {
		disabled: true,
		"aria-label": "Toggle disabled",
		children: (
			<>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					strokeWidth="2"
					strokeLinecap="round"
					strokeLinejoin="round"
					className="h-4 w-4"
					aria-hidden="true"
				>
					<path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24" />
					<line x1="1" y1="1" x2="23" y2="23" />
				</svg>
				Hidden
			</>
		),
	},
};
