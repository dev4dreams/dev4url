export default function BtnCopy({ onClick }) {
  return (
    <button
      onClick={onClick}
      className={`group bg-white hover:border-[#0abab5] hover:border-2 w-7 sm:w-9 md:w-12 aspect-square rounded-[10px] flex items-center justify-center shadow-lg transition`}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        className="w-[16px] sm:w-[20px] md:w-[26px]"
      >
        <path
          fill="currentColor"
          className={
            "text-[#2188d1] group-hover:text-[#0abab5] transition-colors duration-300"
          }
          d="M18.355 6.54h-1.94V4.69a2.69 2.69 0 0 0-1.646-2.484A2.7 2.7 0 0 0 13.745 2h-8.05a2.68 2.68 0 0 0-2.67 2.69v10.09a2.68 2.68 0 0 0 2.67 2.69h1.94v1.85a2.68 2.68 0 0 0 2.67 2.68h8a2.68 2.68 0 0 0 2.67-2.68V9.23a2.69 2.69 0 0 0-2.62-2.69M7.635 9.23v6.74h-1.94a1.18 1.18 0 0 1-1.17-1.19V4.69a1.18 1.18 0 0 1 1.17-1.19h8.05a1.18 1.18 0 0 1 1.17 1.19v1.85h-4.61a2.69 2.69 0 0 0-2.67 2.69"
        />
      </svg>
    </button>
  );
}
