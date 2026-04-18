export default function Loading() {
  return (
    <article className="max-w-none w-full animate-pulse">
      {/* Title Skeleton */}
      <div className="h-12 w-3/4 bg-secondary/50 rounded-lg mb-6"></div>
      
      {/* Paragraph Skeleton */}
      <div className="space-y-4 mb-10">
        <div className="h-4 w-full bg-secondary/40 rounded"></div>
        <div className="h-4 w-[95%] bg-secondary/40 rounded"></div>
        <div className="h-4 w-5/6 bg-secondary/40 rounded"></div>
      </div>

      {/* Subtitle Skeleton */}
      <div className="h-8 w-1/3 bg-secondary/50 rounded-lg mt-10 mb-4 border-b border-border pb-2"></div>
      
      {/* Code Block Skeleton */}
      <div className="h-48 w-full bg-slate-900 rounded-lg border border-slate-800 mb-8 mt-6"></div>

      {/* Paragraph Skeleton */}
      <div className="space-y-4 mb-8">
        <div className="h-4 w-full bg-secondary/40 rounded"></div>
        <div className="h-4 w-4/5 bg-secondary/40 rounded"></div>
      </div>
      
      {/* Subtitle Skeleton */}
      <div className="h-8 w-1/4 bg-secondary/50 rounded-lg mt-8 mb-4 border-b border-border pb-2"></div>
      
      {/* List or Table Skeleton */}
      <div className="space-y-3 my-6 ml-6">
        <div className="h-4 w-1/2 bg-secondary/40 rounded"></div>
        <div className="h-4 w-2/3 bg-secondary/40 rounded"></div>
        <div className="h-4 w-3/4 bg-secondary/40 rounded"></div>
      </div>
    </article>
  );
}
