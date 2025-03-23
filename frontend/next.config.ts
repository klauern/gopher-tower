import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  ...(process.env.NODE_ENV === "production"
    ? {
        output: "export",
        distDir: "out",
      }
    : {}),
  images: {
    unoptimized: true,
  },
  trailingSlash: true,
  env: {
    // Make environment variables available to the client
    NEXT_PUBLIC_API_BASE_URL: process.env.NEXT_PUBLIC_API_BASE_URL || "",
  },
  ...(process.env.NODE_ENV !== "production"
    ? {
        // Add proxy configuration for development only
        async rewrites() {
          return [
            {
              source: "/api/:path*",
              destination: "http://localhost:8080/api/:path*", // Proxy to your Go backend
            },
          ];
        },
      }
    : {}),
};

export default nextConfig;
