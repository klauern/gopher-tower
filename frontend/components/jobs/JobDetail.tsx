import { getApiUrl } from '@/config';
import { useEffect, useState } from 'react';
import { Job } from '../../types/jobs';
import { formatDate } from '../../utils/date';
import { JobStatusBadge } from './JobStatusBadge';

interface JobDetailProps {
  jobId: string;
  onBack: () => void;
}

export function JobDetail({ jobId, onBack }: JobDetailProps) {
  const [job, setJob] = useState<Job | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchJob = async () => {
      if (!jobId) return;

      try {
        setLoading(true);
        const url = getApiUrl(`jobs/${jobId}`);
        console.log('Fetching job details from:', url);

        const response = await fetch(url);
        console.log('Response status:', response.status);

        if (!response.ok) {
          const errorText = await response.text();
          console.error('Error response:', errorText);
          throw new Error('Failed to fetch job details');
        }

        const data = await response.json();
        console.log('Job data:', data);
        setJob(data);
        setError(null);
      } catch (err) {
        console.error('Error fetching job:', err);
        setError(err instanceof Error ? err.message : 'An error occurred');
      } finally {
        setLoading(false);
      }
    };

    console.log('JobDetail mounted with jobId:', jobId);
    fetchJob();
  }, [jobId]);

  if (loading) {
    return <div className="text-center">Loading...</div>;
  }

  if (error) {
    return <div className="text-red-600 text-center">{error}</div>;
  }

  if (!job) {
    return <div className="text-center">Job not found</div>;
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Job Details</h1>
        <button
          onClick={onBack}
          className="text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
        >
          Back to Jobs
        </button>
      </div>
      <div className="bg-white dark:bg-gray-800 shadow-sm rounded-lg p-6">
        <div className="space-y-6">
          <div>
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">{job.name}</h2>
            <p className="text-gray-600 dark:text-gray-300">{job.description}</p>
          </div>
          <div className="flex items-center space-x-4">
            <JobStatusBadge status={job.status} />
            <span className="text-sm text-gray-500 dark:text-gray-400">
              Created: {formatDate(job.createdAt)}
            </span>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Start Date</p>
              <p className="mt-1">{job.startDate ? formatDate(job.startDate) : 'Not started'}</p>
            </div>
            <div>
              <p className="text-sm font-medium text-gray-500 dark:text-gray-400">End Date</p>
              <p className="mt-1">{job.endDate ? formatDate(job.endDate) : 'Not completed'}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
