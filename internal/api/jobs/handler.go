package jobs

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Handler handles HTTP requests for jobs
type Handler struct {
	service Service
}

// NewHandler creates a new job handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers the job routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/jobs", h.CreateJob)
	r.Get("/jobs/{id}", h.GetJob)
	r.Put("/jobs/{id}", h.UpdateJob)
	r.Delete("/jobs/{id}", h.DeleteJob)
	r.Get("/jobs", h.ListJobs)
}

// CreateJob handles job creation requests
func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req JobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get owner ID from context (assuming it's set by auth middleware)
	ownerID := r.Context().Value("user_id")
	if ownerID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	resp, err := h.service.CreateJob(r.Context(), req, ownerID.(string))
	if err != nil {
		switch err {
		case ErrInvalidJob:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetJob handles job retrieval requests
func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing job ID", http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetJob(r.Context(), id)
	if err != nil {
		switch err {
		case ErrJobNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UpdateJob handles job update requests
func (h *Handler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing job ID", http.StatusBadRequest)
		return
	}

	var req JobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.service.UpdateJob(r.Context(), id, req)
	if err != nil {
		switch err {
		case ErrJobNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case ErrInvalidJob:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteJob handles job deletion requests
func (h *Handler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing job ID", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteJob(r.Context(), id)
	if err != nil {
		switch err {
		case ErrJobNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListJobs handles job listing requests
func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	params := JobListParams{
		Page:     1,
		PageSize: 10,
	}

	if page := r.URL.Query().Get("page"); page != "" {
		p, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			return
		}
		params.Page = p
	}

	if pageSize := r.URL.Query().Get("page_size"); pageSize != "" {
		ps, err := strconv.Atoi(pageSize)
		if err != nil {
			http.Error(w, "Invalid page_size parameter", http.StatusBadRequest)
			return
		}
		params.PageSize = ps
	}

	if status := r.URL.Query().Get("status"); status != "" {
		params.Status = JobStatus(status)
	}

	resp, err := h.service.ListJobs(r.Context(), params)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
