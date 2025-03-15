import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import {
	Calendar,
	CreditCard,
	Github,
	Keyboard,
	LifeBuoy,
	LogOut,
	Mail,
	MessageSquare,
	Plus,
	PlusCircle,
	Settings,
	User,
	UserPlus,
	Users,
} from "lucide-react";
import {
	Command,
	CommandDialog,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
	CommandSeparator,
	CommandShortcut,
} from "../../components/command.js";

const meta: Meta<typeof Command> = {
	title: "Design System/Command",
	component: Command,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Command>;

export const Default: Story = {
	render: () => (
		<Command className="rounded-lg border shadow-md">
			<CommandInput placeholder="Type a command or search..." />
			<CommandList>
				<CommandEmpty>No results found.</CommandEmpty>
				<CommandGroup heading="Suggestions">
					<CommandItem>
						<Calendar className="mr-2" />
						<span>Calendar</span>
					</CommandItem>
					<CommandItem>
						<MessageSquare className="mr-2" />
						<span>Messages</span>
						<CommandShortcut>⌘M</CommandShortcut>
					</CommandItem>
					<CommandItem>
						<Settings className="mr-2" />
						<span>Settings</span>
						<CommandShortcut>⌘S</CommandShortcut>
					</CommandItem>
				</CommandGroup>
			</CommandList>
		</Command>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const input = canvas.getByPlaceholderText("Type a command or search...");
		await expect(input).toBeInTheDocument();

		// Test search functionality
		await userEvent.type(input, "set");
		const settingsItem = canvas.getByText("Settings");
		await expect(settingsItem).toBeInTheDocument();
		await expect(canvas.queryByText("Calendar")).not.toBeInTheDocument();
	},
};

export const WithGroups: Story = {
	render: () => (
		<Command className="rounded-lg border shadow-md">
			<CommandInput placeholder="Type a command or search..." />
			<CommandList>
				<CommandEmpty>No results found.</CommandEmpty>
				<CommandGroup heading="Team">
					<CommandItem>
						<Users className="mr-2" />
						<span>View Team</span>
					</CommandItem>
					<CommandItem>
						<UserPlus className="mr-2" />
						<span>Add Member</span>
					</CommandItem>
				</CommandGroup>
				<CommandSeparator />
				<CommandGroup heading="Settings">
					<CommandItem>
						<User className="mr-2" />
						<span>Profile</span>
						<CommandShortcut>⌘P</CommandShortcut>
					</CommandItem>
					<CommandItem>
						<CreditCard className="mr-2" />
						<span>Billing</span>
						<CommandShortcut>⌘B</CommandShortcut>
					</CommandItem>
					<CommandItem>
						<Settings className="mr-2" />
						<span>Settings</span>
						<CommandShortcut>⌘S</CommandShortcut>
					</CommandItem>
				</CommandGroup>
			</CommandList>
		</Command>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Verify groups are present
		await expect(
			canvas.getByText("Team", { selector: "[cmdk-group-heading]" }),
		).toBeInTheDocument();
		await expect(
			canvas.getByText("Settings", { selector: "[cmdk-group-heading]" }),
		).toBeInTheDocument();

		// Test search across groups
		const input = canvas.getByPlaceholderText("Type a command or search...");
		await userEvent.type(input, "add");

		// Wait for filtering to take effect and verify results
		await expect(canvas.getByText("Add Member")).toBeVisible();

		// Wait a bit for the filtering animation
		await new Promise((resolve) => setTimeout(resolve, 100));

		// Verify filtering works by checking that "Add Member" is visible and other items are removed
		await expect(canvas.getByText("Add Member")).toBeVisible();
		await expect(canvas.queryByText("View Team")).toBeNull();
		await expect(canvas.queryByText("Billing")).toBeNull();
	},
};

export const WithDialog: Story = {
	render: () => (
		<CommandDialog>
			<CommandInput placeholder="Type a command or search..." />
			<CommandList>
				<CommandEmpty>No results found.</CommandEmpty>
				<CommandGroup heading="General">
					<CommandItem>
						<Plus className="mr-2" />
						<span>New Project</span>
						<CommandShortcut>⌘N</CommandShortcut>
					</CommandItem>
					<CommandItem>
						<PlusCircle className="mr-2" />
						<span>Create Team</span>
					</CommandItem>
				</CommandGroup>
				<CommandSeparator />
				<CommandGroup heading="Help">
					<CommandItem>
						<Keyboard className="mr-2" />
						<span>Keyboard Shortcuts</span>
					</CommandItem>
					<CommandItem>
						<LifeBuoy className="mr-2" />
						<span>Support</span>
					</CommandItem>
				</CommandGroup>
			</CommandList>
		</CommandDialog>
	),
};

export const WithCustomEmpty: Story = {
	render: () => (
		<Command className="rounded-lg border shadow-md">
			<CommandInput placeholder="Search..." />
			<CommandList>
				<CommandEmpty className="p-6 text-center">
					<div className="space-y-2">
						<MessageSquare className="mx-auto h-6 w-6 text-muted-foreground" />
						<p className="text-sm text-muted-foreground">
							No results found. Try a different search term.
						</p>
					</div>
				</CommandEmpty>
				<CommandGroup heading="Quick Links">
					<CommandItem>
						<Mail className="mr-2" />
						<span>Email</span>
					</CommandItem>
					<CommandItem>
						<Github className="mr-2" />
						<span>GitHub</span>
					</CommandItem>
				</CommandGroup>
			</CommandList>
		</Command>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const input = canvas.getByPlaceholderText("Search...");

		// Test empty state
		await userEvent.type(input, "nonexistent");
		const emptyMessage = canvas.getByText(
			"No results found. Try a different search term.",
		);
		await expect(emptyMessage).toBeVisible();
	},
};

export const WithFooter: Story = {
	render: () => (
		<Command className="rounded-lg border shadow-md">
			<CommandInput placeholder="Type a command or search..." />
			<CommandList>
				<CommandEmpty>No results found.</CommandEmpty>
				<CommandGroup heading="Actions">
					<CommandItem>
						<Settings className="mr-2" />
						<span>Settings</span>
					</CommandItem>
					<CommandItem>
						<User className="mr-2" />
						<span>Profile</span>
					</CommandItem>
				</CommandGroup>
				<CommandSeparator />
				<CommandGroup>
					<CommandItem className="text-red-600">
						<LogOut className="mr-2" />
						<span>Sign out</span>
						<CommandShortcut>⇧⌘Q</CommandShortcut>
					</CommandItem>
				</CommandGroup>
			</CommandList>
		</Command>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const signOutText = canvas.getByText("Sign out");
		const signOutItem = signOutText.closest("[cmdk-item]");
		await expect(signOutItem).toHaveClass("text-red-600");

		// Test shortcut is visible
		const shortcut = canvas.getByText("⇧⌘Q");
		await expect(shortcut).toBeVisible();
	},
};
