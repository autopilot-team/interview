import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import { Skeleton } from "../../components/skeleton.js";

const meta: Meta<typeof Skeleton> = {
	title: "Design System/Skeleton",
	component: Skeleton,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Skeleton>;

export const Default: Story = {
	render: () => <Skeleton className="h-4 w-[250px]" />,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const skeleton = canvas.getByTestId("skeleton");
		await expect(skeleton).toBeInTheDocument();
		await expect(skeleton).toHaveClass("animate-pulse");
	},
};

export const Card: Story = {
	render: () => (
		<div className="flex flex-col space-y-3">
			<Skeleton className="h-[125px] w-[250px] rounded-xl" />
			<div className="space-y-2">
				<Skeleton className="h-4 w-[250px]" />
				<Skeleton className="h-4 w-[200px]" />
			</div>
		</div>
	),
};

export const Profile: Story = {
	render: () => (
		<div className="flex items-center space-x-4">
			<Skeleton className="h-12 w-12 rounded-full" />
			<div className="space-y-2">
				<Skeleton className="h-4 w-[250px]" />
				<Skeleton className="h-4 w-[200px]" />
			</div>
		</div>
	),
};

export const Table: Story = {
	render: () => (
		<div className="space-y-4">
			<div className="space-y-2">
				<Skeleton className="h-4 w-[250px]" />
				<Skeleton className="h-4 w-[200px]" />
				<Skeleton className="h-4 w-[225px]" />
			</div>
			<div className="space-y-2">
				<Skeleton className="h-4 w-[250px]" />
				<Skeleton className="h-4 w-[200px]" />
				<Skeleton className="h-4 w-[225px]" />
			</div>
			<div className="space-y-2">
				<Skeleton className="h-4 w-[250px]" />
				<Skeleton className="h-4 w-[200px]" />
				<Skeleton className="h-4 w-[225px]" />
			</div>
		</div>
	),
};

export const ComplexCard: Story = {
	render: () => (
		<div className="flex flex-col space-y-3">
			<Skeleton className="h-[200px] w-[400px] rounded-xl" />
			<div className="space-y-2">
				<Skeleton className="h-4 w-[400px]" />
				<Skeleton className="h-4 w-[300px]" />
				<div className="flex space-x-2 pt-4">
					<Skeleton className="h-8 w-20 rounded-full" />
					<Skeleton className="h-8 w-20 rounded-full" />
					<Skeleton className="h-8 w-20 rounded-full" />
				</div>
			</div>
		</div>
	),
};

export const WithCustomStyling: Story = {
	render: () => (
		<div className="space-y-4">
			<Skeleton className="h-4 w-[250px] bg-primary/5" />
			<Skeleton className="h-4 w-[250px] bg-secondary/10" />
			<Skeleton className="h-4 w-[250px] bg-destructive/5" />
			<Skeleton className="h-4 w-[250px] animate-[pulse_2s_ease-in-out_infinite]" />
		</div>
	),
};
