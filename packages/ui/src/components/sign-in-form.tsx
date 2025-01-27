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
import { useRef, useState } from "react";
import { useForm } from "react-hook-form";
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
	t: SignInFormT;
	isLoading?: boolean;
	handleSignIn?: (data: SignInFormData) => void;
}

export function SignInForm({
	cfTurnstileSiteKey,
	className,
	t,
	isLoading = false,
	handleSignIn,
	...props
}: SignInFormProps) {
	const turnstileRef = useRef<TurnstileInstance>(null);
	const [showPassword, setShowPassword] = useState(false);
	const signInSchema = z.object({
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
		register,
		handleSubmit,
		formState: { errors },
	} = useForm<SignInFormData>({
		resolver: zodResolver(signInSchema),
	});
	const onSubmit = handleSubmit((data) => {
		data.cfTurnstileToken = turnstileRef.current?.getResponse() || "";
		handleSignIn?.(data);
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
								<div className="flex items-center justify-between">
									<Label htmlFor="password">{t.password}</Label>
									<Link
										to="/forgot-password"
										className="text-sm text-muted-foreground hover:text-foreground"
										tabIndex={isLoading ? -1 : undefined}
										aria-disabled={isLoading}
									>
										{t.forgotPassword}
									</Link>
								</div>

								<div className="space-y-2">
									<div className="relative">
										<Input
											{...register("password")}
											type={showPassword ? "text" : "password"}
											className={cn(
												errors.password &&
													"ring-2 ring-destructive ring-offset-2",
											)}
											placeholder={t.passwordPlaceholder}
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
								/>

								<Button type="submit" className="w-full" disabled={isLoading}>
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
