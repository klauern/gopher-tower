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
		// Create a root command with server URL
		root := &cli.Command{
			Name: "test",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "server",
					Value: server.URL,
				},
			},
		}

		// Add the jobs command as a subcommand
		jobsCmd := JobsCommand()
		root.Commands = append(root.Commands, jobsCmd)

		// Find the list subcommand
		var listCmd *cli.Command
		for _, c := range jobsCmd.Commands {
			if c.Name == "list" {
				listCmd = c
				break
			}
		}
		require.NotNil(t, listCmd)

		// Set up command context
		ctx := context.Background()

		// Test with default flags
		err := listCmd.Action(ctx, listCmd)
		assert.NoError(t, err)
	})

	t.Run("status command", func(t *testing.T) {
		// Create a root command with server URL
		root := &cli.Command{
			Name: "test",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "server",
					Value: server.URL,
				},
			},
		}

		// Add the jobs command as a subcommand
		jobsCmd := JobsCommand()
		root.Commands = append(root.Commands, jobsCmd)

		// Find the status subcommand
		var statusCmd *cli.Command
		for _, c := range jobsCmd.Commands {
			if c.Name == "status" {
				statusCmd = c
				break
			}
		}
		require.NotNil(t, statusCmd)

		// Set up command context
		ctx := context.Background()

		// Set job ID flag
		err := statusCmd.Set("id", "1")
		require.NoError(t, err)

		// Test command execution
		err = statusCmd.Action(ctx, statusCmd)
		assert.NoError(t, err)
	})

	t.Run("invalid server URL", func(t *testing.T) {
		// Create a root command with invalid server URL
		root := &cli.Command{
			Name: "test",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "server",
					Value: "http://invalid-server",
				},
			},
		}

		// Add the jobs command as a subcommand
		jobsCmd := JobsCommand()
		root.Commands = append(root.Commands, jobsCmd)

		// Find the list subcommand
		var listCmd *cli.Command
		for _, c := range jobsCmd.Commands {
			if c.Name == "list" {
				listCmd = c
				break
			}
		}
		require.NotNil(t, listCmd)

		// Set up command context
		ctx := context.Background()

		// Test command execution
		err := listCmd.Action(ctx, listCmd)
		assert.Error(t, err)
	})

	t.Run("server error response", func(t *testing.T) {
		errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}))
		defer errorServer.Close()

		// Create a root command with error server URL
		root := &cli.Command{
			Name: "test",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "server",
					Value: errorServer.URL,
				},
			},
		}

		// Add the jobs command as a subcommand
		jobsCmd := JobsCommand()
		root.Commands = append(root.Commands, jobsCmd)

		// Find the list subcommand
		var listCmd *cli.Command
		for _, c := range jobsCmd.Commands {
			if c.Name == "list" {
				listCmd = c
				break
			}
		}
		require.NotNil(t, listCmd)

		// Set up command context
		ctx := context.Background()

		// Test command execution
		err := listCmd.Action(ctx, listCmd)
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
