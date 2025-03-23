//go:generate go tool mockgen -destination=mock_querier_test.go -package=jobs github.com/klauern/gopher-tower/internal/api/jobs JobQuerier
//go:generate go tool mockgen -destination=mock_service_test.go -package=jobs github.com/klauern/gopher-tower/internal/api/jobs Service

package jobs

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/klauern/gopher-tower/internal/db"
)

var (
	ErrJobNotFound = errors.New("job not found")
	ErrInvalidJob  = errors.New("invalid job data")
)

// JobQuerier defines the interface for job-related database operations
type JobQuerier interface {
	CreateJob(ctx context.Context, arg db.CreateJobParams) (db.Job, error)
	GetJob(ctx context.Context, id string) (db.Job, error)
	UpdateJob(ctx context.Context, arg db.UpdateJobParams) (db.Job, error)
	DeleteJob(ctx context.Context, id string) error
	ListJobs(ctx context.Context, arg db.ListJobsParams) ([]db.Job, error)
}

// Service provides job management operations
type Service interface {
	CreateJob(ctx context.Context, req JobRequest, ownerID string) (*JobResponse, error)
	GetJob(ctx context.Context, id string) (*JobResponse, error)
	UpdateJob(ctx context.Context, id string, req JobRequest) (*JobResponse, error)
	DeleteJob(ctx context.Context, id string) error
	ListJobs(ctx context.Context, params JobListParams) (*JobListResponse, error)
}

// jobService implements the Service interface
type jobService struct {
	queries JobQuerier
}

// NewService creates a new job service
func NewService(queries JobQuerier) Service {
	return &jobService{queries: queries}
}

// generateID generates a new UUID for job IDs
func generateID() string {
	return uuid.New().String()
}

// CreateJob creates a new job
func (s *jobService) CreateJob(ctx context.Context, req JobRequest, ownerID string) (*JobResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidJob
	}

	job, err := s.queries.CreateJob(ctx, db.CreateJobParams{
		ID:          generateID(),
		Name:        req.Name,
		Description: db.StringToNullString(req.Description),
		Status:      string(req.Status),
		StartDate:   db.TimeToNullTime(req.StartDate),
		EndDate:     db.TimeToNullTime(req.EndDate),
		OwnerID:     db.StringToNullString(ownerID),
	})
	if err != nil {
		return nil, err
	}

	return toJobResponse(job), nil
}

// GetJob retrieves a job by ID
func (s *jobService) GetJob(ctx context.Context, id string) (*JobResponse, error) {
	job, err := s.queries.GetJob(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, ErrJobNotFound
		}
		return nil, err
	}

	return toJobResponse(job), nil
}

// UpdateJob updates an existing job
func (s *jobService) UpdateJob(ctx context.Context, id string, req JobRequest) (*JobResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidJob
	}

	job, err := s.queries.UpdateJob(ctx, db.UpdateJobParams{
		ID:          id,
		Name:        req.Name,
		Description: db.StringToNullString(req.Description),
		Status:      string(req.Status),
		StartDate:   db.TimeToNullTime(req.StartDate),
		EndDate:     db.TimeToNullTime(req.EndDate),
	})
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, ErrJobNotFound
		}
		return nil, err
	}

	return toJobResponse(job), nil
}

// DeleteJob deletes a job by ID
func (s *jobService) DeleteJob(ctx context.Context, id string) error {
	err := s.queries.DeleteJob(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return ErrJobNotFound
		}
		return err
	}
	return nil
}

// ListJobs returns a paginated list of jobs
func (s *jobService) ListJobs(ctx context.Context, params JobListParams) (*JobListResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	jobs, err := s.queries.ListJobs(ctx, db.ListJobsParams{
		Limit:  int64((params.Page - 1) * params.PageSize),
		Offset: int64(params.PageSize),
	})
	if err != nil {
		return nil, err
	}

	// Filter jobs by status if specified
	var filteredJobs []db.Job
	if params.Status != "" {
		for _, job := range jobs {
			if job.Status == string(params.Status) {
				filteredJobs = append(filteredJobs, job)
			}
		}
	} else {
		filteredJobs = jobs
	}

	responses := make([]JobResponse, len(filteredJobs))
	for i, job := range filteredJobs {
		responses[i] = *toJobResponse(job)
	}

	return &JobListResponse{
		Jobs:       responses,
		TotalCount: int64(len(filteredJobs)),
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}

// toJobResponse converts a db.Job to a JobResponse
func toJobResponse(job db.Job) *JobResponse {
	return &JobResponse{
		ID:          job.ID,
		Name:        job.Name,
		Description: job.Description.String,
		Status:      JobStatus(job.Status),
		StartDate:   db.NullTimeToTimePtr(job.StartDate),
		EndDate:     db.NullTimeToTimePtr(job.EndDate),
		CreatedAt:   job.CreatedAt,
		UpdatedAt:   job.UpdatedAt,
		OwnerID:     job.OwnerID.String,
	}
}
