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

export interface NavSubItem {
	title: string;
	to: string;
	resource?: string;
	icon: LucideIcon;
	isActive?: boolean;
	hasNotification?: boolean;
	notificationCount?: string;
}

export interface NavMainItem {
	title: string;
	icon: LucideIcon;
	to?: string;
	resource?: string;
	items?: NavSubItem[];
}

export function NavMain({ items }: { items: NavMainItem[] }) {
	return (
		<SidebarGroup>
			<SidebarMenu>
				{items.map((section) => {
					// If no items array is provided and url exists, render a direct link
					if (!section.items && section.to) {
						return (
							<SidebarMenuItem key={section.title}>
								<SidebarMenuButton asChild tooltip={section.title}>
									<Link to={section.to}>
										<section.icon className="size-4" />
										<span>{section.title}</span>
									</Link>
								</SidebarMenuButton>
							</SidebarMenuItem>
						);
					}

					// Otherwise render the collapsible menu
					return (
						<Collapsible
							key={section.title}
							asChild
							defaultOpen={section.items?.some((item) => item.isActive)}
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
														to={item.to}
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
					);
				})}
			</SidebarMenu>
		</SidebarGroup>
	);
}
