import { ReactNode } from "react";
import { clsx } from "clsx";

export default function SpeechBubble({
  children,
}: {
  readonly children: ReactNode;
}) {
  return (
    <div
      className={clsx(
        "relative",
        "ml-4",
        "px-4",
        "py-2",
        "bg-base-300",
        "rounded-lg",
        "before:content-['']",
        "before:absolute",
        "before:top-[35%]",
        "before:left-0",
        "before:border-solid",
        "before:border-[0_16px_16px_0]",
        "before:border-transparent",
        "before:border-r-base-300",
        "before:-translate-x-full",
        "before:translate-y-[calc(-50%-0.4px)]",
        "before:transform-[skew(0,10deg)]",
        "before:origin-right",
      )}
    >
      {children}
    </div>
  );
}
