import { cn } from "@autopilot/ui/lib/utils";
import { Link } from "react-router";

const sizeVariants = {
	sm: {
		icon: "size-6",
		iconInner: "size-4",
		text: "text-base",
	},
	md: {
		icon: "size-8",
		iconInner: "size-5",
		text: "text-lg",
	},
	lg: {
		icon: "size-10",
		iconInner: "size-6",
		text: "text-2xl",
	},
} as const;

interface BrandProps extends Omit<React.ComponentPropsWithoutRef<"a">, "href"> {
	showName?: boolean;
	size?: keyof typeof sizeVariants;
	href?: string;
	spaNavigation?: boolean;
}

export function Brand({
	className,
	showName = true,
	size = "sm",
	href = "/",
	spaNavigation = true,
	...props
}: BrandProps) {
	const sizes = sizeVariants[size];

	return (
		<Link
			to={href}
			className={cn(
				"flex items-center gap-2 font-medium",
				sizes.text,
				className,
			)}
			onClick={(e) => {
				if (!spaNavigation) {
					e.preventDefault();
					window.location.href = href;
					return;
				}
			}}
			{...props}
		>
			<div
				className={cn(
					"flex items-center justify-center rounded-md bg-primary text-primary-foreground",
					sizes.icon,
				)}
			>
				<img
					src="https://assets.autopilot.is/logo.png"
					alt="Autopilot"
					className={sizes.iconInner}
				/>
			</div>

			{showName && "Autopilot"}
		</Link>
	);
}
