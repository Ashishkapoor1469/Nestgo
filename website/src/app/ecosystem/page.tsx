import Link from "next/link";
import { Navbar } from "@/components/layout/Navbar";
import { Footer } from "@/components/layout/Footer";
import { ArrowRight, Package, Code2, Users, BarChart3, Wrench, Star } from "lucide-react";

const ecosystemItems = [
  {
    category: "Official Packages",
    items: [
      {
        name: "NestGo CLI",
        description: "Command-line interface for project scaffolding, code generation, and development workflows.",
        icon: Code2,
        url: "https://github.com/Ashishkapoor1469/Nestgo/tree/main/cmd/nestgo"
      },
      {
        name: "NestGo Core",
        description: "The core framework with modular architecture, DI container, and HTTP server.",
        icon: Package,
        url: "https://github.com/Ashishkapoor1469/Nestgo"
      },
      {
        name: "NestGo Plugins",
        description: "Extensible plugin system for adding observability, security, and integration features.",
        icon: Wrench,
        url: "https://github.com/Ashishkapoor1469/Nestgo/tree/main/plugins"
      }
    ]
  },
  {
    category: "Recommended Libraries",
    items: [
      {
        name: "Chi Router",
        description: "Lightweight and high-performance HTTP router that NestGo uses internally.",
        icon: BarChart3,
        url: "https://github.com/go-chi/chi"
      },
      {
        name: "GORM",
        description: "Powerful ORM library for database operations and migrations in NestGo apps.",
        icon: Package,
        url: "https://gorm.io"
      },
      {
        name: "JWT-Go",
        description: "Industry-standard JWT authentication implementation for Go applications.",
        icon: Code2,
        url: "https://github.com/golang-jwt/jwt"
      }
    ]
  },
  {
    category: "Community Resources",
    items: [
      {
        name: "GitHub Repository",
        description: "Source code, issues, discussions, and contributions for the NestGo project.",
        icon: Package,
        url: "https://github.com/Ashishkapoor1469/Nestgo"
      },
      {
        name: "Documentation",
        description: "Comprehensive API reference, guides, and tutorials for NestGo development.",
        icon: Code2,
        url: "/docs"
      },
      {
        name: "Examples",
        description: "Real-world example projects demonstrating NestGo features and best practices.",
        icon: Star,
        url: "https://github.com/Ashishkapoor1469/Nestgo/tree/main/examples"
      }
    ]
  }
];

export default function EcosystemPage() {
  return (
    <div className="min-h-screen bg-background text-foreground">
      <Navbar />
      <main className="mt-20">
        {/* Hero Section */}
        <section className="py-16 px-4 sm:px-6 lg:px-8 bg-secondary/30">
          <div className="mx-auto max-w-4xl text-center">
            <h1 className="text-4xl md:text-5xl font-bold tracking-tight text-foreground mb-4">
              NestGo Ecosystem
            </h1>
            <p className="text-lg text-muted-foreground mb-8">
              Explore official packages, recommended libraries, and community resources
            </p>
          </div>
        </section>

        {/* Ecosystem Categories */}
        <section className="py-20 px-4 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-6xl">
            {ecosystemItems.map((category, idx) => (
              <div key={idx} className="mb-16 last:mb-0">
                <h2 className="text-2xl font-bold text-foreground mb-8 flex items-center gap-2">
                  <div className="w-1 h-8 bg-primary rounded"></div>
                  {category.category}
                </h2>
                
                <div className="grid md:grid-cols-3 gap-6">
                  {category.items.map((item, i) => {
                    const Icon = item.icon;
                    return (
                      <a
                        key={i}
                        href={item.url}
                        target={item.url.startsWith('http') ? "_blank" : undefined}
                        rel={item.url.startsWith('http') ? "noopener noreferrer" : undefined}
                        className="p-6 bg-white rounded-lg border border-border shadow-sm hover:shadow-lg hover:border-primary transition-all group"
                      >
                        <div className="flex items-start justify-between mb-4">
                          <div className="p-3 bg-secondary rounded-lg group-hover:bg-primary/10 transition-colors">
                            <Icon className="w-6 h-6 text-primary" />
                          </div>
                          <ArrowRight className="w-4 h-4 text-muted-foreground group-hover:text-primary transition-colors" />
                        </div>
                        
                        <h3 className="font-semibold text-foreground mb-2 group-hover:text-primary transition-colors">
                          {item.name}
                        </h3>
                        <p className="text-sm text-muted-foreground">
                          {item.description}
                        </p>
                      </a>
                    );
                  })}
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Integration Guide */}
        <section className="py-16 px-4 sm:px-6 lg:px-8 bg-secondary/30">
          <div className="mx-auto max-w-4xl">
            <h2 className="text-2xl font-bold text-foreground mb-8">Integration Standards</h2>
            
            <div className="grid md:grid-cols-2 gap-6 mb-8">
              <div className="p-6 bg-white rounded-lg border border-border">
                <h3 className="font-semibold text-foreground mb-3">Database</h3>
                <p className="text-muted-foreground mb-4">
                  NestGo works seamlessly with PostgreSQL, MySQL, MongoDB, and other databases through GORM, BSON, or your preferred driver.
                </p>
                <Link href="/docs" className="text-primary font-semibold hover:underline">
                  Learn more →
                </Link>
              </div>
              
              <div className="p-6 bg-white rounded-lg border border-border">
                <h3 className="font-semibold text-foreground mb-3">Observability</h3>
                <p className="text-muted-foreground mb-4">
                  Built-in support for Prometheus metrics, structured logging with slog, and OpenTelemetry tracing.
                </p>
                <Link href="/docs" className="text-primary font-semibold hover:underline">
                  Learn more →
                </Link>
              </div>
              
              <div className="p-6 bg-white rounded-lg border border-border">
                <h3 className="font-semibold text-foreground mb-3">Deployment</h3>
                <p className="text-muted-foreground mb-4">
                  Deploy anywhere - Docker, Kubernetes, serverless platforms, or traditional VPS providers.
                </p>
                <Link href="/docs" className="text-primary font-semibold hover:underline">
                  Learn more →
                </Link>
              </div>
              
              <div className="p-6 bg-white rounded-lg border border-border">
                <h3 className="font-semibold text-foreground mb-3">Message Queues</h3>
                <p className="text-muted-foreground mb-4">
                  Integrate with RabbitMQ, Kafka, Redis, or SQS for asynchronous processing and distributed systems.
                </p>
                <Link href="/docs" className="text-primary font-semibold hover:underline">
                  Learn more →
                </Link>
              </div>
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className="py-16 px-4 sm:px-6 lg:px-8 bg-primary text-primary-foreground">
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="text-3xl font-bold mb-4">Join the NestGo Community</h2>
            <p className="text-primary-foreground/90 mb-8">
              Contribute to NestGo, share your projects, and connect with other developers
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <a 
                href="https://github.com/Ashishkapoor1469/Nestgo"
                target="_blank"
                rel="noopener noreferrer"
                className="px-6 py-3 bg-primary-foreground text-primary rounded-lg font-semibold hover:bg-white transition-colors"
              >
                Star on GitHub
              </a>
              <a 
                href="https://github.com/Ashishkapoor1469/Nestgo/discussions"
                target="_blank"
                rel="noopener noreferrer"
                className="px-6 py-3 border border-primary-foreground rounded-lg font-semibold hover:bg-primary-foreground/10 transition-colors"
              >
                Join Discussions
              </a>
            </div>
          </div>
        </section>
      </main>
      <Footer />
    </div>
  );
}
