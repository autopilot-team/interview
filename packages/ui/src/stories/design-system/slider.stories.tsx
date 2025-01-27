import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import * as React from "react";
import { Slider } from "../../components/slider.js";

const meta: Meta<typeof Slider> = {
	title: "Design System/Slider",
	component: Slider,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Slider>;

export const Default: Story = {
	args: {
		defaultValue: [50],
		max: 100,
		step: 1,
		className: "w-[60vw] max-w-md",
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const slider = canvas.getByRole("slider");
		await expect(slider).toBeInTheDocument();
		await expect(slider).toHaveValue(50);
	},
};

export const Range: Story = {
	render: () => (
		<Slider
			defaultValue={[25, 75]}
			max={100}
			step={1}
			className="w-[60vw] max-w-md"
		/>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Wait for the slider to be mounted
		await new Promise((resolve) => setTimeout(resolve, 200));

		// Get both thumbs of the range slider
		const thumbs = canvas.getAllByRole("slider");
		expect(thumbs).toHaveLength(2);

		// Check first thumb (lower value)
		await expect(thumbs[0]).toHaveAttribute("aria-valuenow", "25");
		await expect(thumbs[0]).toBeVisible();

		// Check second thumb (upper value)
		await expect(thumbs[1]).toHaveAttribute("aria-valuenow", "75");
		await expect(thumbs[1]).toBeVisible();
	},
};

export const Steps: Story = {
	args: {
		defaultValue: [50],
		max: 100,
		step: 10,
		className: "w-[60vw] max-w-md",
	},
};

export const Disabled: Story = {
	args: {
		defaultValue: [50],
		max: 100,
		step: 1,
		disabled: true,
		className: "w-[60vw] max-w-md",
	},
};

export const WithLabels: Story = {
	render: () => (
		<div className="w-[60vw] max-w-md space-y-2">
			<div className="flex justify-between text-sm text-muted-foreground">
				<span>0%</span>
				<span>50%</span>
				<span>100%</span>
			</div>
			<Slider defaultValue={[50]} max={100} step={1} />
		</div>
	),
};

export const WithValue: Story = {
	render: () => {
		const [value, setValue] = React.useState(50);
		return (
			<div className="w-[60vw] max-w-md space-y-2">
				<div className="mb-4 flex items-center justify-center text-center text-sm text-muted-foreground">
					Value: {value}%
				</div>
				<Slider
					value={[value]}
					onValueChange={(values) => setValue(values[0] ?? 50)}
					max={100}
					step={1}
				/>
			</div>
		);
	},
};
