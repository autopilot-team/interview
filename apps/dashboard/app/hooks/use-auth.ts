interface User {
	id: string;
	email: string;
	name: string;
}

interface Entity {
	id: string;
	name: string;
}

interface Organization extends Entity {
	platformId: string;
}

interface Account extends Entity {
	organizationId: string;
}

interface ActiveContext {
	platformId: string | null;
	organizationId: string | null;
	accountId: string | null;
	permissions: Permission[];
	role: Role;
}

interface UserContext {
	user: User;
	activeContext: ActiveContext;
	available: {
		platforms: Entity[];
		organizations: Organization[];
		accounts: Account[];
	};
}

type Permission =
	| "read:payments"
	| "write:payments"
	| "read:accounts"
	| "write:accounts";

type Role = "owner" | "admin" | "member" | "viewer";

export function useAuth() {}
