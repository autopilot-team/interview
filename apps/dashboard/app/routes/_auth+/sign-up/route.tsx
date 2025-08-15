import { Brand } from "@autopilot/ui/components/brand";
import { SignUpForm } from "@autopilot/ui/components/sign-up-form";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useIdentity } from "@/components/identity-provider";
import { config } from "@/libs/config";

export default function Component() {
	const { t } = useTranslation(["identity"]);
	const { isSigningUp, signUp } = useIdentity();
	const [generalError, setGeneralError] = useState("");
	const [successMessage, setSuccessMessage] = useState("");

	return (
		<div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
			<div className="flex w-full max-w-sm flex-col gap-6">
				<Brand className="self-center" />

				<SignUpForm
					cfTurnstileSiteKey={config.cfTurnstileSiteKey}
					generalError={generalError}
					isLoading={isSigningUp}
					successMessage={successMessage}
					t={{
						...t("identity:signUpForm", { returnObjects: true }),
						passwordStrength: t("identity:passwordStrength", {
							returnObjects: true,
						}),
					}}
					handleSignUp={async (data, reset) => {
						const { confirmPassword, ...rest } = data;
						setGeneralError("");
						setSuccessMessage("");

						try {
							await signUp(rest);
							setSuccessMessage(t("identity:success.signUp.description"));
							reset();
						} catch (error) {
							console.error(error);
							setGeneralError(t("identity:errors.signUp.description"));
						}
					}}
				/>
			</div>
		</div>
	);
}
