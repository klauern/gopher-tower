import { useRouter } from 'next/router';
import React, { useState } from 'react';
import { JobFilters } from '../../components/jobs/JobFilters';
import { JobList } from '../../components/jobs/JobList';
import { JobStatus } from '../../types/jobs';

const JobsPage: React.FC = () => {
  const router = useRouter();
  const [selectedStatus, setSelectedStatus] = useState<JobStatus | ''>('');
  const [page, setPage] = useState(1);
  const pageSize = 10;

  const handleStatusChange = (status: JobStatus | '') => {
    setSelectedStatus(status);
    setPage(1); // Reset to first page when filter changes
  };

  const handleCreateJob = () => {
    router.push('/jobs/new');
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Jobs</h1>
        <button
          onClick={handleCreateJob}
          className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium"
        >
          Create Job
        </button>
      </div>
      <div className="mb-6">
        <JobFilters
          selectedStatus={selectedStatus}
          onStatusChange={handleStatusChange}
        />
      </div>
      <JobList
        status={selectedStatus}
        page={page}
        pageSize={pageSize}
        onPageChange={setPage}
      />
    </div>
  );
};

export default JobsPage;
