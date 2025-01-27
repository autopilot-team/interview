import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import {
	ResizableHandle,
	ResizablePanel,
	ResizablePanelGroup,
} from "../../components/resizable.js";

const meta: Meta<typeof ResizablePanelGroup> = {
	title: "Design System/Resizable",
	component: ResizablePanelGroup,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof ResizablePanelGroup>;

export const Default: Story = {
	render: () => (
		<ResizablePanelGroup
			direction="horizontal"
			className="h-[400px] max-w-3xl rounded-lg border"
		>
			<ResizablePanel defaultSize={25}>
				<div className="flex h-full items-center justify-center p-6">
					<span className="font-semibold">Sidebar</span>
				</div>
			</ResizablePanel>
			<ResizableHandle />
			<ResizablePanel defaultSize={75}>
				<div className="flex h-full items-center justify-center p-6">
					<span className="font-semibold">Main Content</span>
				</div>
			</ResizablePanel>
		</ResizablePanelGroup>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const sidebar = canvas.getByText("Sidebar");
		const mainContent = canvas.getByText("Main Content");

		await expect(sidebar).toBeInTheDocument();
		await expect(mainContent).toBeInTheDocument();
	},
};

export const WithHandle: Story = {
	render: () => (
		<ResizablePanelGroup
			direction="horizontal"
			className="h-[400px] max-w-3xl rounded-lg border"
		>
			<ResizablePanel defaultSize={30}>
				<div className="flex h-full items-center justify-center p-6">
					<span className="font-semibold">Panel 1</span>
				</div>
			</ResizablePanel>
			<ResizableHandle withHandle />
			<ResizablePanel defaultSize={40}>
				<div className="flex h-full items-center justify-center p-6">
					<span className="font-semibold">Panel 2</span>
				</div>
			</ResizablePanel>
			<ResizableHandle withHandle />
			<ResizablePanel defaultSize={30}>
				<div className="flex h-full items-center justify-center p-6">
					<span className="font-semibold">Panel 3</span>
				</div>
			</ResizablePanel>
		</ResizablePanelGroup>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const handles = canvas.getAllByRole("separator");
		await expect(handles).toHaveLength(2);
	},
};

export const Vertical: Story = {
	render: () => (
		<ResizablePanelGroup
			direction="vertical"
			className="h-[500px] max-w-3xl rounded-lg border"
		>
			<ResizablePanel defaultSize={25}>
				<div className="flex h-full items-center justify-center p-6">
					<span className="font-semibold">Top Panel</span>
				</div>
			</ResizablePanel>
			<ResizableHandle withHandle />
			<ResizablePanel defaultSize={75}>
				<div className="flex h-full items-center justify-center p-6">
					<span className="font-semibold">Bottom Panel</span>
				</div>
			</ResizablePanel>
		</ResizablePanelGroup>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const topPanel = canvas.getByText("Top Panel");
		const bottomPanel = canvas.getByText("Bottom Panel");

		await expect(topPanel).toBeInTheDocument();
		await expect(bottomPanel).toBeInTheDocument();
	},
};

export const NestedPanels: Story = {
	render: () => (
		<ResizablePanelGroup
			direction="horizontal"
			className="h-[500px] max-w-4xl rounded-lg border"
		>
			<ResizablePanel defaultSize={25}>
				<div className="flex h-full items-center justify-center p-6">
					<span className="font-semibold">Navigation</span>
				</div>
			</ResizablePanel>
			<ResizableHandle withHandle />
			<ResizablePanel defaultSize={75}>
				<ResizablePanelGroup direction="vertical">
					<ResizablePanel defaultSize={70}>
						<div className="flex h-full items-center justify-center p-6">
							<span className="font-semibold">Main Content</span>
						</div>
					</ResizablePanel>
					<ResizableHandle withHandle />
					<ResizablePanel defaultSize={30}>
						<div className="flex h-full items-center justify-center p-6">
							<span className="font-semibold">Preview</span>
						</div>
					</ResizablePanel>
				</ResizablePanelGroup>
			</ResizablePanel>
		</ResizablePanelGroup>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const navigation = canvas.getByText("Navigation");
		const mainContent = canvas.getByText("Main Content");
		const preview = canvas.getByText("Preview");

		await expect(navigation).toBeInTheDocument();
		await expect(mainContent).toBeInTheDocument();
		await expect(preview).toBeInTheDocument();

		const handles = canvas.getAllByRole("separator");
		await expect(handles).toHaveLength(2);
	},
};
