import { Button } from "@autopilot/ui/components/button";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@autopilot/ui/components/card";
import { Input } from "@autopilot/ui/components/input";
import { Label } from "@autopilot/ui/components/label";
import { cn } from "@autopilot/ui/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { Turnstile, type TurnstileInstance } from "@marsidev/react-turnstile";
import { EyeIcon, EyeOffIcon } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { type UseFormReset, useForm } from "react-hook-form";
import { Link } from "react-router";
import { z } from "zod";

interface PasswordStrength {
	score: number;
	color: string;
	label: string;
	requirements: {
		minLength: boolean;
		hasUppercase: boolean;
		hasLowercase: boolean;
		hasNumber: boolean;
		hasSpecial: boolean;
	};
}

function PasswordStrengthIndicator({
	password,
	className,
	onStrengthChange,
	t,
	disabled,
}: {
	password: string;
	className?: string;
	onStrengthChange?: (strength: PasswordStrength) => void;
	t: SignUpFormT["passwordStrength"];
	disabled?: boolean;
}) {
	const [strength, setStrength] = useState<PasswordStrength>({
		score: 0,
		color: "bg-muted",
		label: "",
		requirements: {
			minLength: false,
			hasUppercase: false,
			hasLowercase: false,
			hasNumber: false,
			hasSpecial: false,
		},
	});

	useEffect(() => {
		const requirements = {
			minLength: password.length >= 8,
			hasUppercase: /[A-Z]/.test(password),
			hasLowercase: /[a-z]/.test(password),
			hasNumber: /[0-9]/.test(password),
			hasSpecial: /[^A-Za-z0-9]/.test(password),
		};

		const meetsAllRequirements = Object.values(requirements).every(Boolean);
		const score = Object.values(requirements).filter(Boolean).length;

		let label = "";
		let normalizedScore = 1;

		if (meetsAllRequirements) {
			label = t.secure;
			normalizedScore = 3;
		} else if (score >= 3) {
			label = t.moderate;
			normalizedScore = 2;
		} else if (score > 0) {
			label = t.weak;
			normalizedScore = 1;
		}

		const newStrength = {
			score: normalizedScore,
			color: "bg-muted",
			label,
			requirements,
		};

		setStrength(newStrength);
		onStrengthChange?.(newStrength);
	}, [password, onStrengthChange, t]);

	if (!password) return null;

	const getHint = () => {
		if (!strength.requirements.minLength) {
			return t.requirements.minLength;
		}
		if (
			!strength.requirements.hasUppercase ||
			!strength.requirements.hasLowercase
		) {
			return t.requirements.mixCase;
		}
		if (!strength.requirements.hasNumber) {
			return t.requirements.number;
		}
		if (!strength.requirements.hasSpecial) {
			return t.requirements.special;
		}
		return null;
	};

	const hint = getHint();

	return (
		<div className={cn("space-y-2", className, disabled && "opacity-50")}>
			<div className="flex-1 flex gap-0.5">
				<div
					className={cn(
						"h-1 flex-1 rounded-sm transition-all duration-200",
						strength.score >= 1 ? "bg-red-500" : "bg-muted",
					)}
				/>

				<div
					className={cn(
						"h-1 flex-1 rounded-sm transition-all duration-200",
						strength.score >= 2 ? "bg-yellow-500" : "bg-muted",
					)}
				/>

				<div
					className={cn(
						"h-1 flex-1 rounded-sm transition-all duration-200",
						strength.score >= 3 ? "bg-green-500" : "bg-muted",
					)}
				/>
			</div>

			{hint && <p className="text-xs text-muted-foreground">{hint}</p>}
		</div>
	);
}

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
	formRef?: React.Ref<{ reset: () => void }>;
	handleSignUp?: <T extends SignUpFormData>(
		data: T,
		reset: UseFormReset<T>,
	) => void;
	isLoading?: boolean;
	t: SignUpFormT;
}

