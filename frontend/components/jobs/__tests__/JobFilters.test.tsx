import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { JobStatus } from '../../../types/jobs';
import { JobFilters } from '../JobFilters';

describe('JobFilters', () => {
  it('renders with default selection', () => {
    render(<JobFilters selectedStatus="" onStatusChange={() => { }} />);
    // Target the native select by its associated label
    const select = screen.getByLabelText('Filter by Status:');
    expect(select).toBeInTheDocument();
    // Check the initial value of the native select
    expect(select).toHaveValue('all'); // Assuming '' maps to 'all' value
  });

  it('renders all status options', async () => {
    render(<JobFilters selectedStatus="" onStatusChange={() => { }} />);

    // Target the native select
    const select = screen.getByTestId('native-select') as HTMLSelectElement;
    expect(select).toBeInTheDocument();

    // Get options directly from the native select
    const options = Array.from(select.options);
    const statuses: JobStatus[] = ['pending', 'active', 'complete', 'failed'];

    // Check values and text content of native options
    expect(options[0]).toHaveValue('all'); // First value is 'all'
    expect(options[0]).toHaveTextContent('All Status');

    statuses.forEach((status, index) => {
      const option = options[index + 1]; // +1 because "All Status" is first
      expect(option).toHaveValue(status);
      expect(option).toHaveTextContent(status.charAt(0).toUpperCase() + status.slice(1));
    });
  });

  it('calls onStatusChange when selection changes', async () => {
    const handleStatusChange = vi.fn();
    render(<JobFilters selectedStatus="" onStatusChange={handleStatusChange} />);

    // Target the native select element
    const select = screen.getByTestId('native-select');

    // Use userEvent to change the value
    await userEvent.selectOptions(select, 'active');

    expect(handleStatusChange).toHaveBeenCalledWith('active');
  });

  it('shows correct selected option', () => {
    render(<JobFilters selectedStatus="active" onStatusChange={() => { }} />);
    // Target the native select element
    const select = screen.getByTestId('native-select') as HTMLSelectElement;
    // Check its value directly
    expect(select.value).toBe('active');
  });
});
