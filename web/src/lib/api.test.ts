import { http, HttpResponse } from "msw";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { mswServer } from "@/test/msw";

const BASE_URL = "https://api.internal.test";

// The Supabase server client reads cookies through next/headers, which only
// works inside a request. Stubbing the module keeps these tests focused on the
// request building in apiFetch.
const getSession = vi.fn();
vi.mock("@/lib/supabase/server", () => ({
  createClient: async () => ({ auth: { getSession } }),
}));

async function importApi() {
  return await import("@/lib/api");
}

describe("apiFetch", () => {
  beforeEach(() => {
    vi.stubEnv("API_URL", BASE_URL);
    getSession.mockResolvedValue({
      data: { session: { access_token: "the.access.token" } },
    });
  });

  afterEach(() => {
    vi.unstubAllEnvs();
    vi.clearAllMocks();
  });

  it("resolves the path against the binding URL", async () => {
    const { apiFetch } = await importApi();
    let seen: URL | undefined;
    mswServer.use(
      http.get(`${BASE_URL}/v1/me`, ({ request }) => {
        seen = new URL(request.url);
        return HttpResponse.json({ id: "user-1" });
      }),
    );

    const response = await apiFetch("/v1/me");

    expect(response.status).toBe(200);
    expect(seen?.pathname).toBe("/v1/me");
  });

  it("forwards the Supabase access token as a bearer token", async () => {
    const { apiFetch } = await importApi();
    let authorization: string | null = null;
    mswServer.use(
      http.get(`${BASE_URL}/v1/me`, ({ request }) => {
        authorization = request.headers.get("Authorization");
        return HttpResponse.json({});
      }),
    );

    await apiFetch("/v1/me");

    expect(authorization).toBe("Bearer the.access.token");
  });

  it("omits the header when there is no session", async () => {
    const { apiFetch } = await importApi();
    getSession.mockResolvedValue({ data: { session: null } });
    let authorization: string | null = "unset";
    mswServer.use(
      http.get(`${BASE_URL}/healthz`, ({ request }) => {
        authorization = request.headers.get("Authorization");
        return HttpResponse.json({ status: "ok" });
      }),
    );

    await apiFetch("/healthz");

    expect(authorization).toBeNull();
  });

  // A binding URL with no trailing slash would otherwise drop its last path
  // segment when used as a URL base.
  it("keeps a path prefix on the binding URL", async () => {
    const { apiFetch } = await importApi();
    vi.stubEnv("API_URL", `${BASE_URL}/base`);
    let seen: URL | undefined;
    mswServer.use(
      http.get(`${BASE_URL}/base/v1/me`, ({ request }) => {
        seen = new URL(request.url);
        return HttpResponse.json({});
      }),
    );

    await apiFetch("/v1/me");

    expect(seen?.pathname).toBe("/base/v1/me");
  });

  // The backend response is per-user; a cached one would be served to the
  // wrong person.
  it("never caches the response", async () => {
    const { apiFetch } = await importApi();
    const spy = vi.spyOn(globalThis, "fetch");
    mswServer.use(http.get(`${BASE_URL}/healthz`, () => HttpResponse.json({})));

    await apiFetch("/healthz");

    expect(spy.mock.calls[0]?.[1]).toMatchObject({ cache: "no-store" });
    spy.mockRestore();
  });

  it("explains itself when the binding is missing", async () => {
    const { apiFetch } = await importApi();
    vi.stubEnv("API_URL", "");

    await expect(apiFetch("/healthz")).rejects.toThrow(/API_URL is not set/);
  });
});

describe("apiJson", () => {
  beforeEach(() => {
    vi.stubEnv("API_URL", BASE_URL);
    getSession.mockResolvedValue({ data: { session: null } });
  });

  afterEach(() => {
    vi.unstubAllEnvs();
    vi.clearAllMocks();
  });

  it("parses a successful body", async () => {
    const { apiJson } = await importApi();
    mswServer.use(
      http.get(`${BASE_URL}/healthz`, () =>
        HttpResponse.json({ status: "ok", env: "test" }),
      ),
    );

    await expect(apiJson("/healthz")).resolves.toEqual({
      status: "ok",
      env: "test",
    });
  });

  it("throws on a non-2xx status", async () => {
    const { apiJson } = await importApi();
    mswServer.use(
      http.get(`${BASE_URL}/readyz`, () =>
        HttpResponse.json({ error: "database unreachable" }, { status: 503 }),
      ),
    );

    await expect(apiJson("/readyz")).rejects.toThrow(/503/);
  });
});
