import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { Button } from "../../components/button.js";
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from "../../components/tooltip.js";

const meta: Meta<typeof Tooltip> = {
	title: "Design System/Tooltip",
	component: Tooltip,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
	decorators: [
		(Story) => (
			<TooltipProvider>
				<Story />
			</TooltipProvider>
		),
	],
};

export default meta;
type Story = StoryObj<typeof Tooltip>;

export const Default: Story = {
	render: () => (
		<Tooltip>
			<TooltipTrigger asChild>
				<Button variant="outline">Hover me</Button>
			</TooltipTrigger>
			<TooltipContent>
				<p>Add to library</p>
			</TooltipContent>
		</Tooltip>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const button = canvas.getByRole("button");
		await expect(button).toBeInTheDocument();
		await userEvent.hover(button);

		// Wait for tooltip to appear in the document body since it's rendered in a portal
		const tooltip = await within(document.body).findByRole("tooltip");
		await expect(tooltip).toBeInTheDocument();

		// Wait for animation to complete
		await new Promise((resolve) => setTimeout(resolve, 300));
		await expect(tooltip).toBeVisible();
		await expect(tooltip).toHaveTextContent("Add to library");
	},
};

export const WithIcon: Story = {
	render: () => (
		<Tooltip>
			<TooltipTrigger asChild>
				<Button size="icon" variant="ghost">
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
					>
						<title>Help</title>
						<circle cx="12" cy="12" r="10" />
						<path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3" />
						<path d="M12 17h.01" />
					</svg>
					<span className="sr-only">Help</span>
				</Button>
			</TooltipTrigger>
			<TooltipContent>
				<p>Click for more information</p>
			</TooltipContent>
		</Tooltip>
	),
};

export const WithCustomStyling: Story = {
	render: () => (
		<Tooltip>
			<TooltipTrigger asChild>
				<Button variant="secondary">Custom Tooltip</Button>
			</TooltipTrigger>
			<TooltipContent className="bg-secondary text-secondary-foreground">
				<p>Custom styled tooltip</p>
			</TooltipContent>
		</Tooltip>
	),
};

export const WithDelay: Story = {
	render: () => (
		<Tooltip delayDuration={700}>
			<TooltipTrigger asChild>
				<Button variant="outline">Delayed Tooltip</Button>
			</TooltipTrigger>
			<TooltipContent>
				<p>Tooltip with 700ms delay</p>
			</TooltipContent>
		</Tooltip>
	),
};

export const WithMultipleLines: Story = {
	render: () => (
		<Tooltip>
			<TooltipTrigger asChild>
				<Button variant="outline">Detailed Info</Button>
			</TooltipTrigger>
			<TooltipContent className="max-w-[200px] text-center">
				<p>
					This is a longer tooltip with multiple lines of text to demonstrate
					how longer content is displayed.
				</p>
			</TooltipContent>
		</Tooltip>
	),
};
