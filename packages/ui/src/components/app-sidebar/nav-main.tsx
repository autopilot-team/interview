"use client";

import {
	Collapsible,
	CollapsibleContent,
	CollapsibleTrigger,
} from "@autopilot/ui/components/collapsible";
import { ChevronRight, type LucideIcon } from "@autopilot/ui/components/icons";
import {
	SidebarGroup,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarMenuSub,
	SidebarMenuSubButton,
	SidebarMenuSubItem,
} from "@autopilot/ui/components/sidebar";
import { Link } from "react-router";

export interface NavMainItem {
	title: string;
	icon: LucideIcon;
	items: {
		title: string;
		url: string;
		icon: LucideIcon;
		isActive?: boolean;
		hasNotification?: boolean;
		notificationCount?: string;
	}[];
}

export function NavMain({
	items,
}: {
	items: NavMainItem[];
}) {
	return (
		<SidebarGroup>
			<SidebarMenu>
				{items.map((section) => (
					<Collapsible
						key={section.title}
						asChild
						defaultOpen={section.items.some((item) => item.isActive)}
						className="group/collapsible"
					>
						<SidebarMenuItem>
							<CollapsibleTrigger asChild>
								<SidebarMenuButton tooltip={section.title}>
									<section.icon className="size-4" />
									<span>{section.title}</span>
									<ChevronRight className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
								</SidebarMenuButton>
							</CollapsibleTrigger>

							<CollapsibleContent>
								<SidebarMenuSub>
									{section.items?.map((item) => (
										<SidebarMenuSubItem key={item.title}>
											<SidebarMenuSubButton asChild>
												<Link
													to={item.url}
													data-active={item.isActive}
													data-notification={item.hasNotification}
													data-notification-count={item.notificationCount}
												>
													{item.icon && <item.icon className="size-4" />}
													<span>{item.title}</span>
												</Link>
											</SidebarMenuSubButton>
										</SidebarMenuSubItem>
									))}
								</SidebarMenuSub>
							</CollapsibleContent>
						</SidebarMenuItem>
					</Collapsible>
				))}
			</SidebarMenu>
		</SidebarGroup>
	);
}
