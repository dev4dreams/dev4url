export default function Footer() {
  return (
    <footer className="w-full py-2 flex justify-center items-center">
      <p className="font-ebGaramond text-[12px] md:text-[1em] xl:text-[1.1em]">
        © {new Date().getFullYear()} <span>Dev4url.</span> All rights reserved.
      </p>
    </footer>
  );
}
