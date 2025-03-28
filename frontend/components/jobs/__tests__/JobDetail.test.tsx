import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { JobDetail } from '../JobDetail';

// Mock the config module
vi.mock('@/config', () => ({
  getApiUrl: (path: string) => `/api/${path}`,
}));

describe('JobDetail', () => {
  const mockJob = {
    id: '1',
    name: 'Test Job',
    description: 'Test Description',
    status: 'active',
    createdAt: '2024-03-23T00:00:00Z',
    startDate: '2024-03-23T00:00:00Z',
    endDate: null,
  };

  const mockOnBack = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  it('shows loading state initially', () => {
    render(<JobDetail jobId="1" onBack={mockOnBack} />);
    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });

  it('displays job details when fetch is successful', async () => {
    vi.mocked(global.fetch).mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockJob),
    } as Response);

    render(<JobDetail jobId="1" onBack={mockOnBack} />);

    await waitFor(() => {
      expect(screen.getByText(mockJob.name)).toBeInTheDocument();
    });

    expect(screen.getByText(mockJob.description)).toBeInTheDocument();
    expect(screen.getByText(/Created:/)).toBeInTheDocument();
    expect(screen.getByText('Not completed')).toBeInTheDocument();
  });

  it('displays error message when fetch fails', async () => {
    const errorMessage = 'Failed to fetch job details';
    vi.mocked(global.fetch).mockResolvedValueOnce({
      ok: false,
      text: () => Promise.resolve('Error occurred'),
    } as Response);

    render(<JobDetail jobId="1" onBack={mockOnBack} />);

    await waitFor(() => {
      expect(screen.getByText(errorMessage)).toBeInTheDocument();
    });
  });

  it('displays network error message when fetch throws', async () => {
    vi.mocked(global.fetch).mockRejectedValueOnce(new Error('Network error'));

    render(<JobDetail jobId="1" onBack={mockOnBack} />);

    await waitFor(() => {
      expect(screen.getByText('Network error')).toBeInTheDocument();
    });
  });

  it('calls onBack when back button is clicked', async () => {
    vi.mocked(global.fetch).mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockJob),
    } as Response);

    render(<JobDetail jobId="1" onBack={mockOnBack} />);

    await waitFor(() => {
      expect(screen.getByText('Back to Jobs')).toBeInTheDocument();
    });

    await userEvent.click(screen.getByText('Back to Jobs'));
    expect(mockOnBack).toHaveBeenCalledTimes(1);
  });

  it('shows "Job not found" when job data is null', async () => {
    vi.mocked(global.fetch).mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(null),
    } as Response);

    render(<JobDetail jobId="1" onBack={mockOnBack} />);

    await waitFor(() => {
      expect(screen.getByText('Job not found')).toBeInTheDocument();
    });
  });

  it('formats dates correctly', async () => {
    const jobWithDates = {
      ...mockJob,
      startDate: '2024-03-23T00:00:00Z',
      endDate: '2024-03-24T00:00:00Z',
    };

    vi.mocked(global.fetch).mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(jobWithDates),
    } as Response);

    render(<JobDetail jobId="1" onBack={mockOnBack} />);

    await waitFor(() => {
      // Check that the created date text contains both the label and the formatted date
      expect(screen.getByText(/Created: Mar 23, 2024, 12:00 AM/)).toBeInTheDocument();

      // Check start and end dates
      expect(screen.getByText('Mar 23, 2024, 12:00 AM')).toBeInTheDocument(); // Start date
      expect(screen.getByText('Mar 24, 2024, 12:00 AM')).toBeInTheDocument(); // End date
    });
  });
});
