import "server-only";

import { createClient } from "@/lib/supabase/server";

/**
 * Calls the Go backend over the Vercel service binding.
 *
 * API_URL is injected by Vercel from the binding declared in vercel.json. It
 * resolves only inside a running function — not at build time and not in
 * middleware — so this helper must not be used during static generation.
 */
export async function apiFetch(
  path: string,
  init: RequestInit = {},
): Promise<Response> {
  const baseUrl = process.env.API_URL;
  if (!baseUrl) {
    throw new Error(
      "API_URL is not set. It is injected by the Vercel service binding at runtime; " +
        "locally, run the stack with `vercel dev` or set it in web/.env.local.",
    );
  }

  const supabase = await createClient();
  const {
    data: { session },
  } = await supabase.auth.getSession();

  const headers = new Headers(init.headers);
  if (session?.access_token) {
    headers.set("Authorization", `Bearer ${session.access_token}`);
  }

  return fetch(new URL(path.replace(/^\//, ""), ensureTrailingSlash(baseUrl)), {
    ...init,
    headers,
    // The backend is dynamic per user; never let a response be reused.
    cache: "no-store",
  });
}

/**
 * Same as apiFetch but parses a JSON body and throws on a non-2xx status.
 */
export async function apiJson<T>(
  path: string,
  init: RequestInit = {},
): Promise<T> {
  const response = await apiFetch(path, init);
  if (!response.ok) {
    throw new Error(`API ${init.method ?? "GET"} ${path} failed: ${response.status}`);
  }
  return (await response.json()) as T;
}

function ensureTrailingSlash(url: string): string {
  return url.endsWith("/") ? url : `${url}/`;
}
