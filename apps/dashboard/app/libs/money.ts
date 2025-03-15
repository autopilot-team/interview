/**
 * Formats a currency amount for display. Use dinero.js for more complex
 * currency/amount handling.
 *
 * @param amount - The amount to format.
 * @param currency - The currency code.
 * @returns The formatted currency amount.
 */
export function formatCurrency(
	amount: number | string | undefined,
	currency: string,
): string {
	if (amount === undefined) return "";

	const amountInMinorUnits =
		typeof amount === "string"
			? Math.round(Number.parseFloat(amount))
			: Math.round(amount);
	const zeroDecimalCurrencies = [
		"BIF",
		"CLP",
		"DJF",
		"GNF",
		"JPY",
		"KMF",
		"KRW",
		"MGA",
		"PYG",
		"RWF",
		"UGX",
		"VND",
		"VUV",
		"XAF",
		"XOF",
		"XPF",
	];

	// Special case currencies that need custom handling
	const specialCaseCurrencies = {
		ISK: { isZeroDecimal: true, specialHandling: true },
		HUF: { isZeroDecimal: false, specialHandling: true },
		TWD: { isZeroDecimal: false, specialHandling: true },
	};

	const upperCurrency = currency.toUpperCase();
	const isZeroDecimal = zeroDecimalCurrencies.includes(upperCurrency);
	const specialCase =
		specialCaseCurrencies[upperCurrency as keyof typeof specialCaseCurrencies];

	let displayAmount: number;
	let fractionDigits = 2;

	// Handle standard zero-decimal currencies
	if (isZeroDecimal) {
		displayAmount = amountInMinorUnits;
		fractionDigits = 0;
	} else if (specialCase) {
		// Handle special cases
		if (specialCase.isZeroDecimal) {
			// ISK case: zero-decimal but needs to be represented as two-decimal
			displayAmount = amountInMinorUnits / 100;
			fractionDigits = 0;
		} else {
			// HUF and TWD: treated as zero-decimal for payouts but can have decimals for display
			const majorUnits = Math.floor(amountInMinorUnits / 100);
			const minorUnits = amountInMinorUnits % 100;

			displayAmount = Number(
				`${majorUnits}.${minorUnits.toString().padStart(2, "0")}`,
			);
		}
	} else {
		// Standard case for currencies with decimals
		const majorUnits = Math.floor(amountInMinorUnits / 100);
		const minorUnits = amountInMinorUnits % 100;

		displayAmount = Number(
			`${majorUnits}.${minorUnits.toString().padStart(2, "0")}`,
		);
	}

	return new Intl.NumberFormat(navigator.language || "en-US", {
		style: "currency",
		currency: currency,
		minimumFractionDigits: fractionDigits,
		maximumFractionDigits: fractionDigits,
	}).format(displayAmount);
}
