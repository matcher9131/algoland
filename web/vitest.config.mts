import { fileURLToPath } from "node:url";

import react from "@vitejs/plugin-react";
import { defineConfig } from "vitest/config";

export default defineConfig({
  // Vite resolves the `@/*` paths from tsconfig.json natively; the
  // vite-tsconfig-paths plugin is no longer needed.
  resolve: { tsconfigPaths: true },
  plugins: [react()],
  test: {
    coverage: {
      provider: "v8",
      reportsDirectory: "./coverage",
      reporter: ["text", "html", "lcov"],
      include: ["src/**/*.{ts,tsx}"],
      exclude: ["src/**/*.test.{ts,tsx}", "src/**/*.d.ts", "src/test/**"],
    },
    // Inline projects inherit the plugins and resolve options declared above.
    projects: [
      {
        // Components and anything that touches the DOM.
        test: {
          name: "browser",
          environment: "jsdom",
          globals: true,
          setupFiles: ["./vitest.setup.ts"],
          include: ["src/**/*.test.tsx"],
        },
      },
      {
        // Server-side modules.
        resolve: {
          alias: {
            // `server-only` throws on import unless the bundler picks its
            // `react-server` export condition. Node resolves it directly here
            // and would always take the throwing branch, so point it at the
            // empty module that condition selects.
            "server-only": fileURLToPath(
              new URL("./node_modules/server-only/empty.js", import.meta.url),
            ),
          },
        },
        test: {
          name: "server",
          environment: "node",
          globals: true,
          setupFiles: ["./vitest.setup.server.ts"],
          include: ["src/**/*.test.ts"],
        },
      },
    ],
  },
});
