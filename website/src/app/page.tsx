import { Navbar } from "@/components/layout/Navbar";
import { Footer } from "@/components/layout/Footer";
import { Hero } from "@/components/home/Hero";
import { CorePillars } from "@/components/home/CorePillars";
import { ExpandedFeatures } from "@/components/home/ExpandedFeatures";
import { ProductivityTools } from "@/components/home/ProductivityTools";
import { WhyNestGo } from "@/components/home/WhyNestGo";
import { DeploySection } from "@/components/home/DeploySection";
import { Newsletter } from "@/components/home/Newsletter";

export default function Home() {
  return (
    <main className="min-h-screen bg-background text-foreground">
      <Navbar />
      <Hero />
      <CorePillars />
      <ExpandedFeatures />
      <ProductivityTools />
      <WhyNestGo />
      <DeploySection />
      <Newsletter />
      <Footer />
    </main>
  );
}
