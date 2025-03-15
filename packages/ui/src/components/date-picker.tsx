import { Button } from "@autopilot/ui/components/button";
import { Calendar } from "@autopilot/ui/components/calendar";
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from "@autopilot/ui/components/popover";
import { cn } from "@autopilot/ui/lib/utils";
import { format } from "date-fns";
import { Calendar as CalendarIcon, X as ClearIcon } from "lucide-react";
import { useState } from "react";
import type {
	DateRange,
	PropsMulti,
	PropsRange,
	PropsSingle,
} from "react-day-picker";

export type DatePickerMode = "single" | "multiple" | "range";

type DatePickerBaseProps = {
	/**
	 * The placeholder text to show when no date is selected
	 */
	placeholder?: string;
	/**
	 * The format to display the selected date in
	 * @default "PPP"
	 */
	dateFormat?: string;
	/**
	 * The className to apply to the trigger button
	 */
	className?: string;
	/**
	 * The earliest date that can be selected
	 */
	fromDate?: Date;
	/**
	 * The latest date that can be selected
	 */
	toDate?: Date;
	/**
	 * Whether to show the clear button
	 * @default true
	 */
	allowClear?: boolean;
};

interface DatePickerSingleProps
	extends DatePickerBaseProps,
		Omit<
			PropsSingle,
			"mode" | "selected" | "onSelect" | "fromDate" | "toDate"
		> {
	mode?: "single";
	selected?: Date;
	onSelect?: (date: Date | undefined) => void;
}

interface DatePickerMultipleProps
	extends DatePickerBaseProps,
		Omit<PropsMulti, "mode" | "selected" | "onSelect" | "fromDate" | "toDate"> {
	mode: "multiple";
	selected?: Date[];
	onSelect?: (dates: Date[] | undefined) => void;
}

interface DatePickerRangeProps
	extends DatePickerBaseProps,
		Omit<PropsRange, "mode" | "selected" | "onSelect" | "fromDate" | "toDate"> {
	mode: "range";
	selected?: DateRange;
	onSelect?: (range: DateRange | undefined) => void;
}

export type DatePickerProps =
	| DatePickerSingleProps
	| DatePickerMultipleProps
	| DatePickerRangeProps;

function SingleCalendar({
	selected,
	onSelect,
	fromDate,
	toDate,
	...props
}: Omit<PropsSingle, "mode"> &
	Pick<DatePickerBaseProps, "fromDate" | "toDate">) {
	return (
		<Calendar
			mode="single"
			selected={selected}
			onSelect={onSelect}
			fromDate={fromDate}
			toDate={toDate}
			{...props}
		/>
	);
}

function MultipleCalendar({
	selected,
	onSelect,
	fromDate,
	toDate,
	...props
}: Omit<PropsMulti, "mode"> &
	Pick<DatePickerBaseProps, "fromDate" | "toDate">) {
	return (
		<Calendar
			mode="multiple"
			selected={selected}
			onSelect={onSelect}
			fromDate={fromDate}
			toDate={toDate}
			{...props}
		/>
	);
}

function RangeCalendar({
	selected,
	onSelect,
	fromDate,
	toDate,
	...props
}: Omit<PropsRange, "mode"> &
	Pick<DatePickerBaseProps, "fromDate" | "toDate">) {
	return (
		<Calendar
			mode="range"
			selected={selected}
			onSelect={onSelect}
			fromDate={fromDate}
			toDate={toDate}
			{...props}
		/>
	);
}

export function DatePicker({
	placeholder = "Pick a date",
	dateFormat = "PPP",
	className,
	mode = "single",
	selected,
	onSelect,
	fromDate,
	toDate,
	allowClear = true,
	...calendarProps
}: DatePickerProps) {
	const [date, setDate] = useState<Date | undefined>(
		mode === "single" ? (selected as Date) : undefined,
	);
	const [dates, setDates] = useState<Date[]>(
		mode === "multiple" ? (selected as Date[]) || [] : [],
	);
	const [range, setRange] = useState<DateRange | undefined>(
		mode === "range" ? (selected as DateRange) : undefined,
	);

	const handleSingleSelect = (value: Date | undefined) => {
		setDate(value);
		if (mode === "single" && onSelect) {
			(onSelect as (date: Date | undefined) => void)(value);
		}
	};

	const handleMultipleSelect = (value: Date[] | undefined) => {
		const selectedDates = value || [];
		setDates(selectedDates);
		if (mode === "multiple" && onSelect) {
			(onSelect as (dates: Date[] | undefined) => void)(selectedDates);
		}
	};

	const handleRangeSelect = (value: DateRange | undefined) => {
		setRange(value);
		if (mode === "range" && onSelect) {
			(onSelect as (range: DateRange | undefined) => void)(value);
		}
	};

	const handleClear = (e: React.MouseEvent) => {
		e.stopPropagation();
		if (mode === "single") {
			handleSingleSelect(undefined);
		} else if (mode === "multiple") {
			handleMultipleSelect([]);
		} else {
			handleRangeSelect(undefined);
		}
	};

	const formatDate = (date: Date | undefined) => {
		if (!date) return "";
		return format(date, dateFormat);
	};

	const getDisplayText = () => {
		if (mode === "single" && date) {
			return formatDate(date);
		}

		if (mode === "multiple" && dates.length > 0) {
			return dates.length === 1
				? formatDate(dates[0])
				: `${dates.length} dates selected`;
		}

		if (mode === "range" && range?.from) {
			if (!range.to) {
				return `From ${formatDate(range.from)}`;
			}

			return `${formatDate(range.from)} - ${formatDate(range.to)}`;
		}

		return <span className="text-muted-foreground">{placeholder}</span>;
	};

	const hasValue =
		mode === "single"
			? !!date
			: mode === "multiple"
				? dates.length > 0
				: !!range?.from;

	const renderCalendar = () => {
		if (mode === "multiple") {
			return (
				<MultipleCalendar
					selected={dates}
					onSelect={handleMultipleSelect}
					fromDate={fromDate}
					toDate={toDate}
					{...calendarProps}
				/>
			);
		}

		if (mode === "range") {
			return (
				<RangeCalendar
					selected={range}
					onSelect={handleRangeSelect}
					fromDate={fromDate}
					toDate={toDate}
					{...calendarProps}
				/>
			);
		}

		return (
			<SingleCalendar
				selected={date}
				onSelect={handleSingleSelect}
				fromDate={fromDate}
				toDate={toDate}
				{...calendarProps}
			/>
		);
	};

	return (
		<Popover>
			<PopoverTrigger asChild>
				<Button
					variant="outline"
					className={cn(
						"flex h-10 w-full items-center justify-between text-left font-normal",
						!hasValue && "text-muted-foreground",
						className,
					)}
				>
					<span className="flex items-center gap-2">
						<CalendarIcon className="h-4 w-4" />
						{getDisplayText()}
					</span>

					{allowClear && hasValue && (
						<Button
							variant="ghost"
							size="icon"
							className="h-4 w-4 p-0 hover:bg-transparent"
							onClick={handleClear}
						>
							<ClearIcon className="h-3 w-3" />
						</Button>
					)}
				</Button>
			</PopoverTrigger>

			<PopoverContent className="w-auto p-0" align="start">
				{renderCalendar()}
			</PopoverContent>
		</Popover>
	);
}
