/**
 * Returns a country flag emoji for the given country code.
 *
 * @param countryCode - The country code.
 * @returns The country flag emoji.
 */
export function getCountryFlag(countryCode: string): string {
	if (!countryCode) return "🏳️";

	const codePoints = [...countryCode.toUpperCase()]
		.map((char) => char.charCodeAt(0) + 127397)
		.map((codePoint) => String.fromCodePoint(codePoint))
		.join("");

	return codePoints || "🏳️";
}
