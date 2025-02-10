import { Button } from "@autopilot/ui/components/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@autopilot/ui/components/card";
import { Label } from "@autopilot/ui/components/label";
import { Password } from "@autopilot/ui/components/password";
import { cn } from "@autopilot/ui/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { Turnstile, type TurnstileInstance } from "@marsidev/react-turnstile";
import { useEffect, useRef, useState } from "react";
import { type UseFormReset, useForm } from "react-hook-form";
import { Link } from "react-router";
import { z } from "zod";

type ResetPasswordFormData = {
	cfTurnstileToken: string;
	password: string;
	confirmPassword: string;
};

export interface ResetPasswordFormT {
	title: string;
	description: string;
	password: string;
	passwordPlaceholder: string;
	confirmPassword: string;
	confirmPasswordPlaceholder: string;
	resetPassword: string;
	backToSignIn: string;
	passwordStrength: {
		secure: string;
		moderate: string;
		weak: string;
		requirements: {
			minLength: string;
			mixCase: string;
			number: string;
			special: string;
		};
	};
	errors: {
		passwordRequired: string;
		passwordMinLength: string;
		passwordUppercase: string;
		passwordLowercase: string;
		passwordNumber: string;
		passwordSpecial: string;
		confirmPasswordRequired: string;
		confirmPasswordMatch: string;
	};
}

export interface ResetPasswordFormProps
	extends React.ComponentPropsWithoutRef<"div"> {
	cfTurnstileSiteKey: string;
	generalError?: string;
	handleResetPassword?: <T extends ResetPasswordFormData>(
		data: T,
		reset: UseFormReset<T>,
	) => Promise<void>;
	isLoading?: boolean;
	successMessage?: string;
	t: ResetPasswordFormT;
}

export function ResetPasswordForm({
	cfTurnstileSiteKey,
	className,
	generalError,
	handleResetPassword,
	isLoading = false,
	successMessage,
	t,
	...props
}: ResetPasswordFormProps) {
	const resetPasswordSchema = z
		.object({
			password: z
				.string()
				.min(1, t.errors.passwordRequired)
				.min(8, t.errors.passwordMinLength)
				.regex(/[A-Z]/, t.errors.passwordUppercase)
				.regex(/[a-z]/, t.errors.passwordLowercase)
				.regex(/[0-9]/, t.errors.passwordNumber)
				.regex(/[^A-Za-z0-9]/, t.errors.passwordSpecial),
			confirmPassword: z.string().min(1, t.errors.confirmPasswordRequired),
		})
		.refine((data) => data.password === data.confirmPassword, {
			message: t.errors.confirmPasswordMatch,
			path: ["confirmPassword"],
		});
	const {
		register,
		handleSubmit,
		formState: { errors },
		reset,
	} = useForm<ResetPasswordFormData>({
		resolver: zodResolver(resetPasswordSchema),
	});
	const onSubmit = handleSubmit(async (data) => {
		data.cfTurnstileToken = turnstileRef.current?.getResponse() || "";
		setPassword("");
		await handleResetPassword?.(data, reset);
	});
	const [currentError, setCurrentError] = useState(generalError);
	const [currentSuccess, setCurrentSuccess] = useState(successMessage);
	const turnstileRef = useRef<TurnstileInstance>(null);
	const [password, setPassword] = useState("");

	// Update messages when props change
	useEffect(() => {
		setCurrentError(generalError);
		setCurrentSuccess(successMessage);
	}, [generalError, successMessage]);

	const clearMessages = () => {
		setCurrentError("");
		setCurrentSuccess("");
	};

	return (
		<div className={cn("flex flex-col gap-6", className)} {...props}>
			<Card>
				<CardHeader className="text-center">
					<CardTitle className="text-xl">{t.title}</CardTitle>
					<CardDescription>{t.description}</CardDescription>
				</CardHeader>

				<CardContent>
					<form onSubmit={onSubmit} className="space-y-6">
						{currentError && (
							<div className="rounded-md bg-destructive/15 p-3 text-sm text-destructive">
								{currentError}
							</div>
						)}

						{currentSuccess && (
							<div className="rounded-md bg-green-500/15 p-3 text-sm text-green-600 dark:text-green-500">
								{currentSuccess}
							</div>
						)}

						<div className="space-y-4">
							<div>
								<Label htmlFor="password">{t.password}</Label>
								<div className="mt-1 space-y-2">
									<Password
										placeholder={t.passwordPlaceholder}
										registration={register("password", {
											onChange: (e) => {
												setPassword(e.target.value);
												clearMessages();
											},
										})}
										error={!!errors.password}
										value={password}
										onValueChange={setPassword}
										showStrength
										strengthT={t.passwordStrength}
										disabled={isLoading}
										aria-invalid={!!errors.password}
										aria-describedby={
											errors.password ? "password-error" : undefined
										}
									/>

									{errors.password && (
										<p className="text-sm font-medium text-destructive mt-2">
											{errors.password.message}
										</p>
									)}
								</div>
							</div>

							<div className="space-y-2">
								<Label htmlFor="confirmPassword">{t.confirmPassword}</Label>
								<div className="mt-1 space-y-2">
									<Password
										placeholder={t.confirmPasswordPlaceholder}
										registration={register("confirmPassword", {
											onChange: clearMessages,
										})}
										error={!!errors.confirmPassword}
										disabled={isLoading}
										aria-invalid={!!errors.confirmPassword}
										aria-describedby={
											errors.confirmPassword
												? "confirm-password-error"
												: undefined
										}
									/>

									{errors.confirmPassword && (
										<p className="text-sm font-medium text-destructive mt-2">
											{errors.confirmPassword.message}
										</p>
									)}
								</div>
							</div>

							<div className="pt-4">
								<Turnstile
									ref={turnstileRef}
									className="text-center"
									options={{ size: "invisible" }}
									siteKey={cfTurnstileSiteKey}
								/>

								<Button type="submit" className="w-full" disabled={isLoading}>
									{t.resetPassword}
								</Button>
							</div>
						</div>

						<div className="text-center text-sm">
							<Link
								aria-disabled={isLoading}
								className="text-muted-foreground hover:text-foreground"
								tabIndex={isLoading ? -1 : undefined}
								to="/sign-in"
							>
								{t.backToSignIn}
							</Link>
						</div>
					</form>
				</CardContent>
			</Card>
		</div>
	);
}
