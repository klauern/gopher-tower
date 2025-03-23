import { Navbar } from '@/components/Navbar';
import { useRouter } from 'next/router';
import React from 'react';
import { JobDetail } from '../../components/jobs/JobDetail';

const JobDetailPage: React.FC = () => {
  const router = useRouter();
  const { id } = router.query;

  const handleBack = () => {
    router.push('/jobs');
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <Navbar />
      <div className="container mx-auto px-4 py-8">
        {id && <JobDetail jobId={id as string} onBack={handleBack} />}
      </div>
    </div>
  );
};

export default JobDetailPage;
