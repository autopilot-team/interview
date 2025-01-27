import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import { ScrollArea } from "../../components/scroll-area.js";

const meta: Meta<typeof ScrollArea> = {
	title: "Design System/ScrollArea",
	component: ScrollArea,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof ScrollArea>;

const tags = Array.from({ length: 50 }).map(
	(_, i, a) => `v1.2.${a.length - i}`,
);

export const Default: Story = {
	render: () => (
		<ScrollArea className="h-72 w-48 rounded-md border">
			<div className="p-4" data-testid="radix-scroll-area-viewport">
				<h4 className="mb-4 text-sm font-medium leading-none">Tags</h4>
				{tags.map((tag) => (
					<div key={tag} className="text-sm">
						{tag}
					</div>
				))}
			</div>
		</ScrollArea>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const viewport = canvas.getByTestId("radix-scroll-area-viewport");
		await expect(viewport).toBeInTheDocument();

		// Verify some content is visible
		const heading = canvas.getByRole("heading", { name: "Tags" });
		await expect(heading).toBeVisible();

		// Verify some tags are visible
		const tags = canvas.getAllByText(/v1\.2\.\d+/);
		await expect(tags.length).toBeGreaterThan(0);
	},
};

export const Horizontal: Story = {
	render: () => (
		<ScrollArea className="w-96 whitespace-nowrap rounded-md border p-4">
			<div className="flex w-[800px] gap-4">
				{Array.from({ length: 10 }, (_, i) => ({
					id: `item-${i + 1}`,
					title: `Item ${i + 1}`,
				})).map((item) => (
					<div
						key={item.id}
						className="w-[150px] flex-none rounded-md border border-dashed p-4"
					>
						<div className="font-semibold">{item.title}</div>
						<div className="text-sm text-muted-foreground">
							Horizontal scrolling content
						</div>
					</div>
				))}
			</div>
		</ScrollArea>
	),
};

export const WithMaxHeight: Story = {
	render: () => (
		<ScrollArea className="h-[200px] w-[350px] rounded-md border p-4">
			<div className="space-y-4">
				<h4 className="font-medium leading-none">Lorem Ipsum</h4>
				<p className="text-sm text-muted-foreground">
					Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do
					eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad
					minim veniam, quis nostrud exercitation ullamco laboris nisi ut
					aliquip ex ea commodo consequat.
				</p>
				<p className="text-sm text-muted-foreground">
					Duis aute irure dolor in reprehenderit in voluptate velit esse cillum
					dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non
					proident, sunt in culpa qui officia deserunt mollit anim id est
					laborum.
				</p>
				<p className="text-sm text-muted-foreground">
					Sed ut perspiciatis unde omnis iste natus error sit voluptatem
					accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae
					ab illo inventore veritatis et quasi architecto beatae vitae dicta
					sunt explicabo.
				</p>
			</div>
		</ScrollArea>
	),
};

export const WithCard: Story = {
	render: () => (
		<ScrollArea className="h-[400px] w-[300px] rounded-md border">
			<div className="p-4">
				<h4 className="mb-4 text-sm font-medium">Messages</h4>
				{Array.from({ length: 20 }, (_, i) => ({
					id: `message-${i + 1}`,
					title: `Message ${i + 1}`,
				})).map((message) => (
					<div key={message.id} className="mb-4 rounded-md border p-3">
						<div className="font-semibold">{message.title}</div>
						<div className="mt-1 text-sm text-muted-foreground">
							This is a sample message content that demonstrates a card within a
							scrollable area.
						</div>
					</div>
				))}
			</div>
		</ScrollArea>
	),
};

export const WithNestedContent: Story = {
	render: () => (
		<ScrollArea className="h-[400px] w-[600px] rounded-md border p-4">
			<div className="space-y-8">
				{Array.from({ length: 3 }, (_, i) => ({
					id: `section-${i + 1}`,
					title: `Section ${i + 1}`,
					items: Array.from({ length: 4 }, (_, j) => ({
						id: `item-${i + 1}-${j + 1}`,
						title: `Item ${j + 1}`,
					})),
				})).map((section) => (
					<div key={section.id} className="space-y-4">
						<h4 className="text-lg font-medium">{section.title}</h4>
						<div className="grid gap-4 md:grid-cols-2">
							{section.items.map((item) => (
								<div key={item.id} className="rounded-md border p-4">
									<h5 className="font-medium">{item.title}</h5>
									<p className="mt-2 text-sm text-muted-foreground">
										This is a nested item within {section.title}. It
										demonstrates complex content structure within a scrollable
										area.
									</p>
								</div>
							))}
						</div>
					</div>
				))}
			</div>
		</ScrollArea>
	),
};
