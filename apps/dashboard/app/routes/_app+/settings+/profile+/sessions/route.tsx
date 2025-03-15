import { getCountryFlag } from "@/libs/country";
import type { RouteHandle } from "@/routes/_app+/_layout";
import { api, type v1 } from "@autopilot/api";
import { Button } from "@autopilot/ui/components/button";
import {
	LaptopIcon,
	MapPinIcon,
	SmartphoneIcon,
	TabletIcon,
	TvIcon,
	XIcon,
} from "@autopilot/ui/components/icons";
import { Spinner } from "@autopilot/ui/components/spinner";
import { toast } from "@autopilot/ui/lib/sonner";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { UAParser } from "ua-parser-js";

export const handle = {
	breadcrumb: "common:breadcrumbs.settings.profile.sessions",
} satisfies RouteHandle;

type APISession = v1.components["schemas"]["Session"];

interface Session {
	id: string;
	ua: UAParser.IResult;
	ip: string;
	country: string;
	firstSeen: Date;
	isCurrent: boolean;
}

const parse = (session: APISession): Session => {
	return {
		id: session.id,
		ua: UAParser(session.userAgent || ""),
		ip: session.ipAddress || "Unknown",
		country: session.country || "Unknown",
		isCurrent: session.current,
		firstSeen: new Date(session.createdAt),
	};
};

export default function Component() {
	const { t } = useTranslation(["common"]);

	const apiSessionsQuery = api.v1.useQuery("get", "/v1/identity/sessions");
	const [sessions, setSessions] = useState<Session[]>([]);
	const { isPending: isSessionsLoading, data: sessionsResponse } =
		apiSessionsQuery;
	const { mutateAsync: deleteAllSessionsMutate } = api.v1.useMutation(
		"delete",
		"/v1/identity/sessions",
		{
			onSuccess: () => {
				apiSessionsQuery.refetch();
			},
		},
	);
	const { mutateAsync: deleteSessionMutate } = api.v1.useMutation(
		"delete",
		"/v1/identity/sessions/{id}",
		{
			onSuccess: () => {
				apiSessionsQuery.refetch();
			},
		},
	);

	useEffect(() => {
		const fetchSessions = async () => {
			try {
				if (isSessionsLoading) {
					return;
				}

				if (!sessionsResponse) {
					throw new Error("Failed to sessions keys");
				}

				const sessions = sessionsResponse.sessions?.map(parse) || [];
				setSessions(sessions);
			} catch (error) {
				toast.error(t("common:error.sessions.load.title"), {
					id: "load-api-keys",
					description: t("common:error.sessions.load.description"),
					icon: <XIcon className="size-4 text-destructive" />,
				});
			}
		};

		fetchSessions();
	}, [sessionsResponse, isSessionsLoading, t]);

	const handleRevokeSession = async (id: string) => {
		try {
			await deleteSessionMutate({ params: { path: { id } } });
			toast.success(t("settings.sessions.success"));
		} catch (error) {
			toast.error(t("settings.sessions.error"));
		}
	};

	const handleRevokeAllOtherSessions = async () => {
		try {
			await deleteAllSessionsMutate({});
		} catch (error) {
			toast.error(t("settings.sessions.revokeAllError"));
		}
	};

	const getDeviceIcon = (ua: UAParser.IResult) => {
		switch (ua.device.type) {
			case "wearable":
			case "mobile":
				return <SmartphoneIcon className="h-4 w-4" />;
			case "tablet":
				return <TabletIcon className="h-4 w-4" />;
			case "console":
			case "smarttv":
				return <TvIcon className="h-4 w-4" />;
			default:
				return <LaptopIcon className="h-4 w-4" />;
		}
	};

	return (
		<div className="space-y-8">
			<div className="grid grid-cols-1 gap-x-8 gap-y-4 lg:grid-cols-2">
				<div>
					<h3 className="text-lg font-medium">
						{t("common:settings.sessions.title")}
					</h3>

					<p className="text-sm text-muted-foreground">
						{t("common:settings.sessions.description")}
					</p>
				</div>

				<div className="flex justify-end py-2">
					<Button
						className="text-sm"
						onClick={handleRevokeAllOtherSessions}
						disabled={sessions.length <= 1}
						size="sm"
						variant="destructive"
					>
						{t("common:settings.sessions.revokeAll")}
					</Button>
				</div>
			</div>

			{isSessionsLoading ? (
				<div className="flex h-[200px] items-center justify-center">
					<Spinner className="size-8" />
				</div>
			) : (
				<div className="mt-6 space-y-6">
					{sessions.map((session) => (
						<div key={session.id} className="rounded-lg border p-4">
							<div className="flex items-start justify-between">
								<div className="space-y-1">
									<div className="flex items-center gap-2">
										{getDeviceIcon(session.ua)}
										<p className="font-medium">
											{session.ua.device.model
												? `${session.ua.browser.name} on ${session.ua.device.model}`
												: `${session.ua.browser.name} on ${session.ua.os.name}`}
											{session.isCurrent && (
												<span className="ml-2 text-xs text-green-500">
													{t("common:settings.sessions.currentSession")}
												</span>
											)}
										</p>
									</div>

									<div className="flex flex-wrap gap-x-4 gap-y-2 text-sm text-muted-foreground">
										<div className="flex items-center gap-1 text-xs">
											<MapPinIcon className="h-3 w-3" />
											<span>
												{session.country.length === 2
													? `${getCountryFlag(session.country)}`
													: session.country}
											</span>
											<span>{session.ip}</span>
										</div>
									</div>

									<div className="mt-2 flex flex-wrap gap-4 text-xs text-muted-foreground">
										<div>
											<span className="font-medium">
												{t("common:settings.sessions.browser")}:{" "}
											</span>
											{session.ua.browser.name || "Unknown"}
											&nbsp;
											{session.ua.browser.version || ""}
										</div>
										<div>
											<span className="font-medium">
												{t("common:settings.sessions.os")}:{" "}
											</span>
											{session.ua.os.name || "Unknown"}
											&nbsp;
											{session.ua.os.version || ""}
										</div>
										<div>
											<span className="font-medium">
												{t("common:settings.sessions.firstSeen")}:{" "}
											</span>
											{session.firstSeen.toLocaleDateString()}
										</div>
									</div>
								</div>

								{!session.isCurrent && (
									<Button
										variant="destructive"
										size="sm"
										onClick={() => handleRevokeSession(session.id)}
									>
										{t("common:settings.sessions.revoke")}
									</Button>
								)}
							</div>
						</div>
					))}
				</div>
			)}
		</div>
	);
}
