"use client";

import { motion } from "framer-motion";
import { Check, X } from "lucide-react";

const comparisons = [
  { feature: "Modular Architecture", nestgo: true, gin: false, fiber: false },
  { feature: "Advanced DI System", nestgo: true, gin: false, fiber: false },
  { feature: "Global CLI Tooling", nestgo: "Advanced", gin: "None", fiber: "None" },
  { feature: "Architecture Enforcement", nestgo: "Built-in", gin: "Manual", fiber: "Manual" },
  { feature: "Route Explorer (AST)", nestgo: true, gin: false, fiber: false },
  { feature: "OpenAPI Auto-Gen", nestgo: true, gin: false, fiber: false },
];

export function WhyNestGo() {
  return (
    <section className="py-24 bg-gradient-to-b from-background via-primary/5 to-background">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <motion.div className="text-center mb-16" initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }}>
          <h2 className="text-3xl font-bold tracking-tight text-foreground mb-4 sm:text-4xl">Why Architect with NestGo?</h2>
          <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
            Complete feature parity with the enterprise frameworks you know, while keeping Go's simplicity and performance.
          </p>
        </motion.div>
        
        <motion.div 
          className="overflow-hidden rounded-2xl border border-border bg-white shadow-xl"
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.2 }}
        >
          <div className="overflow-x-auto">
            <table className="w-full divide-y divide-border text-left text-sm">
              <thead>
                <tr className="bg-gradient-to-r from-primary/10 to-primary/5 border-b-2 border-primary/20">
                  <th scope="col" className="px-6 py-5 font-bold text-foreground text-base">Feature</th>
                  <th scope="col" className="px-6 py-5 font-bold text-primary text-base text-center">
                    <span className="inline-flex items-center gap-2">
                      <span className="w-2 h-2 rounded-full bg-primary"></span>
                      NestGo
                    </span>
                  </th>
                  <th scope="col" className="px-6 py-5 font-bold text-muted-foreground text-base text-center">Gin</th>
                  <th scope="col" className="px-6 py-5 font-bold text-muted-foreground text-base text-center">Fiber</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {comparisons.map((row, idx) => (
                  <motion.tr 
                    key={idx}
                    className="hover:bg-primary/5 transition-colors"
                    initial={{ opacity: 0, x: -20 }}
                    whileInView={{ opacity: 1, x: 0 }}
                    viewport={{ once: true }}
                    transition={{ delay: idx * 0.05 }}
                  >
                    <td className="px-6 py-4 font-semibold text-foreground">{row.feature}</td>
                    <td className="px-6 py-4 text-center">
                      {row.nestgo === true ? (
                        <div className="inline-flex items-center justify-center w-6 h-6 rounded-full bg-green-100">
                          <Check className="text-green-600 w-5 h-5" />
                        </div>
                      ) : (
                        <span className="text-primary font-bold">{row.nestgo}</span>
                      )}
                    </td>
                    <td className="px-6 py-4 text-center">
                      {row.gin === true ? (
                        <div className="inline-flex items-center justify-center w-6 h-6 rounded-full bg-green-100">
                          <Check className="text-green-600 w-5 h-5" />
                        </div>
                      ) : row.gin === false ? (
                        <div className="inline-flex items-center justify-center w-6 h-6 rounded-full bg-red-100">
                          <X className="text-red-500 w-5 h-5" />
                        </div>
                      ) : (
                        <span className="text-muted-foreground font-medium">{row.gin}</span>
                      )}
                    </td>
                    <td className="px-6 py-4 text-center">
                      {row.fiber === true ? (
                        <div className="inline-flex items-center justify-center w-6 h-6 rounded-full bg-green-100">
                          <Check className="text-green-600 w-5 h-5" />
                        </div>
                      ) : row.fiber === false ? (
                        <div className="inline-flex items-center justify-center w-6 h-6 rounded-full bg-red-100">
                          <X className="text-red-500 w-5 h-5" />
                        </div>
                      ) : (
                        <span className="text-muted-foreground font-medium">{row.fiber}</span>
                      )}
                    </td>
                  </motion.tr>
                ))}
              </tbody>
            </table>
          </div>
        </motion.div>

        {/* Stats/CTA Section */}
        <motion.div 
          className="mt-16 grid md:grid-cols-3 gap-8"
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          viewport={{ once: true }}
          transition={{ delay: 0.4 }}
        >
          {[
            { label: "Enterprise Features", value: "100%" },
            { label: "Zero Overhead", value: "Native Go" },
            { label: "Developer Experience", value: "Premium" },
          ].map((stat, idx) => (
            <motion.div 
              key={idx}
              className="p-6 rounded-xl bg-white border border-border text-center hover:border-primary/50 transition-colors"
              initial={{ opacity: 0, scale: 0.9 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              transition={{ delay: 0.5 + idx * 0.1 }}
            >
              <p className="text-muted-foreground text-sm font-medium">{stat.label}</p>
              <p className="text-2xl font-bold text-primary mt-2">{stat.value}</p>
            </motion.div>
          ))}
        </motion.div>
      </div>
    </section>
  );
}
