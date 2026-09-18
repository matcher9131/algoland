import { expect, test } from "@playwright/test";

test.describe("home page", () => {
  test("renders the project name", async ({ page }) => {
    await page.goto("/");

    await expect(page.getByRole("heading", { level: 1 })).toHaveText("algoland");
  });

  test("navigates to the status page", async ({ page }) => {
    await page.goto("/");

    await page.getByRole("link").click();

    await expect(page).toHaveURL(/\/status$/);
    await expect(page.getByRole("heading", { level: 1 })).toHaveText("Backend status");
  });
});
