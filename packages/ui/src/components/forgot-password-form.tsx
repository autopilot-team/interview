import { Button } from "@autopilot/ui/components/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@autopilot/ui/components/card";
import { Input } from "@autopilot/ui/components/input";
import { Label } from "@autopilot/ui/components/label";
import { cn } from "@autopilot/ui/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { Turnstile, type TurnstileInstance } from "@marsidev/react-turnstile";
import { useRef } from "react";
import { useForm } from "react-hook-form";
import { Link } from "react-router";
import { z } from "zod";

type ForgotPasswordFormData = {
	cfTurnstileToken: string;
	email: string;
};

export interface ForgotPasswordFormT {
	title: string;
	description: string;
	email: string;
	emailPlaceholder: string;
	resetPassword: string;
	backToSignIn: string;
	errors: {
		emailRequired: string;
		emailInvalid: string;
	};
}

export interface ForgotPasswordFormProps
	extends React.ComponentPropsWithoutRef<"div"> {
	cfTurnstileSiteKey: string;
	t: ForgotPasswordFormT;
	isLoading?: boolean;
	handleForgotPassword?: (data: ForgotPasswordFormData) => void;
}

export function ForgotPasswordForm({
	cfTurnstileSiteKey,
	className,
	t,
	isLoading = false,
	handleForgotPassword,
	...props
}: ForgotPasswordFormProps) {
	const forgotPasswordSchema = z.object({
		email: z
			.string()
			.min(1, t.errors.emailRequired)
			.email(t.errors.emailInvalid),
	});
	const {
		register,
		handleSubmit,
		formState: { errors },
	} = useForm<ForgotPasswordFormData>({
		resolver: zodResolver(forgotPasswordSchema),
	});
	const onSubmit = handleSubmit((data) => {
		data.cfTurnstileToken = turnstileRef.current?.getResponse() || "";
		handleForgotPassword?.(data);
	});
	const turnstileRef = useRef<TurnstileInstance>(null);

	return (
		<div className={cn("flex flex-col gap-6", className)} {...props}>
			<Card>
				<CardHeader className="text-center">
					<CardTitle className="text-xl">{t.title}</CardTitle>
					<CardDescription>{t.description}</CardDescription>
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
								to="/sign-in"
								className="text-muted-foreground hover:text-foreground"
								tabIndex={isLoading ? -1 : undefined}
								aria-disabled={isLoading}
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
