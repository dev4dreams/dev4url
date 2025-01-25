"use client";

import { useEffect } from "react";
import { usePathname, useRouter } from "next/navigation";

import { POST } from "@/app/_utils/api/post";

export default function Shorten() {
  const pathname = usePathname();
  const router = useRouter();
  const apiUrl = "http://localhost:8080/";

  useEffect(() => {
    handleRedirect();
  }, [pathname]);

  const handleRedirect = async () => {
    try {
      const res = await POST({
        url: apiUrl + "redirect",
        body: { ShortenUrl: pathname.slice(1) },
      });

      if (res.original_url) {
        window.location.replace(res.original_url);
      }

      router.push("nonexistent");
    } catch (err) {
      console.log("handleRedirect Error: ", err);
    }
  };
  return (
    <div className="w-full h-full pt-4 flex justify-center items-center">
      <p>Pathname: {pathname}</p>
      <div className=""></div>
    </div>
  );
}
