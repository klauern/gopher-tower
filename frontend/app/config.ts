/**
 * Application configuration
 *
 * This file contains environment-specific configuration settings.
 */

// Get API base URL from environment variable or use default
const getBaseUrl = () => {
  if (typeof window === "undefined") {
    return process.env.NEXT_PUBLIC_API_BASE_URL || "";
  }
  // For static sites, we either use the environment variable or derive from window.location
  return process.env.NEXT_PUBLIC_API_BASE_URL || window.location.origin;
};

// API configuration
export const API_CONFIG = {
  // Endpoints
  endpoints: {
    events: "/events",
    jobs: "/api/jobs",
  },
};

// Export a function to get a full API URL
export function getApiUrl(endpoint: keyof typeof API_CONFIG.endpoints): string {
  const baseUrl = getBaseUrl();
  // If baseUrl is empty (during SSR), just return the relative path
  if (!baseUrl) {
    return API_CONFIG.endpoints[endpoint];
  }
  return `${baseUrl}${API_CONFIG.endpoints[endpoint]}`;
}
