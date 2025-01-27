import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { Building2, Globe, Store } from "lucide-react";
import {
	type Entity,
	EntitySwitcher,
} from "../../components/app-sidebar/entity-switcher.js";
import { SidebarProvider } from "../../components/sidebar.js";

const testEntities: Entity[] = [
	// Platform A with organizations and accounts
	{
		id: "p1",
		name: "Platform A",
		type: "platform",
		logo: Globe,
	},
	{
		id: "o1",
		name: "Organization A1",
		type: "organization",
		parentId: "p1",
		logo: Building2,
	},
	{
		id: "a1",
		name: "Account A1-1",
		type: "account",
		parentId: "o1",
		logo: Store,
	},
	{
		id: "a2",
		name: "Account A1-2",
		type: "account",
		parentId: "o1",
		logo: Store,
	},
	{
		id: "o2",
		name: "Organization A2",
		type: "organization",
		parentId: "p1",
		logo: Building2,
	},
	{
		id: "a3",
		name: "Account A2-1",
		type: "account",
		parentId: "o2",
		logo: Store,
	},

	// Platform B with organizations and accounts
	{
		id: "p2",
		name: "Platform B",
		type: "platform",
		logo: Globe,
	},
	{
		id: "o3",
		name: "Organization B1",
		type: "organization",
		parentId: "p2",
		logo: Building2,
	},
	{
		id: "a4",
		name: "Account B1-1",
		type: "account",
		parentId: "o3",
		logo: Store,
	},

	// Standalone organizations with accounts
	{
		id: "o9",
		name: "Burger King",
		type: "organization",
		logo: Building2,
	},
	{
		id: "a19",
		name: "Times Square",
		type: "account",
		parentId: "o9",
		logo: Store,
	},
	{
		id: "a20",
		name: "Chicago Downtown",
		type: "account",
		parentId: "o9",
		logo: Store,
	},
	{
		id: "a21",
		name: "Miami Beach",
		type: "account",
		parentId: "o9",
		logo: Store,
	},
	{
		id: "o10",
		name: "Dunkin Donuts",
		type: "organization",
		logo: Building2,
	},
	{
		id: "a22",
		name: "Boston Central",
		type: "account",
		parentId: "o10",
		logo: Store,
	},
	{
		id: "a23",
		name: "Manhattan West",
		type: "account",
		parentId: "o10",
		logo: Store,
	},
	{
		id: "a24",
		name: "Brooklyn Heights",
		type: "account",
		parentId: "o10",
		logo: Store,
	},
];

const meta = {
	title: "Application/EntitySwitcher",
	component: EntitySwitcher,
	parameters: {
		layout: "centered",
		viewport: {
			defaultViewport: "desktop",
		},
	},
	args: {
		t: {
			searchPlaceholder: "Search entities...",
			noEntities: "No entities found.",
			noMatchingEntities: "No matching entities found.",
			selectEntity: "Select entity",
			platforms: "Platforms",
			organizations: "Organizations",
			accounts: "Accounts",
			addEntity: "Add entity",
		},
		entities: testEntities,
	},
	argTypes: {
		onEntityChange: { action: "onEntityChange" },
		onCreateClick: { action: "onCreateClick" },
	},
	decorators: [
		(Story) => (
			<SidebarProvider>
				<div className="w-[280px]">
					<Story />
				</div>
			</SidebarProvider>
		),
	],
} satisfies Meta<typeof EntitySwitcher>;

export default meta;
type Story = StoryObj<typeof meta>;

// Test default state and dropdown behavior
export const Default: Story = {
	args: {
		...meta.args,
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Verify initial state
		const trigger = canvas.getByRole("button", {
			name: meta.args.t.selectEntity,
		});
		expect(trigger).toBeInTheDocument();
		expect(trigger).toHaveAttribute("aria-expanded", "false");

		// Open dropdown
		await userEvent.click(trigger);
		await new Promise((resolve) => setTimeout(resolve, 100)); // Wait for animation

		// Verify dropdown is open
		expect(trigger).toHaveAttribute("aria-expanded", "true");

		// Get the portal content
		const portal = within(document.body);

		// Find the dropdown content
		const content = portal.getByRole("listbox");
		expect(content).toBeInTheDocument();

		// Find the command input
		const input = portal.getByPlaceholderText(meta.args.t.searchPlaceholder);
		expect(input).toBeInTheDocument();

		// Find the command groups
		const groups = content.querySelectorAll("[cmdk-group]");
		expect(groups.length).toBeGreaterThan(0);

		// Find the group headings
		const platformsGroup = portal.getByText(meta.args.t.platforms);
		expect(platformsGroup).toBeInTheDocument();

		const organizationsGroup = portal.getByText(meta.args.t.organizations);
		expect(organizationsGroup).toBeInTheDocument();

		// Find the items using flexible text matchers
		const platformA = portal.getByText((content, element) => {
			if (!element) return false;
			const hasText = (node: Element) => node.textContent === "Platform A";
			const nodeHasText = hasText(element);
			const childrenDontHaveText = Array.from(element.children).every(
				(child) => !hasText(child),
			);
			return nodeHasText && childrenDontHaveText;
		});
		const platformAItem = platformA.closest("[cmdk-item]");
		expect(platformAItem).not.toBeNull();

		const burgerKing = portal.getByText((content, element) => {
			if (!element) return false;
			const hasText = (node: Element) => node.textContent === "Burger King";
			const nodeHasText = hasText(element);
			const childrenDontHaveText = Array.from(element.children).every(
				(child) => !hasText(child),
			);
			return nodeHasText && childrenDontHaveText;
		});
		const burgerKingItem = burgerKing.closest("[cmdk-item]");
		expect(burgerKingItem).not.toBeNull();

		// Test keyboard navigation
		await userEvent.keyboard("{Escape}");
		expect(trigger).toHaveAttribute("aria-expanded", "false");
	},
};

