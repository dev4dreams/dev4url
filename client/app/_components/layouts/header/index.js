import Link from "next/link";
import CustomLink from "../../ui/links";

export default function Header() {
  return (
    <nav className="w-dvw px-7 pl-3 sm:px-10 md:px-20  lg:px-40 py-8 flex justify-between items-center fixed top-0">
      <Link
        href="/"
        className="w-[200px] sm:w-[250px] md:w-[300px] xl:w-[400px]"
      >
        <video autoPlay={true} playsInline muted width="100%">
          <source
            src={`https://utfs.io/f/8EU31ZztQwBP2g48Wo6GtJi9QmNVcRCloOTYAsZkS57vqn0y`}
            type="video/mp4"
          ></source>
        </video>
      </Link>

      <div className="flex sm:gap-2 md:gap-10">
        <CustomLink
          to="/customLink"
          label="Custom url"
          disabled="true"
          tooltipText="Coming Soon!"
        />
        <CustomLink
          to="/login"
          label="Login"
          disabled="true"
          tooltipText="Coming Soon!"
        />
      </div>
    </nav>
  );
}
