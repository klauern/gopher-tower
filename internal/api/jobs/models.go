package jobs

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// JobStatus represents the current state of a job
type JobStatus string

const (
	JobStatusPending  JobStatus = "pending"
	JobStatusActive   JobStatus = "active"
	JobStatusComplete JobStatus = "complete"
	JobStatusFailed   JobStatus = "failed"
)

// IsValid checks if the job status is valid
func (s JobStatus) IsValid() bool {
	switch s {
	case JobStatusPending, JobStatusActive, JobStatusComplete, JobStatusFailed:
		return true
	default:
		return false
	}
}

// JobRequest represents the request to create or update a job
type JobRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      JobStatus  `json:"status"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

// Validate checks if the job request is valid
func (r *JobRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}

	switch r.Status {
	case JobStatusPending, JobStatusActive, JobStatusComplete, JobStatusFailed:
		return nil
	default:
		return errors.New("invalid status")
	}
}

// JobResponse represents a job in responses
type JobResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      JobStatus  `json:"status"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	OwnerID     string     `json:"owner_id,omitempty"`
}

// JobListParams represents parameters for listing jobs
type JobListParams struct {
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	Status   JobStatus `json:"status,omitempty"`
}

// Validate checks if the list parameters are valid
func (p *JobListParams) Validate() error {
	if p.Page < 1 {
		return errors.New("page must be greater than 0")
	}
	if p.PageSize < 1 {
		return errors.New("page_size must be greater than 0")
	}
	if p.Status != "" {
		switch p.Status {
		case JobStatusPending, JobStatusActive, JobStatusComplete, JobStatusFailed:
			return nil
		default:
			return errors.New("invalid status")
		}
	}
	return nil
}

// JobListResponse represents the response for listing jobs
type JobListResponse struct {
	Jobs       []JobResponse `json:"jobs"`
	TotalCount int64         `json:"total_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
}

// NewJobID generates a new UUID for a job
func NewJobID() string {
	return uuid.New().String()
}
