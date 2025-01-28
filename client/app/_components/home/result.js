"use client";
import Image from "next/image";

import IconLink from "@/app/_assets/icon/link.svg";
import BtnCopy from "../ui/btn-copy";

export default function UrlResult({ shortenUrl }) {
  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(shortenUrl);
    } catch (err) {
      console.error("handleCopy Error: ", err);
    }
  };

  return (
    <div
      className={`flex items-center gap-3 transition-all duration-500 ${
        shortenUrl ? " translate-y-0 opacity-100" : " translate-y-3 opacity-0"
      } `}
    >
      <div className="flex items-center gap-2 px-2 pr-1 sm:px-4 md:px-5 py-1.5 md:py-3 my-10 rounded-full shadow-lg bg-gradient-to-br from-teal-300 to-teal-400">
        <Image
          src={IconLink}
          alt="link icon"
          className="w-[20px] sm:w-[22px] md:w-[26px] text-white"
        />{" "}
        <span className=" text-white text-[14px] sm:text-[18px] md:text-lg lg:text-xl font-semibold pr-4">
          {shortenUrl}
        </span>
      </div>
      <BtnCopy onClick={handleCopy} hover={"hover:bg-red"} />
    </div>
  );
}
