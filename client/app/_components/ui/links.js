"use client";

import Link from "next/link";

export default function CustomLink({
  label,
  to = "/",
  disabled = false,
  tooltipText,
}) {
  const handleClick = (e) => {
    if (disabled) {
      e.preventDefault();
    }
  };
  const tooltipClass = tooltipText ? "tooltip-before" : "";
  return (
    <Link
      href={to}
      onClick={handleClick}
      className={`${tooltipClass} text-[15px] sm:text-[24px] md:text-2xl xl:text-3xl  font-ebGaramond cursor-pointer px-2 sm:px-4  md:px-7 py-2 `}
    >
      {label}
    </Link>
  );
}
