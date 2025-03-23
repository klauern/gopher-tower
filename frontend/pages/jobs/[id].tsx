import { useRouter } from 'next/router';
import React from 'react';
import { JobForm } from '../../components/jobs/JobForm';

const JobPage: React.FC = () => {
  const router = useRouter();
  const { id } = router.query;

  const handleClose = () => {
    router.push('/jobs');
  };

  const handleSuccess = () => {
    router.push('/jobs');
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Job Details</h1>
      {id && (
        <JobForm
          onClose={handleClose}
          onSuccess={handleSuccess}
        />
      )}
    </div>
  );
};

export default JobPage;
