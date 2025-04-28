import { Navbar } from "@/components/Navbar";
import { useRouter } from "next/router";
import React from "react";
import { JobDetail } from "../../components/jobs/JobDetail";

const uuidRegex =
  /^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/;

const isValidJobId = (id: string) => uuidRegex.test(id);

const JobDetailPage: React.FC = () => {
  const router = useRouter();
  const { id } = router.query;

  const handleBack = () => {
    router.push("/jobs");
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <Navbar />
      <div className="container mx-auto px-4 py-8">
        {id && isValidJobId(id as string) ? (
          <JobDetail jobId={id as string} onBack={handleBack} />
        ) : id ? (
          <div className="text-center text-destructive">
            Invalid job ID format.
          </div>
        ) : null}
      </div>
    </div>
  );
};

export default JobDetailPage;
