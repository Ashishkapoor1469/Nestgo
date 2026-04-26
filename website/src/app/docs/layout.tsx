"use client";

import Link from "next/link";
import Image from "next/image";
import { usePathname } from "next/navigation";
import { useState, useEffect } from "react";
import { AnimatePresence, motion } from "framer-motion";
import {
  ChevronRight,
  Menu,
  X,
  BookOpen,
  Cpu,
  Shield,
  Wrench,
  Home,
  ExternalLink,
} from "lucide-react";

const navItems = [
  {
    category: "Getting Started",
    icon: BookOpen,
    items: [
      { label: "Installation", href: "/docs/installation" },
      { label: "Quick Start", href: "/docs/quick-start" },
      { label: "Building A Simple App", href: "/docs/building-a-simple-app" },
    ],
  },
  {
    category: "Core Concepts",
    icon: Cpu,
    items: [
      { label: "Modules", href: "/docs/modules" },
      { label: "Controllers", href: "/docs/controllers" },
      { label: "Services", href: "/docs/services" },
    ],
  },
  {
    category: "Advanced Features",
    icon: Shield,
    items: [
      { label: "Guards", href: "/docs/guards" },
      { label: "Interceptors", href: "/docs/interceptors" },
      { label: "Authentication", href: "/docs/authentication" },
      { label: "Authorization", href: "/docs/authorization" },
    ],
  },
  {
    category: "Tools",
    icon: Wrench,
    items: [{ label: "CLI", href: "/docs/cli" }],
  },
];

