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

export type SignUpFormData = {
	cfTurnstileToken: string;
	confirmPassword: string;
	email: string;
	name: string;
	password: string;
};

export interface SignUpFormT {
	title: string;
	name: string;
	namePlaceholder: string;
	email: string;
	emailPlaceholder: string;
	password: string;
	passwordPlaceholder: string;
	confirmPassword: string;
	confirmPasswordPlaceholder: string;
	signUp: string;
	haveAccount: string;
	signIn: string;
	termsText: string;
	termsButton: string;
	termsOfService: string;
	privacyPolicy: string;
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
		nameRequired: string;
		nameMinLength: string;
		emailRequired: string;
		emailInvalid: string;
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

export interface SignUpFormProps extends React.ComponentPropsWithoutRef<"div"> {
	cfTurnstileSiteKey: string;
	generalError?: string;
	handleSignUp?: <T extends SignUpFormData>(
		data: T,
		reset: UseFormReset<T>,
	) => Promise<void>;
	isLoading?: boolean;
	successMessage?: string;
	t: SignUpFormT;
}

export function SignUpForm({
	cfTurnstileSiteKey,
	className,
	generalError,
	handleSignUp,
	isLoading = false,
	successMessage,
	t,
	...props
}: SignUpFormProps) {
	const turnstileRef = useRef<TurnstileInstance>(null);
	const [password, setPassword] = useState("");
	const [currentError, setCurrentError] = useState(generalError);
	const [currentSuccess, setCurrentSuccess] = useState(successMessage);
	const signUpSchema = z
		.object({
			cfTurnstileToken: z.string().optional(),
			name: z
				.string()
				.min(1, t.errors.nameRequired)
				.min(2, t.errors.nameMinLength),
			email: z
				.string()
				.min(1, t.errors.emailRequired)
				.email(t.errors.emailInvalid),
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
	} = useForm({
		resolver: zodResolver(signUpSchema),
	});
	const onSubmit = handleSubmit(async (data) => {
		data.cfTurnstileToken = turnstileRef.current?.getResponse() || "";
		setPassword("");
		await handleSignUp?.(
			data as SignUpFormData,
			reset as UseFormReset<SignUpFormData>,
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
								<Label htmlFor="name">{t.name}</Label>
								<div className="mt-1 space-y-2">
									<Input
										{...register("name", {
											onChange: clearMessages,
										})}
										type="text"
										placeholder={t.namePlaceholder}
										className={cn(
											errors.name && "ring-2 ring-destructive ring-offset-2",
										)}
										disabled={isLoading}
										aria-invalid={!!errors.name}
										aria-describedby={errors.name ? "name-error" : undefined}
									/>

									{errors.name && (
										<p className="text-sm font-medium text-destructive mt-2">
											{errors.name.message}
										</p>
									)}
								</div>
							</div>

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

							<div>
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
									onSuccess={() => setIsCfTurnstileLoading(false)}
								/>

								<Button
									type="submit"
									className="w-full"
									disabled={isFormLoading}
								>
									{t.signUp}
								</Button>
							</div>
						</div>

						<div className="text-center text-sm">
							<span className="text-muted-foreground">{t.haveAccount}</span>{" "}
							<Link
								to="/sign-in"
								className="font-medium text-primary hover:text-primary/90"
								tabIndex={isLoading ? -1 : undefined}
								aria-disabled={isLoading}
							>
								{t.signIn}
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
