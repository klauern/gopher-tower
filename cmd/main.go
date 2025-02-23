package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func main() {
	mux := http.NewServeMux()

	// Handle CORS preflight
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for all responses
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only allow GET requests for SSE
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handleSSE(w, r)
	})

	port := 8080
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
		// Add timeouts to prevent connection issues
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 2 * time.Hour, // Long timeout for SSE
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting SSE server on port %d...\n", port)
	log.Fatal(server.ListenAndServe())
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	log.Printf("New client connected from %s\n", r.RemoteAddr)

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // Disable buffering in Nginx if present

	// Ensure the connection supports flushing
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Send an initial comment to establish connection
	fmt.Fprintf(w, ": ping\n\n")
	flusher.Flush()

	// Create a channel to notify of client disconnect
	notify := r.Context().Done()

	// Send events every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	counter := 0
	for {
		select {
		case <-notify:
			log.Printf("Client %s disconnected\n", r.RemoteAddr)
			return
		case <-ticker.C:
			counter++
			var event Event

			// Alternate between different test events
			switch counter % 3 {
			case 0:
				event = Event{
					Type: "counter",
					Payload: map[string]interface{}{
						"count": counter,
					},
				}
			case 1:
				event = Event{
					Type: "time",
					Payload: map[string]interface{}{
						"timestamp": time.Now().Format(time.RFC3339),
					},
				}
			case 2:
				event = Event{
					Type: "random",
					Payload: map[string]interface{}{
						"value": fmt.Sprintf("Test message %d", counter),
					},
				}
			}

			// Convert event to JSON
			data, err := json.Marshal(event)
			if err != nil {
				log.Printf("Error marshaling event: %v\n", err)
				continue
			}

			// Format as SSE data field (ensuring proper line endings)
			sseData := fmt.Sprintf("data: %s\n\n", string(data))

			// Write and flush
			if _, err := fmt.Fprint(w, sseData); err != nil {
				log.Printf("Error writing to client %s: %v\n", r.RemoteAddr, err)
				return
			}
			flusher.Flush()

			log.Printf("Sent event to %s: %s", r.RemoteAddr, strings.TrimSpace(string(data)))
		}
	}
}
