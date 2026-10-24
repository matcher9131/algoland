import "react";

declare module "react" {
  namespace JSX {
    interface IntrinsicElements {
      math: React.DetailedHTMLProps<
        React.HTMLAttributes<HTMLElement>,
        HTMLElement
      > & {
        display?: "block" | "inline";
      };
      mi: React.DetailedHTMLProps<
        React.HTMLAttributes<HTMLElement>,
        HTMLElement
      >;
    }
  }
}
