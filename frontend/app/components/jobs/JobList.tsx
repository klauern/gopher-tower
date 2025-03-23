'use client';

import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { JobListResponse, JobStatus } from '../../types/jobs';
import { formatDate } from '../../utils/date';
import { JobStatusBadge } from './JobStatusBadge';

interface JobListProps {
  status?: JobStatus | '';
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
}

export function JobList({ status, page, pageSize, onPageChange }: JobListProps) {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [data, setData] = useState<JobListResponse | null>(null);

  useEffect(() => {
    const fetchJobs = async () => {
      try {
        setLoading(true);
        const params = new URLSearchParams({
          page: page.toString(),
          page_size: pageSize.toString(),
        });
        if (status) {
          params.append('status', status);
        }

        // Use the environment variable for the API base URL
        const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || '';
        const response = await fetch(`${apiBaseUrl}/api/jobs?${params.toString()}`);
        if (!response.ok) {
          throw new Error('Failed to fetch jobs');
        }

        const jobData: JobListResponse = await response.json();
        setData(jobData);
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An error occurred');
      } finally {
        setLoading(false);
      }
    };

    fetchJobs();
  }, [page, pageSize, status]);

  const handleJobClick = (jobId: string) => {
    router.push(`/jobs/detail?id=${jobId}`);
  };

  if (loading) {
    return <div className="text-center py-4">Loading jobs...</div>;
  }

  if (error) {
    return <div className="text-red-600 text-center py-4">{error}</div>;
  }

  if (!data || data.jobs.length === 0) {
    return <div className="text-center py-4">No jobs found</div>;
  }

  const totalPages = Math.ceil(data.totalCount / pageSize);

  return (
    <div>
      <div className="overflow-x-auto">
        <table className="min-w-full bg-white shadow-sm rounded-lg">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Name
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Start Date
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                End Date
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Created At
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {data.jobs.map((job) => (
              <tr
                key={job.id}
                className="hover:bg-gray-50 cursor-pointer transition-colors"
                onClick={() => handleJobClick(job.id)}
              >
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm font-medium text-gray-900">
                    {job.name}
                  </div>
                  <div className="text-sm text-gray-500">{job.description}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <JobStatusBadge status={job.status} />
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {job.startDate ? formatDate(job.startDate) : '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {job.endDate ? formatDate(job.endDate) : '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatDate(job.createdAt)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="flex justify-between items-center mt-4">
        <div className="text-sm text-gray-700">
          Showing {(page - 1) * pageSize + 1} to{' '}
          {Math.min(page * pageSize, data.totalCount)} of {data.totalCount} jobs
        </div>
        <div className="flex space-x-2">
          <button
            onClick={() => onPageChange(page - 1)}
            disabled={page === 1}
            className="px-3 py-1 border rounded disabled:opacity-50"
          >
            Previous
          </button>
          <button
            onClick={() => onPageChange(page + 1)}
            disabled={page >= totalPages}
            className="px-3 py-1 border rounded disabled:opacity-50"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  );
}
