import Link from "next/link";
import { Navbar } from "@/components/layout/Navbar";
import { Footer } from "@/components/layout/Footer";
import { ArrowRight, BookOpen, Database, Shield, Zap, Cloud, TestTube } from "lucide-react";

const guides = [
  {
    title: "Building Your First REST API",
    description: "Learn how to create a complete REST API with NestGo, including routing, controllers, and services.",
    icon: BookOpen,
    topics: ["Controllers", "Services", "Routing", "Request/Response"]
  },
  {
    title: "Database Integration",
    description: "Connect your NestGo application to PostgreSQL, MySQL, or MongoDB and manage your data effectively.",
    icon: Database,
    topics: ["Database Setup", "ORM Patterns", "Migrations", "Transactions"]
  },
  {
    title: "Authentication & Security",
    description: "Implement JWT authentication, RBAC, and security best practices in your NestGo applications.",
    icon: Shield,
    topics: ["JWT Tokens", "RBAC", "Guards", "Encryption"]
  },
  {
    title: "Real-time Features",
    description: "Build real-time applications using WebSockets and event emitters with NestGo.",
    icon: Zap,
    topics: ["WebSockets", "Event Bus", "Real-time Updates", "Broadcasting"]
  },
  {
    title: "Deployment & DevOps",
    description: "Deploy NestGo applications to production using Docker, Kubernetes, and cloud platforms.",
    icon: Cloud,
    topics: ["Docker", "Kubernetes", "CI/CD", "Environment Config"]
  },
  {
    title: "Testing & Quality",
    description: "Write unit tests, integration tests, and ensure code quality in your NestGo projects.",
    icon: TestTube,
    topics: ["Unit Testing", "Mocking", "Integration Tests", "Coverage"]
  }
];

export default function GuidesPage() {
  return (
    <div className="min-h-screen bg-background text-foreground">
      <Navbar />
      <main className="mt-20">
        {/* Hero Section */}
        <section className="py-16 px-4 sm:px-6 lg:px-8 bg-secondary/30">
          <div className="mx-auto max-w-4xl text-center">
            <h1 className="text-4xl md:text-5xl font-bold tracking-tight text-foreground mb-4">
              NestGo Guides & Tutorials
            </h1>
            <p className="text-lg text-muted-foreground mb-8">
              Step-by-step guides to help you master NestGo and build production-ready applications
            </p>
            <Link 
              href="/docs"
              className="inline-flex items-center gap-2 px-6 py-3 bg-primary text-primary-foreground rounded-lg font-semibold hover:bg-primary/90 transition-colors"
            >
              Read the Docs <ArrowRight className="w-4 h-4" />
            </Link>
          </div>
        </section>

        {/* Guides Grid */}
        <section className="py-20 px-4 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-6xl">
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
              {guides.map((guide, idx) => {
                const Icon = guide.icon;
                return (
                  <div 
                    key={idx}
                    className="p-6 bg-white rounded-lg border border-border shadow-sm hover:shadow-md transition-all"
                  >
                    <div className="flex items-center gap-3 mb-4">
                      <div className="p-3 bg-secondary rounded-lg">
                        <Icon className="w-6 h-6 text-primary" />
                      </div>
                      <h3 className="text-xl font-semibold text-foreground">
                        {guide.title}
                      </h3>
                    </div>
                    <p className="text-muted-foreground mb-4">
                      {guide.description}
                    </p>
                    <div className="flex flex-wrap gap-2">
                      {guide.topics.map((topic, i) => (
                        <span 
                          key={i}
                          className="px-3 py-1 bg-secondary text-sm text-foreground rounded-full"
                        >
                          {topic}
                        </span>
                      ))}
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className="py-16 px-4 sm:px-6 lg:px-8 bg-primary text-primary-foreground">
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="text-3xl font-bold mb-4">Ready to learn NestGo?</h2>
            <p className="text-primary-foreground/90 mb-8">
              Start with our comprehensive guides and build your next backend application
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Link 
                href="/docs/installation"
                className="px-6 py-3 bg-primary-foreground text-primary rounded-lg font-semibold hover:bg-white transition-colors"
              >
                Get Started
              </Link>
              <a 
                href="https://github.com/Ashishkapoor1469/Nestgo"
                target="_blank"
                rel="noopener noreferrer"
                className="px-6 py-3 border border-primary-foreground rounded-lg font-semibold hover:bg-primary-foreground/10 transition-colors"
              >
                View on GitHub
              </a>
            </div>
          </div>
        </section>
      </main>
      <Footer />
    </div>
  );
}
