'use client';

import { JobStatusBadge } from '@/app/components/jobs/JobStatusBadge';
import { useJob } from '@/app/hooks/useJob';
import { formatDate } from '@/app/utils/date';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';

export default function JobDetailPage() {
  const searchParams = useSearchParams();
  const jobId = searchParams.get('id');
  const { job, loading, error } = useJob(jobId || '');

  if (!jobId) {
    return (
      <div className="max-w-4xl mx-auto p-6">
        <div className="mb-4">
          <Link
            href="/jobs"
            className="text-blue-500 hover:text-blue-600 transition-colors"
          >
            ← Back to Jobs
          </Link>
        </div>
        <div className="bg-white shadow-sm rounded-lg p-6">
          <div className="text-center py-4">No job ID provided</div>
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="max-w-4xl mx-auto p-6">
        <div className="mb-4">
          <Link
            href="/jobs"
            className="text-blue-500 hover:text-blue-600 transition-colors"
          >
            ← Back to Jobs
          </Link>
        </div>
        <div className="bg-white shadow-sm rounded-lg p-6">
          <div className="animate-pulse">
            <div className="h-8 bg-gray-200 rounded w-1/3 mb-4"></div>
            <div className="h-4 bg-gray-200 rounded w-2/3 mb-6"></div>
            <div className="grid grid-cols-2 gap-6">
              {[1, 2, 3, 4].map((i) => (
                <div key={i}>
                  <div className="h-4 bg-gray-200 rounded w-1/4 mb-2"></div>
                  <div className="h-4 bg-gray-200 rounded w-1/2"></div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-4xl mx-auto p-6">
        <div className="mb-4">
          <Link
            href="/jobs"
            className="text-blue-500 hover:text-blue-600 transition-colors"
          >
            ← Back to Jobs
          </Link>
        </div>
        <div className="bg-white shadow-sm rounded-lg p-6">
          <div className="text-red-600 text-center py-4">{error}</div>
        </div>
      </div>
    );
  }

  if (!job) {
    return (
      <div className="max-w-4xl mx-auto p-6">
        <div className="mb-4">
          <Link
            href="/jobs"
            className="text-blue-500 hover:text-blue-600 transition-colors"
          >
            ← Back to Jobs
          </Link>
        </div>
        <div className="bg-white shadow-sm rounded-lg p-6">
          <div className="text-center py-4">Job not found</div>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto p-6">
      <div className="mb-4">
        <Link
          href="/jobs"
          className="text-blue-500 hover:text-blue-600 transition-colors"
        >
          ← Back to Jobs
        </Link>
      </div>
      <div className="bg-white shadow-sm rounded-lg p-6">
        <div className="flex justify-between items-start mb-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{job.name}</h1>
            <p className="mt-2 text-gray-600">{job.description}</p>
          </div>
          <JobStatusBadge status={job.status} />
        </div>

        <div className="grid grid-cols-2 gap-6 mt-6">
          <div>
            <h2 className="text-sm font-medium text-gray-500">Start Date</h2>
            <p className="mt-1 text-gray-900">{job.startDate ? formatDate(job.startDate) : '-'}</p>
          </div>
          <div>
            <h2 className="text-sm font-medium text-gray-500">End Date</h2>
            <p className="mt-1 text-gray-900">{job.endDate ? formatDate(job.endDate) : '-'}</p>
          </div>
          <div>
            <h2 className="text-sm font-medium text-gray-500">Created At</h2>
            <p className="mt-1 text-gray-900">{formatDate(job.createdAt)}</p>
          </div>
          <div>
            <h2 className="text-sm font-medium text-gray-500">Last Updated</h2>
            <p className="mt-1 text-gray-900">{formatDate(job.updatedAt)}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
