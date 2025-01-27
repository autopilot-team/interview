import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import { Button } from "../../components/button.js";
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "../../components/card.js";

const meta = {
	title: "Design System/Card",
	component: Card,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
} satisfies Meta<typeof Card>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Simple: Story = {
	args: {
		className: "w-[350px]",
		children: (
			<>
				<CardHeader>
					<CardTitle>Card Title</CardTitle>
					<CardDescription>Card Description</CardDescription>
				</CardHeader>
				<CardContent>
					<p>Card Content</p>
				</CardContent>
			</>
		),
	},
};

export const WithFooter: Story = {
	args: {
		className: "w-[350px]",
		children: (
			<>
				<CardHeader>
					<CardTitle>Newsletter</CardTitle>
					<CardDescription>Get updates on our latest features.</CardDescription>
				</CardHeader>
				<CardContent>
					<p>Subscribe to our newsletter to get the latest updates.</p>
				</CardContent>
				<CardFooter className="flex justify-between">
					<Button variant="outline">Cancel</Button>
					<Button>Subscribe</Button>
				</CardFooter>
			</>
		),
	},
};

export const Notification: Story = {
	args: {
		className: "w-[350px]",
		children: (
			<>
				<CardHeader>
					<CardTitle>Notifications</CardTitle>
					<CardDescription>You have 3 unread messages</CardDescription>
				</CardHeader>
				<CardContent>
					<div className="space-y-4">
						<div className="flex items-center space-x-4">
							<div className="w-2 h-2 bg-blue-500 rounded-full" />
							<div>
								<p className="text-sm font-medium">New message from John</p>
								<p className="text-sm text-gray-500">2 minutes ago</p>
							</div>
						</div>
						<div className="flex items-center space-x-4">
							<div className="w-2 h-2 bg-blue-500 rounded-full" />
							<div>
								<p className="text-sm font-medium">Meeting reminder</p>
								<p className="text-sm text-gray-500">1 hour ago</p>
							</div>
						</div>
					</div>
				</CardContent>
			</>
		),
	},
};

export const CardTest: Story = {
	args: {
		className: "w-[350px]",
		children: (
			<>
				<CardHeader>
					<CardTitle>Test Card</CardTitle>
					<CardDescription>Testing card components</CardDescription>
				</CardHeader>
				<CardContent>
					<p>Content for testing</p>
				</CardContent>
				<CardFooter>
					<Button>Action</Button>
				</CardFooter>
			</>
		),
	},
	play: async ({ canvasElement }: { canvasElement: HTMLElement }) => {
		const canvas = within(canvasElement);

		// Test card title
		const title = canvas.getByText("Test Card");
		await expect(title).toBeInTheDocument();

		// Test card description
		const description = canvas.getByText("Testing card components");
		await expect(description).toBeInTheDocument();

		// Test content
		const content = canvas.getByText("Content for testing");
		await expect(content).toBeInTheDocument();

		// Test button
		const button = canvas.getByRole("button", { name: "Action" });
		await expect(button).toBeInTheDocument();
	},
};
