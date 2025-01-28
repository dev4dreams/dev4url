"use client";

import { useRouter } from "next/navigation";

export default function NonExist() {
  const router = useRouter();
  return (
    <div className="w-full h-full flex flex-col items-center pt-[40%] sm:pt-[20%] lg:pt-[3.5em] md:gap-10">
      <div className="w-[80%] sm:w-[400px] md:w-[500px] lg:w-[580px]">
        <video autoPlay={true} playsInline muted width="100%">
          <source
            src={`https://as85m4vyio.ufs.sh/f/SawKekFykXenLsPhMNayeTEc2nDqrgKXAPS3yxsvQpjMCdla`}
            type="video/mp4"
          ></source>
        </video>
      </div>
      <div className="mt-10 sm:mt-5 flex flex-col gap-8 items-center">
        <p className="text-[#4255c3] font-ebGaramond font-bold text-[1.3em] md:text-[2em] lg:text-[3em]">
          Oops! This short url doesn't exist.
        </p>
        <p className="text-primary font-ebGaramond text-center text-[1.1em] md:text-[1.5em] lg:text-[2em]">
          How about create a short url now?
        </p>
        <button
          type="button"
          onClick={() => router.push("/")}
          className={`flex justify-center items-center lg:w-[300px] py-1 sm:py-2 md:py-2 lg:py-3 px-4 md:px-5 lg:px-2 border-[3px] border-[#1D5D53] bg-white rounded-[35px] shadow-sm text-[14px] sm:text-[20px] md:text-md lg:text-xl font-bold text-[#319B8B] transition-all duration-[250ms] hover:w-[200px] hover:text-bg hover:bg-primary group hover:tracking-[-0.13em] btn`}
        >
          Shorten Your URL !
        </button>
      </div>
    </div>
  );
}
