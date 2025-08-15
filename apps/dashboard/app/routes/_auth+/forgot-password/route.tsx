import { Brand } from "@autopilot/ui/components/brand";
import { ForgotPasswordForm } from "@autopilot/ui/components/forgot-password-form";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useIdentity } from "@/components/identity-provider";
import { config } from "@/libs/config";

interface ForgotPasswordFormData {
	cfTurnstileToken: string;
	email: string;
}

export default function Component() {
	const { t } = useTranslation(["identity"]);
	const { isForgottingPassword, forgotPassword } = useIdentity();
	const [generalError, setGeneralError] = useState("");
	const [successMessage, setSuccessMessage] = useState("");

	return (
		<div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
			<div className="flex w-full max-w-sm flex-col gap-6">
				<Brand className="self-center" />

				<ForgotPasswordForm
					cfTurnstileSiteKey={config.cfTurnstileSiteKey}
					generalError={generalError}
					isLoading={isForgottingPassword}
					successMessage={successMessage}
					t={t("identity:forgotPasswordForm", { returnObjects: true })}
					handleForgotPassword={async (data, reset) => {
						setGeneralError("");
						setSuccessMessage("");

						try {
							await forgotPassword({
								cfTurnstileToken: data.cfTurnstileToken,
								email: data.email,
							});
							reset();
							setSuccessMessage(t("success.forgotPassword.description"));
						} catch (error) {
							console.error(error);
							setGeneralError(t("errors.forgotPassword.description"));
						}
					}}
				/>
			</div>
		</div>
	);
}
