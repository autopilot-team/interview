import { GalleryVerticalEnd } from "@autopilot/ui/components/icons";
import { cn } from "@autopilot/ui/lib/utils";
import { Link } from "react-router";

const sizeVariants = {
	sm: {
		icon: "h-6 w-6",
		iconInner: "size-4",
		text: "text-base",
	},
	md: {
		icon: "h-8 w-8",
		iconInner: "size-5",
		text: "text-lg",
	},
	lg: {
		icon: "h-10 w-10",
		iconInner: "size-6",
		text: "text-2xl",
	},
} as const;

interface BrandProps extends Omit<React.ComponentPropsWithoutRef<"a">, "href"> {
	showName?: boolean;
	size?: keyof typeof sizeVariants;
	href?: string;
}

export function Brand({
	className,
	showName = true,
	size = "sm",
	href = "/",
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
			{...props}
		>
			<div
				className={cn(
					"flex items-center justify-center rounded-md bg-primary text-primary-foreground",
					sizes.icon,
				)}
			>
				<GalleryVerticalEnd className={sizes.iconInner} />
			</div>
			{showName && "Autopilot"}
		</Link>
	);
}
