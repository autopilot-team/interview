import { ApiReferenceReact } from "@scalar/api-reference-react";
import "@scalar/api-reference-react/style.css";

export default function Component() {
	return (
		<ApiReferenceReact
			configuration={{
				hideModels: true,
				spec: {
					url: `${import.meta.env.VITE_API_URL || "http://localhost:3001"}/v1/openapi.json`,
				},
			}}
		/>
	);
}
