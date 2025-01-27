import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { Label } from "../../components/label.js";
import { Switch } from "../../components/switch.js";

const meta: Meta<typeof Switch> = {
	title: "Design System/Switch",
	component: Switch,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Switch>;

export const Default: Story = {
	render: () => (
		<div className="flex items-center space-x-2">
			<Switch id="airplane-mode" />
			<Label htmlFor="airplane-mode">Airplane Mode</Label>
		</div>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const switchElement = canvas.getByRole("switch");
		await expect(switchElement).toBeInTheDocument();
		await expect(switchElement).not.toBeChecked();
		await userEvent.click(switchElement);
		await expect(switchElement).toBeChecked();
	},
};

export const Checked: Story = {
	render: () => (
		<div className="flex items-center space-x-2">
			<Switch id="notifications" defaultChecked />
			<Label htmlFor="notifications">Enable Notifications</Label>
		</div>
	),
};

export const WithDescription: Story = {
	render: () => (
		<div className="flex flex-row items-center justify-between rounded-lg border p-4">
			<div className="space-y-0.5">
				<Label htmlFor="dark-mode">Dark Mode</Label>
				<p className="text-sm text-muted-foreground">
					Enable dark mode for a better viewing experience at night
				</p>
			</div>
			<Switch id="dark-mode" />
		</div>
	),
};

export const Disabled: Story = {
	render: () => (
		<div className="space-y-4">
			<div className="flex items-center space-x-2">
				<Switch id="disabled-unchecked" disabled />
				<Label htmlFor="disabled-unchecked" className="text-muted-foreground">
					Disabled
				</Label>
			</div>
			<div className="flex items-center space-x-2">
				<Switch id="disabled-checked" disabled defaultChecked />
				<Label htmlFor="disabled-checked" className="text-muted-foreground">
					Disabled (checked)
				</Label>
			</div>
		</div>
	),
};

export const WithForm: Story = {
	render: () => (
		<form className="w-full max-w-sm space-y-4 rounded-lg border p-4">
			<div className="space-y-2">
				<h4 className="font-medium">Email Preferences</h4>
				<p className="text-sm text-muted-foreground">
					Manage your email notification preferences
				</p>
			</div>
			<div className="space-y-4">
				<div className="flex items-center justify-between">
					<Label htmlFor="marketing">Marketing emails</Label>
					<Switch id="marketing" defaultChecked />
				</div>
				<div className="flex items-center justify-between">
					<Label htmlFor="social">Social notifications</Label>
					<Switch id="social" defaultChecked />
				</div>
				<div className="flex items-center justify-between">
					<Label htmlFor="security">Security updates</Label>
					<Switch id="security" defaultChecked />
				</div>
			</div>
		</form>
	),
};

export const WithCustomStyling: Story = {
	render: () => (
		<div className="flex items-center space-x-2">
			<Switch
				id="custom"
				className="data-[state=checked]:bg-green-500 data-[state=unchecked]:bg-red-500"
			/>
			<Label htmlFor="custom">Custom Colors</Label>
		</div>
	),
};
