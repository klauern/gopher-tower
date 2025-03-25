package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeDatabase(t *testing.T) {
	// Create an in-memory SQLite database for testing
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	// Test database initialization
	err = initializeDatabase(db)
	require.NoError(t, err)

	// Verify that tables were created by trying to insert and query data
	_, err = db.Exec(`INSERT INTO jobs (name, description, status) VALUES (?, ?, ?)`,
		"Test Job", "Test Description", "pending")
	require.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM jobs").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestHandleSSE(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		validateBody   bool
	}{
		{
			name:           "GET request",
			method:         "GET",
			expectedStatus: http.StatusOK,
			validateBody:   true,
		},
		{
			name:           "POST request",
			method:         "POST",
			expectedStatus: http.StatusMethodNotAllowed,
			validateBody:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/events", nil)
			w := httptest.NewRecorder()

			handleSSE(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.validateBody {
				// Verify headers
				assert.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))
				assert.Equal(t, "no-cache", resp.Header.Get("Cache-Control"))
				assert.Equal(t, "keep-alive", resp.Header.Get("Connection"))
				assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))

				// Read and validate events
				scanner := bufio.NewScanner(resp.Body)
				var events []Event
				for scanner.Scan() {
					line := scanner.Text()
					if strings.HasPrefix(line, "data: ") {
						var event Event
						err := json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &event)
						require.NoError(t, err)
						events = append(events, event)

						// Break after receiving first event to avoid waiting for more
						break
					}
				}

				// Verify at least one event was received
				require.NotEmpty(t, events)
				event := events[0]

				// Verify event structure based on type
				switch event.Type {
				case "counter":
					payload, ok := event.Payload.(map[string]interface{})
					require.True(t, ok)
					assert.Contains(t, payload, "count")
				case "time":
					payload, ok := event.Payload.(map[string]interface{})
					require.True(t, ok)
					assert.Contains(t, payload, "timestamp")
					_, err := time.Parse(time.RFC3339, payload["timestamp"].(string))
					assert.NoError(t, err)
				case "random":
					payload, ok := event.Payload.(map[string]interface{})
					require.True(t, ok)
					assert.Contains(t, payload, "value")
					assert.Contains(t, payload["value"].(string), "Test message")
				default:
					t.Errorf("Unexpected event type: %s", event.Type)
				}
			}
		})
	}
}

func TestEvent_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		event    Event
		expected string
	}{
		{
			name: "counter event",
			event: Event{
				Type:    "counter",
				Payload: map[string]interface{}{"count": 1},
			},
			expected: `{"type":"counter","payload":{"count":1}}`,
		},
		{
			name: "time event",
			event: Event{
				Type:    "time",
				Payload: map[string]interface{}{"timestamp": "2024-03-24T00:00:00Z"},
			},
			expected: `{"type":"time","payload":{"timestamp":"2024-03-24T00:00:00Z"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.event)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(data))
		})
	}
}
