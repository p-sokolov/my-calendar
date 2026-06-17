package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"my-calendar/internal/calendar"
	e "my-calendar/internal/error"
)

type Handler struct {
	storage *calendar.Calendar
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewHandler(calendar *calendar.Calendar) *Handler {
	return &Handler {
		storage: calendar,
	}
}

// GetEvents godoc
//
// @Summary Получить события
// @Tags events
// @Produce json
// @Param user_id query int false "User ID"
// @Param date query string false "Date"
// @Param period query string false "day|week|month"
// @Success 200 {array} calendar.Event
// @Failure 400 {object} handler.ErrorResponse
// @Router /events [get]
func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
    date := r.URL.Query().Get("date")
    period := r.URL.Query().Get("period")

    events, err := h.storage.GetEvents(userID, date, period)
    if errors.Is(err, e.ErrEventNotFound) {
   		writeJson(w, http.StatusBadRequest, ErrorResponse{ Error: err.Error() })
        return
    }

    writeJson(w, http.StatusOK, events)
}

func (h *Handler) GetDaily(w http.ResponseWriter, r *http.Request) {
	// parse query params
	// ?user_id=1&date=2023-12-31
	date := r.URL.Query().Get("date")
	userId := r.URL.Query().Get("user_id")

    events, err := h.storage.GetDaily(userId, date)
    if err != nil {
    	writeJson(w, http.StatusBadRequest, ErrorResponse{ Error: err.Error() })
     	return
    }
    
	writeJson(w, http.StatusOK, events)
}

func (h *Handler) GetWeekly(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	userId := r.URL.Query().Get("user_id")

    events, err := h.storage.GetWeekly(userId, date)
    if err != nil {
   		writeJson(w, http.StatusBadRequest, ErrorResponse{ Error: err.Error() })
     	return
    }
    
	writeJson(w, http.StatusOK, events)
}

func (h *Handler) GetMonthly(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	userId := r.URL.Query().Get("user_id")
	
    events, err := h.storage.GetMonthly(userId, date)
    if err != nil {
   		writeJson(w, http.StatusBadRequest, ErrorResponse{ Error: err.Error() })
    	return
    }

    writeJson(w, http.StatusOK, events)
}

// CreateEvent godoc
//
// @Summary Создать событие
// @Description Создает новое событие календаря
// @Tags events
// @Accept json
// @Produce json
// @Param event body calendar.Event true "Event"
// @Success 200 {object} calendar.Event
// @Failure 400 {object} handler.ErrorResponse
// @Router /events [post]
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event calendar.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        writeJson(w, http.StatusBadRequest, ErrorResponse{ Error: err.Error() })
        return
    }

    created, err := h.storage.CreateEvent(event)
    if err != nil {
   		writeJson(w, http.StatusBadRequest, ErrorResponse{ Error: err.Error() })
     	return
    }

    writeJson(w, http.StatusOK, created)
}

// UpdateEvent godoc
//
// @Summary Обновить событие
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param event body calendar.Event true "Event"
// @Success 200 {object} calendar.Event
// @Failure 400 {object} handler.ErrorResponse
// @Router /events/{id} [put]
func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	
	id, err := strconv.Atoi(idStr)
    if err != nil {
        writeJson(w, http.StatusBadRequest, ErrorResponse{ Error: err.Error() })
        return
    }
    
    var event calendar.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        writeJson(w, http.StatusBadRequest, ErrorResponse{ Error: err.Error() })
        return
    }

    event.EventID = id
    updated, err := h.storage.UpdateEvent(event)
    
    if err != nil {
    	writeJson(w, http.StatusBadRequest, ErrorResponse{ Error: err.Error() })
     	return
    }

    writeJson(w, http.StatusOK, updated)
}

// DeleteEvent godoc
//
// @Summary Удалить событие
// @Tags events
// @Produce json
// @Param id path int true "Event ID"
// @Success 204
// @Failure 400 {object} handler.ErrorResponse
// @Router /events/{id} [delete]
func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	
	idStr := mux.Vars(r)["id"]
	
	id, err := strconv.Atoi(idStr)
    if err != nil {
        writeJson(w, http.StatusBadRequest, ErrorResponse{ Error: err.Error() })
        return
    }

    if err := h.storage.DeleteEvent(id); err != nil {
    	writeJson(w, http.StatusBadRequest, ErrorResponse{ Error: err.Error() })
     	return
    }

    w.WriteHeader(http.StatusNoContent)
}

// helper for writing json
func writeJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

