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
import { NavUser, type NavUserT } from "./nav-user.js";

export interface NavigationT {
	overview: {
		title: string;
		activity: string;
		balances: string;
		performance: string;
		settings: string;
	};
	moneyIn: {
		title: string;
		transactions: string;
		paymentMethods: string;
		checkout: string;
		subscriptions: string;
		products: string;
	};
	moneyFlow: {
		title: string;
		accounts: string;
		transfers: string;
		treasury: string;
		fx: string;
	};
	moneyOut: {
		title: string;
		payouts: string;
		expenses: string;
		vendors: string;
		tax: string;
	};
	operations: {
		title: string;
		processors: string;
		ledger: string;
		reconciliation: string;
		reports: string;
	};
	risk: {
		title: string;
		monitoring: string;
		prevention: string;
		disputes: string;
		compliance: string;
		vault: string;
	};
	intelligence: {
		title: string;
		analytics: string;
		insights: string;
		reports: string;
		models: string;
	};
	developer: {
		title: string;
		apiKeys: string;
		webhooks: string;
		documentation: string;
		status: string;
	};
}

export interface User {
	name: string;
	email: string;
	avatar: string;
}

export interface AppSidebarProps extends React.ComponentProps<typeof Sidebar> {
	t: {
		nav: NavigationT;
		entitySwitcher: EntitySwitcherT;
		navUser: NavUserT;
	};
	entities: Entity[];
	currentEntity?: Entity;
	onEntityChange?: (entity: Entity | undefined) => void;
	onCreateEntity?: () => void;
	user: User;
	navigation: NavMainItem[];
}

export function AppSidebar({
	t,
	entities,
	currentEntity,
	onEntityChange,
	onCreateEntity,
	user,
	navigation,
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
				<NavUser user={user} t={t.navUser} />
			</SidebarFooter>

			<SidebarRail />
		</Sidebar>
	);
}
