import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { JobRequest } from '../../../types/jobs';
import { JobForm } from '../JobForm';

const mockInitialData: JobRequest = {
  name: 'Test Job',
  description: 'Test Description',
  status: 'pending',
  startDate: '2024-03-23T00:00:00Z',
  endDate: '2024-03-24T00:00:00Z',
};

describe('JobForm', () => {
  it('renders create form by default', () => {
    render(
      <JobForm
        onClose={() => {}}
        onSuccess={() => {}}
      />
    );

    expect(screen.getByText('Create New Job')).toBeInTheDocument();
    expect(screen.getByLabelText('Name')).toBeInTheDocument();
    expect(screen.getByLabelText('Description')).toBeInTheDocument();
  });

  it('renders edit form with initial data', () => {
    render(
      <JobForm
        onClose={() => {}}
        onSuccess={() => {}}
        initialData={mockInitialData}
      />
    );

    expect(screen.getByText('Edit Job')).toBeInTheDocument();
    expect(screen.getByDisplayValue('Test Job')).toBeInTheDocument();
    expect(screen.getByDisplayValue('Test Description')).toBeInTheDocument();
  });

  it('handles form submission', async () => {
    const onSuccess = vi.fn();
    global.fetch = vi.fn().mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({}),
    });

    render(
      <JobForm
        onClose={() => {}}
        onSuccess={onSuccess}
      />
    );

    fireEvent.change(screen.getByLabelText('Name'), {
      target: { value: 'New Job' },
    });
    fireEvent.change(screen.getByLabelText('Description'), {
      target: { value: 'New Description' },
    });
    fireEvent.click(screen.getByText('Save'));

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledWith('/api/jobs', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: expect.stringContaining('New Job'),
      });
      expect(onSuccess).toHaveBeenCalled();
    });
  });

  it('handles submission errors', async () => {
    global.fetch = vi.fn().mockRejectedValueOnce(new Error('Failed to create job'));

    render(
      <JobForm
        onClose={() => {}}
        onSuccess={() => {}}
      />
    );

    fireEvent.change(screen.getByLabelText('Name'), {
      target: { value: 'New Job' },
    });
    fireEvent.click(screen.getByText('Save'));

    await waitFor(() => {
      expect(screen.getByText('Failed to create job')).toBeInTheDocument();
    });
  });

  it('calls onClose when cancel is clicked', () => {
    const onClose = vi.fn();
    render(
      <JobForm
        onClose={onClose}
        onSuccess={() => {}}
      />
    );

    fireEvent.click(screen.getByText('Cancel'));
    expect(onClose).toHaveBeenCalled();
  });
});
