import fs from "fs";
import path from "path";
import { compileMDX } from "next-mdx-remote/rsc";
import rehypePrettyCode from "rehype-pretty-code";
import { notFound } from "next/navigation";
import { CodeBlock } from "@/components/CodeBlock";

// MDX component overrides — base styling comes from .docs-article CSS
const components = {
  h1: (props: any) => <h1 {...props} />,
  h2: (props: any) => <h2 {...props} />,
  h3: (props: any) => <h3 {...props} />,
  h4: (props: any) => <h4 {...props} />,
  p: (props: any) => <p {...props} />,
  ul: (props: any) => <ul {...props} />,
  ol: (props: any) => <ol {...props} />,
  li: (props: any) => <li {...props} />,
  a: (props: any) => <a {...props} />,
  strong: (props: any) => <strong {...props} />,
  em: (props: any) => <em {...props} />,
  hr: (props: any) => <hr {...props} />,
  blockquote: (props: any) => <blockquote {...props} />,
  code: (props: any) => (
    <code
      className="relative rounded px-2 py-1 font-mono text-sm font-semibold bg-slate-900 text-cyan-400 border border-slate-700"
      {...props}
    />
  ),
  pre: (props: any) => (
    <CodeBlock>
      <pre className="p-6 overflow-x-auto text-slate-100 leading-relaxed" {...props} />
    </CodeBlock>
  ),
  table: (props: any) => (
    <div className="table-wrapper my-6 rounded-xl border border-border overflow-hidden">
      <table className="w-full text-sm" {...props} />
    </div>
  ),
  thead: (props: any) => <thead {...props} />,
  tr: (props: any) => <tr {...props} />,
  th: (props: any) => <th {...props} />,
  td: (props: any) => <td {...props} />,
};

interface Props {
  params: Promise<{ slug: string[] }>;
}

export async function generateStaticParams() {
  const contentDir = path.join(process.cwd(), "src", "content", "docs");
  let files: string[] = [];
  try {
    files = fs.readdirSync(contentDir);
  } catch (error) {
    console.error("Could not read docs directory", error);
  }
  
  return files
    .filter((file) => file.endsWith(".mdx"))
    .map((file) => ({
      slug: [file.replace(/\.mdx$/, "")],
    }));
}

export async function generateMetadata({ params }: Props) {
  const { slug } = await params;
  const contentDir = path.join(process.cwd(), "src", "content", "docs");
  const slugPath = slug ? slug.join("/") : "installation";
  const filePath = path.join(contentDir, `${slugPath}.mdx`);

  try {
    const fileContent = fs.readFileSync(filePath, "utf8");
    const match = fileContent.match(/^---\r?\n([\s\S]*?)\r?\n---/);
    if (match) {
      const frontmatterContent = match[1];
      const lines = frontmatterContent.split('\n');
      let title = "Docs";
      let description = "NestGo Documentation";
      lines.forEach(line => {
        const [key, value] = line.split(':').map(s => s.trim());
        if (key === "title" && value) title = value.replace(/^["']|["']$/g, '');
        if (key === "description" && value) description = value.replace(/^["']|["']$/g, '');
      });
      return { title: `${title} | NestGo Docs`, description };
    }
  } catch (error) {
    // Ignore error, will fallback
  }
  
  return {
    title: "NestGo Documentation",
    description: "NestGo comprehensive docs."
  };
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
  
  const frontmatterRegex = /^---\r?\n([\s\S]*?)\r?\n---\r?\n([\s\S]*)$/;
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
