import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import { Badge } from "../../components/badge.js";

const meta: Meta<typeof Badge> = {
	title: "Design System/Badge",
	component: Badge,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Badge>;

export const Default: Story = {
	render: () => <Badge>Default</Badge>,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const badge = canvas.getByText("Default");
		await expect(badge).toBeInTheDocument();
		await expect(badge).toHaveClass("bg-primary");
	},
};

export const Secondary: Story = {
	render: () => <Badge variant="secondary">Secondary</Badge>,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const badge = canvas.getByText("Secondary");
		await expect(badge).toBeInTheDocument();
		await expect(badge).toHaveClass("bg-secondary");
	},
};

export const Destructive: Story = {
	render: () => <Badge variant="destructive">Destructive</Badge>,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const badge = canvas.getByText("Destructive");
		await expect(badge).toBeInTheDocument();
		await expect(badge).toHaveClass("bg-destructive");
	},
};

export const Outline: Story = {
	render: () => <Badge variant="outline">Outline</Badge>,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const badge = canvas.getByText("Outline");
		await expect(badge).toBeInTheDocument();
		await expect(badge).toHaveClass("text-foreground");
	},
};

export const AllVariants: Story = {
	render: () => (
		<div className="flex gap-4">
			<Badge>Default</Badge>
			<Badge variant="secondary">Secondary</Badge>
			<Badge variant="destructive">Destructive</Badge>
			<Badge variant="outline">Outline</Badge>
		</div>
	),
};
