import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import Home from "./page";

describe("Home", () => {
  it("renders the project name as the heading", () => {
    render(<Home />);

    expect(screen.getByRole("heading", { level: 1 })).toHaveTextContent(
      "algoland",
    );
  });

  it("links to the backend status page", () => {
    render(<Home />);

    expect(screen.getByRole("link")).toHaveAttribute("href", "/status");
  });

  it("lists every layer of the stack", () => {
    render(<Home />);

    for (const layer of ["Frontend", "Backend", "Database", "Auth"]) {
      expect(screen.getByText(layer)).toBeInTheDocument();
    }
  });
});
