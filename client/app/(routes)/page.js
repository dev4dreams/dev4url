"use client";
import Image from "next/image";

import iconPushPrimary from "@/app/_assets/icon/push-primary.svg";
import { useState } from "react";
import { actions } from "../_actions/request";
import { GET } from "../_utils/api/get";
import { POST } from "../_utils/api/post";
export default function Home() {
  const [url, setUrl] = useState("");
  const apiUrl = "http://localhost:8080/";
  const fetchUsers = async () => {
    try {
      const res = await actions.url.GET.all();

      console.log("data: ", res);
    } catch (err) {
      console.error("TEST failed", err);
    } finally {
      //  setLoading(false);
    }
  };

  // const testServerConnect = async (e) => {
  //   e.preventDefault();
  //   try {
  //     const res = await GET("http://localhost:8080/api/health");
  //     console.log("testServerConnect res: ", res);
  //   } catch (err) {
  //     console.error("testServerConnect Error: ", err);
  //   }
  // };
  const testPost = async (e) => {
    e.preventDefault();
    try {
      const res = await POST({
        url: apiUrl + "api/shorten",
        body: { original_url: "https://www.example.com" },
      });
      console.log("RES: ", res);
    } catch (err) {
      console.error("TestPost Error: ", err);
    }
  };
  return (
    <main className="flex h-full flex-col items-center justify-center">
      <div className="text-center flex flex-col gap-4">
        <h1 className="text-[3.2em] font-bold tracking-tight text-gray-900 font-ebGaramond">
          Shorten your URL
        </h1>
        <p className="text-[1.3em] text-primary">
          Enter your long URL below to create a shorter, shareable link
        </p>
      </div>

      <div className="mt-8 space-y-4 flex flex-col items-center justify-center gap-10">
        <input
          type="url"
          required
          className="w-[700px] px-6 py-3 rounded-[35px] text-[1.1em] border-2 border-accent focus:outline-none focus:ring-2 focus:ring-secondary shadow-lg"
          placeholder="Enter your long url and wait for the magic!"
          onChange={(e) => setUrl(e.target.value)}
        />

        <button
          type="button"
          // onClick={testServerConnect}
          onClick={testPost}
          className="w-[400px] flex justify-between items-center py-3 px-4 border-[3px] border-[#1D5D53] bg-white rounded-[35px] shadow-sm text-xl font-bold text-[#319B8B] transition-all duration-[250ms] hover:w-[180px] hover:text-bg hover:bg-primary group hover:tracking-[-0.13em]"
        >
          <Image
            src={iconPushPrimary}
            alt="push icon"
            width={24}
            height={24}
            className="transition-[filter] duration-[250ms] group-hover:[filter:brightness(0)_invert(1)]"
          />
          Squeeze it
          <Image
            src={iconPushPrimary}
            alt="push icon"
            width={24}
            height={24}
            className="rotate-180 transition-[filter] duration-[250ms] group-hover:[filter:brightness(0)_invert(1)]"
          />
        </button>
      </div>

      {/* Result section - Initially hidden, show when URL is shortened */}
      <div className="mt-8 hidden">
        <div className="bg-white p-4 rounded-lg border border-gray-200">
          <p className="text-sm font-medium text-gray-700">Shortened URL:</p>
          <div className="mt-2 flex items-center justify-between">
            <code className="text-sm text-indigo-600">
              https://dev4url/abcd123
            </code>
            <button
              type="button"
              className="ml-4 px-3 py-1 text-sm text-indigo-600 hover:text-indigo-700 focus:outline-none"
            >
              Copy
            </button>
          </div>
        </div>
      </div>
    </main>
  );
}
