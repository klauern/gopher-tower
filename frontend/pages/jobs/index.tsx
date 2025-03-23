import React, { useState } from 'react';
import { JobFilters } from '../../components/jobs/JobFilters';
import { JobList } from '../../components/jobs/JobList';
import { JobStatus } from '../../types/jobs';

const JobsPage: React.FC = () => {
  const [selectedStatus, setSelectedStatus] = useState<JobStatus | ''>('');
  const [page, setPage] = useState(1);
  const pageSize = 10;

  const handleStatusChange = (status: JobStatus | '') => {
    setSelectedStatus(status);
    setPage(1); // Reset to first page when filter changes
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Jobs</h1>
      <JobFilters
        selectedStatus={selectedStatus}
        onStatusChange={handleStatusChange}
      />
      <JobList
        page={page}
        pageSize={pageSize}
        onPageChange={setPage}
      />
    </div>
  );
};

export default JobsPage;
