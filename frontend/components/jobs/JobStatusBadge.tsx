'use client';

import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { JobStatus } from '../../types/jobs';

interface JobStatusBadgeProps {
  status: JobStatus;
}

const statusStyles: Record<JobStatus, { variant: "default" | "secondary" | "destructive" | "outline", className: string }> = {
  pending: {
    variant: "secondary",
    className: "bg-yellow-100 hover:bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-500"
  },
  active: {
    variant: "default",
    className: "bg-blue-100 hover:bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-500"
  },
  complete: {
    variant: "default",
    className: "bg-green-100 hover:bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-500"
  },
  failed: {
    variant: "destructive",
    className: "bg-red-100 hover:bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-500"
  },
};

export function JobStatusBadge({ status }: JobStatusBadgeProps) {
  const style = statusStyles[status];
  const label = status.charAt(0).toUpperCase() + status.slice(1);

  return (
    <Badge
      variant={style.variant}
      className={cn(
        "font-medium",
        style.className
      )}
    >
      {label}
    </Badge>
  );
}