// Test entity selection and visual states
export const WithCurrentEntity: Story = {
	args: {
		...meta.args,
		currentEntity: testEntities[0], // Platform A
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Verify selected entity is displayed in trigger
		const trigger = canvas.getByRole("button", {
			name: meta.args.t.selectEntity,
		});
		expect(trigger).toBeInTheDocument();
		expect(trigger.querySelector("span.font-semibold")).toHaveTextContent(
			"Platform A",
		);

		// Open dropdown
		await userEvent.click(trigger);
		await new Promise((resolve) => setTimeout(resolve, 100)); // Wait for animation

		// Get the portal content
		const portal = within(document.body);

		// Find the dropdown content
		const content = portal.getByRole("listbox");
		expect(content).toBeInTheDocument();

		// Verify selected item has correct visual state
		const selectedItem = portal.getByText((content, element) => {
			if (!element) return false;
			const hasText = (node: Element) => node.textContent === "Platform A";
			const nodeHasText = hasText(element);
			const childrenDontHaveText = Array.from(element.children).every(
				(child) => !hasText(child),
			);
			const isCommandItem = element.closest("[cmdk-item]") !== null;
			return nodeHasText && childrenDontHaveText && isCommandItem;
		});
		const selectedItemWrapper = selectedItem.closest("[cmdk-item]");
		expect(selectedItemWrapper).not.toBeNull();
		expect(selectedItemWrapper).toHaveClass("bg-accent");

		// Test selection change
		const newItem = portal.getByText((content, element) => {
			if (!element) return false;
			const hasText = (node: Element) => node.textContent === "Burger King";
			const nodeHasText = hasText(element);
			const childrenDontHaveText = Array.from(element.children).every(
				(child) => !hasText(child),
			);
			const isCommandItem = element.closest("[cmdk-item]") !== null;
			return nodeHasText && childrenDontHaveText && isCommandItem;
		});
		await userEvent.click(newItem);
	},
};

// Test search functionality and filtering
export const SearchFunctionality: Story = {
	args: {
		...meta.args,
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Open dropdown
		const trigger = canvas.getByRole("button", {
			name: meta.args.t.selectEntity,
		});
		await userEvent.click(trigger);
		await new Promise((resolve) => setTimeout(resolve, 100)); // Wait for animation

		// Get the portal content
		const portal = within(document.body);

		// Get search input and wait for focus
		const searchInput = portal.getByPlaceholderText(
			meta.args.t.searchPlaceholder,
		);
		await userEvent.click(searchInput);

		// Test exact match
		await userEvent.type(searchInput, "Burger King");
		await new Promise((resolve) => setTimeout(resolve, 100)); // Wait for filtering
		expect(portal.getByText("Burger King")).toBeVisible();
		expect(portal.queryByText("Platform A")).toBeNull();

		// Clear and test partial match
		await userEvent.clear(searchInput);
		await userEvent.type(searchInput, "Manhattan");
		await new Promise((resolve) => setTimeout(resolve, 100));
		expect(portal.getByText("Manhattan West")).toBeVisible();

		// Test no results
		await userEvent.clear(searchInput);
		await userEvent.type(searchInput, "nonexistent");
		await new Promise((resolve) => setTimeout(resolve, 100));
		expect(
			portal.getByText(meta.args.t.noMatchingEntities),
		).toBeInTheDocument();
	},
};

