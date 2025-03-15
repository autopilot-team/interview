import { useIdentity } from "@/components/identity-provider";
import { config } from "@/libs/config";
import type { v1 } from "@autopilot/api";
import { Brand } from "@autopilot/ui/components/brand";
import { SignInForm } from "@autopilot/ui/components/sign-in-form";
import { TwoFactorForm } from "@autopilot/ui/components/two-factor-form";
import { useState } from "react";
import { useTranslation } from "react-i18next";

type ApiError = v1.components["schemas"]["Error"];

export default function Component() {
	const { t } = useTranslation(["identity"]);
	const {
		isSigningIn,
		signIn,
		isVerifyingTwoFactor,
		verifyTwoFactor,
		isTwoFactorPending,
	} = useIdentity();
	const [generalError, setGeneralError] = useState("");
	const [successMessage, setSuccessMessage] = useState("");

	return (
		<div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
			<div className="flex w-full max-w-sm flex-col gap-6">
				<Brand className="self-center" />

				{isTwoFactorPending ? (
					<TwoFactorForm
						generalError={generalError}
						isLoading={isVerifyingTwoFactor}
						t={t("identity:twoFactorForm", { returnObjects: true })}
						handleVerify={async (code) => {
							setGeneralError("");
							try {
								await verifyTwoFactor(code);
							} catch (error) {
								setGeneralError(t("errors.twoFactor.description"));
							}
						}}
					/>
				) : (
					<SignInForm
						cfTurnstileSiteKey={config.cfTurnstileSiteKey}
						generalError={generalError}
						isLoading={isSigningIn}
						successMessage={successMessage}
						t={t("identity:signInForm", { returnObjects: true })}
						handleSignIn={async (data, reset) => {
							setGeneralError("");
							setSuccessMessage("");

							try {
								if (await signIn(data)) {
									reset();
								}
							} catch (error) {
								const err = error as ApiError;

								// Handle specific error codes
								switch (err?.code) {
									case "AccountLocked":
										setGeneralError(
											t("errors.signIn.accountLocked").split("\\n").join("\n"),
										);
										break;
									case "EmailNotVerified":
										setGeneralError(t("errors.signIn.emailNotVerified"));
										break;
									case "InvalidCredentials":
										setGeneralError(t("errors.signIn.invalidCredentials"));
										break;
									default:
										setGeneralError(t("errors.signIn.description"));
								}
							}
						}}
					/>
				)}
			</div>
		</div>
	);
}
