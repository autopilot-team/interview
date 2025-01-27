import { API_BASE_URL } from "@autopilot/api";
import { LogIn } from "@autopilot/ui/components/icons";
import { MessageCard } from "@autopilot/ui/components/message-card";
import { useTranslation } from "react-i18next";
import { useLoaderData } from "react-router";

export async function clientLoader({ request }: { request: Request }) {
	const url = new URL(request.url);
	const token = url.searchParams.get("token");

	if (!token) {
		throw new Response("Token is required", { status: 404 });
	}

	const response = await fetch(`${API_BASE_URL}/v1/identity/verify-email`, {
		headers: {
			"Content-Type": "application/json",
		},
		method: "POST",
		body: JSON.stringify({ token }),
	});

	return {
		isSuccess: response.ok,
	};
}

export default function Component() {
	const { t } = useTranslation(["identity"]);
	const { isSuccess } = useLoaderData<typeof clientLoader>();

	return (
		<div className="flex min-h-svh flex-col items-center justify-center bg-muted p-2">
			{isSuccess ? (
				<MessageCard
					title={t("verifyEmail.success")}
					description={t("verifyEmail.successDescription")}
					backButton={{ show: false }}
					homeButton={{
						icon: <LogIn className="mr-2 h-4 w-4" />,
						label: t("signInForm.signIn"),
						show: true,
						to: "/sign-in",
					}}
					variant="success"
				/>
			) : (
				<MessageCard
					title={t("verifyEmail.error")}
					description={t("verifyEmail.errorDescription")}
					backButton={{ show: false }}
					homeButton={{
						icon: <LogIn className="mr-2 h-4 w-4" />,
						label: t("signInForm.signIn"),
						show: true,
						to: "/sign-in",
					}}
					variant="error"
				/>
			)}
		</div>
	);
}
