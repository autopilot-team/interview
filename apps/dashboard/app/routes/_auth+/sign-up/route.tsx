import { config } from "@/libs/config";
import { Brand } from "@autopilot/ui/components/brand";
import { SignUpForm } from "@autopilot/ui/components/sign-up-form";
import { useTranslation } from "react-i18next";

export default function Component() {
	const { t } = useTranslation(["identity"]);

	return (
		<div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
			<div className="flex w-full max-w-sm flex-col gap-6">
				<Brand className="self-center" />

				<SignUpForm
					cfTurnstileSiteKey={config.cfTurnstileSiteKey}
					isLoading={false}
					t={t("identity:signUpForm", { returnObjects: true })}
					handleSignUp={async (data, reset) => {}}
				/>
			</div>
		</div>
	);
}
