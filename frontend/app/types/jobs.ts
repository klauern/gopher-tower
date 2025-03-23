export type JobStatus = "pending" | "active" | "complete" | "failed";

export interface Job {
  id: string;
  name: string;
  description: string;
  status: JobStatus;
  startDate?: string;
  endDate?: string;
  createdAt: string;
  updatedAt: string;
  ownerId?: string;
}

export interface JobRequest {
  name: string;
  description: string;
  status: JobStatus;
  startDate?: string;
  endDate?: string;
}

export interface JobListParams {
  page: number;
  pageSize: number;
  status?: JobStatus;
}

export interface JobListResponse {
  jobs: Job[];
  totalCount: number;
  page: number;
  pageSize: number;
}
