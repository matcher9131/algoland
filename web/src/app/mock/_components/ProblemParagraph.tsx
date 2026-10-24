import React, { ReactNode } from "react";

export default function ProblemParagraph({
  children,
  ...attributes
}: React.HTMLProps<HTMLHeadingElement> & { children: ReactNode }) {
  return (
    <p className="py-2 leading-normal" {...attributes}>
      {children}
    </p>
  );
}
