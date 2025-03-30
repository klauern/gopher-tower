'use client';

import { ThemeProvider } from "@/components/theme-provider";
import type { AppProps } from 'next/app';
import { Geist, Geist_Mono } from "next/font/google";
import "../styles/globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export default function App({ Component, pageProps }: AppProps) {
  return (
    <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
      <style jsx global>{`
        :root {
          --font-geist-sans: ${geistSans.variable};
          --font-geist-mono: ${geistMono.variable};
        }
      `}</style>
      <div className={`${geistSans.variable} ${geistMono.variable} antialiased`}>
        <Component {...pageProps} />
      </div>
    </ThemeProvider>
  );
}
