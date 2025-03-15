import { Button } from "@autopilot/ui/components/button";
import { Checkbox } from "@autopilot/ui/components/checkbox";
import {
	DropdownMenu,
	DropdownMenuCheckboxItem,
	DropdownMenuContent,
	DropdownMenuTrigger,
} from "@autopilot/ui/components/dropdown-menu";
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@autopilot/ui/components/table";
import {
	type ColumnDef,
	type ColumnFiltersState,
	type SortingState,
	type Table as TableType,
	type VisibilityState,
	flexRender,
	getCoreRowModel,
	getFacetedRowModel,
	getFacetedUniqueValues,
	getFilteredRowModel,
	getPaginationRowModel,
	getSortedRowModel,
	useReactTable,
} from "@tanstack/react-table";
import { LayoutGridIcon } from "lucide-react";
import { ChevronLeft, ChevronRight } from "lucide-react";
import * as React from "react";
import {
	DataTablePagination,
	type DataTablePaginationTranslations,
} from "./data-table-pagination.js";

interface DataTableT {
	noResults: string;
	loading: string;
	columns: string;
	cursorPagination?: DataTablePaginationTranslations;
	offsetPagination?: DataTablePaginationTranslations;
}

interface OffsetPaginationConfig {
	pageSizeOptions?: number[];
	defaultPageSize?: number;
}

interface CursorPaginationConfig {
	after?: string | null;
	before?: string | null;
}

interface DataTableProps<TData, TValue> {
	className?: string;
	columns: ColumnDef<TData, TValue>[];
	cursorPagination?: CursorPaginationConfig;
	data: TData[];
	emptyState?: React.ReactNode;
	enableRowSelection?: boolean;
	enableHiding?: boolean;
	enableSorting?: boolean;
	enableFiltering?: boolean;
	isLoading?: boolean;
	offsetPagination?: OffsetPaginationConfig;
	searchParams?: URLSearchParams;
	setSearchParams?: (params: URLSearchParams) => void;
	t: DataTableT;
}

function DataTableViewOptions<TData>({
	table,
	t,
}: {
	table: TableType<TData>;
	t: string;
}) {
	return (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button variant="outline" size="sm" className="h-8 px-3 border-dashed">
					<LayoutGridIcon className="mr-2 size-3.5" />
					<span className="text-sm">{t}</span>
				</Button>
			</DropdownMenuTrigger>

			<DropdownMenuContent align="end" className="w-[200px]">
				{table
					.getAllColumns()
					.filter((column) => column.getCanHide())
					.map((column) => {
						const headerContent =
							column.columnDef.header?.toString() || column.id;

						return (
							<DropdownMenuCheckboxItem
								key={column.id}
								className="capitalize"
								checked={column.getIsVisible()}
								onCheckedChange={(value) => column.toggleVisibility(!!value)}
							>
								{headerContent}
							</DropdownMenuCheckboxItem>
						);
					})}
			</DropdownMenuContent>
		</DropdownMenu>
	);
}

