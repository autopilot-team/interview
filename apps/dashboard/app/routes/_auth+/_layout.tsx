import { useIdentity } from "@/components/identity-provider";
import { Navigate, Outlet } from "react-router";

export default function Component() {
	const { user } = useIdentity();

	if (user) {
		return <Navigate to="/" replace />;
	}

	return <Outlet />;
}
