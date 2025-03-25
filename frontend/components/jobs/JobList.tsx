'use client';

import { Button } from "@/components/ui/Button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { getApiUrl } from '@/config';
import { useRouter } from 'next/router';
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

        const url = getApiUrl('jobs');
        console.log('Fetching jobs from:', url);

        const response = await fetch(`${url}?${params.toString()}`);
        console.log('Jobs response status:', response.status);

        if (!response.ok) {
          const errorText = await response.text();
          console.error('Error response:', errorText);
          throw new Error('Failed to fetch jobs');
        }

        const jobData: JobListResponse = await response.json();
        console.log('Jobs data:', jobData);
        setData(jobData);
        setError(null);
      } catch (err) {
        console.error('Error fetching jobs:', err);
        setError(err instanceof Error ? err.message : 'An error occurred');
      } finally {
        setLoading(false);
      }
    };

    fetchJobs();
  }, [page, pageSize, status]);

  const handleJobClick = (jobId: string) => {
    console.log('Navigating to job:', jobId);
    router.push(`/jobs/${jobId}`);
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

  const hasNextPage = data.totalCount > page * pageSize;

  return (
    <div>
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Start Date</TableHead>
              <TableHead>End Date</TableHead>
              <TableHead>Created At</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.jobs.map((job) => (
              <TableRow
                key={job.id}
                className="cursor-pointer hover:bg-muted/50"
                onClick={() => handleJobClick(job.id)}
              >
                <TableCell>
                  <div className="font-medium">{job.name}</div>
                  <div className="text-sm text-muted-foreground">{job.description}</div>
                </TableCell>
                <TableCell>
                  <JobStatusBadge status={job.status} />
                </TableCell>
                <TableCell className="text-muted-foreground">
                  {job.startDate ? formatDate(job.startDate) : '-'}
                </TableCell>
                <TableCell className="text-muted-foreground">
                  {job.endDate ? formatDate(job.endDate) : '-'}
                </TableCell>
                <TableCell className="text-muted-foreground">
                  {formatDate(job.createdAt)}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      <div className="flex items-center justify-between space-x-2 py-4">
        <div className="text-sm text-muted-foreground">
          Showing {data.jobs.length === 0 ? 0 : (page - 1) * pageSize + 1} to{' '}
          {(page - 1) * pageSize + data.jobs.length} of {data.totalCount} jobs
        </div>
        <div className="flex space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onPageChange(page - 1)}
            disabled={page === 1}
          >
            Previous
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => onPageChange(page + 1)}
            disabled={!hasNextPage}
          >
            Next
          </Button>
        </div>
      </div>
    </div>
  );
}
