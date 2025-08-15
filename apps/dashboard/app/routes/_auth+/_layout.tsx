import { Navigate, Outlet } from "react-router";
import { useIdentity } from "@/components/identity-provider";

export default function Component() {
	const { user } = useIdentity();

	if (user) {
		return <Navigate to="/" replace />;
	}

	return <Outlet />;
}
