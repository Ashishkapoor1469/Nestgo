import { MetadataRoute } from "next";

export const dynamic = "force-static";

export default function sitemap(): MetadataRoute.Sitemap {
  return [
    {
      url: "https://nestgo.vercel.app",
      lastModified: new Date(),
      priority: 1,
    },
  ];
}