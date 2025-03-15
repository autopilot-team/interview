import { Button } from "@autopilot/ui/components/button";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@autopilot/ui/components/card";
import { Input } from "@autopilot/ui/components/input";
import { Label } from "@autopilot/ui/components/label";
import { Password } from "@autopilot/ui/components/password";
import { cn } from "@autopilot/ui/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { Turnstile, type TurnstileInstance } from "@marsidev/react-turnstile";
import { useEffect, useRef, useState } from "react";
import { type UseFormReset, useForm } from "react-hook-form";
import { Link } from "react-router";
import { z } from "zod";

type SignInFormData = {
	cfTurnstileToken: string;
	email: string;
	password: string;
};

export interface SignInFormT {
	title: string;
	email: string;
	emailPlaceholder: string;
	password: string;
	passwordPlaceholder: string;
	forgotPassword: string;
	signIn: string;
	noAccount: string;
	signUp: string;
	termsText: string;
	termsButton: string;
	termsOfService: string;
	privacyPolicy: string;
	errors: {
		emailRequired: string;
		emailInvalid: string;
		passwordRequired: string;
		passwordMinLength: string;
	};
}

export interface SignInFormProps extends React.ComponentPropsWithoutRef<"div"> {
	cfTurnstileSiteKey: string;
	generalError?: string;
	handleSignIn?: <T extends SignInFormData>(
		data: T,
		reset: UseFormReset<T>,
	) => Promise<void>;
	isLoading?: boolean;
	successMessage?: string;
	t: SignInFormT;
}

export function SignInForm({
	cfTurnstileSiteKey,
	className,
	generalError,
	handleSignIn,
	isLoading = false,
	successMessage,
	t,
	...props
}: SignInFormProps) {
	const turnstileRef = useRef<TurnstileInstance>(null);
	const [currentError, setCurrentError] = useState(generalError);
	const [currentSuccess, setCurrentSuccess] = useState(successMessage);
	const signInSchema = z.object({
		cfTurnstileToken: z.string().optional(),
		email: z
			.string()
			.min(1, t.errors.emailRequired)
			.email(t.errors.emailInvalid),
		password: z
			.string()
			.min(1, t.errors.passwordRequired)
			.min(8, t.errors.passwordMinLength),
	});
	const {
		formState: { errors },
		handleSubmit,
		register,
		reset,
	} = useForm({
		resolver: zodResolver(signInSchema),
	});
	const onSubmit = handleSubmit(async (data) => {
		data.cfTurnstileToken = turnstileRef.current?.getResponse() || "";
		await handleSignIn?.(
			data as SignInFormData,
			reset as UseFormReset<SignInFormData>,
		);
	});
	const [isCfTurnstileLoading, setIsCfTurnstileLoading] = useState(true);
	const isFormLoading = isLoading || isCfTurnstileLoading;

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
								<Label htmlFor="email">{t.email}</Label>
								<div className="mt-1 space-y-2">
									<Input
										{...register("email", {
											onChange: clearMessages,
										})}
										type="email"
										placeholder={t.emailPlaceholder}
										className={cn(
											errors.email && "ring-2 ring-destructive ring-offset-2",
										)}
										disabled={isLoading}
										aria-invalid={!!errors.email}
										aria-describedby={errors.email ? "email-error" : undefined}
									/>

									{errors.email && (
										<p className="text-sm font-medium text-destructive mt-2">
											{errors.email.message}
										</p>
									)}
								</div>
							</div>

							<div className="space-y-2">
								<div className="flex items-center justify-between">
									<Label htmlFor="password">{t.password}</Label>
									<Link
										aria-disabled={isLoading}
										className="text-sm text-muted-foreground hover:text-foreground"
										tabIndex={-1}
										to="/forgot-password"
									>
										{t.forgotPassword}
									</Link>
								</div>

								<div>
									<Password
										aria-describedby={
											errors.password ? "password-error" : undefined
										}
										aria-invalid={!!errors.password}
										disabled={isLoading}
										error={!!errors.password}
										isNewPassword={false}
										placeholder={t.passwordPlaceholder}
										registration={register("password", {
											onChange: clearMessages,
										})}
									/>

									{errors.password && (
										<p className="text-sm font-medium text-destructive mt-2">
											{errors.password.message}
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
									onSuccess={() => setIsCfTurnstileLoading(false)}
								/>

								<Button
									type="submit"
									className="w-full"
									disabled={isFormLoading}
								>
									{t.signIn}
								</Button>
							</div>
						</div>

						<div className="text-center text-sm">
							<span className="text-muted-foreground">{t.noAccount}</span>{" "}
							<Link
								to="/sign-up"
								className="font-medium text-primary hover:text-primary/90"
								tabIndex={isLoading ? -1 : undefined}
								aria-disabled={isLoading}
							>
								{t.signUp}
							</Link>
						</div>
					</form>
				</CardContent>
			</Card>

			<div className="text-balance text-center text-xs text-muted-foreground">
				{t.termsText.split(/\{(terms|privacy|button)\}/).map((part, index) => {
					if (part === "terms") {
						return (
							<Link
								key="terms"
								to="/terms-of-service"
								className="text-primary hover:text-primary/90"
								tabIndex={isLoading ? -1 : undefined}
								aria-disabled={isLoading}
							>
								{t.termsOfService}
							</Link>
						);
					}

					if (part === "privacy") {
						return (
							<Link
								key="privacy"
								to="/privacy-policy"
								className="text-primary hover:text-primary/90"
								tabIndex={isLoading ? -1 : undefined}
								aria-disabled={isLoading}
							>
								{t.privacyPolicy}
							</Link>
						);
					}

					if (part === "button") {
						return (
							<span key="button" className="font-medium text-primary">
								{t.termsButton}
							</span>
						);
					}

					return <span key={`text-${part}`}>{part}</span>;
				})}
			</div>
		</div>
	);
}
