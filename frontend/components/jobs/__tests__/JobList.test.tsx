import { act, fireEvent, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { renderWithRouter } from '../../../test/utils/test-utils';
import { JobListResponse } from '../../../types/jobs';
import { JobList } from '../JobList';

const mockJobsResponse: JobListResponse = {
  jobs: [
    {
      id: '1',
      name: 'Test Job 1',
      description: 'Test Description 1',
      status: 'pending',
      createdAt: '2024-03-23T00:00:00Z',
      updatedAt: '2024-03-23T00:00:00Z',
    },
    {
      id: '2',
      name: 'Test Job 2',
      description: 'Test Description 2',
      status: 'active',
      startDate: '2024-03-23T00:00:00Z',
      createdAt: '2024-03-23T00:00:00Z',
      updatedAt: '2024-03-23T00:00:00Z',
    },
  ],
  totalCount: 15,
  page: 1,
  pageSize: 10,
};

describe('JobList', () => {
  beforeEach(() => {
    vi.spyOn(global, 'fetch').mockImplementation(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(mockJobsResponse),
      } as Response)
    );
  });

  it('renders loading state initially', async () => {
    let resolveFetch: (value: Response) => void;
    const fetchPromise = new Promise<Response>((resolve) => {
      resolveFetch = resolve;
    });

    vi.spyOn(global, 'fetch').mockImplementationOnce(() => fetchPromise);

    renderWithRouter(
      <JobList
        page={1}
        pageSize={10}
        onPageChange={() => {}}
      />
    );

    expect(screen.getByText('Loading jobs...')).toBeInTheDocument();

    await act(async () => {
      resolveFetch!({
        ok: true,
        json: () => Promise.resolve(mockJobsResponse),
      } as Response);
    });
  });

  it('renders jobs after loading', async () => {
    await act(async () => {
      renderWithRouter(
        <JobList
          page={1}
          pageSize={10}
          onPageChange={() => {}}
        />
      );
    });

    await waitFor(() => {
      expect(screen.getByText('Test Job 1')).toBeInTheDocument();
      expect(screen.getByText('Test Job 2')).toBeInTheDocument();
    });
  });

  it('handles pagination correctly', async () => {
    const onPageChange = vi.fn();
    await act(async () => {
      renderWithRouter(
        <JobList
          page={1}
          pageSize={10}
          onPageChange={onPageChange}
        />
      );
    });

    await waitFor(() => {
      expect(screen.getByText('Test Job 1')).toBeInTheDocument();
      expect(screen.getByText('Test Job 2')).toBeInTheDocument();
    });

    const nextButton = screen.getByText('Next');
    expect(nextButton).toBeInTheDocument();
    expect(nextButton).not.toBeDisabled();

    await act(async () => {
      fireEvent.click(nextButton);
    });

    expect(onPageChange).toHaveBeenCalledWith(2);
  });

  it('handles status filtering', async () => {
    await act(async () => {
      renderWithRouter(
        <JobList
          page={1}
          pageSize={10}
          status="active"
          onPageChange={() => {}}
        />
      );
    });

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('status=active')
      );
    });
  });

  it('handles error state', async () => {
    vi.spyOn(global, 'fetch').mockImplementationOnce(() =>
      Promise.reject(new Error('Failed to fetch'))
    );

    await act(async () => {
      renderWithRouter(
        <JobList
          page={1}
          pageSize={10}
          onPageChange={() => {}}
        />
      );
    });

    await waitFor(() => {
      expect(screen.getByText('Failed to fetch')).toBeInTheDocument();
    });
  });
});
