import type { Meta, StoryObj } from "@storybook/react";
import { expect, within } from "@storybook/test";
import {
	Area,
	Bar,
	CartesianGrid,
	Cell,
	Line,
	Pie,
	AreaChart as RechartsAreaChart,
	BarChart as RechartsBarChart,
	LineChart as RechartsLineChart,
	PieChart as RechartsPieChart,
	XAxis,
	YAxis,
} from "recharts";
import { ChartContainer, ChartTooltip } from "../../components/chart.js";

const meta: Meta<typeof ChartContainer> = {
	title: "Design System/Chart",
	component: ChartContainer,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof ChartContainer>;

const lineData = [
	{ name: "Jan", value: 100, previousValue: 80 },
	{ name: "Feb", value: 120, previousValue: 100 },
	{ name: "Mar", value: 170, previousValue: 140 },
	{ name: "Apr", value: 140, previousValue: 160 },
	{ name: "May", value: 200, previousValue: 170 },
	{ name: "Jun", value: 180, previousValue: 190 },
];

const barData = [
	{ name: "Q1", revenue: 400, expenses: 300 },
	{ name: "Q2", revenue: 300, expenses: 250 },
	{ name: "Q3", revenue: 500, expenses: 400 },
	{ name: "Q4", revenue: 280, expenses: 220 },
];

const pieData = [
	{ name: "Desktop", value: 400 },
	{ name: "Mobile", value: 300 },
	{ name: "Tablet", value: 200 },
];

const areaData = [
	{ name: "Mon", users: 10, activeUsers: 8 },
	{ name: "Tue", users: 30, activeUsers: 20 },
	{ name: "Wed", users: 20, activeUsers: 15 },
	{ name: "Thu", users: 40, activeUsers: 30 },
	{ name: "Fri", users: 25, activeUsers: 18 },
];

export const LineChartStory: Story = {
	args: {
		config: {
			value: {
				label: "Current",
				color: "hsl(var(--primary))",
			},
			previousValue: {
				label: "Previous",
				color: "hsl(var(--muted-foreground))",
			},
		},
		className: "w-[600px] h-[400px]",
		children: (
			<div className="w-[600px] h-[400px]">
				<section aria-label="Line chart">
					<RechartsLineChart
						width={600}
						height={400}
						data={lineData}
						margin={{ top: 16, right: 16, bottom: 16, left: 16 }}
					>
						<XAxis
							dataKey="name"
							stroke="hsl(var(--muted-foreground))"
							fontSize={12}
							tickLine={false}
							axisLine={false}
							dy={10}
						/>
						<YAxis
							stroke="hsl(var(--muted-foreground))"
							fontSize={12}
							tickLine={false}
							axisLine={false}
							dx={-10}
						/>
						<CartesianGrid
							strokeDasharray="4"
							className="stroke-border"
							vertical={false}
						/>
						<Line
							type="monotone"
							dataKey="previousValue"
							stroke="hsl(var(--muted-foreground))"
							strokeWidth={1.5}
							dot={false}
							strokeDasharray="4"
							opacity={0.5}
						/>
						<Line
							type="monotone"
							dataKey="value"
							stroke="hsl(var(--primary))"
							strokeWidth={1.5}
							dot={false}
							activeDot={{
								r: 4,
								fill: "hsl(var(--primary))",
								stroke: "hsl(var(--background))",
								strokeWidth: 2,
							}}
						/>
						<ChartTooltip />
					</RechartsLineChart>
				</section>
			</div>
		),
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const chart = await canvas.findByRole("region", { name: /line chart/i });
		await expect(chart).toBeInTheDocument();

		// Check if axes are visible
		await expect(await canvas.findByText("Jan")).toBeVisible();
		await expect(await canvas.findByText("100")).toBeVisible();
	},
};

export const BarChartStory: Story = {
	args: {
		config: {
			revenue: {
				label: "Revenue",
				color: "hsl(var(--primary))",
			},
			expenses: {
				label: "Expenses",
				color: "hsl(var(--muted-foreground))",
			},
		},
		className: "w-[600px] h-[400px]",
		children: (
			<div className="w-[600px] h-[400px]">
				<section aria-label="Bar chart">
					<RechartsBarChart
						width={600}
						height={400}
						data={barData}
						margin={{ top: 16, right: 16, bottom: 16, left: 16 }}
					>
						<XAxis
							dataKey="name"
							stroke="hsl(var(--muted-foreground))"
							fontSize={12}
							tickLine={false}
							axisLine={false}
							dy={10}
						/>
						<YAxis
							stroke="hsl(var(--muted-foreground))"
							fontSize={12}
							tickLine={false}
							axisLine={false}
							dx={-10}
						/>
						<CartesianGrid
							strokeDasharray="4"
							className="stroke-border"
							vertical={false}
						/>
						<Bar
							dataKey="expenses"
							fill="hsl(var(--muted-foreground))"
							radius={[4, 4, 0, 0]}
							opacity={0.5}
						/>
						<Bar
							dataKey="revenue"
							fill="hsl(var(--primary))"
							radius={[4, 4, 0, 0]}
						/>
						<ChartTooltip />
					</RechartsBarChart>
				</section>
			</div>
		),
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const chart = await canvas.findByRole("region", { name: /bar chart/i });
		await expect(chart).toBeInTheDocument();

		// Wait for chart to render
		await new Promise((resolve) => setTimeout(resolve, 1000));

		// Check if axes are visible
		await expect(await canvas.findByText("Q1")).toBeVisible();

		// Find the SVG element and check its contents
		const svgElement = canvasElement.querySelector("svg");
		expect(svgElement).toBeTruthy();

		// Look for text elements with the value we expect
		const yAxisTexts = Array.from(
			svgElement?.querySelectorAll(
				".recharts-cartesian-axis-tick-value tspan",
			) ?? [],
		);
		const hasExpectedValue = yAxisTexts.some((node) => {
			const value = node.textContent?.trim();
			return value === "400" || value === "300" || value === "500";
		});
		expect(hasExpectedValue).toBe(true);
	},
};

export const PieChartStory: Story = {
	args: {
		config: {
			Desktop: {
				label: "Desktop Users",
				color: "hsl(var(--primary))",
			},
			Mobile: {
				label: "Mobile Users",
				color: "hsl(var(--muted-foreground))",
			},
			Tablet: {
				label: "Tablet Users",
				color: "hsl(var(--ring))",
			},
		},
		className: "w-[600px] h-[400px]",
		children: (
			<div className="w-[600px] h-[400px]">
				<section aria-label="Pie chart">
					<RechartsPieChart
						width={600}
						height={400}
						margin={{ top: 16, right: 16, bottom: 16, left: 16 }}
					>
						<Pie
							data={pieData}
							dataKey="value"
							nameKey="name"
							cx="50%"
							cy="50%"
							innerRadius={60}
							outerRadius={80}
							paddingAngle={4}
							fill="currentColor"
							stroke="hsl(var(--background))"
							strokeWidth={2}
						>
							{pieData.map((entry) => (
								<Cell
									key={`cell-${entry.name}`}
									fill={
										entry.name === "Desktop"
											? "hsl(var(--primary))"
											: entry.name === "Mobile"
												? "hsl(var(--muted-foreground))"
												: "hsl(var(--ring))"
									}
									opacity={entry.name === "Mobile" ? 0.5 : 1}
								/>
							))}
						</Pie>
						<ChartTooltip />
					</RechartsPieChart>
				</section>
			</div>
		),
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const chart = await canvas.findByRole("region", { name: /pie chart/i });
		await expect(chart).toBeInTheDocument();
	},
};

export const AreaChartStory: Story = {
	args: {
		config: {
			users: {
				label: "Total Users",
				color: "hsl(var(--primary))",
			},
			activeUsers: {
				label: "Active Users",
				color: "hsl(var(--muted-foreground))",
			},
		},
		className: "w-[600px] h-[400px]",
		children: (
			<div className="w-[600px] h-[400px]">
				<section aria-label="Area chart">
					<RechartsAreaChart
						width={600}
						height={400}
						data={areaData}
						margin={{ top: 16, right: 16, bottom: 16, left: 16 }}
					>
						<XAxis
							dataKey="name"
							stroke="hsl(var(--muted-foreground))"
							fontSize={12}
							tickLine={false}
							axisLine={false}
							dy={10}
						/>
						<YAxis
							stroke="hsl(var(--muted-foreground))"
							fontSize={12}
							tickLine={false}
							axisLine={false}
							dx={-10}
						/>
						<CartesianGrid
							strokeDasharray="4"
							className="stroke-border"
							vertical={false}
						/>
						<Area
							type="monotone"
							dataKey="activeUsers"
							stackId="1"
							stroke="hsl(var(--muted-foreground))"
							fill="hsl(var(--muted-foreground))"
							fillOpacity={0.1}
							strokeWidth={1.5}
						/>
						<Area
							type="monotone"
							dataKey="users"
							stackId="1"
							stroke="hsl(var(--primary))"
							fill="hsl(var(--primary))"
							fillOpacity={0.1}
							strokeWidth={1.5}
						/>
						<ChartTooltip />
					</RechartsAreaChart>
				</section>
			</div>
		),
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const chart = await canvas.findByRole("region", { name: /area chart/i });
		await expect(chart).toBeInTheDocument();

		// Check if axes are visible
		await expect(await canvas.findByText("Mon")).toBeVisible();
		await expect(await canvas.findByText("40")).toBeVisible();
	},
};
