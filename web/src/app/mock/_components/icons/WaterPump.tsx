import { clsx } from "clsx";
import { HTMLAttributes } from "react";
import { twMerge } from "tailwind-merge";

export default function WaterPump({
  className,
  ...otherProps
}: HTMLAttributes<SVGSVGElement>) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 64 64"
      fill="none"
      strokeWidth="2.5"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={twMerge(clsx("stroke-red-400", className))}
      {...otherProps}
    >
      <title>ポンプ</title>
      <rect x="4" y="56" width="56" height="4" rx="1.5" />
      <rect x="18" y="50" width="12" height="6" rx="1" />
      <rect x="44" y="44" width="8" height="12" rx="1" />
      <circle cx="24" cy="36" r="14" />
      <path d="M30 23.4 V14 M38 36 V14" />
      <rect x="27" y="9" width="14" height="5" rx="1" />
      <path d="M10.9 31 H7 M10.9 41 H7" />
      <rect x="3" y="27" width="4" height="18" rx="1" />
      <path d="M38 36 H42" />
      <rect x="42" y="28" width="16" height="16" rx="2" />
      <path d="M47 31 V41 M52 31 V41" />
      <g strokeWidth="2.2">
        <path d="M24 33 Q28 30 29 27" />
        <path d="M24 33 Q28 30 29 27" transform="rotate(90 24 36)" />
        <path d="M24 33 Q28 30 29 27" transform="rotate(180 24 36)" />
        <path d="M24 33 Q28 30 29 27" transform="rotate(270 24 36)" />
      </g>
      <circle cx="24" cy="36" r="2.5" />
    </svg>
  );
}
