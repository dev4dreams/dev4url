"use client";

import Image from "next/image";
import { useEffect } from "react";
import { usePathname, useRouter } from "next/navigation";
import * as Sentry from "@sentry/nextjs";

import { POST } from "@/app/_utils/api/post";

import IconLink from "@/app/_assets/icon/link-black.svg";

export default function Shorten() {
  const pathname = usePathname();
  const router = useRouter();
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;

  useEffect(() => {
    handleRedirect();
  }, [pathname]);

  const handleRedirect = async () => {
    try {
      const res = await POST({
        url: apiUrl + "shortUrl/get",
        body: { ShortenUrl: pathname.slice(1) },
      });
      console.log("RES: ", res);
      if (res.original_url) {
        window.location.replace(res.original_url);
      }
      setTimeout(() => {
        // router.push("404");
      }, 1000);
    } catch (err) {
      console.error("handleRedirect Error: ", err);
      Sentry.captureException(err, "handleRedirect");
      router.push("maintenance");
    }
  };
  return (
    <div className="w-full h-full pt-4 flex justify-center items-center">
      <div className="w-[400px] sm:w-[450px] md:w-[500px] lg:w-[750px] aspect-[5/3] flex flex-col gap-5 sm:gap-8 items-center justify-center bg-accent drop-shadow-lg rounded-xl ">
        <div className="flex items-center gap-1 sm:gap-2 font-zenDots text-[0.5em] sm:text-[1.1em] md:text-[1.3em] lg:text-[1.8em] p-2 md:p-3 sm:pr-4 md:pr-5 rounded-lg bg-white">
          <Image
            src={IconLink}
            alt="link icon"
            className="w-[16px] sm:w-[26px] md:w-[1.5em] text-[#000]"
          />{" "}
          <span className="">https://dev4url.cc/</span>{" "}
          <div className="loader " />
        </div>
        <p className="text-[1.1em] sm:text-[1.5em] lg:text-[2em] text-white font-ebGaramond font-semibold text-center">
          Hold on a sec. <br />
          Redirecting you to the page...
        </p>
      </div>
    </div>
  );
}
