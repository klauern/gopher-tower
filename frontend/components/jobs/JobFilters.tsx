import { JobStatus } from '../../types/jobs';

interface JobFiltersProps {
  selectedStatus: JobStatus | '';
  onStatusChange: (status: JobStatus | '') => void;
}

const statusOptions: { value: JobStatus | ''; label: string }[] = [
  { value: '', label: 'All Status' },
  { value: 'pending', label: 'Pending' },
  { value: 'active', label: 'Active' },
  { value: 'complete', label: 'Complete' },
  { value: 'failed', label: 'Failed' },
];

export function JobFilters({ selectedStatus, onStatusChange }: JobFiltersProps) {
  return (
    <div className="flex items-center space-x-4">
      <label htmlFor="status" className="text-sm font-medium text-gray-700">
        Filter by Status:
      </label>
      <select
        id="status"
        value={selectedStatus}
        onChange={(e) => onStatusChange(e.target.value as JobStatus | '')}
        className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
      >
        {statusOptions.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </div>
  );
}
