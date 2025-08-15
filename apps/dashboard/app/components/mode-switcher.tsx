import type { OperationMode } from "@autopilot/api";
import { getOperationMode, setOperationMode } from "@autopilot/api";
import { Label } from "@autopilot/ui/components/label";
import { Switch } from "@autopilot/ui/components/switch";
import { cn } from "@autopilot/ui/lib/utils";
import { createContext, useContext, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

export interface ModeContextType {
	isTestMode: boolean;
	mode: OperationMode;
	setIsTestMode: (isTestMode: boolean) => void;
	setMode: (mode: OperationMode) => void;
}

const ModeContext = createContext<ModeContextType | undefined>(undefined);

export function ModeSwitcherProvider({
	children,
}: {
	children: React.ReactNode;
}) {
	const [mode, setModeState] = useState<OperationMode>(getOperationMode());
	const [isTestMode, setIsTestMode] = useState<boolean>(mode === "test");

	const setMode = async (newMode: OperationMode) => {
		try {
			setIsTestMode(newMode === "test");
			setModeState(newMode);
			await setOperationMode(newMode);
			window.location.href = window.location.pathname;
		} catch (error) {
			setIsTestMode(mode === "test");
			setModeState(mode);
			throw error;
		} finally {
		}
	};

	// Sync initial mode
	useEffect(() => {
		const currentMode = getOperationMode();
		if (currentMode !== mode) {
			setModeState(currentMode);
		}
	}, [mode]);

	return (
		<ModeContext.Provider value={{ isTestMode, mode, setMode, setIsTestMode }}>
			{children}
		</ModeContext.Provider>
	);
}

export function useMode() {
	const context = useContext(ModeContext);
	if (context === undefined) {
		throw new Error("useMode must be used within a ModeSwitcherProvider");
	}

	return context;
}

export interface ModeSwitcherProps
	extends React.HTMLAttributes<HTMLDivElement> {}

export function ModeSwitcher({ className }: ModeSwitcherProps) {
	const { mode, setMode } = useMode();
	const { t } = useTranslation(["common"]);

	return (
		<section
			className={cn("flex items-center gap-2 my-1.5", className)}
			aria-label={t("common:mode.switcher", "Mode Switcher")}
		>
			<Switch
				checked={mode === "live"}
				className="data-[state=checked]:bg-green-500"
				onCheckedChange={(checked: boolean) =>
					setMode(checked ? "live" : "test")
				}
			/>

			<Label
				htmlFor="mode-switch"
				className={cn(
					"text-sm font-medium leading-none",
					"peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
				)}
			>
				{mode === "live"
					? t("common:mode.live", "Live Mode")
					: t("common:mode.test", "Test Mode")}
			</Label>
		</section>
	);
}
