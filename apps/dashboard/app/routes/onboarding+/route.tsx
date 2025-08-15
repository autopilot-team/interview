import { Link, Navigate } from "react-router";
import { useIdentity } from "@/components/identity-provider";

export default function Component() {
	const { user, entity } = useIdentity();

	if (!user) {
		return <Navigate to="/sign-in" replace />;
	}
	if (entity) {
		return <Navigate to={`/${entity.slug}`} replace />;
	}

	return (
		<div className="p-4">
			<h1 className="text-xl">Show Onboarding</h1>

			<Link to={"sign-out"} className="text-blue-600 underline">
				Sign out
			</Link>
		</div>
	);
}
