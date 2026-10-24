import { expect, test } from "@playwright/test";

/**
 * Exercises the path Vitest cannot reach: an async Server Component that calls
 * the Go backend over the service binding during a real request.
 */
test.describe("status page", () => {
  test.skip(
    !process.env.DATABASE_URL || !process.env.SUPABASE_URL,
    "needs the Go backend, which needs a local Supabase stack",
  );

  test("reports the backend as reachable", async ({ page }) => {
    await page.goto("/status");

    await expect(page.getByText("ok")).toBeVisible();
    await expect(page.getByText("development")).toBeVisible();
  });

  test("renders on every request rather than from a prerender", async ({
    page,
  }) => {
    const first = await page.goto("/status");

    // A statically prerendered page would be served from the CDN cache; this
    // one must be produced by a function per request.
    expect(first?.headers()["cache-control"]).not.toContain("immutable");
    await expect(page.getByText("ok")).toBeVisible();
  });
});
