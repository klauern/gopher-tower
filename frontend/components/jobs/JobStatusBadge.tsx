import { JobStatus } from '../../types/jobs';

interface JobStatusBadgeProps {
  status: JobStatus;
}

const statusColors: Record<JobStatus, { bg: string; text: string }> = {
  pending: { bg: 'bg-yellow-100', text: 'text-yellow-800' },
  active: { bg: 'bg-blue-100', text: 'text-blue-800' },
  complete: { bg: 'bg-green-100', text: 'text-green-800' },
  failed: { bg: 'bg-red-100', text: 'text-red-800' },
};

export function JobStatusBadge({ status }: JobStatusBadgeProps) {
  const colors = statusColors[status];

  return (
    <span
      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${colors.bg} ${colors.text}`}
    >
      {status.charAt(0).toUpperCase() + status.slice(1)}
    </span>
  );
}
