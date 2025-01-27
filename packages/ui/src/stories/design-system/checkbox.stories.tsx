import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { Checkbox } from "../../components/checkbox.js";

const meta = {
	title: "Design System/Checkbox",
	component: Checkbox,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
	argTypes: {
		disabled: {
			control: "boolean",
		},
		defaultChecked: {
			control: "boolean",
		},
	},
} satisfies Meta<typeof Checkbox>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const checkbox = canvas.getByRole("checkbox");
		await expect(checkbox).toBeInTheDocument();
		await userEvent.click(checkbox);
		await expect(checkbox).toBeChecked();
	},
};

export const Checked: Story = {
	args: {
		defaultChecked: true,
	},
};

export const Disabled: Story = {
	args: {
		disabled: true,
	},
};

export const DisabledChecked: Story = {
	args: {
		disabled: true,
		defaultChecked: true,
	},
};

export const WithLabel: Story = {
	args: {
		id: "terms",
	},
	decorators: [
		(Story) => (
			<div className="flex items-center space-x-2">
				<Story />
				<label
					htmlFor="terms"
					className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
				>
					Accept terms and conditions
				</label>
			</div>
		),
	],
};

export const WithDescription: Story = {
	args: {
		id: "marketing",
	},
	decorators: [
		(Story) => (
			<div className="grid gap-1.5">
				<div className="flex items-center space-x-2">
					<Story />
					<label
						htmlFor="marketing"
						className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
					>
						Marketing emails
					</label>
				</div>
				<p className="text-sm text-muted-foreground">
					Receive emails about new products, features, and more.
				</p>
			</div>
		),
	],
};
