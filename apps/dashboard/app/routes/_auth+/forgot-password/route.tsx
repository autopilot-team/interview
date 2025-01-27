import { config } from "@/libs/config";
import { Brand } from "@autopilot/ui/components/brand";
import { ForgotPasswordForm } from "@autopilot/ui/components/forgot-password-form";
import { useTranslation } from "react-i18next";

export default function Component() {
	const { t } = useTranslation(["identity"]);

	return (
		<div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
			<div className="flex w-full max-w-sm flex-col gap-6">
				<Brand className="self-center" />

				<ForgotPasswordForm
					cfTurnstileSiteKey={config.cfTurnstileSiteKey}
					t={t("identity:forgotPasswordForm", { returnObjects: true })}
					handleForgotPassword={async (data: { email: string }) => {}}
				/>
			</div>
		</div>
	);
}
