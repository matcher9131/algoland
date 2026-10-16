import clsx from "clsx";
import React, { ReactNode } from "react";

export default function ProblemSectionHeader({ children, className, ...attributes }: React.HTMLProps<HTMLHeadingElement> & { children: ReactNode }) {
    return (
        <h3 className={clsx("text-xl", "py-2", className)} {...attributes}>
            {children}
        </h3>
    )
}
