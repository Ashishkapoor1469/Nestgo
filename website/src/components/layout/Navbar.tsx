"use client";

import { motion } from "framer-motion";
import Link from "next/link";
import Image from "next/image";
import { ExternalLink, Menu } from "lucide-react";
import { useState, useEffect } from "react";

export function Navbar() {
  const [isScrolled, setIsScrolled] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      const scrollPercentage = (window.scrollY / window.innerHeight) * 100;
      setIsScrolled(scrollPercentage > 5);
    };

    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  return (
    <motion.header
      initial={{ y: -100 }}
      animate={{ y: 0 }}
      transition={{ duration: 0.8, ease: [0.16, 1, 0.3, 1] }}
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
        isScrolled 
          ? "bg-white/90 backdrop-blur-md border-b border-border shadow-sm" 
          : "bg-transparent"
      }`}
    >
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          
          {/* Logo Section */}
          <div className="flex items-center gap-2">
            <Link href="/" className="flex items-center gap-2.5 hover:opacity-80 transition-opacity">
              <Image 
                src="/logo.webp" 
                alt="NestGo Logo" 
                width={60} 
                height={60}
                className="md:w-full md:h-25"
                priority
                unoptimized
              />
            </Link>
          </div>

          {/* Desktop Nav */}
          <nav className="hidden md:flex items-center gap-8 text-sm font-medium">
            <Link href="/docs" className={`transition-colors ${isScrolled ? "text-muted-foreground hover:text-foreground" : "text-white/80 hover:text-white"}`}>
              Docs
            </Link>
            <Link href="/guides" className={`transition-colors ${isScrolled ? "text-muted-foreground hover:text-foreground" : "text-white/80 hover:text-white"}`}>
              Guides
            </Link>
            <Link href="/ecosystem" className={`transition-colors ${isScrolled ? "text-muted-foreground hover:text-foreground" : "text-white/80 hover:text-white"}`}>
              Ecosystem
            </Link>
            <div className={`h-4 w-px ${isScrolled ? "bg-border" : "bg-white/20"}`}></div>
            <a 
              href="https://github.com/Ashishkapoor1469/Nestgo" 
              target="_blank"
              rel="noopener noreferrer"
              className={`flex items-center gap-2 transition-colors ${isScrolled ? "text-muted-foreground hover:text-foreground" : "text-white/80 hover:text-white"}`}
            >
              <ExternalLink className="w-4 h-4" />
              <span>GitHub</span>
            </a>
          </nav>

          {/* Mobile Menu Button */}
          <div className="md:hidden">
            <button className={`p-2 transition-colors ${isScrolled ? "text-muted-foreground hover:text-foreground" : "text-white/80 hover:text-white"}`}>
              <Menu className="w-5 h-5" />
            </button>
          </div>

        </div>
      </div>
    </motion.header>
  );
}
