import { Navbar } from '@/components/Navbar';
import { useRouter } from 'next/router';
import React from 'react';
import { JobForm } from '../../components/jobs/JobForm';

const CreateJobPage: React.FC = () => {
  const router = useRouter();

  const handleClose = () => {
    router.push('/jobs');
  };

  const handleSuccess = () => {
    router.push('/jobs');
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <Navbar />
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-6">Create New Job</h1>
        <div className="bg-white dark:bg-gray-800 shadow-sm rounded-lg p-6">
          <JobForm
            onClose={handleClose}
            onSuccess={handleSuccess}
          />
        </div>
      </div>
    </div>
  );
};

export default CreateJobPage;
