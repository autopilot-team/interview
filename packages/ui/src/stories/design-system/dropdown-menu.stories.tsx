import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import {
	Cloud,
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

import { Button } from "../../components/button.js";
import {
	DropdownMenu,
	DropdownMenuCheckboxItem,
	DropdownMenuContent,
	DropdownMenuGroup,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuPortal,
	DropdownMenuRadioGroup,
	DropdownMenuRadioItem,
	DropdownMenuSeparator,
	DropdownMenuShortcut,
	DropdownMenuSub,
	DropdownMenuSubContent,
	DropdownMenuSubTrigger,
	DropdownMenuTrigger,
} from "../../components/dropdown-menu.js";

const meta: Meta<typeof DropdownMenu> = {
	title: "Design System/DropdownMenu",
	component: DropdownMenu,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof DropdownMenu>;

export const Default: Story = {
	render: () => (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button variant="outline">Open Menu</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent className="w-56">
				<DropdownMenuLabel>My Account</DropdownMenuLabel>
				<DropdownMenuSeparator />
				<DropdownMenuGroup>
					<DropdownMenuItem>
						<User className="mr-2 h-4 w-4" />
						Profile
						<DropdownMenuShortcut>⇧⌘P</DropdownMenuShortcut>
					</DropdownMenuItem>
					<DropdownMenuItem>
						<CreditCard className="mr-2 h-4 w-4" />
						Billing
						<DropdownMenuShortcut>⌘B</DropdownMenuShortcut>
					</DropdownMenuItem>
					<DropdownMenuItem>
						<Settings className="mr-2 h-4 w-4" />
						Settings
						<DropdownMenuShortcut>⌘S</DropdownMenuShortcut>
					</DropdownMenuItem>
					<DropdownMenuItem>
						<Keyboard className="mr-2 h-4 w-4" />
						Keyboard shortcuts
						<DropdownMenuShortcut>⌘K</DropdownMenuShortcut>
					</DropdownMenuItem>
				</DropdownMenuGroup>
			</DropdownMenuContent>
		</DropdownMenu>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByRole("button", { name: "Open Menu" });
		await expect(trigger).toBeInTheDocument();

		// Test menu interaction
		await userEvent.click(trigger);

		// Wait for the dropdown to be mounted in the DOM
		await new Promise((resolve) => setTimeout(resolve, 100));

		// Look for menu items in document.body since they're rendered in a portal
		const menuContent = within(document.body);
		await expect(menuContent.getByText("My Account")).toBeVisible();
		await expect(menuContent.getByText("Profile")).toBeVisible();
		await expect(menuContent.getByText("Settings")).toBeVisible();
	},
};

export const WithSubmenus: Story = {
	render: () => (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button>Team Settings</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent className="w-56">
				<DropdownMenuLabel>Team</DropdownMenuLabel>
				<DropdownMenuSeparator />
				<DropdownMenuGroup>
					<DropdownMenuItem>
						<Users className="mr-2 h-4 w-4" />
						Team Members
					</DropdownMenuItem>
					<DropdownMenuSub>
						<DropdownMenuSubTrigger>
							<UserPlus className="mr-2 h-4 w-4" />
							Invite Users
						</DropdownMenuSubTrigger>
						<DropdownMenuPortal>
							<DropdownMenuSubContent>
								<DropdownMenuItem>
									<Mail className="mr-2 h-4 w-4" />
									Email Invite
								</DropdownMenuItem>
								<DropdownMenuItem>
									<MessageSquare className="mr-2 h-4 w-4" />
									Message Invite
								</DropdownMenuItem>
								<DropdownMenuSeparator />
								<DropdownMenuItem>
									<PlusCircle className="mr-2 h-4 w-4" />
									More Options...
								</DropdownMenuItem>
							</DropdownMenuSubContent>
						</DropdownMenuPortal>
					</DropdownMenuSub>
					<DropdownMenuItem>
						<Plus className="mr-2 h-4 w-4" />
						New Team
						<DropdownMenuShortcut>⌘+T</DropdownMenuShortcut>
					</DropdownMenuItem>
				</DropdownMenuGroup>
				<DropdownMenuSeparator />
				<DropdownMenuItem>
					<Github className="mr-2 h-4 w-4" />
					GitHub
				</DropdownMenuItem>
				<DropdownMenuItem>
					<LifeBuoy className="mr-2 h-4 w-4" />
					Support
				</DropdownMenuItem>
				<DropdownMenuItem disabled>
					<Cloud className="mr-2 h-4 w-4" />
					API Reference
				</DropdownMenuItem>
				<DropdownMenuSeparator />
				<DropdownMenuItem>
					<LogOut className="mr-2 h-4 w-4" />
					Log out
					<DropdownMenuShortcut>⇧⌘Q</DropdownMenuShortcut>
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByRole("button", { name: "Team Settings" });
		await expect(trigger).toBeInTheDocument();

		// Test menu interaction
		await userEvent.click(trigger);

		// Wait for the dropdown to be mounted in the DOM
		await new Promise((resolve) => setTimeout(resolve, 200));

		// Look for menu items in document.body since they're rendered in a portal
		const menuContent = within(document.body);
		await expect(menuContent.getByText("Team")).toBeVisible();
		await expect(menuContent.getByText("Team Members")).toBeVisible();

		// Test submenu interaction
		const inviteUsersTrigger = menuContent.getByRole("menuitem", {
			name: /Invite Users/i,
		});
		await expect(inviteUsersTrigger).toBeVisible();
		await userEvent.click(inviteUsersTrigger);

		// Wait for submenu animation
		await new Promise((resolve) => setTimeout(resolve, 200));

		// Check submenu items
		const emailInviteItem = await menuContent.findByRole("menuitem", {
			name: /Email Invite/i,
		});
		await expect(emailInviteItem).toBeVisible();
		const messageInviteItem = await menuContent.findByRole("menuitem", {
			name: /Message Invite/i,
		});
		await expect(messageInviteItem).toBeVisible();
	},
};

export const WithCheckboxAndRadio: Story = {
	render: () => (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button variant="outline">View Options</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent className="w-56">
				<DropdownMenuLabel>Appearance</DropdownMenuLabel>
				<DropdownMenuSeparator />
				<DropdownMenuCheckboxItem checked>
					Show Status Bar
					<DropdownMenuShortcut>⌘B</DropdownMenuShortcut>
				</DropdownMenuCheckboxItem>
				<DropdownMenuCheckboxItem>Show Full Path</DropdownMenuCheckboxItem>
				<DropdownMenuCheckboxItem>Show Line Numbers</DropdownMenuCheckboxItem>
				<DropdownMenuSeparator />
				<DropdownMenuLabel>Theme</DropdownMenuLabel>
				<DropdownMenuRadioGroup value="system">
					<DropdownMenuRadioItem value="light">Light</DropdownMenuRadioItem>
					<DropdownMenuRadioItem value="dark">Dark</DropdownMenuRadioItem>
					<DropdownMenuRadioItem value="system">System</DropdownMenuRadioItem>
				</DropdownMenuRadioGroup>
			</DropdownMenuContent>
		</DropdownMenu>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByRole("button", { name: "View Options" });
		await expect(trigger).toBeInTheDocument();

		// Test menu interaction
		await userEvent.click(trigger);

		// Wait for the dropdown to be mounted in the DOM
		await new Promise((resolve) => setTimeout(resolve, 200));

		// Look for menu items in document.body since they're rendered in a portal
		const menuContent = within(document.body);

		// Check for checkbox items
		const statusBarItem = await menuContent.findByRole("menuitemcheckbox", {
			name: /Show Status Bar/i,
		});
		await expect(statusBarItem).toBeVisible();
		await expect(statusBarItem).toBeChecked();

		const fullPathItem = await menuContent.findByRole("menuitemcheckbox", {
			name: /Show Full Path/i,
		});
		await expect(fullPathItem).toBeVisible();

		// Check for radio items
		const lightThemeItem = await menuContent.findByRole("menuitemradio", {
			name: /Light/i,
		});
		await expect(lightThemeItem).toBeVisible();

		const darkThemeItem = await menuContent.findByRole("menuitemradio", {
			name: /Dark/i,
		});
		await expect(darkThemeItem).toBeVisible();

		const systemThemeItem = await menuContent.findByRole("menuitemradio", {
			name: /System/i,
		});
		await expect(systemThemeItem).toBeVisible();
		await expect(systemThemeItem).toBeChecked();
	},
};
