import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { Label } from "../../components/label.js";
import { RadioGroup, RadioGroupItem } from "../../components/radio-group.js";

const meta: Meta<typeof RadioGroup> = {
	title: "Design System/RadioGroup",
	component: RadioGroup,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof RadioGroup>;

export const Default: Story = {
	render: () => (
		<RadioGroup defaultValue="option-one">
			<div className="flex items-center space-x-2">
				<RadioGroupItem value="option-one" id="option-one" />
				<Label htmlFor="option-one">Option One</Label>
			</div>
			<div className="flex items-center space-x-2">
				<RadioGroupItem value="option-two" id="option-two" />
				<Label htmlFor="option-two">Option Two</Label>
			</div>
			<div className="flex items-center space-x-2">
				<RadioGroupItem value="option-three" id="option-three" />
				<Label htmlFor="option-three">Option Three</Label>
			</div>
		</RadioGroup>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const radio = canvas.getByLabelText("Option One");
		await expect(radio).toBeInTheDocument();
		await expect(radio).toBeChecked();
		await userEvent.click(canvas.getByLabelText("Option Two"));
		await expect(canvas.getByLabelText("Option Two")).toBeChecked();
	},
};

export const WithDescription: Story = {
	render: () => (
		<RadioGroup defaultValue="comfortable">
			<div className="flex items-start space-x-2">
				<RadioGroupItem value="default" id="default" className="mt-1" />
				<div>
					<Label htmlFor="default">Default</Label>
					<p className="text-sm text-muted-foreground">
						System default spacing for controls and content.
					</p>
				</div>
			</div>
			<div className="flex items-start space-x-2">
				<RadioGroupItem value="comfortable" id="comfortable" className="mt-1" />
				<div>
					<Label htmlFor="comfortable">Comfortable</Label>
					<p className="text-sm text-muted-foreground">
						Additional spacing for better readability.
					</p>
				</div>
			</div>
			<div className="flex items-start space-x-2">
				<RadioGroupItem value="compact" id="compact" className="mt-1" />
				<div>
					<Label htmlFor="compact">Compact</Label>
					<p className="text-sm text-muted-foreground">
						Reduced spacing to show more content.
					</p>
				</div>
			</div>
		</RadioGroup>
	),
};

export const Disabled: Story = {
	render: () => (
		<RadioGroup defaultValue="option-two" disabled>
			<div className="flex items-center space-x-2">
				<RadioGroupItem value="option-one" id="disabled-one" />
				<Label htmlFor="disabled-one">Option One</Label>
			</div>
			<div className="flex items-center space-x-2">
				<RadioGroupItem value="option-two" id="disabled-two" />
				<Label htmlFor="disabled-two">Option Two</Label>
			</div>
		</RadioGroup>
	),
};

export const WithCustomStyles: Story = {
	render: () => (
		<RadioGroup
			defaultValue="option-one"
			className="bg-secondary/20 p-4 rounded-lg"
		>
			<div className="flex items-center space-x-2">
				<RadioGroupItem
					value="option-one"
					id="styled-one"
					className="border-secondary text-secondary"
				/>
				<Label htmlFor="styled-one" className="font-semibold">
					Option One
				</Label>
			</div>
			<div className="flex items-center space-x-2">
				<RadioGroupItem
					value="option-two"
					id="styled-two"
					className="border-secondary text-secondary"
				/>
				<Label htmlFor="styled-two" className="font-semibold">
					Option Two
				</Label>
			</div>
		</RadioGroup>
	),
};

export const Inline: Story = {
	render: () => (
		<RadioGroup defaultValue="option-one" className="flex space-x-4">
			<div className="flex items-center space-x-2">
				<RadioGroupItem value="option-one" id="inline-one" />
				<Label htmlFor="inline-one">Option One</Label>
			</div>
			<div className="flex items-center space-x-2">
				<RadioGroupItem value="option-two" id="inline-two" />
				<Label htmlFor="inline-two">Option Two</Label>
			</div>
			<div className="flex items-center space-x-2">
				<RadioGroupItem value="option-three" id="inline-three" />
				<Label htmlFor="inline-three">Option Three</Label>
			</div>
		</RadioGroup>
	),
};
