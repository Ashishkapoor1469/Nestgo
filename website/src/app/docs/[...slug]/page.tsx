import fs from "fs";
import path from "path";
import { compileMDX } from "next-mdx-remote/rsc";
import rehypePrettyCode from "rehype-pretty-code";
import { notFound } from "next/navigation";
import { CodeBlock } from "@/components/CodeBlock";

// Syntax highlighting configuration
const highlightTheme = {
  keyword: "text-purple-400",
  string: "text-green-400",
  function: "text-cyan-400",
  number: "text-blue-400",
  comment: "text-gray-500",
  operator: "text-slate-300",
};

// Define components allowed in MDX natively
const components = {
  h1: (props: any) => <h1 className="text-4xl font-extrabold tracking-tight lg:text-5xl mb-6 text-foreground" {...props} />,
  h2: (props: any) => <h2 className="mt-10 border-b border-border pb-2 text-3xl font-semibold tracking-tight transition-colors first:mt-0 mb-4 text-foreground" {...props} />,
  h3: (props: any) => <h3 className="mt-8 text-2xl font-semibold tracking-tight mb-4 text-foreground" {...props} />,
  p: (props: any) => <p className="leading-7 [&:not(:first-child)]:mt-6 text-muted-foreground" {...props} />,
  ul: (props: any) => <ul className="my-6 ml-6 list-disc [&>li]:mt-2 text-muted-foreground" {...props} />,
  li: (props: any) => <li className="text-muted-foreground" {...props} />,
  code: (props: any) => {
    // Check if this is inline code (not inside pre)
    return (
      <code 
        className="relative rounded px-2 py-1 font-mono text-sm font-semibold bg-slate-900 text-cyan-400 border border-slate-700" 
        {...props} 
      />
    );
  },
  pre: (props: any) => (
    <CodeBlock>
      <pre className="p-6 overflow-x-auto text-slate-100 leading-relaxed" {...props} />
    </CodeBlock>
  ),
  a: (props: any) => <a className="font-medium text-primary underline underline-offset-4 hover:text-primary/80" {...props} />,
  table: (props: any) => (
    <div className="my-6 rounded-lg border border-border overflow-hidden">
      <table className="w-full text-sm" {...props} />
    </div>
  ),
  thead: (props: any) => <thead className="bg-secondary/50 border-b border-border" {...props} />,
  tr: (props: any) => <tr className="border-b border-border last:border-b-0 hover:bg-secondary/30" {...props} />,
  th: (props: any) => <th className="px-6 py-3 text-left font-semibold text-foreground" {...props} />,
  td: (props: any) => <td className="px-6 py-3 text-muted-foreground" {...props} />,
  blockquote: (props: any) => (
    <blockquote className="mt-6 border-l-4 border-primary pl-6 italic text-muted-foreground" {...props} />
  ),
};

interface Props {
  params: Promise<{ slug: string[] }>;
}

export default async function DocPage({ params }: Props) {
  const { slug } = await params;
  const contentDir = path.join(process.cwd(), "src", "content", "docs");
  const slugPath = slug ? slug.join("/") : "installation";
  const filePath = path.join(contentDir, `${slugPath}.mdx`);

  let fileContent: string;
  try {
    fileContent = fs.readFileSync(filePath, "utf8");
  } catch (error) {
    notFound();
  }

  // Extract frontmatter (YAML between --- markers)
  let frontmatter: any = {};
  let contentWithoutFrontmatter = fileContent;
  
  const frontmatterRegex = /^---\n([\s\S]*?)\n---\n([\s\S]*)$/;
  const match = fileContent.match(frontmatterRegex);
  
  if (match) {
    const frontmatterContent = match[1];
    contentWithoutFrontmatter = match[2];
    
    // Parse YAML frontmatter
    const lines = frontmatterContent.split('\n');
    lines.forEach(line => {
      const [key, value] = line.split(':').map(s => s.trim());
      if (key && value) {
        frontmatter[key] = value.replace(/^["']|["']$/g, '');
      }
    });
  }

  // Pre-configure the syntax highlighter options
  const options = {
    theme: "poimandres",
  };

  const { content } = await compileMDX({
    source: contentWithoutFrontmatter,
    components,
    options: {
      mdxOptions: {
        rehypePlugins: [[rehypePrettyCode, options]],
      },
    },
  });

  return (
    <article className="max-w-none">
      {content}
    </article>
  );
}
