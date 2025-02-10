"use client";

import {
	Avatar,
	AvatarFallback,
	AvatarImage,
} from "@autopilot/ui/components/avatar";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuGroup,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@autopilot/ui/components/dropdown-menu";
import { ChevronsUpDown, LogOut } from "@autopilot/ui/components/icons";
import type { LucideIcon } from "@autopilot/ui/components/icons";
import {
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	useSidebar,
} from "@autopilot/ui/components/sidebar";
import { Link } from "react-router";

export interface NavUserItem {
	title: string;
	icon: LucideIcon;
	url?: string;
	onClick?: () => void;
}

export interface NavUserT {
	signOut: string;
}

export function NavUser({
	user,
	t,
	onSignOutClick,
	items = [],
}: {
	user: {
		name: string;
		email: string;
		image?: string;
	};
	t: NavUserT;
	onSignOutClick?: () => void;
	items?: NavUserItem[];
}) {
	const { isMobile } = useSidebar();

	return (
		<SidebarMenu>
			<SidebarMenuItem>
				<DropdownMenu>
					<DropdownMenuTrigger asChild>
						<SidebarMenuButton
							size="lg"
							className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
						>
							<Avatar className="h-8 w-8 rounded-lg">
								<AvatarImage src={user.image} alt={user.name} />
								<AvatarFallback className="rounded-lg">CN</AvatarFallback>
							</Avatar>

							<div className="grid flex-1 text-left text-sm leading-tight">
								<span className="truncate font-semibold">{user.name}</span>
								<span className="truncate text-xs">{user.email}</span>
							</div>
							<ChevronsUpDown className="ml-auto size-4" />
						</SidebarMenuButton>
					</DropdownMenuTrigger>

					<DropdownMenuContent
						className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
						side={isMobile ? "bottom" : "right"}
						align="end"
						sideOffset={4}
					>
						<DropdownMenuLabel className="p-0 font-normal">
							<div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
								<Avatar className="h-8 w-8 rounded-lg">
									<AvatarImage src={user.image} alt={user.name} />
									<AvatarFallback className="rounded-lg">CN</AvatarFallback>
								</Avatar>
								<div className="grid flex-1 text-left text-sm leading-tight">
									<span className="truncate font-semibold">{user.name}</span>
									<span className="truncate text-xs">{user.email}</span>
								</div>
							</div>
						</DropdownMenuLabel>

						<DropdownMenuSeparator />

						<DropdownMenuGroup>
							{items.map((item) => (
								<DropdownMenuItem
									key={item.title}
									onClick={item.onClick}
									asChild={Boolean(item.url)}
								>
									{item.url ? (
										<Link to={item.url}>
											<item.icon className="size-4" />
											{item.title}
										</Link>
									) : (
										<>
											<item.icon className="size-4" />
											{item.title}
										</>
									)}
								</DropdownMenuItem>
							))}
						</DropdownMenuGroup>

						<DropdownMenuSeparator />

						<DropdownMenuItem onClick={onSignOutClick}>
							<LogOut />
							{t.signOut}
						</DropdownMenuItem>
					</DropdownMenuContent>
				</DropdownMenu>
			</SidebarMenuItem>
		</SidebarMenu>
	);
}
