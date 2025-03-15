"use client";

import {
	Sidebar,
	SidebarContent,
	SidebarFooter,
	SidebarHeader,
	SidebarRail,
} from "@autopilot/ui/components/sidebar";
import type * as React from "react";
import {
	type Entity,
	EntitySwitcher,
	type EntitySwitcherT,
} from "./entity-switcher.js";
import { NavMain, type NavMainItem } from "./nav-main.js";
import { NavUser, type NavUserItem, type NavUserT } from "./nav-user.js";

export interface User {
	name: string;
	email: string;
	image?: string;
}

export interface AppSidebarProps extends React.ComponentProps<typeof Sidebar> {
	t: {
		entitySwitcher: EntitySwitcherT;
		navUser: NavUserT;
	};
	entities: Entity[];
	currentEntity?: Entity;
	onEntityChange?: (entity: Entity | undefined) => void;
	onCreateEntity?: () => Promise<void>;
	onSignOutClick?: () => Promise<void>;
	user: User;
	navigation: NavMainItem[];
	userNavigation?: NavUserItem[];
}

export function AppSidebar({
	t,
	entities,
	currentEntity,
	onEntityChange,
	onCreateEntity,
	onSignOutClick,
	user,
	navigation,
	userNavigation,
	...props
}: AppSidebarProps) {
	if (!t) return null;

	return (
		<Sidebar collapsible="icon" {...props}>
			<SidebarHeader>
				<EntitySwitcher
					t={t.entitySwitcher}
					entities={entities}
					currentEntity={currentEntity}
					onEntityChange={onEntityChange}
					canCreate={true}
					onCreateClick={onCreateEntity}
				/>
			</SidebarHeader>

			<SidebarContent>
				<NavMain key={currentEntity?.id} items={navigation} />
			</SidebarContent>

			<SidebarFooter>
				<NavUser
					items={userNavigation}
					onSignOutClick={onSignOutClick}
					t={t.navUser}
					user={user}
				/>
			</SidebarFooter>

			<SidebarRail />
		</Sidebar>
	);
}
