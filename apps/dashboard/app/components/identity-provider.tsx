import { api, type v1 } from "@autopilot/api";
import {
	type ReactNode,
	createContext,
	useContext,
	useEffect,
	useState,
} from "react";

export type SignInBody =
	v1.paths["/v1/identity/sign-in"]["post"]["requestBody"]["content"]["application/json"];
export type SessionUser = v1.components["schemas"]["SessionUser"];
export type Me =
	v1.paths["/v1/identity/me"]["get"]["responses"]["200"]["content"]["application/json"];

export interface IdentityContextType {
	isInitializing: boolean;
	isSigningIn: boolean;
	isSigningOut: boolean;
	signIn: (body: SignInBody) => Promise<boolean>;
	signOut: () => Promise<void>;
	user?: SessionUser | null;
}

const IdentityContext = createContext<IdentityContextType | null>(null);

interface IdentityProviderProps {
	children: ReactNode;
}

export function IdentityProvider({ children }: IdentityProviderProps) {
	const [user, setUser] = useState<SessionUser | null | undefined>(undefined);
	const { isPending: isSigningIn, mutateAsync: signInAsync } =
		api.v1.useMutation("post", "/v1/identity/sign-in");
	const { isPending: isSigningOut, mutateAsync: signOutAsync } =
		api.v1.useMutation("delete", "/v1/identity/sign-out");
	const { data: session, isLoading: isInitializing } = api.v1.useQuery(
		"get",
		"/v1/identity/me",
		{},
		{ retry: false },
	);

	// Update user when session changes
	useEffect(() => {
		setUser(session?.user?.id ? session.user : null);
	}, [session]);

	async function signIn(body: SignInBody) {
		const response = await signInAsync({
			body,
		});

		if (response.user?.id) {
			setUser(response.user);
			return true;
		}

		return false;
	}

	async function signOut() {
		await signOutAsync({});
		setUser(null);
	}

	return (
		<IdentityContext.Provider
			value={{
				isInitializing,
				isSigningIn,
				isSigningOut,
				signIn,
				signOut,
				user,
			}}
		>
			{children}
		</IdentityContext.Provider>
	);
}
export function useIdentity(): IdentityContextType {
	const context = useContext(IdentityContext);
	if (!context) {
		throw new Error("useIdentity must be used within an IdentityProvider");
	}

	return context;
}
