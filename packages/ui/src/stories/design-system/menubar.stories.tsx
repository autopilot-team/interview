import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import {
	Menubar,
	MenubarCheckboxItem,
	MenubarContent,
	MenubarItem,
	MenubarMenu,
	MenubarRadioGroup,
	MenubarRadioItem,
	MenubarSeparator,
	MenubarShortcut,
	MenubarSub,
	MenubarSubContent,
	MenubarSubTrigger,
	MenubarTrigger,
} from "../../components/menubar.js";

const meta: Meta<typeof Menubar> = {
	title: "Design System/Menubar",
	component: Menubar,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Menubar>;

export const Default: Story = {
	render: () => (
		<Menubar>
			<MenubarMenu>
				<MenubarTrigger>File</MenubarTrigger>
				<MenubarContent>
					<MenubarItem>
						New Tab <MenubarShortcut>⌘T</MenubarShortcut>
					</MenubarItem>
					<MenubarItem>
						New Window <MenubarShortcut>⌘N</MenubarShortcut>
					</MenubarItem>
					<MenubarSeparator />
					<MenubarItem>Share</MenubarItem>
					<MenubarSeparator />
					<MenubarItem>
						Print... <MenubarShortcut>⌘P</MenubarShortcut>
					</MenubarItem>
				</MenubarContent>
			</MenubarMenu>
			<MenubarMenu>
				<MenubarTrigger>Edit</MenubarTrigger>
				<MenubarContent>
					<MenubarItem>
						Undo <MenubarShortcut>⌘Z</MenubarShortcut>
					</MenubarItem>
					<MenubarItem>
						Redo <MenubarShortcut>⇧⌘Z</MenubarShortcut>
					</MenubarItem>
					<MenubarSeparator />
					<MenubarItem>Cut</MenubarItem>
					<MenubarItem>Copy</MenubarItem>
					<MenubarItem>Paste</MenubarItem>
				</MenubarContent>
			</MenubarMenu>
		</Menubar>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Test initial menu presence
		const fileMenu = canvas.getByRole("menuitem", { name: "File" });
		const editMenu = canvas.getByRole("menuitem", { name: "Edit" });
		await expect(fileMenu).toBeInTheDocument();
		await expect(editMenu).toBeInTheDocument();

		// Get initial menu trigger texts for verification
		const menuTriggers = canvas
			.getAllByRole("menuitem")
			.map((item) => item.textContent);
		expect(menuTriggers).toEqual(["File", "Edit"]);

		// Test File menu interaction
		await userEvent.click(fileMenu);
		await new Promise((resolve) => setTimeout(resolve, 200));

		// Menu items should be in a portal
		const menuContent = within(document.body);

		// Wait for menu content to be mounted and verify File menu items
		const fileMenuItems = await menuContent.findAllByRole("menuitem", {
			hidden: true,
		});
		const fileMenuTexts = fileMenuItems
			.filter((item) => !["File", "Edit"].includes(item.textContent || ""))
			.map((item) => item.textContent?.replace(/\s+/g, " ").trim());

		expect(fileMenuTexts).toEqual([
			"New Tab ⌘T",
			"New Window ⌘N",
			"Share",
			"Print... ⌘P",
		]);

		// Test menu switching
		await userEvent.click(editMenu);
		await userEvent.click(editMenu);
		await new Promise((resolve) => setTimeout(resolve, 500));

		// Menu items should be in a portal
		const editMenuContent = within(document.body);

		// Wait for menu content to be mounted and verify Edit menu items
		const editMenuItems = await editMenuContent.findAllByRole("menuitem", {
			hidden: true,
		});
		const editMenuTexts = editMenuItems
			.filter((item) => !["File", "Edit"].includes(item.textContent || ""))
			.map((item) => item.textContent?.replace(/\s+/g, " ").trim());

		expect(editMenuTexts).toEqual([
			"Undo ⌘Z",
			"Redo ⇧⌘Z",
			"Cut",
			"Copy",
			"Paste",
		]);
	},
};

export const WithSubmenus: Story = {
	render: () => (
		<Menubar>
			<MenubarMenu>
				<MenubarTrigger>View</MenubarTrigger>
				<MenubarContent>
					<MenubarSub>
						<MenubarSubTrigger>Zoom</MenubarSubTrigger>
						<MenubarSubContent>
							<MenubarItem>
								Zoom In <MenubarShortcut>⌘+</MenubarShortcut>
							</MenubarItem>
							<MenubarItem>
								Zoom Out <MenubarShortcut>⌘-</MenubarShortcut>
							</MenubarItem>
							<MenubarItem>
								Reset Zoom <MenubarShortcut>⌘0</MenubarShortcut>
							</MenubarItem>
						</MenubarSubContent>
					</MenubarSub>
					<MenubarSeparator />
					<MenubarItem>
						Enter Full Screen <MenubarShortcut>⌘F</MenubarShortcut>
					</MenubarItem>
				</MenubarContent>
			</MenubarMenu>
		</Menubar>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Test initial menu presence
		const viewMenu = canvas.getByRole("menuitem", { name: "View" });
		await expect(viewMenu).toBeInTheDocument();

		// Open View menu
		await userEvent.click(viewMenu);
		await new Promise((resolve) => setTimeout(resolve, 200));

		// Menu items should be in a portal
		const menuContent = within(document.body);
		const zoomTrigger = menuContent.getByRole("menuitem", { name: "Zoom" });
		await expect(zoomTrigger).toBeVisible();

		// Test submenu interaction
		await userEvent.hover(zoomTrigger);
		await new Promise((resolve) => setTimeout(resolve, 200));

		// Verify submenu items
		await expect(menuContent.getByText("Zoom In")).toBeVisible();
		await expect(menuContent.getByText("Zoom Out")).toBeVisible();
		await expect(menuContent.getByText("Reset Zoom")).toBeVisible();

		// Verify main menu item is still visible
		await expect(menuContent.getByText("Enter Full Screen")).toBeVisible();
	},
};

export const WithCheckboxAndRadio: Story = {
	render: () => (
		<Menubar>
			<MenubarMenu>
				<MenubarTrigger>Options</MenubarTrigger>
				<MenubarContent>
					<MenubarCheckboxItem>Show Status Bar</MenubarCheckboxItem>
					<MenubarCheckboxItem checked>Show Full Path</MenubarCheckboxItem>
					<MenubarSeparator />
					<MenubarRadioGroup value="betty">
						<MenubarRadioItem value="andy">Andy</MenubarRadioItem>
						<MenubarRadioItem value="betty">Betty</MenubarRadioItem>
						<MenubarRadioItem value="charlie">Charlie</MenubarRadioItem>
					</MenubarRadioGroup>
				</MenubarContent>
			</MenubarMenu>
		</Menubar>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Test initial menu presence
		const optionsMenu = canvas.getByRole("menuitem", { name: "Options" });
		await expect(optionsMenu).toBeInTheDocument();

		// Open Options menu
		await userEvent.click(optionsMenu);
		await new Promise((resolve) => setTimeout(resolve, 200));

		// Menu items should be in a portal
		const menuContent = within(document.body);

		// Test checkbox items
		const statusBarItem = menuContent.getByRole("menuitemcheckbox", {
			name: "Show Status Bar",
		});
		const fullPathItem = menuContent.getByRole("menuitemcheckbox", {
			name: "Show Full Path",
		});

		await expect(statusBarItem).toBeVisible();
		await expect(statusBarItem).not.toBeChecked();
		await expect(fullPathItem).toBeVisible();
		await expect(fullPathItem).toBeChecked();

		// Test radio items
		const andyItem = menuContent.getByRole("menuitemradio", { name: "Andy" });
		const bettyItem = menuContent.getByRole("menuitemradio", { name: "Betty" });
		const charlieItem = menuContent.getByRole("menuitemradio", {
			name: "Charlie",
		});

		await expect(andyItem).toBeVisible();
		await expect(andyItem).not.toBeChecked();
		await expect(bettyItem).toBeVisible();
		await expect(bettyItem).toBeChecked();
		await expect(charlieItem).toBeVisible();
		await expect(charlieItem).not.toBeChecked();
	},
};

export const ComplexExample: Story = {
	render: () => (
		<Menubar className="w-[500px]">
			<MenubarMenu>
				<MenubarTrigger>File</MenubarTrigger>
				<MenubarContent>
					<MenubarItem>
						New File <MenubarShortcut>⌘N</MenubarShortcut>
					</MenubarItem>
					<MenubarSub>
						<MenubarSubTrigger>Share</MenubarSubTrigger>
						<MenubarSubContent>
							<MenubarItem>Email Link</MenubarItem>
							<MenubarItem>Messages</MenubarItem>
							<MenubarItem>Notes</MenubarItem>
						</MenubarSubContent>
					</MenubarSub>
					<MenubarSeparator />
					<MenubarItem>
						Print <MenubarShortcut>⌘P</MenubarShortcut>
					</MenubarItem>
				</MenubarContent>
			</MenubarMenu>
			<MenubarMenu>
				<MenubarTrigger>Edit</MenubarTrigger>
				<MenubarContent>
					<MenubarCheckboxItem>Show Invisibles</MenubarCheckboxItem>
					<MenubarSeparator />
					<MenubarItem disabled>
						Undo <MenubarShortcut>⌘Z</MenubarShortcut>
					</MenubarItem>
					<MenubarItem>
						Redo <MenubarShortcut>⇧⌘Z</MenubarShortcut>
					</MenubarItem>
				</MenubarContent>
			</MenubarMenu>
			<MenubarMenu>
				<MenubarTrigger>View</MenubarTrigger>
				<MenubarContent>
					<MenubarRadioGroup value="compact">
						<MenubarRadioItem value="compact">Compact</MenubarRadioItem>
						<MenubarRadioItem value="default">Default</MenubarRadioItem>
						<MenubarRadioItem value="expanded">Expanded</MenubarRadioItem>
					</MenubarRadioGroup>
				</MenubarContent>
			</MenubarMenu>
		</Menubar>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Test initial menu presence
		const fileMenu = canvas.getByRole("menuitem", { name: "File" });
		const editMenu = canvas.getByRole("menuitem", { name: "Edit" });
		const viewMenu = canvas.getByRole("menuitem", { name: "View" });

		await expect(fileMenu).toBeInTheDocument();
		await expect(editMenu).toBeInTheDocument();
		await expect(viewMenu).toBeInTheDocument();

		// Test File menu and submenu
		await userEvent.click(fileMenu);
		await new Promise((resolve) => setTimeout(resolve, 200));

		const menuContent = within(document.body);
		await expect(menuContent.getByText("New File")).toBeVisible();

		const shareTrigger = menuContent.getByRole("menuitem", { name: "Share" });
		await expect(shareTrigger).toBeVisible();
		await userEvent.hover(shareTrigger);
		await new Promise((resolve) => setTimeout(resolve, 200));

		await expect(menuContent.getByText("Email Link")).toBeVisible();
		await expect(menuContent.getByText("Messages")).toBeVisible();
		await expect(menuContent.getByText("Notes")).toBeVisible();

		// Test Edit menu with disabled item
		await userEvent.click(editMenu);
		await userEvent.click(editMenu);
		await new Promise((resolve) => setTimeout(resolve, 200));

		const undoItem = menuContent.getByRole("menuitem", { name: /Undo/i });
		await expect(undoItem).toBeVisible();
		await expect(undoItem).toHaveAttribute("data-disabled");

		const showInvisiblesItem = menuContent.getByRole("menuitemcheckbox", {
			name: "Show Invisibles",
		});
		await expect(showInvisiblesItem).toBeVisible();
		await expect(showInvisiblesItem).not.toBeChecked();

		// Test View menu with radio items
		await userEvent.click(viewMenu);
		await userEvent.click(viewMenu);
		await new Promise((resolve) => setTimeout(resolve, 200));

		const compactItem = menuContent.getByRole("menuitemradio", {
			name: "Compact",
		});
		await expect(compactItem).toBeVisible();
		await expect(compactItem).toBeChecked();

		const defaultItem = menuContent.getByRole("menuitemradio", {
			name: "Default",
		});
		await expect(defaultItem).toBeVisible();
		await expect(defaultItem).not.toBeChecked();
	},
};
