import { Button } from "@autopilot/ui/components/button";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@autopilot/ui/components/card";
import { Input } from "@autopilot/ui/components/input";
import {
	InputOTP,
	InputOTPGroup,
	InputOTPSlot,
} from "@autopilot/ui/components/input-otp";
import { Label } from "@autopilot/ui/components/label";
import { cn } from "@autopilot/ui/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { REGEXP_ONLY_DIGITS_AND_CHARS } from "input-otp";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

type TwoFactorFormData = {
	code: string;
};

export interface TwoFactorFormT {
	title: string;
	code: string;
	codePlaceholder: string;
	backupCodePlaceholder: string;
	verify: string;
	useBackupCode: string;
	useAuthenticator: string;
	backupCodeInstructions: string;
	errors: {
		codeRequired: string;
		codeLength: string;
		backupCodeLength: string;
	};
}

export interface TwoFactorFormProps
	extends React.ComponentPropsWithoutRef<"div"> {
	generalError?: string;
	handleVerify?: (code: string) => Promise<void>;
	isLoading?: boolean;
	t: TwoFactorFormT;
}

export function TwoFactorForm({
	className,
	generalError,
	handleVerify,
	isLoading = false,
	t,
	...props
}: TwoFactorFormProps) {
	const [currentError, setCurrentError] = useState(generalError);
	const [isUsingBackupCode, setIsUsingBackupCode] = useState(false);

	const twoFactorSchema = z.object({
		code: z
			.string()
			.min(1, t.errors.codeRequired)
			.length(
				isUsingBackupCode ? 10 : 6,
				isUsingBackupCode ? t.errors.backupCodeLength : t.errors.codeLength,
			),
	});

	const {
		formState: { errors },
		handleSubmit,
		register,
		setValue,
		reset,
	} = useForm<TwoFactorFormData>({
		resolver: zodResolver(twoFactorSchema),
	});

	const onSubmit = handleSubmit(async (data) => {
		await handleVerify?.(data.code);
	});

	// Update messages when props change
	useEffect(() => {
		setCurrentError(generalError);
	}, [generalError]);

	const clearMessages = () => {
		setCurrentError("");
	};

	const toggleMode = () => {
		setIsUsingBackupCode(!isUsingBackupCode);
		reset(); // Clear the form when switching modes
		clearMessages();
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

						<div className="space-y-4">
							<div>
								<Label
									htmlFor="code"
									className="whitespace-pre-line leading-4 block"
								>
									{isUsingBackupCode ? t.backupCodeInstructions : t.code}
								</Label>

								<div className="mt-6.5 space-y-2">
									<div className="flex justify-center">
										{isUsingBackupCode ? (
											<Input
												id="code"
												className="text-center"
												placeholder={t.backupCodePlaceholder}
												maxLength={10}
												disabled={isLoading}
												{...register("code")}
												onChange={(e) => {
													setValue("code", e.target.value);
													clearMessages();
												}}
											/>
										) : (
											<InputOTP
												maxLength={6}
												disabled={isLoading}
												pattern={REGEXP_ONLY_DIGITS_AND_CHARS}
												onComplete={(value) => {
													setValue("code", value);
													clearMessages();
												}}
											>
												<InputOTPGroup>
													{Array.from({ length: 6 }).map((_, index) => (
														<InputOTPSlot
															key={`otp-${index.toString()}`}
															index={index}
														/>
													))}
												</InputOTPGroup>
											</InputOTP>
										)}
									</div>

									{errors.code && (
										<p className="text-sm font-medium text-destructive mt-2 text-center">
											{errors.code.message}
										</p>
									)}
								</div>
							</div>

							<div className="pt-4 space-y-4">
								<Button type="submit" className="w-full" disabled={isLoading}>
									{t.verify}
								</Button>

								<Button
									type="button"
									variant="link"
									className="w-full"
									onClick={toggleMode}
									disabled={isLoading}
								>
									{isUsingBackupCode ? t.useAuthenticator : t.useBackupCode}
								</Button>
							</div>
						</div>
					</form>
				</CardContent>
			</Card>
		</div>
	);
}
