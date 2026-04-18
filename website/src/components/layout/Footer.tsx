import Link from "next/link";
import { ExternalLink } from "lucide-react";

export function Footer() {
  const currentYear = new Date().getFullYear();

  const links = [
    { label: "Docs", href: "/docs" },
    { label: "Guides", href: "/guides" },
    { label: "Ecosystem", href: "/ecosystem" },
    { label: "License (MIT)", href: "https://opensource.org/licenses/MIT" },
  ];

  const social = [
    { name: "GitHub", href: "https://github.com/Ashishkapoor1469/Nestgo" },
    { name: "Twitter", href: "https://twitter.com/nestjs" },
    { name: "LinkedIn", href: "https://linkedin.com" },
  ];

  return (
    <footer className="border-t border-border bg-white">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-12">
        {/* Main Footer Content */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-12 mb-12">
          {/* About */}
          <div>
            <h3 className="font-bold text-foreground mb-4">NestGo</h3>
            <p className="text-sm text-muted-foreground leading-relaxed">
              Enterprise backend architecture for Go. Built for massive performance, compile-time type safety, and infinite scalability.
            </p>
          </div>

          {/* Links */}
          <div>
            <h3 className="font-bold text-foreground mb-4">Resources</h3>
            <div className="space-y-2">
              {links.map((link) => (
                <a
                  key={link.label}
                  href={link.href}
                  target={link.href.startsWith("http") ? "_blank" : undefined}
                  rel={link.href.startsWith("http") ? "noopener noreferrer" : undefined}
                  className="block text-sm text-muted-foreground hover:text-foreground transition-colors"
                >
                  {link.label}
                </a>
              ))}
            </div>
          </div>

          {/* Social Links */}
          <div>
            <h3 className="font-bold text-foreground mb-4">Follow Us</h3>
            <div className="space-y-2">
              {social.map(({ name, href }) => (
                <a
                  key={name}
                  href={href}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-2 text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  <ExternalLink className="w-4 h-4" />
                  {name}
                </a>
              ))}
            </div>
          </div>
        </div>

        {/* Divider */}
        <div className="border-t border-border pt-8">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
            <p className="text-xs text-muted-foreground text-center md:text-left">
              &copy; {currentYear} NestGo Community. Released under the MIT License.
            </p>
            <p className="text-xs text-muted-foreground text-center md:text-right">
              Designed for developers. Built for scale.
            </p>
          </div>
        </div>
      </div>
    </footer>
  );
}
