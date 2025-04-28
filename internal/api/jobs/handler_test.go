package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"
)

type testContext struct {
	ctrl        *gomock.Controller
	mockService *MockService
	handler     *Handler
	router      chi.Router
}

func setupTest(t *testing.T) *testContext {
	ctrl := gomock.NewController(t)
	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	return &testContext{
		ctrl:        ctrl,
		mockService: mockService,
		handler:     handler,
		router:      router,
	}
}

func TestCreateJob(t *testing.T) {
	tests := []struct {
		name       string
		req        JobRequest
		setupAuth  func(r *http.Request)
		setupMock  func(*MockService)
		wantStatus int
		wantResp   *JobResponse
	}{
		{
			name: "successful creation",
			req: JobRequest{
				Name:        "Test Job",
				Description: "Test Description",
				Status:      JobStatusPending,
			},
			setupAuth: func(r *http.Request) {
				ctx := context.WithValue(r.Context(), UserIDKey, "test-user")
				*r = *r.WithContext(ctx)
			},
			setupMock: func(ms *MockService) {
				ms.EXPECT().
					CreateJob(gomock.Any(), gomock.Any(), "test-user").
					Return(&JobResponse{
						ID:          "test-id",
						Name:        "Test Job",
						Description: "Test Description",
						Status:      JobStatusPending,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid request",
			req: JobRequest{
				Name:   "",
				Status: JobStatusPending,
			},
			setupAuth: func(r *http.Request) {
				ctx := context.WithValue(r.Context(), UserIDKey, "test-user")
				*r = *r.WithContext(ctx)
			},
			setupMock: func(ms *MockService) {
				ms.EXPECT().
					CreateJob(gomock.Any(), gomock.Any(), "test-user").
					Return(nil, ErrInvalidJob)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "unauthorized",
			req: JobRequest{
				Name:   "Test Job",
				Status: JobStatusPending,
			},
			setupAuth: func(r *http.Request) {
				// No user_id in context
			},
			setupMock:  func(ms *MockService) {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := setupTest(t)
			defer tc.ctrl.Finish()

			tt.setupMock(tc.mockService)

			body, _ := json.Marshal(tt.req)
			req := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewReader(body))
			tt.setupAuth(req)

			w := httptest.NewRecorder()
			tc.router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("CreateJob() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var got JobResponse
				if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if got.Name != tt.req.Name {
					t.Errorf("CreateJob() name = %v, want %v", got.Name, tt.req.Name)
				}
			}
		})
	}
}

func TestGetJob(t *testing.T) {
	tests := []struct {
		name       string
		jobID      string
		setupMock  func(*MockService)
		wantStatus int
	}{
		{
			name:  "job found",
			jobID: "123e4567-e89b-12d3-a456-426614174000",
			setupMock: func(ms *MockService) {
				ms.EXPECT().
					GetJob(gomock.Any(), "123e4567-e89b-12d3-a456-426614174000").
					Return(&JobResponse{
						ID:          "123e4567-e89b-12d3-a456-426614174000",
						Name:        "Test Job",
						Description: "Test Description",
						Status:      JobStatusPending,
					}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "job not found",
			jobID: "123e4567-e89b-12d3-a456-426614174001",
			setupMock: func(ms *MockService) {
				ms.EXPECT().
					GetJob(gomock.Any(), "123e4567-e89b-12d3-a456-426614174001").
					Return(nil, ErrJobNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := setupTest(t)
			defer tc.ctrl.Finish()

			tt.setupMock(tc.mockService)

			req := httptest.NewRequest(http.MethodGet, "/jobs/"+tt.jobID, nil)
			w := httptest.NewRecorder()
			tc.router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetJob() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestListJobs(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		setupMock  func(*MockService)
		wantStatus int
		wantCount  int
	}{
		{
			name:  "list all jobs",
			query: "?page=1&page_size=10",
			setupMock: func(ms *MockService) {
				ms.EXPECT().
					ListJobs(gomock.Any(), gomock.Any()).
					Return(&JobListResponse{
						Jobs: []JobResponse{
							{ID: "1", Name: "Job 1"},
							{ID: "2", Name: "Job 2"},
						},
						TotalCount: 2,
						Page:       1,
						PageSize:   10,
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name:  "filter by status",
			query: "?page=1&page_size=10&status=pending",
			setupMock: func(ms *MockService) {
				ms.EXPECT().
					ListJobs(gomock.Any(), gomock.Any()).
					Return(&JobListResponse{
						Jobs: []JobResponse{
							{ID: "1", Name: "Job 1", Status: JobStatusPending},
						},
						TotalCount: 1,
						Page:       1,
						PageSize:   10,
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name:       "invalid page",
			query:      "?page=invalid",
			setupMock:  func(ms *MockService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := setupTest(t)
			defer tc.ctrl.Finish()

			tt.setupMock(tc.mockService)

			req := httptest.NewRequest(http.MethodGet, "/jobs"+tt.query, nil)
			w := httptest.NewRecorder()
			tc.router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ListJobs() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var got JobListResponse
				if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if len(got.Jobs) != tt.wantCount {
					t.Errorf("ListJobs() count = %v, want %v", len(got.Jobs), tt.wantCount)
				}
			}
		})
	}
}

func TestUpdateJob(t *testing.T) {
	tests := []struct {
		name       string
		jobID      string
		req        JobRequest
		setupMock  func(*MockService)
		wantStatus int
	}{
		{
			name:  "successful update",
			jobID: "123e4567-e89b-12d3-a456-426614174000",
			req: JobRequest{
				Name:        "Updated Job",
				Description: "Updated Description",
				Status:      JobStatusActive,
			},
			setupMock: func(ms *MockService) {
				ms.EXPECT().
					UpdateJob(gomock.Any(), "123e4567-e89b-12d3-a456-426614174000", gomock.Any()).
					Return(&JobResponse{
						ID:          "123e4567-e89b-12d3-a456-426614174000",
						Name:        "Updated Job",
						Description: "Updated Description",
						Status:      JobStatusActive,
					}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "job not found",
			jobID: "123e4567-e89b-12d3-a456-426614174001",
			req: JobRequest{
				Name:   "Updated Job",
				Status: JobStatusActive,
			},
			setupMock: func(ms *MockService) {
				ms.EXPECT().
					UpdateJob(gomock.Any(), "123e4567-e89b-12d3-a456-426614174001", gomock.Any()).
					Return(nil, ErrJobNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := setupTest(t)
			defer tc.ctrl.Finish()

			tt.setupMock(tc.mockService)

			body, _ := json.Marshal(tt.req)
			req := httptest.NewRequest(http.MethodPut, "/jobs/"+tt.jobID, bytes.NewReader(body))
			w := httptest.NewRecorder()
			tc.router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("UpdateJob() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestDeleteJob(t *testing.T) {
	tests := []struct {
		name       string
		jobID      string
		setupMock  func(*MockService)
		wantStatus int
	}{
		{
			name:  "successful deletion",
			jobID: "123e4567-e89b-12d3-a456-426614174000",
			setupMock: func(ms *MockService) {
				ms.EXPECT().
					DeleteJob(gomock.Any(), "123e4567-e89b-12d3-a456-426614174000").
					Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:  "job not found",
			jobID: "123e4567-e89b-12d3-a456-426614174001",
			setupMock: func(ms *MockService) {
				ms.EXPECT().
					DeleteJob(gomock.Any(), "123e4567-e89b-12d3-a456-426614174001").
					Return(ErrJobNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := setupTest(t)
			defer tc.ctrl.Finish()

			tt.setupMock(tc.mockService)

			req := httptest.NewRequest(http.MethodDelete, "/jobs/"+tt.jobID, nil)
			w := httptest.NewRecorder()
			tc.router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("DeleteJob() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}
