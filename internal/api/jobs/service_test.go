package jobs

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/klauern/gopher-tower/internal/db"
	"go.uber.org/mock/gomock"
)

// MockQueries implements JobQuerier for testing
type MockQueries struct {
	jobs map[string]db.Job
}

func NewMockQueries() *MockQueries {
	return &MockQueries{
		jobs: make(map[string]db.Job),
	}
}

func (m *MockQueries) CreateJob(ctx context.Context, arg db.CreateJobParams) (db.Job, error) {
	job := db.Job{
		ID:          arg.ID,
		Name:        arg.Name,
		Description: arg.Description,
		Status:      arg.Status,
		StartDate:   arg.StartDate,
		EndDate:     arg.EndDate,
		OwnerID:     arg.OwnerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.jobs[job.ID] = job
	return job, nil
}

func (m *MockQueries) GetJob(ctx context.Context, id string) (db.Job, error) {
	job, ok := m.jobs[id]
	if !ok {
		return db.Job{}, db.ErrNotFound
	}
	return job, nil
}

func (m *MockQueries) UpdateJob(ctx context.Context, arg db.UpdateJobParams) (db.Job, error) {
	job, ok := m.jobs[arg.ID]
	if !ok {
		return db.Job{}, db.ErrNotFound
	}

	job.Name = arg.Name
	job.Description = arg.Description
	job.Status = arg.Status
	job.StartDate = arg.StartDate
	job.EndDate = arg.EndDate
	job.UpdatedAt = time.Now()

	m.jobs[arg.ID] = job
	return job, nil
}

func (m *MockQueries) DeleteJob(ctx context.Context, id string) error {
	if _, ok := m.jobs[id]; !ok {
		return db.ErrNotFound
	}
	delete(m.jobs, id)
	return nil
}

func (m *MockQueries) ListJobs(ctx context.Context, arg db.ListJobsParams) ([]db.Job, error) {
	jobs := make([]db.Job, 0, len(m.jobs))
	for _, job := range m.jobs {
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// Implement other required methods of db.Queries with no-op implementations
func (m *MockQueries) CreateActivityLog(context.Context, db.CreateActivityLogParams) (db.ActivityLog, error) {
	return db.ActivityLog{}, errors.New("not implemented")
}

func (m *MockQueries) CreateAttachment(context.Context, db.CreateAttachmentParams) (db.Attachment, error) {
	return db.Attachment{}, errors.New("not implemented")
}

func (m *MockQueries) CreateComment(context.Context, db.CreateCommentParams) (db.Comment, error) {
	return db.Comment{}, errors.New("not implemented")
}

func (m *MockQueries) CreateNotification(context.Context, db.CreateNotificationParams) (db.Notification, error) {
	return db.Notification{}, errors.New("not implemented")
}

func (m *MockQueries) CreateTask(context.Context, db.CreateTaskParams) (db.Task, error) {
	return db.Task{}, errors.New("not implemented")
}

func (m *MockQueries) CreateUser(context.Context, db.CreateUserParams) (db.User, error) {
	return db.User{}, errors.New("not implemented")
}

func (m *MockQueries) DeleteAttachment(context.Context, string) error {
	return errors.New("not implemented")
}

func (m *MockQueries) DeleteComment(context.Context, string) error {
	return errors.New("not implemented")
}

func (m *MockQueries) DeleteNotification(context.Context, string) error {
	return errors.New("not implemented")
}

func (m *MockQueries) DeleteTask(context.Context, string) error {
	return errors.New("not implemented")
}

func (m *MockQueries) DeleteUser(context.Context, string) error {
	return errors.New("not implemented")
}

func (m *MockQueries) GetAttachment(context.Context, string) (db.Attachment, error) {
	return db.Attachment{}, errors.New("not implemented")
}

func (m *MockQueries) GetComment(context.Context, string) (db.Comment, error) {
	return db.Comment{}, errors.New("not implemented")
}

func (m *MockQueries) GetNotification(context.Context, string) (db.Notification, error) {
	return db.Notification{}, errors.New("not implemented")
}

func (m *MockQueries) GetTask(context.Context, string) (db.Task, error) {
	return db.Task{}, errors.New("not implemented")
}

func (m *MockQueries) GetUser(context.Context, string) (db.User, error) {
	return db.User{}, errors.New("not implemented")
}

func (m *MockQueries) GetUserByEmail(context.Context, string) (db.User, error) {
	return db.User{}, errors.New("not implemented")
}

func (m *MockQueries) GetUserByUsername(context.Context, string) (db.User, error) {
	return db.User{}, errors.New("not implemented")
}

func (m *MockQueries) ListActivityLogsByEntity(context.Context, db.ListActivityLogsByEntityParams) ([]db.ActivityLog, error) {
	return nil, errors.New("not implemented")
}

func (m *MockQueries) ListActivityLogsByUser(context.Context, db.ListActivityLogsByUserParams) ([]db.ActivityLog, error) {
	return nil, errors.New("not implemented")
}

func (m *MockQueries) ListAttachmentsByJob(context.Context, sql.NullString) ([]db.Attachment, error) {
	return nil, errors.New("not implemented")
}

func (m *MockQueries) ListAttachmentsByTask(context.Context, sql.NullString) ([]db.Attachment, error) {
	return nil, errors.New("not implemented")
}

func (m *MockQueries) ListCommentsByTask(context.Context, db.ListCommentsByTaskParams) ([]db.Comment, error) {
	return nil, errors.New("not implemented")
}

func (m *MockQueries) ListJobsByOwner(context.Context, db.ListJobsByOwnerParams) ([]db.Job, error) {
	return nil, errors.New("not implemented")
}

func (m *MockQueries) ListNotificationsByUser(context.Context, db.ListNotificationsByUserParams) ([]db.Notification, error) {
	return nil, errors.New("not implemented")
}

func (m *MockQueries) ListTasks(context.Context, db.ListTasksParams) ([]db.Task, error) {
	return nil, errors.New("not implemented")
}

func (m *MockQueries) ListTasksByJob(context.Context, db.ListTasksByJobParams) ([]db.Task, error) {
	return nil, errors.New("not implemented")
}

func (m *MockQueries) ListTasksByStatus(context.Context, db.ListTasksByStatusParams) ([]db.Task, error) {
	return nil, errors.New("not implemented")
}

func (m *MockQueries) ListTasksByUser(context.Context, db.ListTasksByUserParams) ([]db.Task, error) {
	return nil, errors.New("not implemented")
}

func (m *MockQueries) ListUnreadNotificationsByUser(context.Context, db.ListUnreadNotificationsByUserParams) ([]db.Notification, error) {
	return nil, errors.New("not implemented")
}

func (m *MockQueries) ListUsers(context.Context, db.ListUsersParams) ([]db.User, error) {
	return nil, errors.New("not implemented")
}

func (m *MockQueries) MarkAllNotificationsAsRead(context.Context, string) error {
	return errors.New("not implemented")
}

func (m *MockQueries) MarkNotificationAsRead(context.Context, string) (db.Notification, error) {
	return db.Notification{}, errors.New("not implemented")
}

func (m *MockQueries) UpdateComment(context.Context, db.UpdateCommentParams) (db.Comment, error) {
	return db.Comment{}, errors.New("not implemented")
}

func (m *MockQueries) UpdateTask(context.Context, db.UpdateTaskParams) (db.Task, error) {
	return db.Task{}, errors.New("not implemented")
}

func (m *MockQueries) UpdateUser(context.Context, db.UpdateUserParams) (db.User, error) {
	return db.User{}, errors.New("not implemented")
}

func (m *MockQueries) WithTx(tx *sql.Tx) *db.Queries {
	return &db.Queries{}
}

func TestJobService_CreateJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockJobQuerier(ctrl)
	svc := NewService(mockQuerier)
	ctx := context.Background()

	tests := []struct {
		name    string
		req     JobRequest
		ownerID string
		setup   func()
		wantErr bool
	}{
		{
			name: "valid job",
			req: JobRequest{
				Name:        "Test Job",
				Description: "Test Description",
				Status:      JobStatusPending,
				StartDate:   nil,
				EndDate:     nil,
			},
			ownerID: "owner123",
			setup: func() {
				mockQuerier.EXPECT().
					CreateJob(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, arg db.CreateJobParams) (db.Job, error) {
						return db.Job{
							ID:          arg.ID,
							Name:        arg.Name,
							Description: arg.Description,
							Status:      arg.Status,
							StartDate:   arg.StartDate,
							EndDate:     arg.EndDate,
							OwnerID:     arg.OwnerID,
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
						}, nil
					})
			},
			wantErr: false,
		},
		{
			name: "empty name",
			req: JobRequest{
				Name:   "",
				Status: JobStatusPending,
			},
			ownerID: "owner123",
			setup:   func() {}, // No mock expectations as validation should fail
			wantErr: true,
		},
		{
			name: "invalid status",
			req: JobRequest{
				Name:   "Test Job",
				Status: "invalid",
			},
			ownerID: "owner123",
			setup:   func() {}, // No mock expectations as validation should fail
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			resp, err := svc.CreateJob(ctx, tt.req, tt.ownerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("CreateJob() returned nil response")
					return
				}
				if resp.Name != tt.req.Name {
					t.Errorf("CreateJob() name = %v, want %v", resp.Name, tt.req.Name)
				}
				if resp.Status != tt.req.Status {
					t.Errorf("CreateJob() status = %v, want %v", resp.Status, tt.req.Status)
				}
			}
		})
	}
}

