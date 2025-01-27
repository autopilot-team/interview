import { Alert, AlertTitle } from "@autopilot/ui/components/alert";
import { Button } from "@autopilot/ui/components/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@autopilot/ui/components/card";
import { cn } from "@autopilot/ui/lib/utils";
import { AlertTriangle, ArrowLeft, Check, Home, X } from "lucide-react";
import { Link } from "react-router";

export type MessageCardVariant = "default" | "error" | "success";

export interface MessageCardProps {
	backButton?: {
		icon?: React.ReactNode;
		label?: string;
		show?: boolean;
	};
	description: string;
	footer?: string;
	homeButton?: {
		icon?: React.ReactNode;
		label?: string;
		show?: boolean;
		to?: string;
	};
	stack?: string;
	title: string;
	variant?: MessageCardVariant;
}

export function MessageCard({
	backButton = { label: "Back", show: true },
	description,
	footer,
	homeButton = { label: "Home", show: true, to: "/" },
	stack,
	title,
	variant = "default",
}: MessageCardProps) {
	return (
		<div className="bg-muted min-h-screen w-full flex flex-col items-center justify-center px-6">
			<div className="flex flex-col gap-6 w-full max-w-md">
				<Card>
					<CardHeader className="text-center space-y-2">
						<div className="mx-auto">
							{variant === "error" ? (
								<X className="size-8 text-destructive" />
							) : variant === "success" ? (
								<Check className="size-8 text-green-500" />
							) : (
								<AlertTriangle className="size-8 text-muted-foreground" />
							)}
						</div>

						<CardTitle className="text-xl">{title}</CardTitle>

						<CardDescription className="text-sm text-muted-foreground">
							{description}
						</CardDescription>
					</CardHeader>

					<CardContent
						className={cn({
							"pb-1": !stack && !homeButton.show && !backButton.show,
						})}
					>
						{stack && (
							<div className="space-y-2 mb-6">
								<Alert variant="destructive" className="text-xs font-mono">
									<AlertTitle className="mb-2 text-[0.8rem]">
										Stack Trace
									</AlertTitle>

									<div className="max-h-[200px] overflow-auto whitespace-pre-wrap">
										{stack}
									</div>
								</Alert>
							</div>
						)}

						{homeButton.show || backButton.show ? (
							<div className="space-y-4">
								<div className="flex items-center justify-center gap-2">
									{homeButton.show && (
										<Button asChild className="w-full">
											<Link to={homeButton.to || "/"}>
												{homeButton.icon ? (
													homeButton.icon
												) : (
													<Home className="mr-2 h-4 w-4" />
												)}
												{homeButton.label}
											</Link>
										</Button>
									)}

									{backButton.show && (
										<Button
											variant="outline"
											type="button"
											className="w-full"
											onClick={() => window.history.back()}
										>
											{backButton.icon ? (
												backButton.icon
											) : (
												<ArrowLeft className="mr-2 h-4 w-4" />
											)}
											{backButton.label}
										</Button>
									)}
								</div>
							</div>
						) : null}
					</CardContent>
				</Card>

				{footer ? (
					<div className="text-balance text-center text-xs text-muted-foreground">
						{footer}
					</div>
				) : null}
			</div>
		</div>
	);
}
