"use client";

import Image from "next/image";
import { useCallback } from "react";

import IconWarning from "@/app/_assets/icon/warning.svg";

export default function ErrorMsg({ error }) {
  const renderErrorMsg = useCallback(() => {
    const baseClasses =
      "h-full text-[#cf073d] text-[8px] sm:text-[12px] md:text-[16px] font-semibold pt-1 pr-2";

    switch (error) {
      case "url":
        return (
          <p className={baseClasses}>
            This URL has been flagged as potentially unsafe. <br />
            We can't process URL that may contain unsafe content.
          </p>
        );
      case "dev4url":
        return (
          <p
            className={`${baseClasses} flex flex-col items-center justify-center`}
          >
            <div className="eyes"></div>
            <br />
            Don't make me redirect to myself, please.
          </p>
        );
      default:
        return (
          <p className={baseClasses}>
            Sorry for the inconvenience. <br />
            Server is during maintenance.
            <br />
            Please try again later.
          </p>
        );
    }
  }, [error]);
  return (
    <div
      className={`${
        error ? "display-default opacity-100" : "display-none opacity-0"
      } ${
        error == "dev4url" ? "bg-accent" : "bg-white"
      } transition-opacity duration-300 flex items-center gap-3 py-1 md:py-4 my-6 border-2 border-red-800 p-2 rounded-lg`}
    >
      {error !== "dev4url" && (
        <Image
          src={IconWarning}
          alt="warning icon"
          className="w-[16px] md:w-[25px]"
        />
      )}{" "}
      {renderErrorMsg()}
    </div>
  );
}
