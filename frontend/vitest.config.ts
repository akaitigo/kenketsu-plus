import path from "node:path";
import { defineConfig } from "vitest/config";

export default defineConfig({
	esbuild: {
		jsx: "automatic",
	},
	test: {
		globals: true,
		environment: "node",
		include: ["src/**/*.test.ts", "src/**/*.test.tsx", "test/**/*.test.ts"],
		exclude: ["test/e2e/**"],
		coverage: {
			provider: "v8",
			reporter: ["text", "lcov"],
			include: ["src/**/*.ts", "src/**/*.tsx"],
			exclude: ["src/**/*.test.*", "src/app/layout.tsx"],
		},
	},
	resolve: {
		alias: {
			"@": path.resolve(__dirname, "./src"),
		},
	},
});
