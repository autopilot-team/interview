import { useIdentity } from "@/components/identity-provider";
import { VerifyPassword } from "@/components/verify-password";
import type { RouteHandle } from "@/routes/_app+/_layout";
import { api, type v1 } from "@autopilot/api";
import { Button } from "@autopilot/ui/components/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@autopilot/ui/components/dialog";
import {
	CheckIcon,
	CopyIcon,
	QrCodeIcon,
	XIcon,
} from "@autopilot/ui/components/icons";
import { Input } from "@autopilot/ui/components/input";
import { Label } from "@autopilot/ui/components/label";
import { Password } from "@autopilot/ui/components/password";
import { Switch } from "@autopilot/ui/components/switch";
import { zodResolver } from "@autopilot/ui/lib/hook-form-zod-resolver";
import { useForm } from "@autopilot/ui/lib/react-hook-form";
import { toast } from "@autopilot/ui/lib/sonner";
import { cn } from "@autopilot/ui/lib/utils";
import { z } from "@autopilot/ui/lib/zod";
import { useState } from "react";
import { useTranslation } from "react-i18next";

export const handle = {
	breadcrumb: "common:breadcrumbs.settings.profile.security",
} satisfies RouteHandle;

interface UpdatePasswordData {
	currentPassword: string;
	newPassword: string;
	confirmPassword: string;
}

type ApiError = v1.components["schemas"]["Error"];

function PasswordSection() {
	const { t } = useTranslation(["common", "identity"]);
	const { updatePassword, isUpdatingPassword } = useIdentity();
	const [password, setPassword] = useState("");

	const updatePasswordSchema = z
		.object({
			currentPassword: z
				.string()
				.min(8, t("identity:security.password.errors.passwordMinLength")),
			newPassword: z
				.string()
				.min(8, t("identity:security.password.errors.passwordMinLength"))
				.regex(
					/[A-Z]/,
					t("identity:security.password.errors.passwordUppercase"),
				)
				.regex(
					/[a-z]/,
					t("identity:security.password.errors.passwordLowercase"),
				)
				.regex(/[0-9]/, t("identity:security.password.errors.passwordNumber"))
				.regex(
					/[^A-Za-z0-9]/,
					t("identity:security.password.errors.passwordSpecial"),
				),
			confirmPassword: z.string(),
		})
		.refine(
			(data: UpdatePasswordData) => data.newPassword === data.confirmPassword,
			{
				message: t("identity:security.password.errors.passwordsDoNotMatch"),
				path: ["confirmPassword"],
			},
		);

	const {
		register,
		handleSubmit,
		formState: { errors },
		reset,
		setError,
	} = useForm<UpdatePasswordData>({
		resolver: zodResolver(updatePasswordSchema),
	});

	const onSubmit = handleSubmit(async (data) => {
		try {
			await updatePassword({
				currentPassword: data.currentPassword,
				newPassword: data.newPassword,
			});
			reset();
			setPassword("");
			toast.success(t("identity:security.password.success"), {
				icon: <CheckIcon className="h-4 w-4 text-green-500" />,
			});
		} catch (error) {
			const err = error as ApiError;

			switch (err?.code) {
				case "InvalidCredentials":
					setError("currentPassword", {
						message: t("identity:security.password.errors.invalidPassword"),
					});
					break;
				default:
					toast.error(t("identity:security.password.error"), {
						icon: <XIcon className="h-4 w-4 text-destructive" />,
					});
			}
		}
	});

	return (
		<div className="border-b pb-8">
			<div className="grid grid-cols-1 gap-x-8 gap-y-4 lg:grid-cols-2">
				<div>
					<h3 className="text-lg font-medium">
						{t("identity:security.password.title")}
					</h3>

					<p className="text-sm text-muted-foreground">
						{t("identity:security.password.description")}
					</p>
				</div>

				<form onSubmit={onSubmit} className="space-y-4">
					<div className="space-y-2">
						<Label htmlFor="currentPassword">
							{t("identity:security.password.currentPassword")}
						</Label>

						<Password
							autoComplete="current-password"
							disabled={isUpdatingPassword}
							error={!!errors.currentPassword}
							isNewPassword={false}
							registration={register("currentPassword")}
						/>

						{errors.currentPassword?.message && (
							<p className="text-sm text-destructive">
								{errors.currentPassword.message}
							</p>
						)}
					</div>

					<div className="space-y-2">
						<Label htmlFor="newPassword">
							{t("identity:security.password.newPassword")}
						</Label>

						<Password
							disabled={isUpdatingPassword}
							error={!!errors.newPassword}
							onValueChange={setPassword}
							registration={register("newPassword", {
								onChange: (e) => {
									setPassword(e.target.value);
								},
							})}
							showStrength
							strengthT={t("identity:passwordStrength", {
								returnObjects: true,
							})}
							value={password}
						/>

						{errors.newPassword?.message && (
							<p className="text-sm text-destructive">
								{errors.newPassword.message}
							</p>
						)}
					</div>

					<div className="space-y-2">
						<Label htmlFor="confirmPassword">
							{t("identity:security.password.confirmPassword")}
						</Label>

						<Password
							disabled={isUpdatingPassword}
							error={!!errors.confirmPassword}
							registration={register("confirmPassword")}
						/>

						{errors.confirmPassword?.message && (
							<p className="text-sm text-destructive">
								{errors.confirmPassword.message}
							</p>
						)}
					</div>

					<div className="flex justify-end mt-8">
						<Button type="submit" disabled={isUpdatingPassword}>
							{t("identity:security.password.updatePassword")}
						</Button>
					</div>
				</form>
			</div>
		</div>
	);
}

