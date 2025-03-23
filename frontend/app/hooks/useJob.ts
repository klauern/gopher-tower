import { getApiUrl } from "@/app/config";
import { Job } from "@/app/types/jobs";
import { useEffect, useState } from "react";

export function useJob(id: string) {
  const [job, setJob] = useState<Job | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchJob = async () => {
      try {
        setLoading(true);
        const baseUrl = getApiUrl("jobs");
        // Use the current origin if we got a relative URL
        const url = baseUrl.startsWith("/")
          ? `${window.location.origin}${baseUrl}/${id}`
          : `${baseUrl}/${id}`;

        const response = await fetch(url);

        if (!response.ok) {
          throw new Error("Failed to fetch job details");
        }

        const jobData: Job = await response.json();
        setJob(jobData);
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : "An error occurred");
      } finally {
        setLoading(false);
      }
    };

    fetchJob();
  }, [id]);

  return { job, loading, error };
}
