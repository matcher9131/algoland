import { afterAll, afterEach, beforeAll } from "vitest";

import { mswServer } from "@/test/msw";

beforeAll(() => {
  // An unhandled request is almost always a test that would otherwise reach
  // the real Supabase project or the real backend, so fail loudly.
  mswServer.listen({ onUnhandledRequest: "error" });
});

afterEach(() => {
  mswServer.resetHandlers();
});

afterAll(() => {
  mswServer.close();
});
