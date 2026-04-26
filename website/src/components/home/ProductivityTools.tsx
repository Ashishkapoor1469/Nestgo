"use client";

import { motion } from "framer-motion";
import { Terminal, Activity, Waypoints, Network, CheckCircle2 } from "lucide-react";

const tools = [
  {
    id: "routes",
    icon: Waypoints,
    name: "AST Route Explorer",
    command: "nestgo routes",
    description: "Instantly index and visualize every API endpoint via fast AST parsing—without booting the server.",
    terminal: `╔══════════════════════════════════════════╗
║         🌐 API Route Explorer          ║  
╚══════════════════════════════════════════╝

  METHOD  PATH                SUMMARY                     
  ────────────────────────────────────────────────────────
  GET     /todos              List all todos              
  POST    /auth/login         Authenticate User              
  GET     /users/{id}         Fetch User Profile              `
  },
  {
    id: "lint-arch",
    icon: Network,
    name: "Clean Architecture Linter",
    command: "nestgo lint-arch",
    description: "Enforce strict dependency rules to guarantee massive teams don't accidentally couple domains.",
    terminal: `[LINT] Scanning module boundaries...
[ERROR] Cross-module boundary violation detected!
  File: internal/modules/users/service.go
  Imports: "sample-app/internal/modules/auth"

  Fix: Domains must communicate through DI interfaces.
  Exit Code: 1`
  },
  {
    id: "doctor",
    icon: Activity,
    name: "Zero-Config Diagnostics",
    command: "nestgo doctor",
    description: "Scan your entire project configuration, Go toolchain, and dependency graph health in milliseconds.",
    terminal: `⚕️  NestGo Diagnostics Report:

[PASS] Go toolchain v1.22.0 detected
[PASS] Project module found: github.com/my-api
[PASS] go.mod dependencies are synchronized
[PASS] Configuration valid (nestgo.yaml)
[PASS] Environment mappings correctly binded
[PASS] Port 8080 is available

✓ Workspace health is Excellent.`
  }
];

export function ProductivityTools() {
  return (
    <section className="py-24 bg-slate-950 text-white overflow-hidden relative">
      <div className="absolute inset-0 bg-[url('/grid.svg')] bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))]" />
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 relative z-10">
        <motion.div 
          className="text-center mb-16"
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
        >
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/20 text-primary mb-6">
            <CheckCircle2 className="w-4 h-4" />
            <span className="text-sm font-medium tracking-wide uppercase">New in v0.3.0</span>
          </div>
          <h2 className="text-3xl md:text-5xl font-bold mb-4 tracking-tight">Radical Developer Productivity</h2>
          <p className="text-lg text-slate-400 max-w-2xl mx-auto">
            NestGo ships with a powerful unified CLI suite that automates the hardest parts of scaling enterprise codebases.
          </p>
        </motion.div>

        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6 sm:gap-8">
          {tools.map((tool, idx) => {
            const Icon = tool.icon;
            return (
              <motion.div
                key={tool.id}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: idx * 0.15 }}
                className="flex flex-col h-full"
              >
                <div className="mb-6 flex items-center gap-4">
                  <div className="p-3 bg-white/5 rounded-xl">
                    <Icon className="w-6 h-6 text-primary" />
                  </div>
                  <div>
                    <h3 className="text-xl font-bold">{tool.name}</h3>
                    <code className="text-xs text-primary font-mono bg-primary/10 px-2 py-1 rounded mt-1 inline-block">
                      $ {tool.command}
                    </code>
                  </div>
                </div>
                <p className="text-slate-400 text-sm mb-6 flex-grow">{tool.description}</p>
                
                {/* Terminal Window */}
                <div className="rounded-xl border border-white/10 bg-slate-900 overflow-hidden shadow-2xl">
                  <div className="flex items-center gap-2 px-4 py-3 bg-slate-800/50 border-b border-white/5">
                    <div className="w-2.5 h-2.5 rounded-full bg-red-500/80" />
                    <div className="w-2.5 h-2.5 rounded-full bg-yellow-500/80" />
                    <div className="w-2.5 h-2.5 rounded-full bg-emerald-500/80" />
                  </div>
                  <div className="p-3 sm:p-4 overflow-x-auto min-h-[180px]">
                    <pre className="text-xs font-mono text-slate-300 whitespace-pre">
                      {tool.terminal}
                    </pre>
                  </div>
                </div>
              </motion.div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
