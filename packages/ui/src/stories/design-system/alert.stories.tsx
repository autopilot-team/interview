import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import { Alert, AlertDescription, AlertTitle } from "../../components/alert.js";

const meta = {
	title: "Design System/Alert",
	component: Alert,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
	argTypes: {
		variant: {
			control: "select",
			options: ["default", "destructive"],
		},
	},
} satisfies Meta<typeof Alert>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {
		children: (
			<>
				<AlertTitle>Heads up!</AlertTitle>
				<AlertDescription>
					You can add components to your app using the cli.
				</AlertDescription>
			</>
		),
	},
};

export const Destructive: Story = {
	args: {
		variant: "destructive",
		children: (
			<>
				<AlertTitle>Error</AlertTitle>
				<AlertDescription>
					Your session has expired. Please log in again.
				</AlertDescription>
			</>
		),
	},
};

export const WithIcon: Story = {
	args: {
		children: (
			<>
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
					className="h-4 w-4"
					aria-label="Info icon"
					role="img"
				>
					<circle cx="12" cy="12" r="10" />
					<line x1="12" y1="16" x2="12" y2="12" />
					<line x1="12" y1="8" x2="12.01" y2="8" />
				</svg>
				<AlertTitle>Note</AlertTitle>
				<AlertDescription>
					This is an important message with an icon.
				</AlertDescription>
			</>
		),
	},
};

export const AlertTest: Story = {
	args: {
		variant: "destructive",
		children: (
			<>
				<AlertTitle>Test Alert</AlertTitle>
				<AlertDescription>This is a test alert message.</AlertDescription>
			</>
		),
	},
	play: async ({ canvasElement }: { canvasElement: HTMLElement }) => {
		const canvas = within(canvasElement);

		// Test alert title
		const title = canvas.getByText("Test Alert");
		await expect(title).toBeInTheDocument();

		// Test alert description
		const description = canvas.getByText("This is a test alert message.");
		await expect(description).toBeInTheDocument();

		// Test alert role
		const alert = canvas.getByRole("alert");
		await expect(alert).toBeInTheDocument();

		// Test destructive variant
		await expect(alert).toHaveClass("border-destructive/50");
	},
};
