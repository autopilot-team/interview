import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import {
	Accordion,
	AccordionContent,
	AccordionItem,
	AccordionTrigger,
} from "../../components/accordion.js";

const meta: Meta<typeof Accordion> = {
	title: "Design System/Accordion",
	component: Accordion,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Accordion>;

export const Default: Story = {
	render: () => (
		<Accordion type="single" collapsible className="w-[400px]">
			<AccordionItem value="item-1">
				<AccordionTrigger>Is it accessible?</AccordionTrigger>
				<AccordionContent>
					Yes. It adheres to the WAI-ARIA design pattern.
				</AccordionContent>
			</AccordionItem>
		</Accordion>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByRole("button", { name: "Is it accessible?" });
		await expect(trigger).toBeInTheDocument();

		// Test accordion interaction
		await userEvent.click(trigger);
		const content = canvas.getByText(
			"Yes. It adheres to the WAI-ARIA design pattern.",
		);
		await expect(content).toBeVisible();
	},
};

export const Multiple: Story = {
	render: () => (
		<Accordion type="multiple" className="w-[400px]">
			<AccordionItem value="item-1">
				<AccordionTrigger>What is your refund policy?</AccordionTrigger>
				<AccordionContent>
					If you're unhappy with your purchase for any reason, email us within
					90 days and we'll refund you in full, no questions asked.
				</AccordionContent>
			</AccordionItem>
			<AccordionItem value="item-2">
				<AccordionTrigger>Do you offer technical support?</AccordionTrigger>
				<AccordionContent>
					Yes, we offer email and phone support 24 hours a day, 7 days a week.
				</AccordionContent>
			</AccordionItem>
		</Accordion>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const refundTrigger = canvas.getByRole("button", {
			name: "What is your refund policy?",
		});
		const supportTrigger = canvas.getByRole("button", {
			name: "Do you offer technical support?",
		});

		// Test multiple items can be open simultaneously
		await userEvent.click(refundTrigger);
		await userEvent.click(supportTrigger);

		const refundPolicy = canvas.getByText(
			/If you're unhappy with your purchase/,
		);
		const support = canvas.getByText(/Yes, we offer email and phone support/);

		await expect(refundPolicy).toBeVisible();
		await expect(support).toBeVisible();
	},
};

export const WithCustomStyles: Story = {
	render: () => (
		<Accordion type="single" collapsible className="w-[400px]">
			<AccordionItem value="item-1" className="border-b-primary">
				<AccordionTrigger className="text-primary hover:text-primary/80">
					Can I customize the styling?
				</AccordionTrigger>
				<AccordionContent className="text-muted-foreground">
					Yes, the component accepts custom className props for styling.
				</AccordionContent>
			</AccordionItem>
		</Accordion>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByRole("button", {
			name: "Can I customize the styling?",
		});
		await expect(trigger).toHaveClass("text-primary");

		await userEvent.click(trigger);
		const content = canvas.getByText(/Yes, the component accepts/);
		await expect(content).toBeVisible();
		await expect(content.parentElement).toHaveClass("text-muted-foreground");
	},
};

export const FAQ: Story = {
	render: () => (
		<Accordion type="single" collapsible className="w-[600px]">
			<AccordionItem value="item-1">
				<AccordionTrigger>What payment methods do you accept?</AccordionTrigger>
				<AccordionContent>
					We accept all major credit cards, PayPal, and bank transfers. For
					business customers, we also offer invoice-based payments with net-30
					terms.
				</AccordionContent>
			</AccordionItem>
			<AccordionItem value="item-2">
				<AccordionTrigger>How long does shipping take?</AccordionTrigger>
				<AccordionContent>
					Domestic shipping typically takes 3-5 business days. International
					shipping can take 7-14 business days depending on the destination.
				</AccordionContent>
			</AccordionItem>
			<AccordionItem value="item-3">
				<AccordionTrigger>Do you ship internationally?</AccordionTrigger>
				<AccordionContent>
					Yes, we ship to over 100 countries worldwide. Shipping costs and
					delivery times vary by location.
				</AccordionContent>
			</AccordionItem>
		</Accordion>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const questions = [
			"What payment methods do you accept?",
			"How long does shipping take?",
			"Do you ship internationally?",
		];

		// Test each accordion item
		for (const question of questions) {
			const trigger = canvas.getByRole("button", { name: question });
			await userEvent.click(trigger);
			await expect(trigger).toHaveAttribute("data-state", "open");

			// Click again to close
			await userEvent.click(trigger);
			await expect(trigger).toHaveAttribute("data-state", "closed");
		}
	},
};
