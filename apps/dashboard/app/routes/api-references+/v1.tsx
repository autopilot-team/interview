import { ApiReferenceReact } from "@scalar/api-reference-react";

export default function Component() {
	return (
		<ApiReferenceReact
			configuration={{
				customCss: ":root { --scalar-custom-header-height: 4.15rem; }",
				hideModels: true,
				url: `${import.meta.env.VITE_API_URL || "http://localhost:3001"}/v1/openapi.json`,
			}}
		/>
	);
}
