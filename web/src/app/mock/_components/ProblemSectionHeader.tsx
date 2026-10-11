import React, { ReactNode } from "react";

export default function ProblemSectionHeader({ children, ...attributes }: React.HTMLProps<HTMLHeadingElement> & { children: ReactNode }) {

    return (
        <h3 className="text-xl py-2" {...attributes}>
            {children}
        </h3>
    )
}
