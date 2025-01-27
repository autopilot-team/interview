import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import { Checkbox } from "../../components/checkbox.js";
import {
	Table,
	TableBody,
	TableCaption,
	TableCell,
	TableFooter,
	TableHead,
	TableHeader,
	TableRow,
} from "../../components/table.js";

const meta = {
	title: "Design System/Table",
	component: Table,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
} satisfies Meta<typeof Table>;

export default meta;
type Story = StoryObj<typeof Table>;

const invoices = [
	{
		invoice: "INV001",
		paymentStatus: "Paid",
		totalAmount: "$250.00",
		paymentMethod: "Credit Card",
	},
	{
		invoice: "INV002",
		paymentStatus: "Pending",
		totalAmount: "$150.00",
		paymentMethod: "PayPal",
	},
	{
		invoice: "INV003",
		paymentStatus: "Unpaid",
		totalAmount: "$350.00",
		paymentMethod: "Bank Transfer",
	},
	{
		invoice: "INV004",
		paymentStatus: "Paid",
		totalAmount: "$450.00",
		paymentMethod: "Credit Card",
	},
	{
		invoice: "INV005",
		paymentStatus: "Paid",
		totalAmount: "$550.00",
		paymentMethod: "PayPal",
	},
];

export const Default: Story = {
	render: () => (
		<Table>
			<TableCaption>A list of your recent invoices.</TableCaption>
			<TableHeader>
				<TableRow>
					<TableHead>Invoice</TableHead>
					<TableHead>Status</TableHead>
					<TableHead>Method</TableHead>
					<TableHead className="text-right">Amount</TableHead>
				</TableRow>
			</TableHeader>
			<TableBody>
				{invoices.map((invoice) => (
					<TableRow key={invoice.invoice}>
						<TableCell className="font-medium">{invoice.invoice}</TableCell>
						<TableCell>{invoice.paymentStatus}</TableCell>
						<TableCell>{invoice.paymentMethod}</TableCell>
						<TableCell className="text-right">{invoice.totalAmount}</TableCell>
					</TableRow>
				))}
			</TableBody>
			<TableFooter>
				<TableRow>
					<TableCell colSpan={3}>Total</TableCell>
					<TableCell className="text-right">$1,750.00</TableCell>
				</TableRow>
			</TableFooter>
		</Table>
	),
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const table = canvas.getByRole("table");
		await expect(table).toBeInTheDocument();
		const rows = canvas.getAllByRole("row");
		await expect(rows).toHaveLength(7); // Header + 5 rows + footer
	},
};

export const WithSelection: Story = {
	render: () => (
		<Table>
			<TableHeader>
				<TableRow>
					<TableHead className="w-[50px]">
						<Checkbox />
					</TableHead>
					<TableHead>Invoice</TableHead>
					<TableHead>Status</TableHead>
					<TableHead>Method</TableHead>
					<TableHead className="text-right">Amount</TableHead>
				</TableRow>
			</TableHeader>
			<TableBody>
				{invoices.map((invoice) => (
					<TableRow key={invoice.invoice}>
						<TableCell>
							<Checkbox />
						</TableCell>
						<TableCell className="font-medium">{invoice.invoice}</TableCell>
						<TableCell>{invoice.paymentStatus}</TableCell>
						<TableCell>{invoice.paymentMethod}</TableCell>
						<TableCell className="text-right">{invoice.totalAmount}</TableCell>
					</TableRow>
				))}
			</TableBody>
		</Table>
	),
};

const tasks = [
	{
		id: "TASK-8782",
		title:
			"You can't compress the program without quantifying the open-source SSD pixel!",
		status: "in progress",
		label: "documentation",
		priority: "high",
	},
	{
		id: "TASK-7878",
		title:
			"Try to calculate the EXE feed, maybe it will index the multi-byte pixel!",
		status: "backlog",
		label: "bug",
		priority: "medium",
	},
	{
		id: "TASK-7839",
		title: "We need to bypass the neural TCP card!",
		status: "todo",
		label: "feature",
		priority: "high",
	},
];

export const WithBadges: Story = {
	render: () => (
		<Table>
			<TableHeader>
				<TableRow>
					<TableHead className="w-[100px]">Task</TableHead>
					<TableHead>Title</TableHead>
					<TableHead className="w-[100px]">Status</TableHead>
					<TableHead className="w-[100px]">Priority</TableHead>
				</TableRow>
			</TableHeader>
			<TableBody>
				{tasks.map((task) => (
					<TableRow key={task.id}>
						<TableCell className="font-medium">{task.id}</TableCell>
						<TableCell>{task.title}</TableCell>
						<TableCell>
							<div className="flex w-[100px] items-center">
								<span
									className={`capitalize ${task.status === "in progress" ? "text-blue-500" : task.status === "todo" ? "text-yellow-500" : "text-slate-500"}`}
								>
									{task.status}
								</span>
							</div>
						</TableCell>
						<TableCell>
							<div className="flex w-[100px] items-center">
								<span
									className={`capitalize ${task.priority === "high" ? "text-red-500" : "text-orange-500"}`}
								>
									{task.priority}
								</span>
							</div>
						</TableCell>
					</TableRow>
				))}
			</TableBody>
		</Table>
	),
};
