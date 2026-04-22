"use client";

import { motion } from "framer-motion";
import { 
  Package, 
  Layers, 
  Lock, 
  Zap, 
  Database, 
  GitBranch, 
  BarChart3, 
  Server 
} from "lucide-react";

const features = [
  {
    icon: Layers,
    title: "Modularity",
    description: "Organize applications into self-contained modules with automatic dependency resolution"
  },
  {
    icon: Zap,
    title: "Scalability",
    description: "Scale seamlessly with efficient, battle-tested components and microservice support"
  },
  {
    icon: Package,
    title: "Dependency Injection",
    description: "Sophisticated DI system for better testability and loose coupling"
  },
  {
    icon: Lock,
    title: "Type Safety",
    description: "Compile-time type checking prevents errors before they reach production"
  },
  {
    icon: Database,
    title: "Rich Ecosystem",
    description: "Database drivers, caching layers, queue systems, and more out of the box"
  },
  {
    icon: GitBranch,
    title: "Route Discovery",
    description: "Built-in AST-based route explorer to visualize your API surface instantly."
  },
  {
    icon: Server,
    title: "Auth Scaffolding",
    description: "Generate complete Auth modules with JWT, Bcrypt, and Protected Guards in one command."
  },
  {
    icon: BarChart3,
    title: "Project Metrics",
    description: "Analyze project health, complexity, and resource distribution with nestgo metrics."
  },
  {
    icon: Layers,
    title: "Arch Linter",
    description: "Enforce clean architecture patterns by rejecting illegal cross-module dependencies."
  },
];

export function ExpandedFeatures() {
  return (
    <section className="py-24 bg-secondary/20">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <motion.div 
          className="text-center mb-16"
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
        >
          <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Everything you need..</h2>
          <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
            Build robust, powerful, and scalable server-side applications. Stop reinventing the wheel.
          </p>
        </motion.div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {features.map((feature, idx) => {
            const Icon = feature.icon;
            return (
              <motion.div
                key={idx}
                className="p-6 rounded-xl bg-white border border-border hover:border-primary/50 hover:shadow-lg transition-all duration-300 group"
                initial={{ opacity: 0, scale: 0.95 }}
                whileInView={{ opacity: 1, scale: 1 }}
                viewport={{ once: true }}
                transition={{ delay: idx * 0.05 }}
              >
                <div className="p-3 bg-primary/10 w-fit rounded-lg mb-4 group-hover:bg-primary/20 transition-colors">
                  <Icon className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-lg font-bold text-foreground mb-2 uppercase tracking-wide">{feature.title}</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">{feature.description}</p>
              </motion.div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
