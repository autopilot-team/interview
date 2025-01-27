import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import {
	Avatar,
	AvatarFallback,
	AvatarImage,
} from "../../components/avatar.js";

const meta = {
	title: "Design System/Avatar",
	component: Avatar,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
} satisfies Meta<typeof Avatar>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {
		children: (
			<>
				<AvatarImage src="https://github.com/shadcn.png" alt="@shadcn" />
				<AvatarFallback>CN</AvatarFallback>
			</>
		),
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// First check if the fallback is visible (it will be until the image loads)
		const fallback = canvas.getByText("CN");
		await expect(fallback).toBeInTheDocument();
	},
};

export const WithFallback: Story = {
	args: {
		children: (
			<>
				<AvatarImage src="/broken-image.jpg" alt="@johndoe" />
				<AvatarFallback>JD</AvatarFallback>
			</>
		),
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const fallback = canvas.getByText("JD");
		await expect(fallback).toBeInTheDocument();
	},
};

export const CustomSize: Story = {
	args: {
		className: "h-16 w-16",
		children: (
			<>
				<AvatarImage src="https://github.com/shadcn.png" alt="@shadcn" />
				<AvatarFallback>CN</AvatarFallback>
			</>
		),
	},
};

export const WithCustomFallbackColor: Story = {
	args: {
		children: (
			<>
				<AvatarImage src="/broken-image.jpg" alt="@johndoe" />
				<AvatarFallback className="bg-primary text-primary-foreground">
					JD
				</AvatarFallback>
			</>
		),
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const fallback = canvas.getByText("JD");
		await expect(fallback).toBeInTheDocument();
		await expect(fallback).toHaveClass("bg-primary", "text-primary-foreground");
	},
};
