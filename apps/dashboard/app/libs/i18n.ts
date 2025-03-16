import i18next, { type InitOptions } from "i18next";
import I18nextBrowserLanguageDetector from "i18next-browser-languagedetector";
import I18nextHttpBackend from "i18next-http-backend";
import { initReactI18next } from "react-i18next";
import type common from "../../public/locales/en/common.json";
import type identity from "../../public/locales/en/identity.json";

const defaultNS = "common";
const ns = [defaultNS, "identity"];
const supportedLngs = ["en", "zh-CN", "zh-TW"];

declare module "i18next" {
	interface CustomTypeOptions {
		defaultNS: typeof defaultNS;
		resources: {
			common: typeof common;
			identity: typeof identity;
		};
	}
}

const config = {
	backend: {
		loadPath: "/locales/{{lng}}/{{ns}}.json",
	},
	defaultNS,
	detection: {
		order: ["navigator"],
	},
	fallbackLng: "en",
	interpolation: {
		escapeValue: false,
	},
	ns,
	react: {
		useSuspense: false,
	},
	supportedLngs,
} satisfies InitOptions;

i18next
	.use(I18nextBrowserLanguageDetector)
	.use(I18nextHttpBackend)
	.use(initReactI18next);

if (import.meta.env.MODE !== "production") {
	const { HMRPlugin } = await import("i18next-hmr/plugin");

	i18next.use(
		new HMRPlugin({
			vite: {
				client: true,
			},
		}),
	);
}

export async function initI18n() {
	await i18next.init(config);
}
