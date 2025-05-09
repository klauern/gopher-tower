package main

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			// Create a context with timeout and test flag
			ctx := context.WithValue(context.Background(), "test", true)
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			req := httptest.NewRequest(tt.method, "/api/events", nil)
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			// Create a done channel to signal test completion
			done := make(chan struct{})
			go func() {
				handleSSE(w, req)
				close(done)
			}()

			// Wait for either context timeout or test completion
			select {
			case <-ctx.Done():
				t.Fatal("Test timed out")
			case <-done:
				// Test completed normally
			}

			resp := w.Result()
			defer resp.Body.Close()
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