function TwoFactorSection() {
	const { user } = useIdentity();
	const { t } = useTranslation(["common", "identity"]);
	const [is2FAEnabled, setIs2FAEnabled] = useState(user?.isTwoFactorEnabled);
	const [is2FADialogOpen, setIs2FADialogOpen] = useState(false);
	const [qrCodeUrl, setQrCodeUrl] = useState<string | null>(null);
	const [backupCodes, setBackupCodes] = useState<string[]>([]);
	const [setupStep, setSetupStep] = useState<"password" | "qr" | "verify">(
		"password",
	);
	const [verifyPasswordOpen, setVerifyPasswordOpen] = useState(false);
	const [verifyAction, setVerifyAction] = useState<"enable" | "disable">(
		"enable",
	);
	const [isLoadingQR, setIsLoadingQR] = useState(false);
	const [qrError, setQrError] = useState<Error | null>(null);
	const [hasCopied, setHasCopied] = useState(false);

	const enable2FASchema = z.object({
		code: z
			.string()
			.length(6, t("identity:twoFactorForm.errors.codeLength"))
			.regex(/^\d+$/, t("identity:twoFactorForm.errors.codeRequired")),
	});

	const {
		register: registerEnable2FA,
		handleSubmit: handleSubmitEnable2FA,
		formState: { errors: enable2FAErrors },
		reset: resetEnable2FA,
		setError: setEnable2FAError,
	} = useForm({
		resolver: zodResolver(enable2FASchema),
	});

	const setupTwoFactorMutation = api.v1.useMutation(
		"post",
		"/v1/identity/setup-two-factor",
	);
	const enableTwoFactorMutation = api.v1.useMutation(
		"post",
		"/v1/identity/enable-two-factor",
	);
	const disableTwoFactorMutation = api.v1.useMutation(
		"delete",
		"/v1/identity/disable-two-factor",
	);

	const loadQRCode = async () => {
		setIsLoadingQR(true);
		setQrError(null);
		try {
			const response = await setupTwoFactorMutation.mutateAsync({});
			setQrCodeUrl(response.qrCode);
			setBackupCodes(response.backupCodes || []);
		} catch (error) {
			setQrError(
				error instanceof Error ? error : new Error("Failed to load QR code"),
			);
			toast.error(t("identity:errors.twoFactor.description"));
		} finally {
			setIsLoadingQR(false);
		}
	};

	const handleEnable2FA = handleSubmitEnable2FA(async (data) => {
		try {
			await enableTwoFactorMutation.mutateAsync({
				body: {
					code: data.code,
				},
			});
			setIs2FAEnabled(true);
			setIs2FADialogOpen(false);
			resetEnable2FA();
			setSetupStep("password");
			setQrCodeUrl(null);
			setBackupCodes([]);
			toast.success(t("identity:security.twoFactor.enableSuccess"), {
				icon: <CheckIcon className="h-4 w-4 text-green-500" />,
				duration: 5000,
			});
		} catch (error: unknown) {
			const err = error as ApiError;

			switch (err?.code) {
				case "InvalidTwoFactorCode":
					setEnable2FAError("code", {
						message: t("identity:errors.twoFactor.invalidCode"),
					});
					break;
				default:
					setEnable2FAError("code", {
						message: t("identity:security.twoFactor.enableError"),
					});
			}
		}
	});

	return (
		<div>
			<div className="grid grid-cols-1 gap-x-8 gap-y-4 lg:grid-cols-2">
				<div>
					<h3 className="text-lg font-medium">
						{t("identity:security.twoFactor.title")}
					</h3>

					<p className="text-sm text-muted-foreground">
						{t("identity:security.twoFactor.description")}
					</p>
				</div>

				<div className="flex items-center justify-between rounded-lg border p-4">
					<div className="space-y-1">
						<p className="font-medium">
							{t("identity:security.twoFactor.label")}
						</p>

						<p className="text-sm text-muted-foreground">
							{is2FAEnabled
								? t("identity:security.twoFactor.enabled")
								: t("identity:security.twoFactor.disabled")}
						</p>
					</div>

					<Switch
						checked={is2FAEnabled}
						onCheckedChange={(checked) => {
							setVerifyAction(checked ? "enable" : "disable");
							setVerifyPasswordOpen(true);
						}}
					/>
				</div>
			</div>

			<VerifyPassword
				open={verifyPasswordOpen}
				onOpenChange={setVerifyPasswordOpen}
				onVerified={async () => {
					if (verifyAction === "enable") {
						setIs2FADialogOpen(true);
						setSetupStep("qr");
						// Reset any existing QR code or backup codes
						setQrCodeUrl(null);
						setBackupCodes([]);
						await loadQRCode();
					} else {
						try {
							await disableTwoFactorMutation.mutateAsync({});
							setIs2FAEnabled(false);
							toast.success(t("identity:security.twoFactor.disableSuccess"), {
								icon: <CheckIcon className="h-4 w-4 text-green-500" />,
								duration: 5000,
							});
						} catch (error) {
							toast.error(t("identity:security.twoFactor.disableError"));
						}
					}
					setVerifyPasswordOpen(false);
				}}
				onError={() => {
					toast.error(t("identity:verifyPassword.errors.invalidPassword"));
				}}
			/>

			<Dialog
				open={is2FADialogOpen}
				onOpenChange={(open) => {
					// Prevent closing by clicking outside or pressing escape during setup
					if (!open && setupStep !== "verify") {
						return;
					}
					setIs2FADialogOpen(open);
				}}
			>
				<DialogContent
					className={cn(
						"sm:max-w-md",
						setupStep !== "verify" && "[&>button]:hidden",
					)}
				>
					<DialogHeader>
						<DialogTitle>
							{setupStep === "verify"
								? t("identity:security.twoFactor.enable.verify")
								: t("identity:security.twoFactor.enable.title")}
						</DialogTitle>

						<DialogDescription>
							{setupStep === "verify"
								? t("identity:security.twoFactor.verify")
								: t("identity:security.twoFactor.enable.description")}
						</DialogDescription>
					</DialogHeader>

					<form onSubmit={handleEnable2FA} className="space-y-4">
						{setupStep === "qr" && (
							<div className="space-y-4">
								<div className="rounded-lg border p-4">
									<div className="flex justify-center">
										{isLoadingQR ? (
											<div className="flex h-32 w-32 items-center justify-center">
												<div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
											</div>
										) : qrError ? (
											<div className="flex h-32 w-32 flex-col items-center justify-center gap-2">
												<QrCodeIcon className="h-12 w-12 text-muted-foreground" />
												<Button
													type="button"
													variant="outline"
													size="sm"
													onClick={loadQRCode}
												>
													{t("identity:security.twoFactor.retry")}
												</Button>
											</div>
										) : qrCodeUrl ? (
											<img
												src={qrCodeUrl}
												alt="QR Code"
												className="h-32 w-32"
											/>
										) : (
											<QrCodeIcon className="h-32 w-32" />
										)}
									</div>

									<p className="mt-2 text-center text-sm text-muted-foreground">
										{qrError
											? t("identity:security.twoFactor.error")
											: t("identity:security.twoFactor.scanQrCode")}
									</p>
								</div>

								{!qrError && (
									<>
										<div className="space-y-2">
											<div className="flex items-center justify-between">
												<h4 className="font-medium">
													{t("identity:security.twoFactor.backupCodes.title")}
												</h4>
												<Button
													type="button"
													variant="ghost"
													size="sm"
													className="h-8 px-2"
													onClick={async () => {
														await navigator.clipboard.writeText(
															backupCodes.join("\n"),
														);
														setHasCopied(true);
														setTimeout(() => setHasCopied(false), 2000);
													}}
												>
													{hasCopied ? (
														<CheckIcon className="h-4 w-4 text-green-500" />
													) : (
														<CopyIcon className="h-4 w-4" />
													)}
													<span>
														{hasCopied
															? t(
																	"identity:security.twoFactor.backupCodes.copied",
																)
															: t(
																	"identity:security.twoFactor.backupCodes.copyAll",
																)}
													</span>
												</Button>
											</div>

											<p className="text-sm text-muted-foreground">
												{t(
													"identity:security.twoFactor.backupCodes.description",
												)}
											</p>

											<div className="grid grid-cols-2 gap-2">
												{backupCodes.map((code) => (
													<code
														key={code}
														className="select-all rounded bg-muted px-2 py-1 text-center font-mono text-sm"
													>
														{code}
													</code>
												))}
											</div>
										</div>

										<Button
											type="button"
											className="mt-4 w-full"
											onClick={() => setSetupStep("verify")}
										>
											{t("identity:security.twoFactor.enable.next")}
										</Button>
									</>
								)}
							</div>
						)}

						{setupStep === "verify" && (
							<div className="space-y-2">
								<Label htmlFor="enable-code">
									{t("identity:security.twoFactor.code")}
								</Label>

								<Input
									id="enable-code"
									{...registerEnable2FA("code")}
									className={cn("mt-1", {
										"ring-2 ring-destructive": enable2FAErrors.code,
									})}
								/>

								{enable2FAErrors.root?.message && (
									<p className="text-sm text-destructive">
										{enable2FAErrors.root?.message}
									</p>
								)}
							</div>
						)}

						{setupStep === "verify" && (
							<DialogFooter>
								<Button
									type="button"
									variant="ghost"
									onClick={() => setIs2FADialogOpen(false)}
								>
									{t("identity:security.twoFactor.cancel")}
								</Button>

								<Button type="submit">
									{t("identity:security.twoFactor.enable.confirm")}
								</Button>
							</DialogFooter>
						)}
					</form>
				</DialogContent>
			</Dialog>
		</div>
	);
}

export default function Component() {
	return (
		<div className="space-y-8">
			<PasswordSection />
			<TwoFactorSection />
		</div>
	);
}