func TestJobService_GetJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockJobQuerier(ctrl)
	svc := NewService(mockQuerier)
	ctx := context.Background()

	testJob := db.Job{
		ID:          "test-job-id",
		Name:        "Test Job",
		Description: db.StringToNullString("Test Description"),
		Status:      string(JobStatusPending),
		OwnerID:     db.StringToNullString("owner123"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tests := []struct {
		name    string
		jobID   string
		setup   func()
		wantErr bool
	}{
		{
			name:  "existing job",
			jobID: testJob.ID,
			setup: func() {
				mockQuerier.EXPECT().
					GetJob(gomock.Any(), testJob.ID).
					Return(testJob, nil)
			},
			wantErr: false,
		},
		{
			name:  "non-existent job",
			jobID: "non-existent-id",
			setup: func() {
				mockQuerier.EXPECT().
					GetJob(gomock.Any(), "non-existent-id").
					Return(db.Job{}, db.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			resp, err := svc.GetJob(ctx, tt.jobID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("GetJob() returned nil response")
					return
				}
				if resp.ID != tt.jobID {
					t.Errorf("GetJob() ID = %v, want %v", resp.ID, tt.jobID)
				}
			}
		})
	}
}

func TestJobService_UpdateJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockJobQuerier(ctrl)
	svc := NewService(mockQuerier)
	ctx := context.Background()

	testJob := db.Job{
		ID:          "test-job-id",
		Name:        "Test Job",
		Description: db.StringToNullString("Test Description"),
		Status:      string(JobStatusPending),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tests := []struct {
		name    string
		jobID   string
		req     JobRequest
		setup   func()
		wantErr bool
	}{
		{
			name:  "valid update",
			jobID: testJob.ID,
			req: JobRequest{
				Name:        "Updated Job",
				Description: "Updated Description",
				Status:      JobStatusActive,
			},
			setup: func() {
				mockQuerier.EXPECT().
					UpdateJob(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, arg db.UpdateJobParams) (db.Job, error) {
						testJob.Name = arg.Name
						testJob.Description = arg.Description
						testJob.Status = arg.Status
						testJob.UpdatedAt = time.Now()
						return testJob, nil
					})
			},
			wantErr: false,
		},
		{
			name:  "non-existent job",
			jobID: "non-existent-id",
			req: JobRequest{
				Name:   "Updated Job",
				Status: JobStatusActive,
			},
			setup: func() {
				mockQuerier.EXPECT().
					UpdateJob(gomock.Any(), gomock.Any()).
					Return(db.Job{}, db.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name:  "invalid request",
			jobID: testJob.ID,
			req: JobRequest{
				Name:   "",
				Status: JobStatusActive,
			},
			setup:   func() {}, // No mock expectations as validation should fail
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			resp, err := svc.UpdateJob(ctx, tt.jobID, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("UpdateJob() returned nil response")
					return
				}
				if resp.Name != tt.req.Name {
					t.Errorf("UpdateJob() name = %v, want %v", resp.Name, tt.req.Name)
				}
				if resp.Status != tt.req.Status {
					t.Errorf("UpdateJob() status = %v, want %v", resp.Status, tt.req.Status)
				}
			}
		})
	}
}

