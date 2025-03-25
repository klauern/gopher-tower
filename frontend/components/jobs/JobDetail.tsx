'use client';

import { Button } from "@/components/ui/Button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { getApiUrl } from '@/config';
import { ArrowLeft } from "lucide-react";
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
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center text-muted-foreground">Loading...</div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center text-destructive">{error}</div>
        </CardContent>
      </Card>
    );
  }

  if (!job) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center text-muted-foreground">Job not found</div>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <Button
            variant="ghost"
            size="icon"
            onClick={onBack}
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <h1 className="text-3xl font-bold">Job Details</h1>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{job.name}</CardTitle>
          {job.description && (
            <CardDescription>{job.description}</CardDescription>
          )}
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="flex items-center space-x-4">
            <JobStatusBadge status={job.status} />
            <span className="text-sm text-muted-foreground">
              Created: {formatDate(job.createdAt)}
            </span>
          </div>

          <div className="grid grid-cols-2 gap-6">
            <div className="space-y-1">
              <p className="text-sm font-medium text-muted-foreground">Start Date</p>
              <p>{job.startDate ? formatDate(job.startDate) : 'Not started'}</p>
            </div>
            <div className="space-y-1">
              <p className="text-sm font-medium text-muted-foreground">End Date</p>
              <p>{job.endDate ? formatDate(job.endDate) : 'Not completed'}</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
