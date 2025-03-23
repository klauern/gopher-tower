import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { JobStatus } from '../../../types/jobs';
import { JobFilters } from '../JobFilters';

describe('JobFilters', () => {
  it('renders with default selection', () => {
    render(
      <JobFilters
        selectedStatus=""
        onStatusChange={() => {}}
      />
    );

    expect(screen.getByLabelText('Filter by Status:')).toBeInTheDocument();
    expect(screen.getByText('All Status')).toBeInTheDocument();
  });

  it('renders all status options', () => {
    render(
      <JobFilters
        selectedStatus=""
        onStatusChange={() => {}}
      />
    );

    const statuses: JobStatus[] = ['pending', 'active', 'complete', 'failed'];
    statuses.forEach(status => {
      expect(screen.getByText(status.charAt(0).toUpperCase() + status.slice(1))).toBeInTheDocument();
    });
  });

  it('calls onStatusChange when selection changes', () => {
    const onStatusChange = vi.fn();
    render(
      <JobFilters
        selectedStatus=""
        onStatusChange={onStatusChange}
      />
    );

    fireEvent.change(screen.getByLabelText('Filter by Status:'), {
      target: { value: 'active' },
    });

    expect(onStatusChange).toHaveBeenCalledWith('active');
  });

  it('shows correct selected option', () => {
    render(
      <JobFilters
        selectedStatus="pending"
        onStatusChange={() => {}}
      />
    );

    const select = screen.getByLabelText('Filter by Status:') as HTMLSelectElement;
    expect(select.value).toBe('pending');
  });
});
