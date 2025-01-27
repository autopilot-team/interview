import { appendFileSync, createReadStream } from "node:fs";
import { createInterface } from "node:readline";

const readInterface = createInterface({
	input: createReadStream(process.argv[2]),
	output: process.stdout,
	terminal: false,
});

const stripQuotes = (str) =>
	str.startsWith('"') || str.startsWith("'") ? str.slice(1, -1) : str;

const replaceEnvVars = (str) => {
	const value = str
		.replaceAll(/\$\{([a-zA-Z0-9_]+):\+:\$[a-zA-Z0-9_]+\}/g, (_, key) =>
			((v) => (v ? `:${v}` : ""))(process.env[key]),
		)
		.replaceAll(/\$\{([a-zA-Z0-9_]+)\}/g, (_, key) => process.env[key] ?? "")
		.replaceAll(/\$([a-zA-Z0-9_]+)/g, (_, key) => process.env[key] ?? "");
	console.error("FOO", str, value);
	return value;
};

readInterface.on("line", (line) => {
	const match = line.match(/^export ([^=]+)=(.*)$/);

	if (match) {
		const [_, key, value_] = match;
		const value = stripQuotes(value_);

		if (key === "PATH") {
			const paths = value
				.replaceAll("${PATH:+:$PATH}", "")
				.replaceAll("$PATH", "")
				.replaceAll("${PATH}", "")
				.split(":");

			for (const path of paths) {
				appendFileSync(process.env.GITHUB_PATH, `${path}\n`);
			}
		} else {
			const v = replaceEnvVars(value);
			appendFileSync(process.env.GITHUB_ENV, `${key}=${v}\n`);
		}
	}
});