export function DataTable<TData, TValue>({
	className,
	columns: userColumns,
	data,
	isLoading = false,
	emptyState,
	offsetPagination = {},
	cursorPagination,
	enableRowSelection = false,
	enableHiding = false,
	enableSorting = true,
	enableFiltering = true,
	searchParams,
	setSearchParams,
	t,
}: DataTableProps<TData, TValue>) {
	const isCursorPagination = !!cursorPagination;
	const {
		defaultPageSize = offsetPagination?.defaultPageSize || 10,
		pageSizeOptions = offsetPagination?.pageSizeOptions || [10, 20, 30, 40, 50],
	} = offsetPagination;
	const [rowSelection, setRowSelection] = React.useState({});
	const [columnVisibility, setColumnVisibility] =
		React.useState<VisibilityState>({});
	const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
		[],
	);
	const [sorting, setSorting] = React.useState<SortingState>([]);
	const [{ pageIndex, pageSize }, setPagination] = React.useState({
		pageIndex: 0,
		pageSize: defaultPageSize,
	});
	const updateSearchParams = React.useCallback(
		(updates: Record<string, string | null>) => {
			if (!setSearchParams || !searchParams) return;

			const newParams = new URLSearchParams(searchParams);
			for (const [key, value] of Object.entries(updates)) {
				if (value === null) {
					newParams.delete(key);
				} else {
					newParams.set(key, value);
				}
			}

			setSearchParams(newParams);
		},
		[searchParams, setSearchParams],
	);
	const paginationState = React.useMemo(
		() => ({
			pageIndex,
			pageSize,
		}),
		[pageIndex, pageSize],
	);
	const columns = React.useMemo(() => {
		if (!enableRowSelection) return userColumns;

		const selectionColumn: ColumnDef<TData, TValue> = {
			id: "select",
			header: ({ table }) => (
				<Checkbox
					checked={table.getIsAllPageRowsSelected()}
					onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
					aria-label="Select all"
				/>
			),
			cell: ({ row }) => (
				<Checkbox
					checked={row.getIsSelected()}
					onCheckedChange={(value) => row.toggleSelected(!!value)}
					aria-label="Select row"
				/>
			),
			enableSorting: false,
			enableHiding: false,
		};

		return [selectionColumn, ...userColumns];
	}, [enableRowSelection, userColumns]);
	const table = useReactTable({
		data,
		columns,
		state: {
			sorting,
			columnVisibility,
			rowSelection,
			columnFilters,
			pagination: paginationState,
		},
		enableRowSelection,
		enableColumnFilters: enableFiltering,
		enableSorting,
		enableHiding,
		onRowSelectionChange: setRowSelection,
		onColumnFiltersChange: setColumnFilters,
		onColumnVisibilityChange: setColumnVisibility,
		onPaginationChange: (updaterOrValue) => {
			setPagination(updaterOrValue);

			if (setSearchParams) {
				const newPagination =
					typeof updaterOrValue === "function"
						? updaterOrValue(paginationState)
						: updaterOrValue;

				if (!isCursorPagination) {
					updateSearchParams({
						page: String(newPagination.pageIndex + 1),
						pageSize: String(newPagination.pageSize),
					});
				}
			}
		},
		onSortingChange: (updaterOrValue) => {
			setSorting(updaterOrValue);

			if (!setSearchParams) return;

			const newSorting =
				typeof updaterOrValue === "function"
					? updaterOrValue(sorting)
					: updaterOrValue;

			if (newSorting.length > 0) {
				updateSearchParams({
					sortColumn: newSorting[0]?.id || null,
					sortDir: newSorting[0]?.desc ? "desc" : "asc",
				});
			} else {
				updateSearchParams({
					sortColumn: null,
					sortDir: null,
				});
			}
		},
		getCoreRowModel: getCoreRowModel(),
		getFilteredRowModel: getFilteredRowModel(),
		getPaginationRowModel: isCursorPagination
			? undefined
			: getPaginationRowModel(),
		getSortedRowModel: getSortedRowModel(),
		getFacetedRowModel: getFacetedRowModel(),
		getFacetedUniqueValues: getFacetedUniqueValues(),
		manualPagination: isCursorPagination,
		pageCount: isCursorPagination ? 2 : undefined, // Need at least 2 for next/prev to work
	});

	React.useEffect(() => {
		if (!searchParams || isCursorPagination) return;

		const pageParam = searchParams.get("page");
		const pageSizeParam = searchParams.get("pageSize");
		if (pageParam || pageSizeParam) {
			let validPageIndex = 0;
			if (pageParam) {
				const parsedPageIndex = Number.parseInt(pageParam, 10);

				if (!Number.isNaN(parsedPageIndex) && parsedPageIndex > 0) {
					validPageIndex = parsedPageIndex - 1;
				}
			}

			let validPageSize = defaultPageSize;
			if (pageSizeParam) {
				const parsedPageSize = Number.parseInt(pageSizeParam, 10);

				if (!Number.isNaN(parsedPageSize) && parsedPageSize > 0) {
					validPageSize =
						pageSizeOptions.length > 0
							? pageSizeOptions.includes(parsedPageSize)
								? parsedPageSize
								: defaultPageSize
							: parsedPageSize;
				}
			}

			// Calculate total pages for offset pagination to cap page index
			const totalPages = Math.ceil(data.length / validPageSize);
			if (totalPages > 0 && validPageIndex >= totalPages) {
				validPageIndex = totalPages - 1;
			}

			setPagination({
				pageIndex: validPageIndex,
				pageSize: validPageSize,
			});

			const needsParamUpdate =
				(pageParam && validPageIndex + 1 !== Number.parseInt(pageParam, 10)) ||
				(pageSizeParam && validPageSize !== Number.parseInt(pageSizeParam, 10));

			if (needsParamUpdate) {
				updateSearchParams({
					page: String(validPageIndex + 1),
					pageSize: String(validPageSize),
				});
			}
		}

		const sortColumn = searchParams.get("sortColumn");
		const sortDir = searchParams.get("sortDir");
		if (sortColumn && sortDir) {
			const validSortDir =
				sortDir === "desc" || sortDir === "asc" ? sortDir : "asc";

			setSorting([
				{
					id: sortColumn,
					desc: validSortDir === "desc",
				},
			]);

			if (validSortDir !== sortDir) {
				updateSearchParams({
					sortColumn,
					sortDir: validSortDir,
				});
			}
		}

		const visibilityParam = searchParams.get("columnVisibility");
		if (visibilityParam) {
			try {
				const visibilityState = JSON.parse(visibilityParam);
				setColumnVisibility(visibilityState);
			} catch (e) {
				console.error("Failed to parse column visibility from URL", e);
			}
		}

		const filtersParam = searchParams.get("filters");
		if (filtersParam) {
			try {
				const filtersState = JSON.parse(filtersParam);
				setColumnFilters(filtersState);
			} catch (e) {
				console.error("Failed to parse filters from URL", e);
			}
		}
	}, [
		defaultPageSize,
		pageSizeOptions,
		searchParams,
		updateSearchParams,
		isCursorPagination,
		data,
	]);

	// Initialize search params if they're not set
	React.useEffect(() => {
		if (!searchParams || !setSearchParams || isCursorPagination) return;

		const updates: Record<string, string | null> = {};
		let hasUpdates = false;

		if (isCursorPagination) {
			if (!searchParams.has("pageSize")) {
				updates.pageSize = String(pageSize);
				hasUpdates = true;
			}

			if (searchParams.has("page")) {
				updates.page = null;
				hasUpdates = true;
			}
		} else {
			if (!searchParams.has("page")) {
				updates.page = String(pageIndex + 1);
				hasUpdates = true;
			}

			if (!searchParams.has("pageSize")) {
				updates.pageSize = String(pageSize);
				hasUpdates = true;
			}

			if (searchParams.has("after") || searchParams.has("before")) {
				updates.after = null;
				updates.before = null;
				hasUpdates = true;
			}
		}

		if (
			sorting.length > 0 &&
			(!searchParams.has("sortColumn") || !searchParams.has("sortDir"))
		) {
			updates.sortColumn = sorting[0]?.id || "";
			updates.sortDir = sorting[0]?.desc ? "desc" : "asc";
			hasUpdates = true;
		}

		if (hasUpdates) {
			updateSearchParams(updates);
		}
	}, [
		searchParams,
		setSearchParams,
		pageIndex,
		pageSize,
		sorting,
		updateSearchParams,
		isCursorPagination,
	]);

	return (
		<div className={className}>
			<div className="flex items-center justify-between py-4">
				<div className="flex-1" />
				{enableHiding && <DataTableViewOptions table={table} t={t.columns} />}
			</div>

			<div className="rounded-md border border-border/90">
				<Table>
					<TableHeader>
						{table.getHeaderGroups().map((headerGroup) => (
							<TableRow key={headerGroup.id} className="bg-muted/50">
								{headerGroup.headers.map((header) => {
									return (
										<TableHead
											key={header.id}
											colSpan={header.colSpan}
											className="bg-zinc-100 h-14 text-sm font-semibold text-gray-700"
											style={{
												width: header.getSize(),
												minWidth: header.getSize(),
											}}
										>
											{header.isPlaceholder
												? null
												: flexRender(
														header.column.columnDef.header,
														header.getContext(),
													)}
										</TableHead>
									);
								})}
							</TableRow>
						))}
					</TableHeader>

					<TableBody>
						{isLoading ? (
							<TableRow>
								<TableCell
									colSpan={columns.length}
									className="h-24 text-center"
								>
									{t.loading}
								</TableCell>
							</TableRow>
						) : table.getRowModel().rows?.length ? (
							table.getRowModel().rows.map((row) => (
								<TableRow
									key={row.id}
									data-state={row.getIsSelected() && "selected"}
									className="hover:bg-muted/40"
								>
									{row.getVisibleCells().map((cell) => (
										<TableCell
											key={cell.id}
											className="py-3"
											style={{
												width: cell.column.getSize(),
												minWidth: cell.column.getSize(),
											}}
										>
											{flexRender(
												cell.column.columnDef.cell,
												cell.getContext(),
											)}
										</TableCell>
									))}
								</TableRow>
							))
						) : (
							<TableRow>
								<TableCell
									colSpan={columns.length}
									className="h-24 text-center"
								>
									{emptyState || t.noResults}
								</TableCell>
							</TableRow>
						)}
					</TableBody>
				</Table>
			</div>

			<div className="mt-4">
				{isCursorPagination ? (
					<div className="flex items-center justify-between px-2">
						<div className="flex-1 text-sm text-muted-foreground">
							{table.options.enableRowSelection &&
							table.getFilteredSelectedRowModel().rows.length > 0
								? (
										t.cursorPagination?.rowsSelected ||
										"{{count}} of {{total}} row(s) selected"
									)
										.replace(
											"{{count}}",
											String(table.getFilteredSelectedRowModel().rows.length),
										)
										.replace(
											"{{total}}",
											String(table.getFilteredRowModel().rows.length),
										)
								: null}
						</div>

						<div className="flex items-center gap-4">
							<div className="flex items-center gap-2">
								<Button
									variant="outline"
									className="h-8 w-8 p-0"
									onClick={() => {
										const params = new URLSearchParams(searchParams);

										if (cursorPagination?.before) {
											params.set("before", cursorPagination.before);
											params.delete("after");
											setSearchParams?.(params);
										}
									}}
									disabled={!cursorPagination?.before}
								>
									<span className="sr-only">
										{t.cursorPagination?.navigation?.previous}
									</span>
									<ChevronLeft />
								</Button>

								<Button
									variant="outline"
									className="h-8 w-8 p-0"
									onClick={() => {
										const params = new URLSearchParams(searchParams);

										if (cursorPagination?.after) {
											params.set("after", cursorPagination.after);
											params.delete("before");
											setSearchParams?.(params);
										}
									}}
									disabled={!cursorPagination?.after}
								>
									<span className="sr-only">
										{t.cursorPagination?.navigation?.next}
									</span>
									<ChevronRight />
								</Button>
							</div>
						</div>
					</div>
				) : (
					<DataTablePagination
						table={table}
						pageSizeOptions={offsetPagination.pageSizeOptions}
						t={
							t.offsetPagination || {
								rowsPerPage: "",
								of: "",
								rowsSelected: "",
								page: "",
								navigation: {
									first: "",
									previous: "",
									next: "",
									last: "",
								},
							}
						}
					/>
				)}
			</div>
		</div>
	);
}
