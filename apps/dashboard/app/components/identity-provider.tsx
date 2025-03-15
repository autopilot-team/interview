import { api, type v1 } from "@autopilot/api";
import {
	type ReactNode,
	createContext,
	useContext,
	useEffect,
	useState,
} from "react";
import { useParams } from "react-router";

export enum AuthState {
	INITIALIZING = "INITIALIZING",
	CHECKING_SESSION = "CHECKING_SESSION",
	AUTHENTICATED = "AUTHENTICATED",
	UNAUTHENTICATED = "UNAUTHENTICATED",
	ERROR = "ERROR",
}

export type SignInBody =
	v1.paths["/v1/identity/sign-in"]["post"]["requestBody"]["content"]["application/json"];
export type SignUpBody =
	v1.paths["/v1/identity/sign-up"]["post"]["requestBody"]["content"]["application/json"];
export type ForgotPasswordBody =
	v1.components["schemas"]["ForgotPasswordRequestBody"];
export type ResetPasswordBody = {
	newPassword: string;
	token: string;
};
export type Entity = v1.components["schemas"]["Entity"];
export type EntityRole = v1.components["schemas"]["EntityRole"];
export type Me =
	v1.paths["/v1/identity/me"]["get"]["responses"]["200"]["content"]["application/json"];
export type Membership = v1.components["schemas"]["Membership"];
export type SessionUser = v1.components["schemas"]["SessionUser"];
export type UpdatePasswordBody = {
	currentPassword: string;
	newPassword: string;
};

export interface IdentityContextType {
	authState: AuthState;
	isForgottingPassword: boolean;
	isInitializing: boolean;
	isResettingPassword: boolean;
	isSigningIn: boolean;
	isSigningOut: boolean;
	isSigningUp: boolean;
	isUpdatingPassword: boolean;
	isVerifyingTwoFactor: boolean;
	isTwoFactorPending: boolean;
	forgotPassword: (body: ForgotPasswordBody) => Promise<void>;
	resetPassword: (body: ResetPasswordBody) => Promise<void>;
	signIn: (body: SignInBody) => Promise<boolean>;
	signOut: () => Promise<void>;
	signUp: (body: SignUpBody) => Promise<void>;
	updatePassword: (body: UpdatePasswordBody) => Promise<void>;
	verifyTwoFactor: (code: string) => Promise<void>;
	switchEntity: (slug: string, id?: string) => Promise<void>;
	refreshUser: () => Promise<void>;
	user?: SessionUser | null;
	entity?: Entity | null;
	role?: EntityRole | null;
}

const IdentityContext = createContext<IdentityContextType | null>(null);

interface IdentityProviderProps {
	children: ReactNode;
}

export function IdentityProvider({ children }: IdentityProviderProps) {
	const { entity: entitySlug } = useParams();
	const [user, setUser] = useState<SessionUser | null | undefined>(undefined);
	const [entity, setEntity] = useState<Entity | null | undefined>();
	const [role, setRole] = useState<EntityRole | null | undefined>();
	const [authState, setAuthState] = useState<AuthState>(AuthState.INITIALIZING);
	const [isTwoFactorPending, setIsTwoFactorPending] = useState(false);
	const { isPending: isSigningIn, mutateAsync: signInAsync } =
		api.v1.useMutation("post", "/v1/identity/sign-in");
	const { isPending: isSigningUp, mutateAsync: signUpAsync } =
		api.v1.useMutation("post", "/v1/identity/sign-up");
	const { isPending: isSigningOut, mutateAsync: signOutAsync } =
		api.v1.useMutation("delete", "/v1/identity/sign-out");
	const { isPending: isForgottingPassword, mutateAsync: forgotPasswordAsync } =
		api.v1.useMutation("post", "/v1/identity/forgot-password");
	const { isPending: isResettingPassword, mutateAsync: resetPasswordAsync } =
		api.v1.useMutation("post", "/v1/identity/reset-password");
	const { isPending: isUpdatingPassword, mutateAsync: updatePasswordAsync } =
		api.v1.useMutation("post", "/v1/identity/update-password");
	const { isPending: isVerifyingTwoFactor, mutateAsync: verifyTwoFactorAsync } =
		api.v1.useMutation("post", "/v1/identity/verify-two-factor");

	useEffect(() => {
		if (entity) {
			window.localStorage.setItem("entityID", entity?.id);
		}
	}, [entity]);

	const {
		data: session,
		isLoading: isInitializing,
		isRefetching,
		refetch: refetchSession,
	} = api.v1.useQuery("get", "/v1/identity/me", {}, { retry: false });

	// Update user and auth state when session changes
	useEffect(() => {
		if (isInitializing) {
			setAuthState(AuthState.CHECKING_SESSION);
			return;
		}

		if (isRefetching) {
			return;
		}

		if (!session?.user?.id) {
			setUser(null);
			setAuthState(AuthState.UNAUTHENTICATED);
			return;
		}

		setUser(session.user);
		setAuthState(AuthState.AUTHENTICATED);
		setEntity(session.activeEntity);
		setRole(session.entityRole);

		if (session.activeEntity) {
			localStorage.setItem("entityID", session.activeEntity.id);
		}
	}, [session, isInitializing, isRefetching]);

	useEffect(() => {
		const handler = () => {
			if (entity) {
				localStorage.setItem("entityID", entity?.id || "");
			} else {
				localStorage.removeItem("entityID");
			}
		};
		window.addEventListener("focus", handler);
		return () => window.removeEventListener("focus", handler);
	}, [entity]);

	// Handle API errors
	useEffect(() => {
		const handleError = () => {
			setAuthState(AuthState.ERROR);
		};

		window.addEventListener("unhandledrejection", handleError);
		return () => window.removeEventListener("unhandledrejection", handleError);
	}, []);

	async function signIn(body: SignInBody) {
		const response = await signInAsync({ body });
		if (response.isTwoFactorPending) {
			setIsTwoFactorPending(true);
			return false;
		}

		await refetchSession();

		return true;
	}

	async function signUp(body: SignUpBody) {
		await signUpAsync({
			body,
		});
	}

	async function signOut() {
		await signOutAsync({});
		setUser(null);
		setEntity(null);
	}

	async function forgotPassword(body: ForgotPasswordBody) {
		await forgotPasswordAsync({
			body,
		});
	}

	async function resetPassword(body: ResetPasswordBody) {
		await resetPasswordAsync({
			body,
		});
	}

	async function updatePassword(body: UpdatePasswordBody) {
		await updatePasswordAsync({
			body,
		});
	}

	async function verifyTwoFactor(code: string) {
		await verifyTwoFactorAsync({ body: { code } });
		setIsTwoFactorPending(false);
		await refetchSession();
	}

	async function switchEntity(slug: string, id?: string) {
		if (id && entity?.id !== id) {
			localStorage.setItem("entityID", id);
		}

		let path = `/${slug}`;
		if (entitySlug) {
			path += window.location.pathname.replace(`/${entitySlug}`, "");
		}

		window.location.href = path;
	}

	async function refreshUser() {
		await refetchSession();
	}

	return (
		<IdentityContext.Provider
			value={{
				authState,
				isInitializing,
				isSigningIn,
				isSigningOut,
				isSigningUp,
				isForgottingPassword,
				isResettingPassword,
				isTwoFactorPending,
				isUpdatingPassword,
				isVerifyingTwoFactor,
				signIn,
				signOut,
				signUp,
				forgotPassword,
				resetPassword,
				updatePassword,
				verifyTwoFactor,
				switchEntity,
				refreshUser,
				user,
				entity,
				role,
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
