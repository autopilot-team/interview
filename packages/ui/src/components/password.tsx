import { Button } from "@autopilot/ui/components/button";
import { EyeIcon, EyeOffIcon } from "@autopilot/ui/components/icons";
import { Input } from "@autopilot/ui/components/input";
import { PasswordStrength } from "@autopilot/ui/components/password-strength";
import { cn } from "@autopilot/ui/lib/utils";
import { useState } from "react";
import type { UseFormRegisterReturn } from "react-hook-form";

export interface PasswordStrengthT {
	secure: string;
	moderate: string;
	weak: string;
	requirements: {
		minLength: string;
		mixCase: string;
		number: string;
		special: string;
	};
}

export interface PasswordProps
	extends React.ComponentPropsWithoutRef<typeof Input> {
	/** Whether the input is disabled */
	disabled?: boolean;
	/** Whether the input is in an error state */
	error?: boolean;
	/** Whether the input is a new password */
	isNewPassword?: boolean;
	/** Called when password value changes */
	onValueChange?: (value: string) => void;
	/** Form registration object from react-hook-form */
	registration?: Partial<UseFormRegisterReturn>;
	/** Whether to show the password strength indicator */
	showStrength?: boolean;
	/** Password strength translations */
	strengthT?: PasswordStrengthT;
	/** Current password value for strength calculation */
	value?: string;
}

export function Password({
	className,
	disabled,
	error,
	isNewPassword = true,
	onValueChange,
	registration,
	showStrength,
	strengthT,
	value = "",
	...props
}: PasswordProps) {
	const [showPassword, setShowPassword] = useState(false);

	return (
		<div className="mt-1 space-y-2">
			<div className="relative">
				<Input
					{...props}
					{...registration}
					autoComplete={isNewPassword ? "off" : undefined}
					data-1p-ignore={isNewPassword ? "true" : undefined}
					className={cn(error && "ring-2 ring-destructive", className)}
					disabled={disabled}
					onChange={(e) => {
						registration?.onChange?.(e);
						onValueChange?.(e.target.value);
					}}
					type={showPassword ? "text" : "password"}
				/>

				<Button
					aria-label={showPassword ? "Hide password" : "Show password"}
					className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
					disabled={disabled}
					onClick={() => setShowPassword(!showPassword)}
					size="sm"
					type="button"
					variant="ghost"
				>
					{showPassword ? (
						<EyeOffIcon className="h-4 w-4" />
					) : (
						<EyeIcon className="h-4 w-4" />
					)}
				</Button>
			</div>

			{showStrength && strengthT && (
				<PasswordStrength
					password={value}
					requirements={strengthT.requirements}
					strength={strengthT}
					disabled={disabled}
				/>
			)}
		</div>
	);
}
