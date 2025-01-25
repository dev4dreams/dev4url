import { Geist, Geist_Mono, EB_Garamond, Zen_Dots } from "next/font/google";
import "./globals.css";
import Header from "../_components/layouts/header";
import iconLogo from "../_assets/icon/dev4url_icon.png";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

const zenDots = Zen_Dots({
  weight: "400",
  subsets: ["latin"],
  display: "swap",
  variable: "--font-zen-dots",
});

const ebGaramond = EB_Garamond({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-eb-garamond",
});

export const metadata = {
  title: "Dev4url",
  description: "Develop for your shorten url",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <head>
        <link rel="icon" href={iconLogo.src} />
      </head>
      <body
        className={`${geistSans.variable} ${geistMono.variable} ${ebGaramond.variable} ${zenDots.variable} antialiased h-dvh bg-bg flex-col items-start pt-[7.5%]`}
      >
        <Header />
        <div className="px-[10%] w-dvw h-[97%] flex-col justify-center items-center ">
          {children}
        </div>
      </body>
    </html>
  );
}
