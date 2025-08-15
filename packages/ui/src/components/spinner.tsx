import { cn } from "@autopilot/ui/lib/utils";
import { cva, type VariantProps } from "class-variance-authority";

const spinnerVariants = cva(
	"inline-block animate-spin rounded-full border-solid border-current border-r-transparent motion-reduce:animate-[spin_1.5s_linear_infinite]",
	{
		variants: {
			size: {
				xs: "h-3 w-3 border",
				sm: "h-4 w-4 border",
				md: "h-6 w-6 border-2",
				lg: "h-8 w-8 border-2",
				xl: "h-12 w-12 border-3",
			},
			variant: {
				default: "text-primary",
				muted: "text-muted-foreground",
				destructive: "text-destructive",
			},
		},
		defaultVariants: {
			size: "md",
			variant: "default",
		},
	},
);

export interface SpinnerProps
	extends React.HTMLAttributes<HTMLDivElement>,
		VariantProps<typeof spinnerVariants> {
	/**
	 * The label for screen readers
	 * @default "Loading..."
	 */
	label?: string;
}

export function Spinner({
	className,
	size,
	variant,
	label = "Loading...",
	...props
}: SpinnerProps) {
	return (
		<div className={cn("inline-flex", className)} {...props}>
			<div className={cn(spinnerVariants({ size, variant }))} />
			<span className="sr-only">{label}</span>
		</div>
	);
}
