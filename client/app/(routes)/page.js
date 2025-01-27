"use client";

import Image from "next/image";
import { useCallback, useState } from "react";
import * as Sentry from "@sentry/nextjs";

import { POST } from "../_utils/api/post";

import iconPushPrimary from "@/app/_assets/icon/push-primary.svg";

import UrlResult from "../_components/home/result";
import ErrorMsg from "../_components/home/errorMsg";

export default function Home() {
  const [url, setUrl] = useState("");
  const [shortenUrl, setShortenUrl] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;

  const handleCreateShortUrl = async () => {
    try {
      // reset when submit
      setIsLoading(true);
      setErrorMsg("");
      setShortenUrl("");

      const res = await POST({
        url: apiUrl + "shortUrl/post",
        body: { original_url: url },
      });
      setUrl("");
      console.log("res", res);
      if (!res.shortenUrl) {
        console.log(typeof res, JSON.parse(res));
        const error = await JSON.parse(res);

        setErrorMsg(error?.errors[0].includes("dev4url") ? "dev4url" : "url");
      }

      setIsLoading(false);
      setShortenUrl(res.shortenUrl);
    } catch (err) {
      console.error("handleCreateShortUrl Error: ", err);
      Sentry.captureException(err, "handleCreateShortUrl");
      setErrorMsg("server");
      setIsLoading(false);
    }
  };

  const renderSubmitBtn = useCallback(
    () => (
      <button
        type="button"
        onClick={handleCreateShortUrl}
        disabled={!url}
        className={`${
          isLoading
            ? "w-[100px] md:w-[200px]"
            : "w-[220px] sm:w-[250px] md:w-[300px] lg:w-[400px]"
        } flex ${
          isLoading ? "justify-center" : "justify-between"
        } items-center py-1 sm:py-2 md:py-3 px-4 border-[3px] border-[#1D5D53] bg-white rounded-[35px] shadow-sm text-sm sm:text-md md:text-xl font-bold text-[#319B8B] transition-all duration-[250ms] hover:w-[200px] hover:text-bg hover:bg-primary group hover:tracking-[-0.13em] btn`}
      >
        {isLoading ? (
          <div className="loader text-[14px] sm:text-[20px] md:text-[24px]" />
        ) : (
          <>
            <Image
              src={iconPushPrimary}
              alt="push icon"
              className="w-[18px] md:w-[24px] lg:w-[32px] transition-[filter] duration-[250ms] group-hover:[filter:brightness(0)_invert(1)]"
            />
            Squeeze it
            <Image
              src={iconPushPrimary}
              alt="push icon"
              className="w-[18px] md:w-[24px] lg:w-[32px] rotate-180 transition-[filter] duration-[250ms] group-hover:[filter:brightness(0)_invert(1)]"
            />
          </>
        )}
      </button>
    ),
    [isLoading, url, handleCreateShortUrl]
  );
  return (
    <main className="flex h-full flex-col items-center justify-center">
      <div className="text-center flex flex-col gap-4">
        <h1 className="text-[1.6em] sm:text-[2em] lg:text-[3.2em] font-bold tracking-tight text-gray-900 font-ebGaramond">
          Shorten your URL
        </h1>
        <p className="text-[0.9em] sm:text-[1em] md:text-[1.3em] text-primary">
          Enter your <span className="font-extrabold">long</span> URL below to
          create a <span className="font-bold">shorter</span>, shareable link
        </p>
      </div>

      <div className="mt-8 space-y-4 flex flex-col items-center justify-center gap-5 md:gap-10">
        <input
          type="url"
          required
          value={url}
          className="w-[330px] sm:w-[420px] md:w-[500px] lg:w-[700px] px-3 md:px-6 py-1 sm:py-2 md:py-3 rounded-[35px] text-[0.7em] md:text-[1.1em] border-2 border-accent focus:outline-none focus:ring-2 focus:ring-secondary shadow-lg"
          placeholder="Enter your long url and wait for the magic!"
          onChange={(e) => setUrl(e.target.value)}
        />

        {renderSubmitBtn()}
      </div>
      {errorMsg && <ErrorMsg error={errorMsg} />}
      <UrlResult shortenUrl={shortenUrl} />
    </main>
  );
}
