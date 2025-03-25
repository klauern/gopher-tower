'use client';

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
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
      <label htmlFor="status" className="text-sm font-medium">
        Filter by Status:
      </label>
      <Select
        value={selectedStatus || undefined}
        onValueChange={(value) => onStatusChange(value as JobStatus | '')}
      >
        <SelectTrigger className="w-[180px]">
          <SelectValue placeholder="All Status" />
        </SelectTrigger>
        <SelectContent>
          {statusOptions.map((option) => (
            <SelectItem
              key={option.value || 'all'}
              value={option.value || 'all'}
            >
              {option.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
