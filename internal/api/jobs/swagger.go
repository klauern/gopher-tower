package jobs

// @title Gopher Tower Jobs API
// @version 1.0
// @description Job management API for Gopher Tower
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// CreateJob godoc
// @Summary Create a new job
// @Description Create a new job with the provided details
// @Tags jobs
// @Accept json
// @Produce json
// @Param job body JobRequest true "Job details"
// @Success 200 {object} JobResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /jobs [post]
//
//nolint:unused // This function exists only for swagger documentation
func (h *Handler) swaggerCreateJob() {}

// GetJob godoc
// @Summary Get job details
// @Description Get details of a specific job by ID
// @Tags jobs
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} JobResponse
// @Failure 404 {string} string "Job not found"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /jobs/{id} [get]
//
//nolint:unused // This function exists only for swagger documentation
func (h *Handler) swaggerGetJob() {}

// UpdateJob godoc
// @Summary Update job details
// @Description Update details of a specific job
// @Tags jobs
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Param job body JobRequest true "Job details"
// @Success 200 {object} JobResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Job not found"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /jobs/{id} [put]
//
//nolint:unused // This function exists only for swagger documentation
func (h *Handler) swaggerUpdateJob() {}

// DeleteJob godoc
// @Summary Delete a job
// @Description Delete a specific job by ID
// @Tags jobs
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "Job not found"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /jobs/{id} [delete]
//
//nolint:unused // This function exists only for swagger documentation
func (h *Handler) swaggerDeleteJob() {}

// ListJobs godoc
// @Summary List jobs
// @Description Get a paginated list of jobs with optional filters
// @Tags jobs
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 10)"
// @Param status query string false "Filter by status (pending, active, complete, failed)"
// @Success 200 {object} JobListResponse
// @Failure 400 {string} string "Invalid parameters"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /jobs [get]
//
//nolint:unused // This function exists only for swagger documentation
func (h *Handler) swaggerListJobs() {}
