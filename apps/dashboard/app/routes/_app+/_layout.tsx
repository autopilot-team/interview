import { AuthState, useIdentity } from "@/components/identity-provider";
import { ModeSwitcher } from "@/components/mode-switcher";
import { AppSidebar } from "@autopilot/ui/components/app-sidebar/app-sidebar";
import type { Entity } from "@autopilot/ui/components/app-sidebar/entity-switcher";
import type { NavMainItem } from "@autopilot/ui/components/app-sidebar/nav-main";
import {
	Breadcrumb,
	BreadcrumbItem,
	BreadcrumbLink,
	BreadcrumbList,
	BreadcrumbPage,
	BreadcrumbSeparator,
} from "@autopilot/ui/components/breadcrumb";
import { ClientOnly } from "@autopilot/ui/components/client-only";
import {
	Activity,
	Cable,
	Code2,
	CreditCard,
	HandCoins,
	Settings,
} from "@autopilot/ui/components/icons";
import { Separator } from "@autopilot/ui/components/separator";
import {
	SidebarInset,
	SidebarProvider,
	SidebarTrigger,
} from "@autopilot/ui/components/sidebar";
import type { ParseKeys } from "i18next";
import React, { useEffect } from "react";
import { useTranslation } from "react-i18next";
import {
	Navigate,
	Outlet,
	useLocation,
	useMatches,
	useNavigate,
	useParams,
} from "react-router";

export interface RouteHandle {
	breadcrumb?: ParseKeys<["common"]>;
}

// Paths that are not associated with an entity.
const NON_ENTITY_PATHS = ["/settings"];

export default function AppLayout() {
	const { t } = useTranslation(["common"]);
	const { authState, entity, role, signOut, switchEntity, user } =
		useIdentity();
	const location = useLocation();
	const navigate = useNavigate();
	const matches = useMatches();
	const params = useParams();
	const entityPath = (path: string) =>
		entity?.slug ? `/${entity.slug}${path}` : "#";
	const navigation: NavMainItem[] = [];
	const filterSub = (item: NavMainItem): NavMainItem => {
		if (!item.items) return item;
		const filtered = item.items.filter((subItem) => {
			if (!subItem.resource) return true;
			if (!role) return false;
			return subItem.resource in role.access;
		});

		return { ...item, items: filtered };
	};
	const roleNavigation = navigation.reduce((filtered, item) => {
		if (!item.resource) {
			filtered.push(filterSub(item));
			return filtered;
		}

		if (!role) {
			return filtered;
		}

		if (item.resource in role.access) {
			filtered.push(filterSub(item));
		}
		return filtered;
	}, [] as NavMainItem[]);
	const handleCreateEntity = async () => {};
	const onEntityChange = React.useCallback(
		async (entity: Entity | undefined) => {
			if (!entity) return;
			await switchEntity(entity.slug, entity.id);
		},
		[switchEntity],
	);
	const entities = React.useMemo<Entity[]>(() => {
		if (!user?.memberships) return [];

		return user.memberships.map((membership) => ({
			id: membership.entity.id,
			name: membership.entity.name,
			logo: Activity,
			type: membership.entity.type as "platform" | "organization" | "account",
			slug: membership.entity.slug,
			parentId: membership.entity.parentId,
		}));
	}, [user?.memberships]);
	const activeEntity = React.useMemo<Entity | undefined>(() => {
		if (!entity) return undefined;

		return {
			id: entity.id,
			name: entity.name,
			logo: Activity,
			type: entity.type as "platform" | "organization" | "account",
			slug: entity.slug,
			parentId: entity.parentId,
		};
	}, [entity]);

	useEffect(() => {
		if (authState === AuthState.UNAUTHENTICATED) {
			navigate("/sign-in");
		}
	}, [authState, navigate]);

	if (!user) {
		return <Navigate to="/sign-in" replace />;
	}

	if (
		entity &&
		!params.entity &&
		NON_ENTITY_PATHS.every((path) => !location.pathname.startsWith(path))
	) {
		switchEntity(entity.slug, entity.id);
		return;
	}

	if (!entity) {
		return <Navigate to="/onboarding" replace />;
	}

	return (
		<SidebarProvider>
			<ClientOnly>
				<div className="fixed top-4 right-4 z-50">
					<ModeSwitcher />
				</div>
			</ClientOnly>

			<AppSidebar
				currentEntity={activeEntity}
				entities={entities}
				navigation={roleNavigation}
				onCreateEntity={handleCreateEntity}
				onEntityChange={onEntityChange}
				onSignOutClick={async () => {
					await signOut();
				}}
				t={{
					entitySwitcher: t("common:entitySwitcher", { returnObjects: true }),
					navUser: t("common:navUser", { returnObjects: true }),
				}}
				user={user}
				userNavigation={[
					{
						title: t("common:navUser.billing"),
						icon: CreditCard,
						url: "/billing",
					},
					{
						title: t("common:navUser.settings"),
						icon: Settings,
						url: "/settings/profile",
					},
				]}
			/>

			<SidebarInset className="w-full md:w-[calc(100%-16rem)]">
				<header className="sticky top-0 z-10 flex h-16 shrink-0 items-center gap-2 border-b border-gray-200 bg-white px-4 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
					<SidebarTrigger className="-ml-1" />

					<Separator orientation="vertical" className="mr-2 h-4" />

					<Breadcrumb>
						<BreadcrumbList>
							{matches
								.filter(
									(match) =>
										(match.handle as RouteHandle | undefined)?.breadcrumb,
								)
								.map((match, index, matches) => {
									return (
										<React.Fragment key={match.id}>
											<BreadcrumbItem className="hidden md:block">
												{index === matches.length - 1 ? (
													<BreadcrumbPage>
														{/* biome-ignore lint/style/noNonNullAssertion: <explanation> */}
														{t((match.handle as RouteHandle).breadcrumb!)}
													</BreadcrumbPage>
												) : (
													<BreadcrumbLink href={match.pathname}>
														{/* biome-ignore lint/style/noNonNullAssertion: <explanation> */}
														{t((match.handle as RouteHandle).breadcrumb!)}
													</BreadcrumbLink>
												)}
											</BreadcrumbItem>

											{index < matches.length - 1 && (
												<BreadcrumbSeparator className="hidden md:block" />
											)}
										</React.Fragment>
									);
								})}
						</BreadcrumbList>
					</Breadcrumb>
				</header>

				<Outlet />
			</SidebarInset>
		</SidebarProvider>
	);
}