func TestJobService_DeleteJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockJobQuerier(ctrl)
	svc := NewService(mockQuerier)
	ctx := context.Background()

	tests := []struct {
		name    string
		jobID   string
		setup   func()
		wantErr bool
	}{
		{
			name:  "existing job",
			jobID: "test-job-id",
			setup: func() {
				mockQuerier.EXPECT().
					DeleteJob(gomock.Any(), "test-job-id").
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "non-existent job",
			jobID: "non-existent-id",
			setup: func() {
				mockQuerier.EXPECT().
					DeleteJob(gomock.Any(), "non-existent-id").
					Return(db.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := svc.DeleteJob(ctx, tt.jobID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJobService_ListJobs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockJobQuerier(ctrl)
	svc := NewService(mockQuerier)
	ctx := context.Background()

	testJobs := []db.Job{
		{
			ID:          "job-1",
			Name:        "Job 1",
			Description: db.StringToNullString("Description 1"),
			Status:      string(JobStatusPending),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "job-2",
			Name:        "Job 2",
			Description: db.StringToNullString("Description 2"),
			Status:      string(JobStatusActive),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	tests := []struct {
		name    string
		params  JobListParams
		setup   func()
		want    int
		wantErr bool
	}{
		{
			name: "list all jobs",
			params: JobListParams{
				Page:     1,
				PageSize: 10,
			},
			setup: func() {
				mockQuerier.EXPECT().
					ListJobs(gomock.Any(), gomock.Any()).
					Return(testJobs, nil)
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "filter by status",
			params: JobListParams{
				Page:     1,
				PageSize: 10,
				Status:   JobStatusPending,
			},
			setup: func() {
				mockQuerier.EXPECT().
					ListJobs(gomock.Any(), gomock.Any()).
					Return(testJobs, nil)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "invalid page",
			params: JobListParams{
				Page:     0,
				PageSize: 10,
			},
			setup:   func() {}, // No mock expectations as validation should fail
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			resp, err := svc.ListJobs(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("ListJobs() returned nil response")
					return
				}
				if len(resp.Jobs) != tt.want {
					t.Errorf("ListJobs() returned %d jobs, want %d", len(resp.Jobs), tt.want)
				}
			}
		})
	}
}
