import { Button } from "@autopilot/ui/components/button";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@autopilot/ui/components/dropdown-menu";
import type { Row } from "@tanstack/react-table";
import { MoreHorizontalIcon } from "lucide-react";

export interface Action<TData> {
	label: string;
	onClick: (row: Row<TData>) => void;
	icon?: React.ReactNode;
	variant?: "default" | "destructive";
	separator?: "before" | "after";
}

interface DataTableRowActionsProps<TData> {
	row: Row<TData>;
	actions: Action<TData>[];
}

export function DataTableRowActions<TData>({
	row,
	actions,
}: DataTableRowActionsProps<TData>) {
	return (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button
					variant="ghost"
					size="icon"
					className="flex h-8 w-8 p-0 data-[state=open]:bg-muted"
				>
					<MoreHorizontalIcon className="size-4" />
					<span className="sr-only">Open menu</span>
				</Button>
			</DropdownMenuTrigger>

			<DropdownMenuContent align="end" className="w-[160px]">
				{actions.map((action, i) => (
					<div key={`${action.label}-${i}`}>
						{action.separator === "before" && <DropdownMenuSeparator />}

						<DropdownMenuItem
							onClick={() => action.onClick(row)}
							className={
								action.variant === "destructive" ? "text-destructive" : ""
							}
						>
							{action.icon && (
								<span className="mr-2 size-4">{action.icon}</span>
							)}
							{action.label}
						</DropdownMenuItem>

						{action.separator === "after" && <DropdownMenuSeparator />}
					</div>
				))}
			</DropdownMenuContent>
		</DropdownMenu>
	);
}
