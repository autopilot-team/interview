import { Button } from "@autopilot/ui/components/button";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@autopilot/ui/components/select";
import type { Table } from "@tanstack/react-table";
import { ChevronLeft, ChevronRight } from "lucide-react";

export interface DataTablePaginationTranslations {
	rowsPerPage: string;
	of: string;
	rowsSelected: string;
	page: string;
	navigation: {
		first?: string;
		previous: string;
		next: string;
		last?: string;
	};
}

interface DataTablePaginationProps<TData> {
	table: Table<TData>;
	pageSizeOptions?: number[];
	t: DataTablePaginationTranslations;
}

export function DataTablePagination<TData>({
	table,
	pageSizeOptions = [10, 20, 30, 40, 50],
	t,
}: DataTablePaginationProps<TData>) {
	const selectedRows = table.getFilteredSelectedRowModel().rows.length;
	const totalRows = table.getFilteredRowModel().rows.length;

	return (
		<div className="flex items-center justify-between px-2">
			<div className="flex-1 text-sm text-muted-foreground">
				{table.options.enableRowSelection && selectedRows > 0
					? t.rowsSelected
							.replace("{{count}}", String(selectedRows))
							.replace("{{total}}", String(totalRows))
					: null}
			</div>

			<div className="flex items-center space-x-6 lg:space-x-8">
				<div className="flex items-center space-x-2">
					<p className="text-sm font-medium">{t.rowsPerPage}</p>
					<Select
						value={`${table.getState().pagination.pageSize}`}
						onValueChange={(value) => {
							table.setPageSize(Number(value));
						}}
					>
						<SelectTrigger className="h-8 w-[70px]">
							<SelectValue placeholder={table.getState().pagination.pageSize} />
						</SelectTrigger>

						<SelectContent side="top">
							{pageSizeOptions.map((pageSize) => (
								<SelectItem key={pageSize} value={`${pageSize}`}>
									{pageSize}
								</SelectItem>
							))}
						</SelectContent>
					</Select>
				</div>

				<div className="flex w-[100px] items-center justify-center text-sm font-medium">
					{t.page
						.replace(
							"{{current}}",
							String(table.getState().pagination.pageIndex + 1),
						)
						.replace("{{total}}", String(table.getPageCount()))}
				</div>

				<div className="flex items-center space-x-2">
					<Button
						variant="outline"
						className="size-8 p-0"
						onClick={() => table.previousPage()}
						disabled={!table.getCanPreviousPage()}
					>
						<span className="sr-only">{t.navigation.previous}</span>
						<ChevronLeft />
					</Button>

					<Button
						variant="outline"
						className="size-8 p-0"
						onClick={() => table.nextPage()}
						disabled={!table.getCanNextPage()}
					>
						<span className="sr-only">{t.navigation.next}</span>
						<ChevronRight />
					</Button>
				</div>
			</div>
		</div>
	);
}
