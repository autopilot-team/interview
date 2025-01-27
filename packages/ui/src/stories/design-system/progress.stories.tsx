import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import * as React from "react";
import { Progress } from "../../components/progress.js";

const meta: Meta<typeof Progress> = {
	title: "Design System/Progress",
	component: Progress,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Progress>;

export const Default: Story = {
	args: {
		value: 40,
		className: "w-[60vw] max-w-md",
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const progress = canvas.getByRole("progressbar");
		await expect(progress).toBeInTheDocument();

		// Wait for styles to be applied
		await new Promise((resolve) => setTimeout(resolve, 100));

		// Check the progress indicator
		const indicator = progress.querySelector(
			'[class*="bg-primary"]',
		) as HTMLDivElement;
		expect(indicator).toBeInTheDocument();

		if (!indicator) throw new Error("Progress indicator not found");

		// Check the inline transform style
		expect(indicator.style.transform).toBe("translateX(-60%)");

		// Verify data attributes
		await expect(progress).toHaveAttribute("data-state", "indeterminate");
		await expect(progress).toHaveAttribute("data-max", "100");

		// Verify ARIA attributes
		await expect(progress).toHaveAttribute("aria-valuemin", "0");
		await expect(progress).toHaveAttribute("aria-valuemax", "100");
	},
};

export const WithLabel: Story = {
	render: () => {
		const [progress, setProgress] = React.useState(40);

		return (
			<div className="w-[60vw] max-w-md space-y-2">
				<Progress value={progress} />
				<div className="text-sm text-muted-foreground text-center">
					{progress}% complete
				</div>
			</div>
		);
	},
};

const AnimatedProgress = () => {
	const [progress, setProgress] = React.useState(13);

	React.useEffect(() => {
		const timer = setInterval(() => {
			setProgress((prevProgress) =>
				prevProgress >= 100 ? 0 : prevProgress + 1,
			);
		}, 100);
		return () => clearInterval(timer);
	}, []);

	return (
		<div className="w-[60vw] max-w-md space-y-2">
			<Progress value={progress} />
			<div className="text-sm text-muted-foreground text-center">
				Loading... {progress}%
			</div>
		</div>
	);
};

export const Animated: Story = {
	render: () => <AnimatedProgress />,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const progress = canvas.getByRole("progressbar");
		await expect(progress).toBeInTheDocument();

		// Wait for the initial value to be set
		const indicator = progress.querySelector(
			'[class*="bg-primary"]',
		) as HTMLDivElement;
		expect(indicator).toBeInTheDocument();

		if (!indicator) throw new Error("Progress indicator not found");

		// Check initial transform
		expect(indicator.style.transform).toBe("translateX(-87%)"); // 100 - 13 = 87

		// Verify initial attributes
		await expect(progress).toHaveAttribute("data-state", "indeterminate");
		await expect(progress).toHaveAttribute("data-max", "100");
		await expect(progress).toHaveAttribute("aria-valuemin", "0");
		await expect(progress).toHaveAttribute("aria-valuemax", "100");

		// Wait a bit and verify the value increases
		await new Promise((resolve) => setTimeout(resolve, 200));

		// Get the new transform value and verify it changed
		const computedStyle = window.getComputedStyle(indicator);
		expect(computedStyle.transform).not.toBe("translateX(-87%)");

		// Verify loading text is present and updates
		await expect(canvas.getByText(/Loading\.\.\. \d+%/)).toBeInTheDocument();
	},
};

export const CustomColors: Story = {
	render: () => (
		<div className="w-[60vw] max-w-md space-y-4">
			<Progress value={40} className="bg-secondary/20" />
			<Progress value={60} className="bg-secondary/20 [&>div]:bg-green-500" />
			<Progress value={80} className="bg-secondary/20 [&>div]:bg-red-500" />
		</div>
	),
};

export const CustomSizes: Story = {
	render: () => (
		<div className="w-[60vw] max-w-md space-y-4">
			<Progress value={40} className="h-1" />
			<Progress value={40} className="h-2" />
			<Progress value={40} className="h-3" />
			<Progress value={40} className="h-4" />
		</div>
	),
};

export const Indeterminate: Story = {
	render: () => (
		<div className="w-[60vw] max-w-md space-y-2">
			<Progress className="[&>div]:animate-[progress-loading_1s_ease-in-out_infinite]" />
			<div className="text-sm text-muted-foreground text-center">
				Processing...
			</div>
		</div>
	),
};
