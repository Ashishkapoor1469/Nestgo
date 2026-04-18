"use client";

import { CopyButton } from "./CopyButton";

export function CodeBlock({
  children,
  ...props
}: {
  children: React.ReactNode;
} & React.HTMLAttributes<HTMLDivElement>) {
  // Extract text content from the code element
  let codeText = "";
  
  const extractText = (node: any): string => {
    if (typeof node === "string") {
      return node;
    }
    if (node && typeof node === "object" && "props" in node) {
      const { children } = node.props;
      if (Array.isArray(children)) {
        return children.map(extractText).join("");
      }
      return extractText(children);
    }
    return "";
  };

  if (typeof children === "object" && children !== null && "props" in children) {
    codeText = extractText(children);
  }

  return (
    <div className="relative my-6 rounded-xl border border-slate-700 bg-slate-900 overflow-hidden text-sm shadow-2xl group">
      <div className="flex items-center justify-between px-4 py-3 bg-gradient-to-r from-slate-800 to-slate-900 border-b border-slate-700">
        <div className="flex gap-2">
          <div className="w-3 h-3 rounded-full bg-red-500 shadow-lg"></div>
          <div className="w-3 h-3 rounded-full bg-yellow-500 shadow-lg"></div>
          <div className="w-3 h-3 rounded-full bg-emerald-500 shadow-lg"></div>
        </div>
        <div className="flex items-center gap-3">
          <span className="font-mono text-xs text-slate-400 uppercase tracking-widest">Code</span>
          {codeText && <CopyButton code={codeText} />}
        </div>
      </div>
      <div className="overflow-x-auto" {...props}>
        {children}
      </div>
    </div>
  );
}
