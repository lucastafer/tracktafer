import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  env: {
    NEXT_PUBLIC_GOOGLE_MAPS_API_KEY:
      process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY,
    NEXT_PUBLIC_NEST_API_URL: process.env.NEXT_PUBLIC_NEST_API_URL,
    NEXT_PUBLIC_NEXT_API_URL: process.env.NEXT_PUBLIC_NEXT_API_URL,
    NEST_API_URL: process.env.NEST_API_URL,
  },
};

export default nextConfig;