// Test hierarchical data and collapsible sections
export const HierarchicalBehavior: Story = {
	args: {
		...meta.args,
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Open dropdown
		const trigger = canvas.getByRole("button", {
			name: meta.args.t.selectEntity,
		});
		await userEvent.click(trigger);
		await new Promise((resolve) => setTimeout(resolve, 100));

		// Get the portal content
		const portal = within(document.body);

		// Find the command menu
		const commandMenu = portal.getByRole("listbox");

		// Find the group headings
		const headings = commandMenu.querySelectorAll("[cmdk-group-heading]");
		const platformsHeading = Array.from(headings).find(
			(h) => h.textContent === meta.args.t.platforms,
		);
		expect(platformsHeading).toBeInTheDocument();

		// Find the items
		const items = commandMenu.querySelectorAll("[cmdk-item]");
		const itemTexts = Array.from(items).map((i) => i.textContent);
		expect(itemTexts).toContain("Organization A1");
		expect(itemTexts).toContain("Account A1-1");
	},
};

// Test mobile responsiveness
export const MobileLayout: Story = {
	args: {
		...meta.args,
	},
	parameters: {
		viewport: {
			defaultViewport: "mobile1",
		},
		chromatic: {
			viewports: ["mobile1"],
		},
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Open dropdown
		const trigger = canvas.getByRole("button", {
			name: meta.args.t.selectEntity,
		});
		await userEvent.click(trigger);
		await new Promise((resolve) => setTimeout(resolve, 100));

		// Get the portal content
		const portal = within(document.body);

		// Verify mobile positioning
		const popperContent = portal.getByRole("menu").closest("[data-side]");
		expect(popperContent).toHaveAttribute("data-side", "bottom");

		// Test touch interaction
		const platformA = portal.getByText((content, element) => {
			if (!element) return false;
			const hasText = (node: Element) => node.textContent === "Platform A";
			const nodeHasText = hasText(element);
			const childrenDontHaveText = Array.from(element.children).every(
				(child) => !hasText(child),
			);
			return nodeHasText && childrenDontHaveText;
		});
		await userEvent.click(platformA);
	},
};

// Test keyboard navigation and accessibility
export const KeyboardNavigation: Story = {
	args: {
		...meta.args,
		currentEntity: testEntities.find((e) => e.name === "Organization A1"),
	},
	play: async ({ canvasElement, args }) => {
		const canvas = within(canvasElement);

		// Open with keyboard
		const trigger = canvas.getByRole("button", {
			name: meta.args.t.selectEntity,
		});
		await userEvent.tab();
		expect(trigger).toHaveFocus();
		await userEvent.keyboard("{Enter}");
		await new Promise((resolve) => setTimeout(resolve, 100));

		// Get the portal content
		const portal = within(document.body);

		// Verify dropdown is open
		expect(trigger).toHaveAttribute("aria-expanded", "true");

		// Find and click Organization A1 directly
		const orgA1 = portal.getByText((content, element) => {
			if (!element) return false;
			const hasText = (node: Element) => node.textContent === "Organization A1";
			const nodeHasText = hasText(element);
			const childrenDontHaveText = Array.from(element.children).every(
				(child) => !hasText(child),
			);
			const isCommandItem = element.closest("[cmdk-item]") !== null;
			return nodeHasText && childrenDontHaveText && isCommandItem;
		});
		await userEvent.click(orgA1);
		await new Promise((resolve) => setTimeout(resolve, 100));

		// Verify dropdown is closed
		expect(trigger).toHaveAttribute("aria-expanded", "false");

		// Test type to select
		await userEvent.click(trigger);
		await new Promise((resolve) => setTimeout(resolve, 100));

		await userEvent.keyboard("bur");
		expect(portal.getByText("Burger King")).toBeVisible();
	},
};

// Test empty state
export const EmptyState: Story = {
	args: {
		...meta.args,
		entities: [],
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Open dropdown
		const trigger = canvas.getByRole("button", {
			name: meta.args.t.selectEntity,
		});
		await userEvent.click(trigger);
		await new Promise((resolve) => setTimeout(resolve, 100));

		// Get the portal content
		const portal = within(document.body);

		// Verify empty state
		expect(portal.getByText(meta.args.t.noEntities)).toBeInTheDocument();
	},
};

// Test create entity functionality
export const CreateEntity: Story = {
	args: {
		...meta.args,
		canCreate: true,
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);

		// Open dropdown
		const trigger = canvas.getByRole("button", {
			name: meta.args.t.selectEntity,
		});
		await userEvent.click(trigger);
		await new Promise((resolve) => setTimeout(resolve, 100));

		// Get the portal content
		const portal = within(document.body);

		// Verify create button
		const createButton = portal.getByText(meta.args.t.addEntity);
		expect(createButton).toBeVisible();

		// Test create button interaction
		await userEvent.click(createButton);
		expect(trigger).toHaveAttribute("aria-expanded", "false");
	},
};
