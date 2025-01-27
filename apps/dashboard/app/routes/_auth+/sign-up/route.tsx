import { config } from "@/libs/config";
import { api } from "@autopilot/api";
import { Brand } from "@autopilot/ui/components/brand";
import { SignUpForm } from "@autopilot/ui/components/sign-up-form";
import { useToast } from "@autopilot/ui/hooks/use-toast";
import { useTranslation } from "react-i18next";

export default function Component() {
	const { toast } = useToast();
	const { t } = useTranslation(["identity"]);
	const { isPending, mutate } = api.v1.useMutation(
		"post",
		"/v1/identity/sign-up",
	);

	return (
		<div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
			<div className="flex w-full max-w-sm flex-col gap-6">
				<Brand className="self-center" />

				<SignUpForm
					cfTurnstileSiteKey={config.cfTurnstileSiteKey}
					isLoading={isPending}
					t={t("identity:signUpForm", { returnObjects: true })}
					handleSignUp={(data, reset) => {
						const { confirmPassword, ...rest } = data;

						mutate(
							{
								body: rest,
							},
							{
								onError: () => {
									toast({
										title: t("identity:errors.signUp.title"),
										description: t("identity:errors.signUp.description"),
									});
								},
								onSuccess: () => {
									toast({
										duration: 10_000,
										title: t("identity:success.signUp.title"),
										description: t("identity:success.signUp.description"),
									});
									reset();
								},
							},
						);
					}}
				/>
			</div>
		</div>
	);
}