export function SignUpForm({
	cfTurnstileSiteKey,
	className,
	handleSignUp,
	isLoading = false,
	t,
	...props
}: SignUpFormProps) {
	const turnstileRef = useRef<TurnstileInstance>(null);
	const [password, setPassword] = useState("");
	const [showPassword, setShowPassword] = useState(false);
	const [showConfirmPassword, setShowConfirmPassword] = useState(false);
	const [passwordStrength, setPasswordStrength] = useState<PasswordStrength>({
		score: 0,
		color: "bg-muted",
		label: "",
		requirements: {
			minLength: false,
			hasUppercase: false,
			hasLowercase: false,
			hasNumber: false,
			hasSpecial: false,
		},
	});
	const signUpSchema = z
		.object({
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
	} = useForm<SignUpFormData>({
		resolver: zodResolver(signUpSchema),
	});
	const onSubmit = handleSubmit((data) => {
		data.cfTurnstileToken = turnstileRef.current?.getResponse() || "";
		handleSignUp?.(data, reset);
	});

	return (
		<div className={cn("flex flex-col gap-6", className)} {...props}>
			<Card>
				<CardHeader className="text-center">
					<CardTitle className="text-xl">{t.title}</CardTitle>
				</CardHeader>

				<CardContent>
					<form onSubmit={onSubmit} className="space-y-6">
						<div className="space-y-4">
							<div className="space-y-2">
								<Label htmlFor="name">{t.name}</Label>
								<div className="space-y-2">
									<Input
										{...register("name")}
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

							<div className="space-y-2">
								<Label htmlFor="email">{t.email}</Label>
								<div className="space-y-2">
									<Input
										{...register("email")}
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
								<Label htmlFor="password">{t.password}</Label>
								<div className="space-y-2">
									<div className="relative">
										<Input
											autoComplete="off"
											data-1p-ignore
											{...register("password", {
												onChange: (e) => setPassword(e.target.value),
											})}
											type={showPassword ? "text" : "password"}
											placeholder={t.passwordPlaceholder}
											className={cn(
												errors.password &&
													"ring-2 ring-destructive ring-offset-2",
											)}
											disabled={isLoading}
											aria-invalid={!!errors.password}
											aria-describedby={
												errors.password ? "password-error" : undefined
											}
										/>
										<Button
											type="button"
											variant="ghost"
											size="sm"
											className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
											onClick={() => setShowPassword(!showPassword)}
											disabled={isLoading}
											aria-label={
												showPassword ? "Hide password" : "Show password"
											}
										>
											{showPassword ? (
												<EyeOffIcon className="h-4 w-4 text-muted-foreground" />
											) : (
												<EyeIcon className="h-4 w-4 text-muted-foreground" />
											)}
										</Button>
									</div>

									<PasswordStrengthIndicator
										password={password}
										onStrengthChange={setPasswordStrength}
										t={t.passwordStrength}
										disabled={isLoading}
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
								<div className="space-y-2">
									<div className="relative">
										<Input
											autoComplete="off"
											data-1p-ignore
											{...register("confirmPassword")}
											type={showConfirmPassword ? "text" : "password"}
											placeholder={t.confirmPasswordPlaceholder}
											className={cn(
												errors.confirmPassword &&
													"ring-2 ring-destructive ring-offset-2",
											)}
											disabled={isLoading}
											aria-invalid={!!errors.confirmPassword}
											aria-describedby={
												errors.confirmPassword
													? "confirm-password-error"
													: undefined
											}
										/>
										<Button
											type="button"
											variant="ghost"
											size="sm"
											className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
											onClick={() =>
												setShowConfirmPassword(!showConfirmPassword)
											}
											disabled={isLoading}
											aria-label={
												showConfirmPassword ? "Hide password" : "Show password"
											}
										>
											{showConfirmPassword ? (
												<EyeOffIcon className="h-4 w-4 text-muted-foreground" />
											) : (
												<EyeIcon className="h-4 w-4 text-muted-foreground" />
											)}
										</Button>
									</div>
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
