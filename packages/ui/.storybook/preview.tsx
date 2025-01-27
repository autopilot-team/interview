import { withThemeByClassName } from "@storybook/addon-themes";
import type { Decorator, Preview } from "@storybook/react";
import React from "react";
import "../src/styles/globals.css";

const withThemeContainer: Decorator = (Story, context) => (
	<div className="min-h-screen bg-background text-foreground">
		<div className="container flex min-h-screen items-center justify-center">
			<div>
				<Story {...context} />
			</div>
		</div>
	</div>
);

export const decorators = [
	withThemeByClassName({
		themes: {
			light: "light",
			dark: "dark",
		},
		defaultTheme: "light",
	}),
	withThemeContainer,
];

const preview: Preview = {
	parameters: {
		controls: {
			matchers: {
				color: /(background|color)$/i,
				date: /Date$/i,
			},
		},
		options: {
			storySort: {
				method: "alphabetical",
			},
		},
	},
};

export default preview;
