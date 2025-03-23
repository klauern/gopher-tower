package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/urfave/cli/v3"
)

// JobStatus represents the current state of a job
type JobStatus string

const (
	JobStatusPending  JobStatus = "pending"
	JobStatusActive   JobStatus = "active"
	JobStatusComplete JobStatus = "complete"
	JobStatusFailed   JobStatus = "failed"
)

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

// JobListResponse represents the response from listing jobs
type JobListResponse struct {
	Jobs       []JobResponse `json:"jobs"`
	TotalCount int           `json:"total_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
}

// JobsCommand returns the jobs subcommand
func JobsCommand() *cli.Command {
	return &cli.Command{
		Name:    "jobs",
		Usage:   "Manage jobs in Gopher Tower",
		Aliases: []string{"j"},
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all jobs",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "page",
						Usage:   "Page number",
						Value:   1,
						Aliases: []string{"p"},
					},
					&cli.IntFlag{
						Name:    "page-size",
						Usage:   "Number of items per page",
						Value:   10,
						Aliases: []string{"s"},
					},
					&cli.StringFlag{
						Name:  "status",
						Usage: "Filter by status (pending, active, complete, failed)",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					serverURL := cmd.Root().String("server")
					page := cmd.Int("page")
					pageSize := cmd.Int("page-size")
					status := cmd.String("status")

					url := fmt.Sprintf("%s/api/jobs?page=%d&page_size=%d", serverURL, page, pageSize)
					if status != "" {
						url += fmt.Sprintf("&status=%s", status)
					}

					resp, err := http.Get(url)
					if err != nil {
						return fmt.Errorf("failed to fetch jobs: %w", err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						return fmt.Errorf("server returned error: %s", resp.Status)
					}

					var listResp JobListResponse
					if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
						return fmt.Errorf("failed to decode response: %w", err)
					}

					// Handle empty response
					if len(listResp.Jobs) == 0 {
						fmt.Println("No jobs found")
						return nil
					}

					// Pretty print jobs
					fmt.Printf("Total jobs: %d (Page %d of %d)\n\n",
						listResp.TotalCount,
						listResp.Page,
						(listResp.TotalCount+listResp.PageSize-1)/listResp.PageSize,
					)

					for _, job := range listResp.Jobs {
						fmt.Printf("ID: %s\n", job.ID)
						fmt.Printf("Name: %s\n", job.Name)
						fmt.Printf("Description: %s\n", job.Description)
						fmt.Printf("Status: %s\n", job.Status)
						if job.StartDate != nil {
							fmt.Printf("Start Date: %s\n", job.StartDate.Format(time.RFC3339))
						}
						if job.EndDate != nil {
							fmt.Printf("End Date: %s\n", job.EndDate.Format(time.RFC3339))
						}
						fmt.Printf("Created: %s\n", job.CreatedAt.Format(time.RFC3339))
						fmt.Printf("Updated: %s\n", job.UpdatedAt.Format(time.RFC3339))
						if job.OwnerID != "" {
							fmt.Printf("Owner: %s\n", job.OwnerID)
						}
						fmt.Println("---")
					}

					return nil
				},
			},
			{
				Name:  "status",
				Usage: "Get status of a specific job",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "Job ID to check status for",
						Required: true,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					serverURL := cmd.Root().String("server")
					jobID := cmd.String("id")
					url := fmt.Sprintf("%s/api/jobs/%s", serverURL, jobID)

					resp, err := http.Get(url)
					if err != nil {
						return fmt.Errorf("failed to fetch job status: %w", err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						return fmt.Errorf("server returned error: %s", resp.Status)
					}

					var job JobResponse
					if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
						return fmt.Errorf("failed to decode response: %w", err)
					}

					fmt.Printf("ID: %s\n", job.ID)
					fmt.Printf("Name: %s\n", job.Name)
					fmt.Printf("Description: %s\n", job.Description)
					fmt.Printf("Status: %s\n", job.Status)
					if job.StartDate != nil {
						fmt.Printf("Start Date: %s\n", job.StartDate.Format(time.RFC3339))
					}
					if job.EndDate != nil {
						fmt.Printf("End Date: %s\n", job.EndDate.Format(time.RFC3339))
					}
					fmt.Printf("Created: %s\n", job.CreatedAt.Format(time.RFC3339))
					fmt.Printf("Updated: %s\n", job.UpdatedAt.Format(time.RFC3339))
					if job.OwnerID != "" {
						fmt.Printf("Owner: %s\n", job.OwnerID)
					}

					return nil
				},
			},
		},
	}
}
