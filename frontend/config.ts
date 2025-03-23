/**
 * Application configuration
 *
 * This file contains environment-specific configuration settings.
 */

// Get API base URL from environment variable or use default
const getBaseUrl = () => {
  // In development, use relative URLs to let Next.js handle proxying
  if (process.env.NODE_ENV === "development") {
    return "";
  }
  const url = process.env.NEXT_PUBLIC_API_BASE_URL || "";
  console.log("API Base URL:", url);
  return url;
};

// Get the full API base URL, used for SSE and other direct connections
const getFullBaseUrl = () => {
  return process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";
};

// API endpoint configuration
const API_CONFIG = {
  endpoints: {
    events: "/api/events",
    jobs: "/api/jobs",
  },
} as const;

type StaticEndpoint = keyof typeof API_CONFIG.endpoints;
type DynamicEndpoint = `jobs/${string}`;
type Endpoint = StaticEndpoint | DynamicEndpoint;

// Export a function to get a full API URL
export function getApiUrl(endpoint: Endpoint): string {
  // Always use full URL for SSE endpoints
  if (endpoint === "events") {
    const url = `${getFullBaseUrl()}${API_CONFIG.endpoints.events}`;
    console.log("Generated SSE URL:", url);
    return url;
  }

  const baseUrl = getBaseUrl();

  // Handle dynamic job endpoints
  if (endpoint.startsWith("jobs/")) {
    const url = `${baseUrl}/api/${endpoint}`;
    console.log("Generated dynamic API URL:", url);
    return url;
  }

  // Handle static endpoints
  const url = `${baseUrl}${API_CONFIG.endpoints[endpoint as StaticEndpoint]}`;
  console.log("Generated static API URL:", url);
  return url;
}
