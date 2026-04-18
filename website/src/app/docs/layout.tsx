import Link from "next/link";
import { Navbar } from "@/components/layout/Navbar";

export default function DocsLayout({ children }: { children: React.ReactNode }) {
  const navItems = [
    {
      category: "Getting Started",
      items: [
        { label: "Installation", href: "/docs/installation" },
        { label: "Quick Start", href: "/docs/quick-start" },
      ]
    },
    {
      category: "Core Concepts",
      items: [
        { label: "Modules", href: "/docs/modules" },
        { label: "Controllers", href: "/docs/controllers" },
        { label: "Services", href: "/docs/services" },
      ]
    },
    {
      category: "Advanced Features",
      items: [
        { label: "Guards", href: "/docs/guards" },
        { label: "Interceptors", href: "/docs/interceptors" },
        { label: "Authentication", href: "/docs/authentication" },
        { label: "Authorization", href: "/docs/authorization" },
      ]
    },
    {
      category: "Tools",
      items: [
        { label: "CLI", href: "/docs/cli" },
      ]
    }
  ];

  return (
    <div className="min-h-screen bg-background text-foreground">
      <Navbar />
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 mt-20 flex gap-8">
        {/* Sidebar */}
        <aside className="fixed w-64 h-[calc(100vh-5rem)] overflow-y-auto border-r border-border py-8 hidden lg:block">
          <nav className="space-y-8 text-sm px-4">
            {navItems.map((section) => (
              <div key={section.category}>
                <h4 className="font-semibold text-foreground mb-3">{section.category}</h4>
                <ul className="space-y-2">
                  {section.items.map((item) => (
                    <li key={item.href}>
                      <Link 
                        href={item.href} 
                        className="text-muted-foreground hover:text-primary transition-colors"
                      >
                        {item.label}
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </nav>
        </aside>

        {/* Content area */}
        <main className="flex-1 lg:pl-72 py-8 min-w-0">
          {children}
        </main>
      </div>
    </div>
  );
}
