export const config = {
	cfTurnstileSiteKey:
		import.meta.env.VITE_CF_TURNSTILE_SITE_KEY || "1x00000000000000000000AA",
};

export type Config = typeof config;
