import { useIdentity } from "@/components/identity-provider";
import { config } from "@/libs/config";
import { Brand } from "@autopilot/ui/components/brand";
import { ResetPasswordForm } from "@autopilot/ui/components/reset-password-form";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useSearchParams } from "react-router";

interface ResetPasswordFormData {
	cfTurnstileToken: string;
	password: string;
	confirmPassword: string;
}

export default function Component() {
	const { t } = useTranslation(["identity"]);
	const { isResettingPassword, resetPassword } = useIdentity();
	const [generalError, setGeneralError] = useState("");
	const [successMessage, setSuccessMessage] = useState("");
	const [searchParams] = useSearchParams();
	const navigate = useNavigate();
	const token = searchParams.get("token");

	// Redirect to sign in if no token is provided
	useEffect(() => {
		if (!token) {
			navigate("/sign-in");
		}
	}, [token, navigate]);

	return (
		<div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
			<div className="flex w-full max-w-sm flex-col gap-6">
				<Brand className="self-center" />

				<ResetPasswordForm
					cfTurnstileSiteKey={config.cfTurnstileSiteKey}
					generalError={generalError}
					isLoading={isResettingPassword}
					successMessage={successMessage}
					t={{
						...t("identity:resetPasswordForm", { returnObjects: true }),
						passwordStrength: t("identity:passwordStrength", {
							returnObjects: true,
						}),
					}}
					handleResetPassword={async (data, reset) => {
						setGeneralError("");
						setSuccessMessage("");

						try {
							await resetPassword({
								newPassword: data.password,
								token: token || "",
							});
							reset();
							setSuccessMessage(t("success.resetPassword.description"));
							// Redirect to sign in after 3 seconds
							setTimeout(() => {
								navigate("/sign-in");
							}, 3000);
						} catch (error) {
							console.error(error);
							setGeneralError(t("errors.resetPassword.invalidToken"));
						}
					}}
				/>
			</div>
		</div>
	);
}
