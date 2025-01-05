import Link from "next/link";
import CustomLink from "../../ui/links";

export default function Header() {
  return (
    <nav className="w-dvw px-40 py-8 flex justify-between items-center fixed top-0">
      <Link href="/">
        <video autoPlay={true} playsInline muted width="400px">
          <source
            src={`https://utfs.io/f/8EU31ZztQwBP2g48Wo6GtJi9QmNVcRCloOTYAsZkS57vqn0y`}
            type="video/mp4"
          ></source>
        </video>
      </Link>

      <div className="flex gap-10">
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
