import {
	Collapsible,
	CollapsibleContent,
	CollapsibleTrigger,
} from "@autopilot/ui/components/collapsible";
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
	CommandSeparator,
} from "@autopilot/ui/components/command";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@autopilot/ui/components/dropdown-menu";
import {
	Check,
	ChevronRight,
	ChevronsUpDown,
	Plus,
} from "@autopilot/ui/components/icons";
import {
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarMenuSub,
	useSidebar,
} from "@autopilot/ui/components/sidebar";
import { cn } from "@autopilot/ui/lib/utils";
import * as React from "react";

export interface Entity {
	id: string;
	name: string;
	logo: React.ElementType;
	type: "platform" | "organization" | "account";
	parentId?: string;
}

export interface EntitySwitcherT {
	searchPlaceholder: string;
	noEntities: string;
	noMatchingEntities: string;
	selectEntity: string;
	platforms: string;
	organizations: string;
	accounts: string;
	addEntity: string;
}

export interface EntitySwitcherProps {
	entities: Entity[];
	currentEntity?: Entity;
	onEntityChange?: (entity: Entity | undefined) => void;
	canCreate?: boolean;
	onCreateClick?: () => void;
	t: EntitySwitcherT;
}

interface TreeNode {
	entity: Entity;
	children: TreeNode[];
}

