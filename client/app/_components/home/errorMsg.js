import IconWarning from "@/app/_assets/icon/warning.svg";
import Image from "next/image";

export default function ErrorMsg({ error }) {
  return (
    <div
      className={`${
        error ? "display-default opacity-100" : "display-none opacity-0"
      } transition-opacity duration-300 flex items-center gap-3 py-1 md:py-4 my-6 border-2 border-red-800 p-2 rounded-lg`}
    >
      <Image
        src={IconWarning}
        alt="warning icon"
        className="w-[16px] md:w-[25px]"
      />{" "}
      {error == "url" ? (
        <p className="h-full text-[#cf073d] text-[8px] sm:text-[12px] md:text-[16px] font-semibold pt-1 pr-2">
          This URL has been flagged as potentially unsafe. <br />
          We can't process URL that may contain unsafe content.
        </p>
      ) : (
        <p className="h-full text-[#cf073d] text-[8px] sm:text-[12px] md:text-[16px] font-semibold pt-1 pr-2">
          Sorry for the inconvenience. <br /> Server is during maintenance.
          <br />
          Please try again later.
        </p>
      )}
    </div>
  );
}
