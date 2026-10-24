import { existsSync } from "node:fs";
import { homedir } from "node:os";
import path from "node:path";

import { defineConfig, devices } from "@playwright/test";

/**
 * Chromium needs libnss3, libnspr4 and libasound2, which normally come from
 * `playwright install-deps` and root. `make e2e-deps` unpacks them into the
 * home directory instead, so point the loader at them when they are there.
 * On a machine with the system packages installed this is a no-op.
 */
const localDeps = path.join(
  homedir(),
  ".local/playwright-deps/usr/lib/x86_64-linux-gnu",
);
if (existsSync(localDeps)) {
  process.env.LD_LIBRARY_PATH = [localDeps, process.env.LD_LIBRARY_PATH]
    .filter(Boolean)
    .join(path.delimiter);
}

const PORT = Number(process.env.E2E_PORT ?? 3100);
const API_PORT = Number(process.env.E2E_API_PORT ?? 8100);
const baseURL = `http://127.0.0.1:${PORT}`;

/**
 * The Go backend is only started when the database settings are present.
 * Without them the API specs skip themselves, so a contributor can still run
 * the frontend E2E suite without a local Supabase stack.
 */
const hasBackendEnv = Boolean(
  process.env.DATABASE_URL && process.env.SUPABASE_URL,
);

/** Copies the named variables through to a server, skipping any that are unset. */
function passThrough(...names: string[]): Record<string, string> {
  const env: Record<string, string> = {};
  for (const name of names) {
    const value = process.env[name];
    if (value) {
      env[name] = value;
    }
  }
  return env;
}

const apiServer = {
  command: "go run ./cmd/server",
  cwd: "../api",
  url: `http://127.0.0.1:${API_PORT}/healthz`,
  reuseExistingServer: !process.env.CI,
  timeout: 120_000,
  env: {
    PORT: String(API_PORT),
    VERCEL_ENV: "development",
    ...passThrough("DATABASE_URL", "SUPABASE_URL"),
  },
};

const webAppServer = {
  command: `npm run dev -- --port ${PORT}`,
  url: baseURL,
  reuseExistingServer: !process.env.CI,
  timeout: 120_000,
  env: {
    // Stands in for the Vercel service binding, which only exists in a
    // deployment.
    API_URL: `http://127.0.0.1:${API_PORT}`,
    // The proxy builds a Supabase client on every request, so these must be
    // present or every page fails. Locally they come from web/.env.local,
    // which Next.js loads on its own; in CI they are passed through here.
    ...passThrough("NEXT_PUBLIC_SUPABASE_URL", "NEXT_PUBLIC_SUPABASE_ANON_KEY"),
  },
};

const webServer = hasBackendEnv ? [apiServer, webAppServer] : [webAppServer];

export default defineConfig({
  testDir: "./e2e",
  fullyParallel: true,
  // A test left focused with .only should fail CI rather than silently skip
  // the rest of the suite.
  forbidOnly: Boolean(process.env.CI),
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: process.env.CI
    ? [["github"], ["html", { open: "never" }]]
    : [["list"]],
  use: {
    baseURL,
    trace: "on-first-retry",
  },
  projects: [{ name: "chromium", use: { ...devices["Desktop Chrome"] } }],
  webServer,
});
