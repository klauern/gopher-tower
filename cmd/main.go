package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/klauern/gopher-tower/internal/api/jobs"
	"github.com/klauern/gopher-tower/internal/db"
	"github.com/klauern/gopher-tower/internal/static"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func initializeDatabase(db *sql.DB) error {
	// Execute schema
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

func main() {
	// Create a file system with the frontend directory as the root
	fsys, err := fs.Sub(static.Files, "frontend")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
	}

	// Create a file server for static assets
	fileServer := http.FileServer(http.FS(fsys))

	// Initialize SQLite database
	dbConn, err := sql.Open("sqlite", "gopher-tower.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer dbConn.Close()

	// Test database connection
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize database schema
	if err := initializeDatabase(dbConn); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create database queries
	queries := db.New(dbConn)

	// Initialize HTTP server and routes
	router := chi.NewRouter()

	// Initialize jobs service and handler
	jobService := jobs.NewService(queries)
	jobHandler := jobs.NewHandler(jobService)

	// Mount API routes
	apiRouter := chi.NewRouter()
	// Add CORS middleware for API routes
	apiRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})
	jobHandler.RegisterRoutes(apiRouter)

	// Handle SSE endpoint under /api
	apiRouter.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests for SSE
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handleSSE(w, r)
	})

	router.Mount("/api", apiRouter)

	// For everything else, use the file server
	router.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request for: %s", r.URL.Path)
		fileServer.ServeHTTP(w, r)
	}))

	port := 8080
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 2 * time.Hour, // Long timeout for SSE
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting server on port %d...\n", port)
	log.Fatal(server.ListenAndServe())
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests for SSE
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("New client connected from %s\n", r.RemoteAddr)

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("X-Accel-Buffering", "no") // Disable buffering in Nginx if present

	// Ensure the connection supports flushing
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Send an initial comment to establish connection
	fmt.Fprintf(w, "retry: 1000\n") // Tell client to retry every 1 second
	fmt.Fprintf(w, ": ping\n\n")
	flusher.Flush()

	// Create a channel to notify of client disconnect
	ctx := r.Context()

	// Send events every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Keep track of connection state
	connected := true
	defer func() {
		if connected {
			log.Printf("Client %s disconnected\n", r.RemoteAddr)
		}
	}()

	counter := 0
	for connected {
		select {
		case <-ctx.Done():
			connected = false
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
				connected = false
				return
			}
			flusher.Flush()

			log.Printf("Sent event to %s: %s", r.RemoteAddr, strings.TrimSpace(string(data)))

			// In test mode, exit after sending one event
			if _, isTest := ctx.Value("test").(bool); isTest {
				return
			}
		}
	}
}