function SidebarContent({
  pathname,
  onClose,
}: {
  pathname: string;
  onClose?: () => void;
}) {
  return (
    <div className="h-full flex flex-col overflow-hidden">
      {/* Sidebar header */}
      <div className="flex items-center justify-between px-4 py-3 border-b border-border flex-shrink-0">
        <Link
          href="/"
          onClick={onClose}
          className="flex items-center gap-2 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
        >
          <Home className="w-4 h-4" />
          <span>Back to Home</span>
        </Link>
        {onClose && (
          <button
            onClick={onClose}
            className="p-1.5 rounded-lg text-muted-foreground hover:text-foreground hover:bg-secondary transition-colors"
            aria-label="Close sidebar"
          >
            <X className="w-4 h-4" />
          </button>
        )}
      </div>

      {/* Nav items */}
      <nav className="flex-1 overflow-y-auto py-4 px-3 space-y-5">
        {navItems.map((section) => {
          const Icon = section.icon;
          return (
            <div key={section.category}>
              <div className="flex items-center gap-2 mb-2 px-2">
                <Icon className="w-3.5 h-3.5 text-primary flex-shrink-0" />
                <h4 className="text-xs font-bold text-foreground uppercase tracking-wider">
                  {section.category}
                </h4>
              </div>
              <ul className="space-y-0.5 pl-4 border-l border-border ml-1.5">
                {section.items.map((item) => {
                  const isActive = pathname === item.href;
                  return (
                    <li key={item.href}>
                      <Link
                        href={item.href}
                        onClick={onClose}
                        className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm transition-all duration-150 ${
                          isActive
                            ? "bg-primary/10 text-primary font-semibold"
                            : "text-muted-foreground hover:text-foreground hover:bg-secondary"
                        }`}
                      >
                        {isActive && (
                          <ChevronRight className="w-3 h-3 flex-shrink-0" />
                        )}
                        <span>{item.label}</span>
                      </Link>
                    </li>
                  );
                })}
              </ul>
            </div>
          );
        })}
      </nav>

      {/* Footer */}
      <div className="px-4 py-3 border-t border-border flex-shrink-0">
        <a
          href="https://github.com/Ashishkapoor1469/Nestgo"
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center gap-2 text-xs text-muted-foreground hover:text-foreground transition-colors"
        >
          <ExternalLink className="w-3.5 h-3.5" />
          <span>View on GitHub</span>
        </a>
      </div>
    </div>
  );
}

export default function DocsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const [mobileSidebarOpen, setMobileSidebarOpen] = useState(false);
  const [desktopSidebarOpen, setDesktopSidebarOpen] = useState(true);

  useEffect(() => {
    setMobileSidebarOpen(false);
  }, [pathname]);

  useEffect(() => {
    document.body.style.overflow = mobileSidebarOpen ? "hidden" : "";
    return () => { document.body.style.overflow = ""; };
  }, [mobileSidebarOpen]);

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">

      {/* ── Docs Top Bar (replaces main Navbar) ── */}
      <header className="fixed top-0 left-0 right-0 z-50 h-14 bg-white/95 backdrop-blur-md border-b border-border shadow-sm flex items-center">
        <div className="w-full flex items-center px-4 sm:px-6 gap-4">

          {/* Logo */}
          <Link href="/" className="flex items-center hover:opacity-80 transition-opacity flex-shrink-0">
            <Image
              src="/logo.webp"
              alt="NestGo"
              width={44}
              height={44}
              priority
              unoptimized
            />
          </Link>

          {/* Divider */}
          <div className="h-5 w-px bg-border hidden sm:block" />

          {/* Sidebar toggle (desktop) */}
          <button
            onClick={() => setDesktopSidebarOpen((o) => !o)}
            className="hidden lg:flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors px-2 py-1.5 rounded-lg hover:bg-secondary"
            aria-label={desktopSidebarOpen ? "Collapse sidebar" : "Expand sidebar"}
          >
            <Menu className="w-4 h-4" />
            <span className="text-xs font-medium">
              {desktopSidebarOpen ? "Collapse" : "Expand"}
            </span>
          </button>

          {/* Mobile sidebar toggle */}
          <button
            onClick={() => setMobileSidebarOpen(true)}
            className="lg:hidden flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors px-2 py-1.5 rounded-lg hover:bg-secondary"
            aria-label="Open navigation"
          >
            <Menu className="w-4 h-4" />
          </button>

          {/* Breadcrumb */}
          <nav className="flex items-center gap-1 text-sm text-muted-foreground min-w-0">
            <Link href="/" className="hover:text-foreground transition-colors shrink-0">Home</Link>
            <ChevronRight className="w-3.5 h-3.5 shrink-0" />
            <Link href="/docs" className="hover:text-foreground transition-colors shrink-0">Docs</Link>
            {pathname !== "/docs" && (
              <>
                <ChevronRight className="w-3.5 h-3.5 shrink-0" />
                <span className="text-foreground font-medium truncate">
                  {navItems
                    .flatMap((s) => s.items)
                    .find((i) => i.href === pathname)?.label ?? ""}
                </span>
              </>
            )}
          </nav>

          {/* Spacer + GitHub link */}
          <div className="ml-auto flex-shrink-0">
            <a
              href="https://github.com/Ashishkapoor1469/Nestgo"
              target="_blank"
              rel="noopener noreferrer"
              className="hidden sm:flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors px-3 py-1.5 rounded-lg hover:bg-secondary"
            >
              <ExternalLink className="w-3.5 h-3.5" />
              <span>GitHub</span>
            </a>
          </div>
        </div>
      </header>

      {/* ── Body (below the 56px top bar) ── */}
      <div className="flex flex-1 pt-14">

        {/* Desktop Sidebar */}
        <AnimatePresence initial={false}>
          {desktopSidebarOpen && (
            <motion.aside
              key="desktop-sidebar"
              initial={{ width: 0, opacity: 0 }}
              animate={{ width: 260, opacity: 1 }}
              exit={{ width: 0, opacity: 0 }}
              transition={{ duration: 0.22, ease: "easeInOut" }}
              className="hidden lg:block flex-shrink-0 overflow-hidden"
            >
              <div className="w-[260px] fixed top-14 bottom-0 border-r border-border bg-background overflow-hidden">
                <SidebarContent pathname={pathname} />
              </div>
            </motion.aside>
          )}
        </AnimatePresence>

        {/* Main Content */}
        <main className="flex-1 min-w-0">
          <div className="max-w-3xl mx-auto px-5 sm:px-8 py-10">
            <article className="docs-article">{children}</article>
          </div>
        </main>
      </div>

      {/* Mobile Sidebar Drawer */}
      <AnimatePresence>
        {mobileSidebarOpen && (
          <>
            <motion.div
              key="mobile-backdrop"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.2 }}
              onClick={() => setMobileSidebarOpen(false)}
              className="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm lg:hidden"
            />
            <motion.aside
              key="mobile-drawer"
              initial={{ x: "-100%" }}
              animate={{ x: 0 }}
              exit={{ x: "-100%" }}
              transition={{ duration: 0.28, ease: [0.32, 0.72, 0, 1] }}
              className="fixed top-0 left-0 bottom-0 z-50 w-72 bg-background border-r border-border shadow-2xl lg:hidden"
            >
              <SidebarContent
                pathname={pathname}
                onClose={() => setMobileSidebarOpen(false)}
              />
            </motion.aside>
          </>
        )}
      </AnimatePresence>
    </div>
  );
}
