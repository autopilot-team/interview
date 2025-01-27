"use client";

import { buttonVariants } from "@autopilot/ui/components/button";
import { cn } from "@autopilot/ui/lib/utils";
import {
	ChevronDownIcon,
	ChevronLeftIcon,
	ChevronRightIcon,
	ChevronUpIcon,
} from "lucide-react";
import type * as React from "react";
import {
	type ChevronProps,
	DayFlag,
	DayPicker,
	SelectionState,
	UI,
} from "react-day-picker";

export type CalendarProps = React.ComponentProps<typeof DayPicker>;

function Calendar({
	className,
	classNames,
	showOutsideDays = true,
	...props
}: CalendarProps) {
	return (
		<DayPicker
			showOutsideDays={showOutsideDays}
			className={cn("p-3", className)}
			classNames={{
				[UI.Months]: "relative",
				[UI.Month]: "space-y-4 ml-0",
				[UI.MonthCaption]: "flex justify-center items-center h-7",
				[UI.CaptionLabel]: "text-sm font-medium",
				[UI.PreviousMonthButton]: cn(
					buttonVariants({ variant: "outline" }),
					"absolute left-1 top-0 h-7 w-7 bg-transparent p-0 opacity-50 hover:opacity-100",
				),
				[UI.NextMonthButton]: cn(
					buttonVariants({ variant: "outline" }),
					"absolute right-1 top-0 h-7 w-7 bg-transparent p-0 opacity-50 hover:opacity-100",
				),
				[UI.MonthGrid]: "w-full border-collapse space-y-1",
				[UI.Weekdays]: "flex",
				[UI.Weekday]:
					"text-muted-foreground rounded-md w-9 font-normal text-[0.8rem]",
				[UI.Week]: "flex w-full mt-2",
				[UI.Day]:
					"h-9 w-9 text-center text-sm p-0 relative [&:has([aria-selected])]:bg-accent focus-within:relative focus-within:z-20 [&:has([aria-selected]):not(.day-range-end):not(.day-range-start):not(.day-range-middle)]:rounded-md",
				[UI.DayButton]: cn(
					buttonVariants({ variant: "ghost" }),
					"h-9 w-9 p-0 font-normal aria-selected:opacity-100 [&:not([aria-selected])]:hover:bg-primary [&:not([aria-selected])]:hover:text-primary-foreground [&[aria-selected]]:hover:bg-accent [&[aria-selected]]:hover:text-accent-foreground [&:not(.day-range-end):not(.day-range-start):not(.day-range-middle)]:hover:rounded-md [&[aria-selected]:not(.day-range-end):not(.day-range-start):not(.day-range-middle)]:rounded-md",
				),
				[SelectionState.range_end]: "day-range-end rounded-r-md",
				[SelectionState.range_start]: "day-range-start rounded-l-md",
				[SelectionState.selected]:
					"bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground focus:bg-primary focus:text-primary-foreground [&:not(.day-range-end):not(.day-range-start):not(.day-range-middle)]:rounded-md [&_button:not(.day-range-end):not(.day-range-start):not(.day-range-middle)]:rounded-md",
				[SelectionState.range_middle]:
					"aria-selected:bg-accent aria-selected:text-accent-foreground !rounded-none [&_button]:!rounded-none",
				[DayFlag.today]: "bg-accent text-accent-foreground",
				[DayFlag.outside]:
					"day-outside text-muted-foreground opacity-50 aria-selected:bg-accent/50 aria-selected:text-muted-foreground aria-selected:opacity-30",
				[DayFlag.disabled]: "text-muted-foreground opacity-50",
				[DayFlag.hidden]: "invisible",
				...classNames,
			}}
			components={{
				Chevron: ({ ...props }) => <Chevron {...props} />,
			}}
			{...props}
		/>
	);
}

function Chevron({ orientation }: ChevronProps) {
	switch (orientation) {
		case "left":
			return <ChevronLeftIcon className="h-4 w-4" />;
		case "right":
			return <ChevronRightIcon className="h-4 w-4" />;
		case "up":
			return <ChevronUpIcon className="h-4 w-4" />;
		case "down":
			return <ChevronDownIcon className="h-4 w-4" />;
		default:
			return null;
	}
}

Calendar.displayName = "Calendar";

export { Calendar };
