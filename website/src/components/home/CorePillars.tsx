"use client";

import { motion } from "framer-motion";
import { Zap, Shield, Lightbulb } from "lucide-react";

const pillars = [
  {
    icon: Zap,
    title: "Powerful",
    description: "Extensible modular architecture with compile-time dependency injection for excellent developer experience"
  },
  {
    icon: Shield,
    title: "Reliable",
    description: "Enterprise-grade framework with built-in guards, interceptors, and comprehensive lifecycle management"
  },
  {
    icon: Lightbulb,
    title: "Progressive",
    description: "Inspired by proven patterns. Everything you need to scale from MVP to enterprise applications"
  }
];

export function CorePillars() {
  return (
    <section className="py-16 bg-background border-y border-border">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="grid md:grid-cols-3 gap-8">
          {pillars.map((pillar, idx) => {
            const Icon = pillar.icon;
            return (
              <motion.div 
                key={idx}
                className="text-center group"
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: idx * 0.1 }}
              >
                <div className="flex justify-center mb-4">
                  <div className="p-4 rounded-xl bg-primary/10 group-hover:bg-primary/20 transition-colors">
                    <Icon className="w-8 h-8 text-primary" />
                  </div>
                </div>
                <h3 className="text-xl font-bold text-foreground mb-2 uppercase tracking-wide">{pillar.title}</h3>
                <p className="text-muted-foreground leading-relaxed">{pillar.description}</p>
              </motion.div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
