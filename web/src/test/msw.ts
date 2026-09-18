import { setupServer } from "msw/node";

/**
 * A shared MSW server with no default handlers.
 *
 * Intercepting at the network layer rather than stubbing `fetch` means the
 * code under test runs its real request-building logic — headers, URL
 * resolution, response parsing — and only the wire is faked.
 */
export const mswServer = setupServer();
