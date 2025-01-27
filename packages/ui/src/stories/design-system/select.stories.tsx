import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectLabel,
	SelectTrigger,
	SelectValue,
} from "../../components/select.js";

const meta = {
	title: "Design System/Select",
	component: Select,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
} satisfies Meta<typeof Select>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<Select>
			<SelectTrigger className="w-[180px]">
				<SelectValue placeholder="Select a fruit" />
			</SelectTrigger>
			<SelectContent>
				<SelectGroup>
					<SelectItem value="apple">Apple</SelectItem>
					<SelectItem value="banana">Banana</SelectItem>
					<SelectItem value="orange">Orange</SelectItem>
				</SelectGroup>
			</SelectContent>
		</Select>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByRole("combobox");
		await expect(trigger).toBeInTheDocument();
		await trigger.click();

		// Wait for select content to be mounted and visible in the document body
		const content = await within(document.body).findByRole("listbox");
		await expect(content).toBeInTheDocument();

		// Wait for animation to complete
		await new Promise((resolve) => setTimeout(resolve, 300));

		// Find and verify the option
		const option = await within(content).findByRole("option", {
			name: "Apple",
		});
		await expect(option).toBeVisible();
	},
};

export const WithLabel: Story = {
	render: () => (
		<div className="grid w-full max-w-sm gap-1.5">
			<label htmlFor="fruit" className="text-sm font-medium leading-none">
				Favorite fruit
			</label>
			<Select>
				<SelectTrigger id="fruit" className="w-[180px]">
					<SelectValue placeholder="Select a fruit" />
				</SelectTrigger>
				<SelectContent>
					<SelectGroup>
						<SelectItem value="apple">Apple</SelectItem>
						<SelectItem value="banana">Banana</SelectItem>
						<SelectItem value="orange">Orange</SelectItem>
					</SelectGroup>
				</SelectContent>
			</Select>
		</div>
	),
};

export const WithGroups: Story = {
	render: () => (
		<Select>
			<SelectTrigger className="w-[180px]">
				<SelectValue placeholder="Select a food" />
			</SelectTrigger>
			<SelectContent>
				<SelectGroup>
					<SelectLabel>Fruits</SelectLabel>
					<SelectItem value="apple">Apple</SelectItem>
					<SelectItem value="banana">Banana</SelectItem>
					<SelectItem value="orange">Orange</SelectItem>
				</SelectGroup>
				<SelectGroup>
					<SelectLabel>Vegetables</SelectLabel>
					<SelectItem value="carrot">Carrot</SelectItem>
					<SelectItem value="potato">Potato</SelectItem>
					<SelectItem value="tomato">Tomato</SelectItem>
				</SelectGroup>
			</SelectContent>
		</Select>
	),
};

export const Disabled: Story = {
	render: () => (
		<Select disabled>
			<SelectTrigger className="w-[180px]">
				<SelectValue placeholder="Select a fruit" />
			</SelectTrigger>
			<SelectContent>
				<SelectGroup>
					<SelectItem value="apple">Apple</SelectItem>
					<SelectItem value="banana">Banana</SelectItem>
					<SelectItem value="orange">Orange</SelectItem>
				</SelectGroup>
			</SelectContent>
		</Select>
	),
};

export const WithDisabledItems: Story = {
	render: () => (
		<Select>
			<SelectTrigger className="w-[180px]">
				<SelectValue placeholder="Select a fruit" />
			</SelectTrigger>
			<SelectContent>
				<SelectGroup>
					<SelectItem value="apple">Apple</SelectItem>
					<SelectItem value="banana" disabled>
						Banana
					</SelectItem>
					<SelectItem value="orange">Orange</SelectItem>
				</SelectGroup>
			</SelectContent>
		</Select>
	),
};
