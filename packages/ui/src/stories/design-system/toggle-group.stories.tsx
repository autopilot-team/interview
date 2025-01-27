import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import {
	AlignCenter,
	AlignJustify,
	AlignLeft,
	AlignRight,
	Bold,
	Italic,
	LayoutGrid,
	List,
	ListOrdered,
	Table2,
	Underline,
} from "lucide-react";

import { ToggleGroup, ToggleGroupItem } from "../../components/toggle-group.js";

const meta: Meta<typeof ToggleGroup> = {
	title: "Design System/ToggleGroup",
	component: ToggleGroup,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof ToggleGroup>;

export const Default: Story = {
	render: () => (
		<ToggleGroup type="single" defaultValue="center">
			<ToggleGroupItem value="left" aria-label="Align left">
				<AlignLeft className="h-4 w-4" />
			</ToggleGroupItem>
			<ToggleGroupItem value="center" aria-label="Align center">
				<AlignCenter className="h-4 w-4" />
			</ToggleGroupItem>
			<ToggleGroupItem value="right" aria-label="Align right">
				<AlignRight className="h-4 w-4" />
			</ToggleGroupItem>
			<ToggleGroupItem value="justify" aria-label="Justify">
				<AlignJustify className="h-4 w-4" />
			</ToggleGroupItem>
		</ToggleGroup>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const leftButton = canvas.getByLabelText("Align left");
		await expect(leftButton).toBeInTheDocument();

		// Test toggle interaction
		await userEvent.click(leftButton);
		await expect(leftButton).toHaveAttribute("data-state", "on");
	},
};

export const Multiple: Story = {
	render: () => (
		<ToggleGroup type="multiple" defaultValue={["bold"]}>
			<ToggleGroupItem value="bold" aria-label="Toggle bold">
				<Bold className="h-4 w-4" />
			</ToggleGroupItem>
			<ToggleGroupItem value="italic" aria-label="Toggle italic">
				<Italic className="h-4 w-4" />
			</ToggleGroupItem>
			<ToggleGroupItem value="underline" aria-label="Toggle underline">
				<Underline className="h-4 w-4" />
			</ToggleGroupItem>
		</ToggleGroup>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const boldButton = canvas.getByLabelText("Toggle bold");
		const italicButton = canvas.getByLabelText("Toggle italic");

		// Test multiple selection
		await expect(boldButton).toHaveAttribute("data-state", "on");
		await userEvent.click(italicButton);
		await expect(italicButton).toHaveAttribute("data-state", "on");
		await expect(boldButton).toHaveAttribute("data-state", "on");
	},
};

export const WithVariants: Story = {
	render: () => (
		<div className="space-y-4">
			<ToggleGroup type="single" variant="outline" defaultValue="list">
				<ToggleGroupItem value="list" aria-label="List view">
					<List className="h-4 w-4" />
				</ToggleGroupItem>
				<ToggleGroupItem value="grid" aria-label="Grid view">
					<LayoutGrid className="h-4 w-4" />
				</ToggleGroupItem>
				<ToggleGroupItem value="table" aria-label="Table view">
					<Table2 className="h-4 w-4" />
				</ToggleGroupItem>
			</ToggleGroup>

			<ToggleGroup
				type="single"
				variant="default"
				size="lg"
				defaultValue="ordered"
			>
				<ToggleGroupItem value="unordered" aria-label="Unordered list">
					<List className="h-5 w-5" />
				</ToggleGroupItem>
				<ToggleGroupItem value="ordered" aria-label="Ordered list">
					<ListOrdered className="h-5 w-5" />
				</ToggleGroupItem>
			</ToggleGroup>
		</div>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const listViewButton = canvas.getByLabelText("List view");
		const gridViewButton = canvas.getByLabelText("Grid view");

		// Test variant interaction
		await expect(listViewButton).toHaveAttribute("data-state", "on");
		await userEvent.click(gridViewButton);
		await expect(gridViewButton).toHaveAttribute("data-state", "on");
		await expect(listViewButton).toHaveAttribute("data-state", "off");
	},
};

export const Disabled: Story = {
	render: () => (
		<ToggleGroup type="single" defaultValue="center">
			<ToggleGroupItem value="left" aria-label="Align left">
				<AlignLeft className="h-4 w-4" />
			</ToggleGroupItem>
			<ToggleGroupItem value="center" aria-label="Align center">
				<AlignCenter className="h-4 w-4" />
			</ToggleGroupItem>
			<ToggleGroupItem value="right" aria-label="Align right" disabled>
				<AlignRight className="h-4 w-4" />
			</ToggleGroupItem>
		</ToggleGroup>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const disabledButton = canvas.getByLabelText("Align right");

		// Test disabled state
		await expect(disabledButton).toBeDisabled();
		await expect(disabledButton).toHaveAttribute("data-disabled");
	},
};
