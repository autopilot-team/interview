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
import { useEffect, useRef, useState } from "react";
import { type UseFormReset, useForm } from "react-hook-form";
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
	generalError?: string;
	handleForgotPassword?: <T extends ForgotPasswordFormData>(
		data: T,
		reset: UseFormReset<T>,
	) => Promise<void>;
	isLoading?: boolean;
	successMessage?: string;
	t: ForgotPasswordFormT;
}

export function ForgotPasswordForm({
	cfTurnstileSiteKey,
	className,
	generalError,
	handleForgotPassword,
	isLoading = false,
	successMessage,
	t,
	...props
}: ForgotPasswordFormProps) {
	const forgotPasswordSchema = z.object({
		cfTurnstileToken: z.string().optional(),
		email: z
			.string()
			.min(1, t.errors.emailRequired)
			.email(t.errors.emailInvalid),
	});
	const {
		register,
		handleSubmit,
		formState: { errors },
		reset,
	} = useForm({
		resolver: zodResolver(forgotPasswordSchema),
	});
	const onSubmit = handleSubmit(async (data) => {
		data.cfTurnstileToken = turnstileRef.current?.getResponse() || "";
		await handleForgotPassword?.(
			data as ForgotPasswordFormData,
			reset as UseFormReset<ForgotPasswordFormData>,
		);
	});
	const [currentError, setCurrentError] = useState(generalError);
	const [currentSuccess, setCurrentSuccess] = useState(successMessage);
	const turnstileRef = useRef<TurnstileInstance>(null);
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
