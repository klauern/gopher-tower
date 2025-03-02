/**
 * Application configuration
 *
 * This file contains environment-specific configuration settings.
 */

// Determine if we're in development mode
const isDevelopment = process.env.NODE_ENV === "development";

// Get API base URL from environment variable or use default
const apiBaseUrl =
  process.env.NEXT_PUBLIC_API_BASE_URL ||
  (isDevelopment ? "http://localhost:8080" : "");

// API configuration
export const API_CONFIG = {
  // Base URL for API calls
  baseUrl: apiBaseUrl,

  // Endpoints
  endpoints: {
    events: "/events",
  },
};

// Export a function to get a full API URL
export function getApiUrl(endpoint: keyof typeof API_CONFIG.endpoints): string {
  return `${API_CONFIG.baseUrl}${API_CONFIG.endpoints[endpoint]}`;
}
