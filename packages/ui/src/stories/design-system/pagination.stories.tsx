import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import {
	Pagination,
	PaginationContent,
	PaginationEllipsis,
	PaginationItem,
	PaginationLink,
	PaginationNext,
	PaginationPrevious,
} from "../../components/pagination.js";

const meta: Meta<typeof Pagination> = {
	title: "Design System/Pagination",
	component: Pagination,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Pagination>;

export const Default: Story = {
	render: () => (
		<Pagination>
			<PaginationContent>
				<PaginationItem>
					<PaginationPrevious href="#" />
				</PaginationItem>
				<PaginationItem>
					<PaginationLink href="#" isActive>
						1
					</PaginationLink>
				</PaginationItem>
				<PaginationItem>
					<PaginationLink href="#">2</PaginationLink>
				</PaginationItem>
				<PaginationItem>
					<PaginationLink href="#">3</PaginationLink>
				</PaginationItem>
				<PaginationItem>
					<PaginationNext href="#" />
				</PaginationItem>
			</PaginationContent>
		</Pagination>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const pagination = canvas.getByRole("navigation");
		await expect(pagination).toBeInTheDocument();

		const currentPage = canvas.getByRole("link", { current: "page" });
		await expect(currentPage).toHaveTextContent("1");

		const nextButton = canvas.getByRole("link", { name: /next/i });
		await expect(nextButton).toBeInTheDocument();
	},
};

export const WithEllipsis: Story = {
	render: () => (
		<Pagination>
			<PaginationContent>
				<PaginationItem>
					<PaginationPrevious href="#" />
				</PaginationItem>
				<PaginationItem>
					<PaginationLink href="#" isActive>
						1
					</PaginationLink>
				</PaginationItem>
				<PaginationItem>
					<PaginationLink href="#">2</PaginationLink>
				</PaginationItem>
				<PaginationItem>
					<PaginationLink href="#">3</PaginationLink>
				</PaginationItem>
				<PaginationItem>
					<PaginationEllipsis />
				</PaginationItem>
				<PaginationItem>
					<PaginationLink href="#">12</PaginationLink>
				</PaginationItem>
				<PaginationItem>
					<PaginationNext href="#" />
				</PaginationItem>
			</PaginationContent>
		</Pagination>
	),
};

export const Disabled: Story = {
	render: () => (
		<Pagination>
			<PaginationContent>
				<PaginationItem>
					<PaginationPrevious
						href="#"
						aria-disabled="true"
						className="pointer-events-none opacity-50"
					/>
				</PaginationItem>
				<PaginationItem>
					<PaginationLink href="#" isActive>
						1
					</PaginationLink>
				</PaginationItem>
				<PaginationItem>
					<PaginationLink href="#">2</PaginationLink>
				</PaginationItem>
				<PaginationItem>
					<PaginationLink href="#">3</PaginationLink>
				</PaginationItem>
				<PaginationItem>
					<PaginationNext href="#" />
				</PaginationItem>
			</PaginationContent>
		</Pagination>
	),
};

export const CustomStyling: Story = {
	render: () => (
		<Pagination>
			<PaginationContent>
				<PaginationItem>
					<PaginationPrevious
						href="#"
						className="hover:bg-primary hover:text-primary-foreground"
					/>
				</PaginationItem>
				{[1, 2, 3, 4, 5].map((page) => (
					<PaginationItem key={page}>
						<PaginationLink
							href="#"
							isActive={page === 1}
							className={
								page === 1
									? "bg-primary text-primary-foreground hover:bg-primary/90"
									: ""
							}
						>
							{page}
						</PaginationLink>
					</PaginationItem>
				))}
				<PaginationItem>
					<PaginationNext
						href="#"
						className="hover:bg-primary hover:text-primary-foreground"
					/>
				</PaginationItem>
			</PaginationContent>
		</Pagination>
	),
};
