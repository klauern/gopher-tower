'use client';

import { useState } from 'react';
import { JobFilters } from '../components/jobs/JobFilters';
import { JobForm } from '../components/jobs/JobForm';
import { JobList } from '../components/jobs/JobList';
import { Button } from '../components/ui/Button';
import { JobStatus } from '../types/jobs';

export default function JobsPage() {
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [selectedStatus, setSelectedStatus] = useState<JobStatus | ''>('');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(10);

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Jobs Dashboard</h1>
        <Button
          onClick={() => setShowCreateForm(true)}
          className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded"
        >
          Create New Job
        </Button>
      </div>

      <JobFilters
        selectedStatus={selectedStatus}
        onStatusChange={setSelectedStatus}
      />

      <div className="mt-6">
        <JobList
          status={selectedStatus}
          page={currentPage}
          pageSize={pageSize}
          onPageChange={setCurrentPage}
        />
      </div>

      {showCreateForm && (
        <JobForm
          onClose={() => setShowCreateForm(false)}
          onSuccess={() => {
            setShowCreateForm(false);
            // Refresh job list
            setCurrentPage(1);
          }}
        />
      )}
    </div>
  );
}
