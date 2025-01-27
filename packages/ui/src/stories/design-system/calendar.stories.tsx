import type { Meta, StoryObj } from "@storybook/react";
import { expect, userEvent, within } from "@storybook/test";
import { addDays } from "date-fns";
import { useState } from "react";
import type { DateRange as ReactDayPickerDateRange } from "react-day-picker";
import { Calendar } from "../../components/calendar.js";

const meta: Meta<typeof Calendar> = {
	title: "Design System/Calendar",
	component: Calendar,
	parameters: {
		layout: "centered",
	},
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof Calendar>;

const CalendarDemo = () => {
	const [date, setDate] = useState<Date | undefined>(new Date());

	return (
		<Calendar
			mode="single"
			selected={date}
			onSelect={setDate}
			className="rounded-md border"
		/>
	);
};

export const Default: Story = {
	render: () => <CalendarDemo />,
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const calendar = canvas.getByRole("grid");
		await expect(calendar).toBeInTheDocument();

		const today = canvas.getByRole("gridcell", { selected: true });
		await expect(today).toBeInTheDocument();

		const nextMonthButton = canvas.getByRole("button", { name: /next month/i });
		await userEvent.click(nextMonthButton);

		// Test date selection
		const dateCell = canvas.getByRole("gridcell", { name: /15/i });
		const dateButton = within(dateCell).getByRole("button");
		await userEvent.click(dateButton);
		await expect(dateCell).toHaveAttribute("aria-selected", "true");
	},
};

const DateRangeDemo = () => {
	const [dateRange, setDateRange] = useState<
		ReactDayPickerDateRange | undefined
	>({
		from: new Date(),
		to: addDays(new Date(), 7),
	});

	return (
		<Calendar
			mode="range"
			selected={dateRange}
			onSelect={setDateRange}
			className="rounded-md border"
		/>
	);
};

export const DateRange: Story = {
	render: () => <DateRangeDemo />,
};

const MultipleSelectionDemo = () => {
	const [dates, setDates] = useState<Date[]>([
		new Date(),
		addDays(new Date(), 2),
		addDays(new Date(), 5),
	]);

	return (
		<Calendar
			mode="multiple"
			selected={dates}
			onSelect={(dates) => setDates(dates || [])}
			className="rounded-md border"
		/>
	);
};

export const MultipleSelection: Story = {
	render: () => <MultipleSelectionDemo />,
};

export const WithFooter: Story = {
	render: () => {
		const [date, setDate] = useState<Date | undefined>(new Date());

		return (
			<div className="space-y-4">
				<Calendar
					mode="single"
					selected={date}
					onSelect={setDate}
					className="rounded-md border"
					footer={
						<div className="mt-3 flex justify-center text-sm text-muted-foreground">
							Click to select a date
						</div>
					}
				/>
			</div>
		);
	},
};

export const Disabled: Story = {
	args: {
		mode: "single",
		selected: new Date(),
		className: "rounded-md border",
		disabled: true,
	},
};

export const DisabledDates: Story = {
	render: () => {
		const [date, setDate] = useState<Date | undefined>(new Date());

		return (
			<Calendar
				mode="single"
				selected={date}
				onSelect={setDate}
				className="rounded-md border"
				disabled={[
					{ from: addDays(new Date(), 1), to: addDays(new Date(), 4) },
					{ from: addDays(new Date(), 8), to: addDays(new Date(), 11) },
				]}
				defaultMonth={new Date()}
			/>
		);
	},
};
