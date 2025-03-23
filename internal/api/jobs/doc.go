/*
Package jobs provides a complete job management service for the gopher-tower application.

Job Lifecycle:

	Pending -> Active -> Complete/Failed

Core Components:

  - Handler: HTTP endpoints for job operations
  - Service: Business logic and job management
  - Models: Data structures and validation
  - JobQuerier: Database interface for job persistence

Example Usage:

	// Create a new job service
	jobService := jobs.NewService(dbQueries)

	// Create a new job handler
	jobHandler := jobs.NewHandler(jobService)

	// Register routes (using chi router)
	router.Mount("/api", jobHandler.RegisterRoutes)

API Endpoints:

	POST   /jobs       - Create a new job
	GET    /jobs/{id}  - Get job details
	PUT    /jobs/{id}  - Update job status/details
	DELETE /jobs/{id}  - Delete a job
	GET    /jobs       - List jobs with pagination and filters

Request/Response Examples:

Create Job:

	POST /jobs
	{
		"name": "Example Job",
		"description": "Job description",
		"status": "pending",
		"start_date": "2024-03-22T00:00:00Z",
		"end_date": "2024-03-23T00:00:00Z"
	}

List Jobs:

	GET /jobs?page=1&page_size=10&status=active

	Response:
	{
		"jobs": [...],
		"total_count": 42,
		"page": 1,
		"page_size": 10
	}

Error Handling:

The package uses standard HTTP status codes:
  - 200: Success
  - 201: Created
  - 400: Bad Request (validation errors)
  - 404: Not Found
  - 500: Internal Server Error

Custom errors:
  - ErrJobNotFound: Job doesn't exist
  - ErrInvalidJob: Invalid job data

Testing:

The package uses uber-go/mock for mocking the database layer:

	task go:generate  # Generate mocks
	task go:test      # Run tests

For more details on testing, see the test files and examples.
*/
package jobs
