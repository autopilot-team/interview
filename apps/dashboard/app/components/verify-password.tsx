import { api } from "@autopilot/api";
import { Button } from "@autopilot/ui/components/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@autopilot/ui/components/dialog";
import { Label } from "@autopilot/ui/components/label";
import { Password } from "@autopilot/ui/components/password";
import { zodResolver } from "@autopilot/ui/lib/hook-form-zod-resolver";
import { useForm } from "@autopilot/ui/lib/react-hook-form";
import { z } from "@autopilot/ui/lib/zod";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";

interface VerifyPasswordData {
	password: string;
}

interface VerifyPasswordProps {
	/** Whether the dialog is open */
	open: boolean;
	/** Called when the dialog should be closed */
	onOpenChange: (open: boolean) => void;
	/** Called when password verification is successful */
	onVerified: () => void;
	/** Called when password verification fails */
	onError: () => void;
	/** Whether to show loading state */
	isLoading?: boolean;
	/** Custom title for the dialog */
	title?: string;
	/** Custom description for the dialog */
	description?: string;
	/** Custom text for the submit button */
	submitText?: string;
}

export function VerifyPassword({
	open,
	onOpenChange,
	onVerified,
	onError,
	isLoading = false,
	title,
	description,
	submitText,
}: VerifyPasswordProps) {
	const { t } = useTranslation(["common", "identity"]);

	const verifyPasswordSchema = z.object({
		password: z
			.string()
			.min(1, t("identity:verifyPassword.errors.passwordRequired")),
	});

	const {
		register,
		handleSubmit,
		formState: { errors },
		reset,
	} = useForm<VerifyPasswordData>({
		resolver: zodResolver(verifyPasswordSchema),
	});

	const verifyPasswordMutation = api.v1.useMutation(
		"post",
		"/v1/identity/verify-password",
	);

	const onSubmit = handleSubmit(async (data) => {
		try {
			await verifyPasswordMutation.mutateAsync({
				body: {
					password: data.password,
				},
			});
			onVerified();
			reset();
		} catch (error) {
			onError();
		}
	});

	// Reset form when dialog is opened
	useEffect(() => {
		if (open) {
			reset();
		}
	}, [open, reset]);

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>
						{title || t("identity:verifyPassword.title")}
					</DialogTitle>
					<DialogDescription>
						{description || t("identity:verifyPassword.description")}
					</DialogDescription>
				</DialogHeader>

				<form onSubmit={onSubmit}>
					<div className="space-y-2">
						<Label htmlFor="verify-password">
							{t("identity:verifyPassword.label")}
						</Label>

						<Password
							autoComplete="current-password"
							aria-describedby={
								errors.password ? "verify-password-error" : undefined
							}
							disabled={isLoading}
							error={!!errors.password}
							registration={register("password")}
						/>

						{errors.password?.message && (
							<p className="text-sm text-destructive">
								{errors.password.message}
							</p>
						)}
					</div>
				</form>

				<DialogFooter className="mt-8">
					<Button
						type="button"
						variant="ghost"
						onClick={() => onOpenChange(false)}
						disabled={isLoading}
					>
						{t("identity:verifyPassword.cancel")}
					</Button>

					<Button
						type="submit"
						form="verify-password-form"
						disabled={isLoading}
					>
						{submitText || t("identity:verifyPassword.next")}
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
