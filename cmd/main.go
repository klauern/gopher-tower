package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/klauern/gopher-tower/internal/api/jobs"
	"github.com/klauern/gopher-tower/internal/db"
	"github.com/klauern/gopher-tower/internal/db/migrate"
	"github.com/klauern/gopher-tower/internal/db/migrations"
	_ "modernc.org/sqlite"
)

// spaFileSystem handles serving the static frontend, falling back to index.html
type spaFileSystem struct {
	root http.FileSystem
}

func (fs *spaFileSystem) Open(name string) (http.File, error) {
	f, err := fs.root.Open(name)
	// If the file doesn't exist (likely a client-side route), serve index.html
	if os.IsNotExist(err) {
		// Ensure we don't loop indefinitely if index.html is missing
		if name == "index.html" {
			return nil, err
		}
		indexFile, indexErr := fs.root.Open("index.html")
		if indexErr != nil {
			// If index.html itself is missing, return the original error
			return nil, err
		}
		return indexFile, nil
	}
	return f, err
}

type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func main() {
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

	// Run database migrations
	if err := migrate.MigrateDB("gopher-tower.db", migrations.Files); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create database queries
	queries := db.New(dbConn)

	// Initialize HTTP server and routes
	router := chi.NewRouter()

	// Add standard middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger) // Keep or add logging middleware
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second)) // Example timeout

	// Initialize jobs service and handler
	jobService := jobs.NewService(queries)
	jobHandler := jobs.NewHandler(jobService)

	// Mount API routes under /api
	router.Route("/api", func(r chi.Router) {
		// Add CORS middleware specifically for API routes
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust in production if needed
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Add headers your frontend might send

				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusOK)
					return
				}

				next.ServeHTTP(w, r)
			})
		})

		jobHandler.RegisterRoutes(r)
		r.Get("/events", handleSSE) // Keep SSE handler under /api
	})

	// --- Serve Static Frontend Files ---
	// Determine the directory where the frontend build output lives
	// This assumes the Go binary runs from the project root. Adjust if needed.
	workDir, _ := os.Getwd()
	staticPath := filepath.Join(workDir, "frontend", "out") // Path to Next.js static export output

	// Check if the static directory exists
	if _, err := os.Stat(staticPath); os.IsNotExist(err) {
		log.Printf("Warning: Static file directory not found at %s. Frontend will not be served.", staticPath)
		// Optionally, you could choose to log.Fatal here if the frontend is mandatory
	} else {
		log.Printf("Serving static files from %s", staticPath)
		staticFilesDir := http.Dir(staticPath)
		fileServer := http.FileServer(&spaFileSystem{root: staticFilesDir})

		// Handle all other requests by serving static files or index.html
		router.Handle("/*", http.StripPrefix("/", fileServer))
	}
	// --- End Static File Serving ---

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
