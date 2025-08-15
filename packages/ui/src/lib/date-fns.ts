export * from "date-fns";

import { enUS, zhCN, zhTW } from "date-fns/locale";

export function getDateFnsLocale(locale: string) {
	if (locale === "zh-CN") return zhCN;
	if (locale === "zh-TW") return zhTW;

	return enUS;
}
