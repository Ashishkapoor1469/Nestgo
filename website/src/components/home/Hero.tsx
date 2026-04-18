"use client";

import { motion } from "framer-motion";
import Link from "next/link";
import { Copy, Check, Terminal, ArrowRight, Code2 } from "lucide-react";
import { useState } from "react";

export function Hero() {
  const [copied, setCopied] = useState(false);
  const command = "go install github.com/Ashishkapoor1469/Nestgo@latest";

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(command);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (error) {
      console.error("Failed to copy:", error);
    }
  };

  const container = {
    hidden: { opacity: 0 },
    show: {
      opacity: 1,
      transition: {
        staggerChildren: 0.08,
        delayChildren: 0.4,
      },
    },
  };

  const item = {
    hidden: { opacity: 0, y: 10 },
    show: { opacity: 1, y: 0 },
  };

  return (
    <section 
      className="relative min-h-screen w-full flex items-center justify-center overflow-hidden pt-16 pb-32 md:pt-20 md:pb-40 bg-cover bg-center bg-no-repeat"
      style={{
        backgroundImage: "url('/bg.png')",
        backgroundAttachment: 'fixed',
      }}
    >
      {/* Dark overlay for text readability */}
      <div className="absolute inset-0 bg-gradient-to-b from-black/40 via-black/30 to-black/50" />
      
      <div className="mx-auto max-w-5xl px-4 sm:px-6 lg:px-8 text-center flex flex-col items-center relative z-10">
        
        {/* Main Headline - Clean and Simple like NestJS */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.1 }}
          className="mb-8"
        >
          <h1 className="text-5xl md:text-7xl lg:text-8xl font-bold tracking-tight text-white mb-4">
            Hello, NestGo!
          </h1>
          <h2 className="text-xl md:text-2xl lg:text-3xl font-medium text-white/90">
            A progressive Go framework for building efficient, reliable and scalable server-side applications.
          </h2>
        </motion.div>

        {/* Quick Links */}
        <motion.div 
          className="mt-8 flex flex-col sm:flex-row items-center justify-center gap-3"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.2 }}
        >
          <Link
            href="/docs"
            className="text-white/90 font-semibold hover:text-white transition-colors flex items-center gap-1"
          >
            Documentation <ArrowRight className="w-4 h-4" />
          </Link>
          <span className="text-white/50 hidden sm:block">•</span>
          <a
            href="https://github.com/Ashishkapoor1469/Nestgo"
            target="_blank"
            rel="noopener noreferrer"
            className="text-white/90 font-semibold hover:text-white transition-colors flex items-center gap-1"
          >
            Source code <ArrowRight className="w-4 h-4" />
          </a>
        </motion.div>

        {/* CTAs - Larger and more prominent */}
        <motion.div 
          className="mt-12 flex flex-col sm:flex-row items-center gap-4 w-full max-w-2xl justify-center"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.3 }}
        >
          <Link 
            href="/docs/installation" 
            className="px-6 md:px-8 py-3 md:py-4 inline-flex items-center justify-center rounded-lg bg-primary text-primary-foreground font-semibold text-base md:text-lg hover:bg-primary/90 active:scale-95 transition-all duration-200 shadow-lg hover:shadow-xl hover:-translate-y-1 gap-2"
          >
            <Code2 className="w-5 h-5" />
            Start Building
          </Link>
          <button 
            onClick={handleCopy}
            className="px-4 md:px-6 py-3 md:py-4 flex items-center gap-3 rounded-lg border-2 border-white/30 bg-white/10 backdrop-blur-sm text-white font-mono text-sm shadow-lg hover:border-white/50 hover:bg-white/20 active:scale-95 transition-all duration-200 group"
          >
            <span className="text-white/80 font-bold hidden sm:inline">$</span>
            <span className="truncate group-hover:text-white transition-colors text-white/90">
              {copied ? "Copied!" : "go install github.com/nestgo/..."}
            </span>
            {copied ? (
              <Check className="w-5 h-5 text-green-400 flex-shrink-0" />
            ) : (
              <Copy className="w-5 h-5 text-white/60 group-hover:text-white transition-colors flex-shrink-0" />
            )}
          </button>
        </motion.div>

        {/* Terminal Example - Enhanced */}
        <motion.div 
          initial={{ opacity: 0, y: 50, rotateX: -15 }}
          animate={{ opacity: 1, y: 0, rotateX: 0 }}
          transition={{ duration: 0.8, delay: 0.5, type: "spring", stiffness: 80 }}
          style={{ perspective: "1200px" }}
          className="mt-20 md:mt-24 w-full max-w-4xl"
        >
          {/* 3D Shadow/Depth Effect */}
          <div className="absolute inset-0 -z-10 rounded-2xl bg-gradient-to-b from-white/20 via-white/10 to-white/5 blur-3xl transform translate-y-12" />
          
          <div className="rounded-2xl border border-white/20 bg-slate-900/80 backdrop-blur-md shadow-2xl overflow-hidden transform transition-all duration-300 hover:shadow-2xl hover:border-white/30">
            {/* Terminal Header */}
            <div className="flex items-center justify-between px-6 py-4 bg-gradient-to-r from-slate-800/60 via-slate-800/60 to-slate-900/60 backdrop-blur-sm border-b border-white/10">
              <div className="flex items-center gap-3">
                <div className="w-3 h-3 rounded-full bg-red-500/80"></div>
                <div className="w-3 h-3 rounded-full bg-yellow-500/80"></div>
                <div className="w-3 h-3 rounded-full bg-emerald-500/80"></div>
              </div>
              <div className="text-xs font-mono text-white/50 uppercase tracking-widest flex items-center gap-2">
                <Terminal className="w-4 h-4" /> $ nestgo generate resource
              </div>
              <div className="w-8" />
            </div>
            
            {/* Terminal Content */}
            <div className="bg-slate-900/70 backdrop-blur-sm p-8 font-mono text-sm text-left overflow-x-auto min-h-72">
              <motion.div variants={container} initial="hidden" animate="show">
                <motion.p variants={item}>
                  <span className="text-slate-400">~/projects/app</span> <span className="text-purple-400 font-semibold">nestgo</span> <span className="text-green-400">generate</span> <span className="text-cyan-400">resource</span> <span className="text-yellow-300">users</span>
                </motion.p>
                <motion.p variants={item} className="text-emerald-400 mt-4 font-semibold">
                  ✓ Scaffolding resource: users
                </motion.p>
                <motion.p variants={item} className="text-slate-400 mt-2">  ├─ internal/modules/users/service.go</motion.p>
                <motion.p variants={item} className="text-slate-400">  ├─ internal/modules/users/dto.go</motion.p>
                <motion.p variants={item} className="text-slate-400">  ├─ internal/modules/users/entity.go</motion.p>
                <motion.p variants={item} className="text-slate-400">  ├─ internal/modules/users/module.go</motion.p>
                <motion.p variants={item} className="text-slate-400">  └─ internal/modules/users/controller.go</motion.p>
                <motion.p variants={item} className="text-emerald-400 mt-4 font-semibold">
                  ✓ Resources generated successfully in 145ms
                </motion.p>
                <motion.p variants={item} className="text-slate-500 mt-3">
                  Ready for hot reload: <span className="text-cyan-400">nestgo dev --watch</span>
                </motion.p>
              </motion.div>
            </div>
          </div>
        </motion.div>

      </div>
    </section>
  );
}
