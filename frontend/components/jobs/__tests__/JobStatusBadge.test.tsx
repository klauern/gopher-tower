import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { JobStatus } from '../../../types/jobs';
import { JobStatusBadge } from '../JobStatusBadge';

describe('JobStatusBadge', () => {
  const statuses: JobStatus[] = ['pending', 'active', 'complete', 'failed'];

  it.each(statuses)('renders %s status with correct styling', (status) => {
    render(<JobStatusBadge status={status} />);
    const badge = screen.getByText(status.charAt(0).toUpperCase() + status.slice(1));
    expect(badge).toBeInTheDocument();

    // Check for correct color classes based on status
    if (status === 'pending') {
      expect(badge.className).toContain('bg-yellow-100');
      expect(badge.className).toContain('text-yellow-800');
    } else if (status === 'active') {
      expect(badge.className).toContain('bg-blue-100');
      expect(badge.className).toContain('text-blue-800');
    } else if (status === 'complete') {
      expect(badge.className).toContain('bg-green-100');
      expect(badge.className).toContain('text-green-800');
    } else if (status === 'failed') {
      expect(badge.className).toContain('bg-red-100');
      expect(badge.className).toContain('text-red-800');
    }
  });
});
