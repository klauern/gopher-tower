package commands

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestJobsCommand(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/jobs":
			// Test list jobs endpoint
			listResp := JobListResponse{
				Jobs: []JobResponse{
					{
						ID:          "1",
						Name:        "Test Job",
						Description: "Test Description",
						Status:      JobStatusPending,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				},
				TotalCount: 1,
				Page:       1,
				PageSize:   10,
			}
			json.NewEncoder(w).Encode(listResp)

		case "/api/jobs/1":
			// Test get job endpoint
			job := JobResponse{
				ID:          "1",
				Name:        "Test Job",
				Description: "Test Description",
				Status:      JobStatusPending,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			json.NewEncoder(w).Encode(job)

		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	t.Run("list command", func(t *testing.T) {
		// Create root command
		app := &cli.Command{
			Name: "test",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "server",
					Value: server.URL,
				},
			},
		}

		// Add jobs command as subcommand
		jobsCmd := JobsCommand()
		app.Commands = append(app.Commands, jobsCmd)

		// Find the list subcommand
		var listCmd *cli.Command
		for _, c := range jobsCmd.Commands {
			if c.Name == "list" {
				listCmd = c
				break
			}
		}
		require.NotNil(t, listCmd)

		// Test with default flags
		ctx := context.Background()
		err := app.Run(ctx, []string{"test", "--server", server.URL, "jobs", "list"})
		assert.NoError(t, err)
	})

	t.Run("status command", func(t *testing.T) {
		// Create root command
		app := &cli.Command{
			Name: "test",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "server",
					Value: server.URL,
				},
			},
		}

		// Add jobs command as subcommand
		jobsCmd := JobsCommand()
		app.Commands = append(app.Commands, jobsCmd)

		// Find the status subcommand
		var statusCmd *cli.Command
		for _, c := range jobsCmd.Commands {
			if c.Name == "status" {
				statusCmd = c
				break
			}
		}
		require.NotNil(t, statusCmd)

		// Test with job ID
		ctx := context.Background()
		err := app.Run(ctx, []string{"test", "--server", server.URL, "jobs", "status", "--id", "1"})
		assert.NoError(t, err)
	})

	t.Run("invalid server URL", func(t *testing.T) {
		// Create root command
		app := &cli.Command{
			Name: "test",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "server",
					Value: "http://invalid-server",
				},
			},
		}

		// Add jobs command as subcommand
		jobsCmd := JobsCommand()
		app.Commands = append(app.Commands, jobsCmd)

		// Find the list subcommand
		var listCmd *cli.Command
		for _, c := range jobsCmd.Commands {
			if c.Name == "list" {
				listCmd = c
				break
			}
		}
		require.NotNil(t, listCmd)

		// Test with invalid server
		ctx := context.Background()
		err := app.Run(ctx, []string{"test", "--server", "http://invalid-server", "jobs", "list"})
		assert.Error(t, err)
	})

	t.Run("server error response", func(t *testing.T) {
		errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}))
		defer errorServer.Close()

		// Create root command
		app := &cli.Command{
			Name: "test",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "server",
					Value: errorServer.URL,
				},
			},
		}

		// Add jobs command as subcommand
		jobsCmd := JobsCommand()
		app.Commands = append(app.Commands, jobsCmd)

		// Find the list subcommand
		var listCmd *cli.Command
		for _, c := range jobsCmd.Commands {
			if c.Name == "list" {
				listCmd = c
				break
			}
		}
		require.NotNil(t, listCmd)

		// Test with error server
		ctx := context.Background()
		err := app.Run(ctx, []string{"test", "--server", errorServer.URL, "jobs", "list"})
		assert.Error(t, err)
	})
}

func TestJobStatus_String(t *testing.T) {
	tests := []struct {
		status   JobStatus
		expected string
	}{
		{JobStatusPending, "pending"},
		{JobStatusActive, "active"},
		{JobStatusComplete, "complete"},
		{JobStatusFailed, "failed"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}
