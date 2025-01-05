/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  darkMode: false,
  theme: {
    extend: {
      colors: {
        text: "#081513",
        bg: "#f2faf9",
        primary: "#42c3b0",
        secondary: "#95e1d6",
        accent: "#69d9c9",
      },
      fontFamily: {
        ebGaramond: ["EB Garamond", "serif"],
        zenDots: ["Zen Dots", "cursive"],
      },
    },
  },
  plugins: [],
};
