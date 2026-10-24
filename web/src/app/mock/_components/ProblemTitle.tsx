import React, { ReactNode } from "react";

export default function ProblemTitle({
  children,
  ...attributes
}: React.HTMLProps<HTMLHeadingElement> & { children: ReactNode }) {
  return (
    <h2 className="text-2xl mb-4 border-b" {...attributes}>
      {children}
    </h2>
  );
}
