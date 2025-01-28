"use client";

import { useRouter } from "next/navigation";

export default function Maintenance() {
  const router = useRouter();
  return (
    <div className="w-full h-full flex flex-col items-center pt-[40%] sm:pt-[20%] lg:pt-[3.5em] md:gap-10">
      <div className="w-[80%] sm:w-[400px] md:w-[500px] lg:w-[60%]">
        <video autoPlay={true} playsInline muted width="100%">
          <source
            src={`https://as85m4vyio.ufs.sh/f/SawKekFykXenLsPhMNayeTEc2nDqrgKXAPS3yxsvQpjMCdla`}
            type="video/mp4"
          ></source>
        </video>
      </div>
      <div className="mt-10 sm:mt-5 flex flex-col gap-8 items-center">
        <p className="text-[#4255c3] font-ebGaramond font-bold text-[1.3em] md:text-[2em] lg:text-[2em] text-center">
          Sorry for the inconvenience. <br /> Server is during maintenance.
        </p>
        <p className="text-red-500 font-ebGaramond text-center text-[1.1em] md:text-[1.5em] lg:text-[2em]">
          We would be back soon.
        </p>
      </div>
    </div>
  );
}
