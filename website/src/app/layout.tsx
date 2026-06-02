import type { Metadata } from "next";
import { Inter, JetBrains_Mono } from "next/font/google";
import "./globals.css";

const inter = Inter({
  variable: "--font-inter",
  subsets: ["latin"],
});

const jetbrainsMono = JetBrains_Mono({
  variable: "--font-jetbrains-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "NestGo — Enterprise Backend Architecture for Go",
  description:
    "A production-grade, modular, and high-performance backend framework for Go inspired by NestJS.",

  verification: {
    google: "yehPcw5ESWZp48ves1OHjvjHzUY_Vb_QlPVxKPcoWNg",
  },

  icons: {
    icon: "/logo.webp",
    shortcut: "/logo.webp",
    apple: "/logo.webp",
  },

  openGraph: {
    title: "NestGo — Enterprise Backend Architecture for Go",
    description:
      "A production-grade, modular, and high-performance backend framework for Go inspired by NestJS.",
    images: [
      {
        url: "/logo.webp",
        width: 512,
        height: 512,
        alt: "NestGo Logo",
      },
    ],
  },

  twitter: {
    card: "summary_large_image",
    title: "NestGo — Enterprise Backend Architecture for Go",
    description:
      "A production-grade, modular, and high-performance backend framework for Go inspired by NestJS.",
    images: ["/logo.webp"],
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${inter.variable} ${jetbrainsMono.variable} font-sans antialiased overflow-x-hidden`}
      >
        {children}
      </body>
    </html>
  );
}