import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import {
	Copy,
	CreditCard,
	Github,
	Keyboard,
	LogOut,
	Mail,
	MessageSquare,
	Plus,
	Settings,
	User,
	UserPlus,
} from "lucide-react";

import {
	ContextMenu,
	ContextMenuCheckboxItem,
	ContextMenuContent,
	ContextMenuItem,
	ContextMenuLabel,
	ContextMenuRadioGroup,
	ContextMenuRadioItem,
	ContextMenuSeparator,
	ContextMenuShortcut,
	ContextMenuSub,
	ContextMenuSubContent,
	ContextMenuSubTrigger,
	ContextMenuTrigger,
} from "../../components/context-menu.js";

const meta: Meta<typeof ContextMenu> = {
	title: "Design System/ContextMenu",
	component: ContextMenu,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof ContextMenu>;

export const Default: Story = {
	render: () => (
		<ContextMenu>
			<ContextMenuTrigger className="flex h-[150px] w-[300px] items-center justify-center rounded-md border border-dashed text-sm">
				Right click here
			</ContextMenuTrigger>
			<ContextMenuContent className="w-64">
				<ContextMenuItem inset>
					Back
					<ContextMenuShortcut>⌘[</ContextMenuShortcut>
				</ContextMenuItem>
				<ContextMenuItem inset disabled>
					Forward
					<ContextMenuShortcut>⌘]</ContextMenuShortcut>
				</ContextMenuItem>
				<ContextMenuItem inset>
					Reload
					<ContextMenuShortcut>⌘R</ContextMenuShortcut>
				</ContextMenuItem>
				<ContextMenuSeparator />
				<ContextMenuItem>
					<Copy className="mr-2 h-4 w-4" />
					Copy
					<ContextMenuShortcut>⌘C</ContextMenuShortcut>
				</ContextMenuItem>
				<ContextMenuItem>
					<Settings className="mr-2 h-4 w-4" />
					Settings
					<ContextMenuShortcut>⌘,</ContextMenuShortcut>
				</ContextMenuItem>
			</ContextMenuContent>
		</ContextMenu>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByText("Right click here");
		await expect(trigger).toBeInTheDocument();

		// Trigger context menu
		await userEvent.pointer([
			{
				target: trigger,
				keys: "[MouseRight]",
			},
		]);

		// Wait for animation to complete
		await new Promise((resolve) => setTimeout(resolve, 300));

		// Wait for menu to be mounted in document body and animation to complete
		const menu = await within(document.body).findByRole(
			"menu",
			{},
			{ timeout: 2000 },
		);

		// Additional wait to ensure animation is complete
		await new Promise((resolve) => setTimeout(resolve, 100));
		await expect(menu).toBeVisible();

		// Check menu items
		const menuContent = within(menu);
		await expect(
			menuContent.getByRole("menuitem", { name: /back/i }),
		).toBeVisible();
		await expect(
			menuContent.getByRole("menuitem", { name: /forward/i }),
		).toBeVisible();
		await expect(
			menuContent.getByRole("menuitem", { name: /reload/i }),
		).toBeVisible();
		await expect(menuContent.getByRole("separator")).toBeVisible();
		await expect(
			menuContent.getByRole("menuitem", { name: /copy/i }),
		).toBeVisible();
	},
};

export const WithSubmenus: Story = {
	render: () => (
		<ContextMenu>
			<ContextMenuTrigger className="flex h-[150px] w-[300px] items-center justify-center rounded-md border border-dashed text-sm">
				Right click here
			</ContextMenuTrigger>
			<ContextMenuContent className="w-64">
				<ContextMenuLabel>My Account</ContextMenuLabel>
				<ContextMenuItem>
					<User className="mr-2 h-4 w-4" />
					Profile
					<ContextMenuShortcut>⇧⌘P</ContextMenuShortcut>
				</ContextMenuItem>
				<ContextMenuItem>
					<CreditCard className="mr-2 h-4 w-4" />
					Billing
					<ContextMenuShortcut>⌘B</ContextMenuShortcut>
				</ContextMenuItem>
				<ContextMenuSeparator />
				<ContextMenuSub>
					<ContextMenuSubTrigger>
						<UserPlus className="mr-2 h-4 w-4" />
						Invite users
					</ContextMenuSubTrigger>
					<ContextMenuSubContent className="w-48">
						<ContextMenuItem>
							<Mail className="mr-2 h-4 w-4" />
							Email
						</ContextMenuItem>
						<ContextMenuItem>
							<MessageSquare className="mr-2 h-4 w-4" />
							Message
						</ContextMenuItem>
						<ContextMenuSeparator />
						<ContextMenuItem>
							<Plus className="mr-2 h-4 w-4" />
							More...
						</ContextMenuItem>
					</ContextMenuSubContent>
				</ContextMenuSub>
				<ContextMenuSeparator />
				<ContextMenuItem>
					<Github className="mr-2 h-4 w-4" />
					GitHub
				</ContextMenuItem>
				<ContextMenuItem>
					<LogOut className="mr-2 h-4 w-4" />
					Log out
					<ContextMenuShortcut>⇧⌘Q</ContextMenuShortcut>
				</ContextMenuItem>
			</ContextMenuContent>
		</ContextMenu>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByText("Right click here");
		await expect(trigger).toBeInTheDocument();

		// Trigger context menu
		await userEvent.pointer([
			{
				target: trigger,
				keys: "[MouseRight]",
			},
		]);

		// Wait for animation to complete
		await new Promise((resolve) => setTimeout(resolve, 300));

		// Wait for menu to be mounted in document body
		const menu = await within(document.body).findByRole(
			"menu",
			{},
			{ timeout: 2000 },
		);
		await expect(menu).toBeInTheDocument();

		// Check menu items
		const menuContent = within(menu);
		await expect(menuContent.getByText("My Account")).toBeVisible();
		await expect(menuContent.getByText("Profile")).toBeVisible();
		await expect(menuContent.getByText("Billing")).toBeVisible();

		// Check submenu
		const submenuTrigger = menuContent.getByText("Invite users");
		await expect(submenuTrigger).toBeVisible();
		await userEvent.hover(submenuTrigger);

		// Wait for animation to complete
		await new Promise((resolve) => setTimeout(resolve, 300));

		// Wait for submenu to appear
		const submenu = await within(document.body).findByRole(
			"menu",
			{ name: /Invite users/i },
			{ timeout: 2000 },
		);
		await expect(submenu).toBeInTheDocument();

		// Wait for animation to complete
		await new Promise((resolve) => setTimeout(resolve, 300));
		await expect(submenu).toBeVisible();

		// Check submenu items
		const submenuContent = within(submenu);
		await expect(submenuContent.getByText("Email")).toBeVisible();
		await expect(submenuContent.getByText("Message")).toBeVisible();
	},
};

export const WithCheckboxAndRadio: Story = {
	render: () => (
		<ContextMenu>
			<ContextMenuTrigger className="flex h-[150px] w-[300px] items-center justify-center rounded-md border border-dashed text-sm">
				Right click here
			</ContextMenuTrigger>
			<ContextMenuContent className="w-64">
				<ContextMenuLabel>Preferences</ContextMenuLabel>
				<ContextMenuSeparator />
				<ContextMenuCheckboxItem checked>
					Show Full Path
					<ContextMenuShortcut>⌘P</ContextMenuShortcut>
				</ContextMenuCheckboxItem>
				<ContextMenuCheckboxItem>Show Hidden Files</ContextMenuCheckboxItem>
				<ContextMenuSeparator />
				<ContextMenuLabel>View Mode</ContextMenuLabel>
				<ContextMenuRadioGroup value="list">
					<ContextMenuRadioItem value="grid">Grid View</ContextMenuRadioItem>
					<ContextMenuRadioItem value="list">List View</ContextMenuRadioItem>
					<ContextMenuRadioItem value="column">
						Column View
					</ContextMenuRadioItem>
				</ContextMenuRadioGroup>
				<ContextMenuSeparator />
				<ContextMenuItem>
					<Keyboard className="mr-2 h-4 w-4" />
					Keyboard Shortcuts
					<ContextMenuShortcut>⌘K</ContextMenuShortcut>
				</ContextMenuItem>
			</ContextMenuContent>
		</ContextMenu>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const trigger = canvas.getByText("Right click here");
		await expect(trigger).toBeInTheDocument();

		// Trigger context menu
		await userEvent.pointer([
			{
				target: trigger,
				keys: "[MouseRight]",
			},
		]);

		// Wait for animation to complete
		await new Promise((resolve) => setTimeout(resolve, 300));

		// Wait for menu to be mounted in document body and animation to complete
		const menu = await within(document.body).findByRole(
			"menu",
			{},
			{ timeout: 2000 },
		);

		// Additional wait to ensure animation is complete
		await new Promise((resolve) => setTimeout(resolve, 100));
		await expect(menu).toBeVisible();

		// Check menu items
		const menuContent = within(menu);
		await expect(menuContent.getByText("Preferences")).toBeVisible();

		// Check checkbox items
		const showFullPath = menuContent.getByText("Show Full Path");
		await expect(showFullPath).toBeVisible();
		await expect(
			showFullPath.closest('[role="menuitemcheckbox"]'),
		).toHaveAttribute("data-state", "checked");

		const showHiddenFiles = menuContent.getByText("Show Hidden Files");
		await expect(showHiddenFiles).toBeVisible();
		await expect(
			showHiddenFiles.closest('[role="menuitemcheckbox"]'),
		).toHaveAttribute("data-state", "unchecked");

		// Check radio items
		const gridView = menuContent.getByText("Grid View");
		await expect(gridView).toBeVisible();
		await expect(gridView.closest('[role="menuitemradio"]')).toHaveAttribute(
			"data-state",
			"unchecked",
		);

		const listView = menuContent.getByText("List View");
		await expect(listView).toBeVisible();
		await expect(listView.closest('[role="menuitemradio"]')).toHaveAttribute(
			"data-state",
			"checked",
		);

		const columnView = menuContent.getByText("Column View");
		await expect(columnView).toBeVisible();
		await expect(columnView.closest('[role="menuitemradio"]')).toHaveAttribute(
			"data-state",
			"unchecked",
		);
	},
};
