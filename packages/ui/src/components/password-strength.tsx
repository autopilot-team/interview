import { cn } from "@autopilot/ui/lib/utils";
import { useEffect, useState } from "react";

interface PasswordStrengthRequirements {
	minLength: string;
	mixCase: string;
	number: string;
	special: string;
}

interface PasswordStrengthLabels {
	secure: string;
	moderate: string;
	weak: string;
}

interface PasswordStrengthState {
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

export interface PasswordStrengthProps {
	password: string;
	className?: string;
	disabled?: boolean;
	requirements: PasswordStrengthRequirements;
	strength: PasswordStrengthLabels;
	onStrengthChange?: (strength: PasswordStrengthState) => void;
}

export function PasswordStrength({
	password,
	className,
	disabled,
	requirements,
	strength,
	onStrengthChange,
}: PasswordStrengthProps) {
	const [state, setState] = useState<PasswordStrengthState>({
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
		const reqs = {
			minLength: password.length >= 8,
			hasUppercase: /[A-Z]/.test(password),
			hasLowercase: /[a-z]/.test(password),
			hasNumber: /[0-9]/.test(password),
			hasSpecial: /[^A-Za-z0-9]/.test(password),
		};

		const meetsAllRequirements = Object.values(reqs).every(Boolean);
		const score = Object.values(reqs).filter(Boolean).length;

		let label = "";
		let normalizedScore = 1;

		if (meetsAllRequirements) {
			label = strength.secure;
			normalizedScore = 3;
		} else if (score >= 3) {
			label = strength.moderate;
			normalizedScore = 2;
		} else if (score > 0) {
			label = strength.weak;
			normalizedScore = 1;
		}

		const newState = {
			score: normalizedScore,
			color: "bg-muted",
			label,
			requirements: reqs,
		};

		setState(newState);
		onStrengthChange?.(newState);
	}, [password, strength, onStrengthChange]);

	if (!password) return null;

	const getHint = () => {
		if (!state.requirements.minLength) {
			return requirements.minLength;
		}
		if (!state.requirements.hasUppercase || !state.requirements.hasLowercase) {
			return requirements.mixCase;
		}
		if (!state.requirements.hasNumber) {
			return requirements.number;
		}
		if (!state.requirements.hasSpecial) {
			return requirements.special;
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
						state.score >= 1 ? "bg-red-500" : "bg-muted",
					)}
				/>

				<div
					className={cn(
						"h-1 flex-1 rounded-sm transition-all duration-200",
						state.score >= 2 ? "bg-yellow-500" : "bg-muted",
					)}
				/>

				<div
					className={cn(
						"h-1 flex-1 rounded-sm transition-all duration-200",
						state.score >= 3 ? "bg-green-500" : "bg-muted",
					)}
				/>
			</div>

			{hint && <p className="text-xs text-muted-foreground">{hint}</p>}
		</div>
	);
}
