import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import {
	HoverCard,
	HoverCardContent,
	HoverCardTrigger,
} from "../../components/hover-card.js";

const meta: Meta<typeof HoverCard> = {
	title: "Design System/HoverCard",
	component: HoverCard,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof HoverCard>;

export const Default: Story = {
	render: () => (
		<HoverCard>
			<HoverCardTrigger className="cursor-pointer underline">
				Hover over me
			</HoverCardTrigger>
			<HoverCardContent>
				<div className="space-y-2">
					<h4 className="text-sm font-semibold">Information</h4>
					<p className="text-sm">
						This is a hover card. It appears when you hover over the trigger.
					</p>
				</div>
			</HoverCardContent>
		</HoverCard>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByText("Hover over me");
		await expect(trigger).toBeInTheDocument();
		await userEvent.hover(trigger);

		// Wait for hover card to appear in the document body since it's rendered in a portal
		const content = await within(document.body).findByText(
			/This is a hover card/,
		);
		await expect(content).toBeInTheDocument();

		// Wait for animation to complete
		await new Promise((resolve) => setTimeout(resolve, 300));
		await expect(content).toBeVisible();
	},
};

export const WithUserProfile: Story = {
	render: () => (
		<HoverCard>
			<HoverCardTrigger className="cursor-pointer underline">
				@username
			</HoverCardTrigger>
			<HoverCardContent className="w-80">
				<div className="flex justify-between space-x-4">
					<div className="space-y-1">
						<h4 className="text-sm font-semibold">@username</h4>
						<p className="text-sm">
							Software developer and open source contributor.
						</p>
						<div className="flex items-center pt-2">
							<span className="text-xs text-muted-foreground">
								Joined December 2023
							</span>
						</div>
					</div>
				</div>
			</HoverCardContent>
		</HoverCard>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByText("@username");
		await expect(trigger).toBeInTheDocument();
		await userEvent.hover(trigger);

		// Wait for hover card to appear in the document body since it's rendered in a portal
		const content = await within(document.body).findByText(
			/Software developer/,
		);
		await expect(content).toBeInTheDocument();

		// Wait for animation to complete
		await new Promise((resolve) => setTimeout(resolve, 300));
		await expect(content).toBeVisible();
	},
};

export const WithCustomAlignment: Story = {
	render: () => (
		<div className="flex h-40 w-full items-center justify-center">
			<HoverCard>
				<HoverCardTrigger className="cursor-pointer underline">
					Right aligned
				</HoverCardTrigger>
				<HoverCardContent align="end" className="w-[200px]">
					<p className="text-sm">
						This hover card is aligned to the end of the trigger.
					</p>
				</HoverCardContent>
			</HoverCard>
		</div>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByText("Right aligned");
		await expect(trigger).toBeInTheDocument();
		await userEvent.hover(trigger);

		// Wait for hover card to appear in the document body since it's rendered in a portal
		const content = await within(document.body).findByText(
			/This hover card is aligned/,
		);
		await expect(content).toBeInTheDocument();

		// Wait for animation to complete
		await new Promise((resolve) => setTimeout(resolve, 300));
		await expect(content).toBeVisible();
		await expect(content.parentElement).toHaveAttribute("data-align", "end");
	},
};

export const WithCustomOffset: Story = {
	render: () => (
		<HoverCard>
			<HoverCardTrigger className="cursor-pointer underline">
				Offset hover card
			</HoverCardTrigger>
			<HoverCardContent sideOffset={10}>
				<p className="text-sm">This hover card has a custom offset of 10px.</p>
			</HoverCardContent>
		</HoverCard>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByText("Offset hover card");
		await expect(trigger).toBeInTheDocument();
		await userEvent.hover(trigger);

		// Wait for hover card to appear in the document body since it's rendered in a portal
		const content = await within(document.body).findByText(/custom offset/);
		await expect(content).toBeInTheDocument();

		// Wait for animation to complete
		await new Promise((resolve) => setTimeout(resolve, 300));
		await expect(content).toBeVisible();
	},
};
