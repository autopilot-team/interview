import { useIdentity } from "@/components/identity-provider";
import { config } from "@/libs/config";
import { Brand } from "@autopilot/ui/components/brand";
import { SignInForm } from "@autopilot/ui/components/sign-in-form";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router";

interface SignInError {
	code?: string;
}

export default function Component() {
	const navigate = useNavigate();
	const { t } = useTranslation(["identity"]);
	const { isSigningIn, signIn } = useIdentity();
	const [generalError, setGeneralError] = useState("");
	const [successMessage, setSuccessMessage] = useState("");

	return (
		<div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
			<div className="flex w-full max-w-sm flex-col gap-6">
				<Brand className="self-center" />

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
								navigate("/");
							}
						} catch (error) {
							const err = error as SignInError;
							// Handle specific error codes
							switch (err?.code) {
								case "identity.account_locked":
									setGeneralError(
										t("errors.signIn.accountLocked").split("\\n").join("\n"),
									);
									break;
								case "identity.email_not_verified":
									setGeneralError(t("errors.signIn.emailNotVerified"));
									break;
								case "identity.invalid_credentials":
									setGeneralError(t("errors.signIn.invalidCredentials"));
									break;
								case "identity.invalid_session":
									setGeneralError(t("errors.signIn.invalidSession"));
									break;
								default:
									setGeneralError(t("errors.signIn.description"));
							}
						}
					}}
				/>
			</div>
		</div>
	);
}
