package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	err := h.db.Ping(r.Context())
	if err != nil {
		response := struct {
			Status   string `json:"status"`
			Database string `json:"database"`
		}{
			Status:   "degraded",
			Database: "error",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}

		return
	}

	response := struct {
		Status   string `json:"status"`
		Database string `json:"database"`
	}{
		Status:   "ok",
		Database: "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func (h *HealthHandler) HandleLive(w http.ResponseWriter, r *http.Request) {
}

func (h *HealthHandler) HandleReady(w http.ResponseWriter, r *http.Request) {}
