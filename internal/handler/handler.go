package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"my-calendar/internal/calendar"
)

type Handler struct {
	storage *calendar.Calendar
}

func NewHandler(calendar calendar.Calendar) *Handler {
	return &Handler {
		storage: &calendar,
	}
}

func (h *Handler) GetDaily(w http.ResponseWriter, r *http.Request) {

	// parse query params
	// ?user_id=1&date=2023-12-31
	date := r.URL.Query().Get("date")
	userId := r.URL.Query().Get("user_id")

    events, err := h.storage.GetDaily(userId, date)
    if err != nil {
    	writeJson(w, http.StatusBadRequest, err.Error())
    }
    
	writeJson(w, http.StatusOK, events)
}

func (h *Handler) GetMonthly(w http.ResponseWriter, r *http.Request) {

	date := r.URL.Query().Get("date")
	userId := r.URL.Query().Get("user_id")
	
	targetDate, err := time.Parse("2001-01-01", date)
	if err != nil {
		writeJson(w, http.StatusBadRequest, "date must be in YYYY-MM-DD format")
        return
	}

    events, err := h.storage.GetMonthly(userId, date)
    if err != nil {
    	writeJson(w, http.StatusBadRequest, err.Error())
    }

    writeJson(w, http.StatusOK, events)
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {

	
}


// helper for writing json
func writeJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