export function EntitySwitcher({
	entities,
	currentEntity,
	onEntityChange,
	canCreate = false,
	onCreateClick,
	t,
}: EntitySwitcherProps) {
	const { isMobile } = useSidebar();
	const [open, setOpen] = React.useState(false);
	const [search, setSearch] = React.useState("");
	const { entityTrees, matchedItems } = React.useMemo(() => {
		const parentMap = new Map(
			entities.map((e) => [e.id, entities.find((p) => p.id === e.parentId)]),
		);
		const nodeMap = new Map<string, TreeNode>();

		// Build trees without filtering
		const trees = {
			platforms: [] as TreeNode[],
			orgs: [] as TreeNode[],
			accounts: [] as TreeNode[],
		};

		// First pass: create all nodes
		for (const entity of entities) {
			nodeMap.set(entity.id, { entity, children: [] });
		}

		// Second pass: build parent-child relationships
		for (const entity of entities) {
			const node = nodeMap.get(entity.id);
			if (!node) continue;

			const parent = parentMap.get(entity.id);

			if (!parent) {
				// Root level node
				if (entity.type === "platform") {
					trees.platforms.push(node);
				} else if (entity.type === "organization") {
					trees.orgs.push(node);
				} else if (entity.type === "account") {
					trees.accounts.push(node);
				}
			} else {
				// Add as child to parent
				const parentNode = nodeMap.get(parent.id);
				if (parentNode) {
					parentNode.children.push(node);
				}
			}
		}

		// Sort all levels by name
		const sortByName = (a: TreeNode, b: TreeNode) =>
			a.entity.name.localeCompare(b.entity.name);

		trees.platforms.sort(sortByName);
		trees.orgs.sort(sortByName);
		trees.accounts.sort(sortByName);

		// Sort children of each node
		const sortChildren = (node: TreeNode) => {
			node.children.sort(sortByName);
			node.children.forEach(sortChildren);
		};

		trees.platforms.forEach(sortChildren);
		trees.orgs.forEach(sortChildren);
		trees.accounts.forEach(sortChildren);

		// Find matched items when searching
		const searchLower = search.toLowerCase();
		const matched = !search
			? []
			: entities
					.filter((entity) => {
						// Direct name match
						if (entity.name.toLowerCase().includes(searchLower)) {
							return true;
						}
						// If this is an organization, check if any of its child accounts match
						if (entity.type === "organization") {
							return entities.some(
								(e) =>
									e.type === "account" &&
									e.parentId === entity.id &&
									e.name.toLowerCase().includes(searchLower),
							);
						}
						// If this is a platform, check if any of its child organizations or accounts match
						if (entity.type === "platform") {
							const childOrgs = entities.filter(
								(e) => e.parentId === entity.id,
							);
							return childOrgs.some((org) => {
								if (org.name.toLowerCase().includes(searchLower)) {
									return true;
								}
								return entities.some(
									(acc) =>
										acc.type === "account" &&
										acc.parentId === org.id &&
										acc.name.toLowerCase().includes(searchLower),
								);
							});
						}
						return false;
					})
					.sort((a, b) => {
						// Sort by type first (platform -> organization -> account)
						const typeOrder = { platform: 0, organization: 1, account: 2 };
						const typeCompare = typeOrder[a.type] - typeOrder[b.type];
						if (typeCompare !== 0) return typeCompare;
						// Then sort by name
						return a.name.localeCompare(b.name);
					});

		return {
			entityTrees: trees,
			matchedItems: matched,
		};
	}, [entities, search]);

	const EntityTreeItem = React.useCallback(
		({ node, depth = 0 }: { node: TreeNode; depth?: number }) => {
			const hasChildren = node.children?.length > 0;
			const active = currentEntity?.id === node.entity.id;

			if (!hasChildren) {
				return (
					<CommandItem
						value={node.entity.name}
						onSelect={() => {
							onEntityChange?.(node.entity);
							setOpen(false);
						}}
						className={cn(
							"gap-2 w-full px-2",
							active && "bg-accent text-accent-foreground",
							"hover:bg-accent/50 hover:text-accent-foreground",
						)}
					>
						<div
							className={cn(
								"flex size-6 items-center justify-center rounded-sm border",
								active && "border-accent-foreground/50 bg-accent-foreground/10",
								!active && "border-border",
							)}
						>
							<node.entity.logo className="size-4 shrink-0" />
						</div>
						<span className="flex-1 truncate">{node.entity.name}</span>
						{active && <Check className="ml-auto size-4 opacity-70" />}
					</CommandItem>
				);
			}

			return (
				<Collapsible defaultOpen={true} className="group/collapsible ml-0.5">
					<CollapsibleTrigger asChild>
						<CommandItem
							value={node.entity.name}
							onSelect={() => {
								onEntityChange?.(node.entity);
								setOpen(false);
							}}
							className={cn(
								"gap-2 w-full px-2",
								active && "bg-accent text-accent-foreground",
								"hover:bg-accent/50 hover:text-accent-foreground",
							)}
						>
							<ChevronRight className="size-4 transition-transform shrink-0 group-data-[state=open]/collapsible:rotate-90" />
							<div
								className={cn(
									"flex size-6 items-center justify-center rounded-sm border",
									active &&
										"border-accent-foreground/50 bg-accent-foreground/10",
									!active && "border-border",
								)}
							>
								<node.entity.logo className="size-4 shrink-0" />
							</div>
							<span className="flex-1 truncate">{node.entity.name}</span>
							{active && <Check className="ml-auto size-4 opacity-70" />}
						</CommandItem>
					</CollapsibleTrigger>

					<CollapsibleContent>
						<SidebarMenuSub className="mr-0 pr-0">
							{node.children?.map((child) => (
								<EntityTreeItem
									key={child.entity.id}
									node={child}
									depth={depth + 1}
								/>
							))}
						</SidebarMenuSub>
					</CollapsibleContent>
				</Collapsible>
			);
		},
		[currentEntity, onEntityChange],
	);

	const MatchedItem = React.useCallback(
		({ entity }: { entity: Entity }) => {
			const active = currentEntity?.id === entity.id;

			return (
				<CommandItem
					value={entity.name}
					onSelect={() => {
						onEntityChange?.(entity);
						setOpen(false);
					}}
					className={cn(
						"gap-2 w-full px-2",
						active && "bg-accent text-accent-foreground",
						"hover:bg-accent/50 hover:text-accent-foreground",
					)}
				>
					<div
						className={cn(
							"flex size-6 items-center justify-center rounded-sm border",
							active && "border-accent-foreground/50 bg-accent-foreground/10",
							!active && "border-border",
						)}
					>
						<entity.logo className="size-4 shrink-0" />
					</div>
					<div className="flex-1 min-w-0">
						<div className="truncate">{entity.name}</div>
						<div className="text-xs text-muted-foreground capitalize">
							{entity.type}
						</div>
					</div>
					{active && <Check className="ml-auto size-4 opacity-70" />}
				</CommandItem>
			);
		},
		[currentEntity, onEntityChange],
	);

	return (
		<SidebarMenu>
			<SidebarMenuItem>
				<DropdownMenu open={open} onOpenChange={setOpen}>
					<DropdownMenuTrigger asChild>
						<SidebarMenuButton
							size="lg"
							className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
							aria-label={t.selectEntity}
						>
							{currentEntity && (
								<>
									<div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
										<currentEntity.logo className="size-4" />
									</div>

									<div className="grid flex-1 text-left text-sm leading-tight">
										<span className="truncate font-semibold">
											{currentEntity.name}
										</span>

										<span className="truncate text-xs capitalize text-muted-foreground">
											{currentEntity.type}
										</span>
									</div>

									<ChevronsUpDown className="ml-auto size-4 opacity-50" />
								</>
							)}
						</SidebarMenuButton>
					</DropdownMenuTrigger>

					<DropdownMenuContent
						className="w-[--radix-dropdown-menu-trigger-width] min-w-[280px] rounded-lg"
						align="start"
						side={isMobile ? "bottom" : "right"}
						sideOffset={4}
					>
						<Command>
							<CommandInput
								placeholder={t.searchPlaceholder}
								value={search}
								onValueChange={setSearch}
							/>

							<CommandList className="mt-1">
								{!search &&
									entityTrees.platforms.length === 0 &&
									entityTrees.orgs.length === 0 &&
									entityTrees.accounts.length === 0 && (
										<CommandEmpty>{t.noEntities}</CommandEmpty>
									)}

								{search && matchedItems.length === 0 && (
									<CommandEmpty>{t.noMatchingEntities}</CommandEmpty>
								)}

								{!search ? (
									<>
										{entityTrees.platforms.length > 0 && (
											<CommandGroup
												heading={t.platforms}
												className="[&>h3]:mb-2 [&>h3]:ml-2 [&>h3]:select-none [&>div]:list-none p-0 mb-2"
											>
												{entityTrees.platforms.map((node) => (
													<EntityTreeItem key={node.entity.id} node={node} />
												))}
											</CommandGroup>
										)}

										{entityTrees.orgs.length > 0 && (
											<>
												<CommandSeparator />
												<CommandGroup
													heading={t.organizations}
													className="[&>h3]:mb-2 [&>h3]:ml-2 [&>h3]:select-none [&>div]:list-none p-0 mb-2"
												>
													{entityTrees.orgs.map((node) => (
														<EntityTreeItem key={node.entity.id} node={node} />
													))}
												</CommandGroup>
											</>
										)}

										{entityTrees.accounts.length > 0 && (
											<>
												<CommandSeparator />
												<CommandGroup
													heading={t.accounts}
													className="[&>h3]:mb-2 [&>h3]:ml-2 [&>h3]:select-none [&>div]:list-none p-0"
												>
													{entityTrees.accounts.map((node) => (
														<EntityTreeItem key={node.entity.id} node={node} />
													))}
												</CommandGroup>
											</>
										)}
									</>
								) : (
									<CommandGroup className="p-0 [&_[cmdk-group-items]]:space-y-0.5">
										{matchedItems.map((entity) => (
											<MatchedItem key={entity.id} entity={entity} />
										))}
									</CommandGroup>
								)}
							</CommandList>
						</Command>

						{canCreate && (
							<>
								<DropdownMenuSeparator />
								<DropdownMenuItem
									className="gap-2 p-2"
									onClick={() => {
										onCreateClick?.();
										setOpen(false);
									}}
								>
									<div className="flex size-6 items-center justify-center rounded-md border bg-background">
										<Plus className="size-4" />
									</div>
									<div className="font-medium text-muted-foreground">
										{t.addEntity}
									</div>
								</DropdownMenuItem>
							</>
						)}
					</DropdownMenuContent>
				</DropdownMenu>
			</SidebarMenuItem>
		</SidebarMenu>
	);
}
