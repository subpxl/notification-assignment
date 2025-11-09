package handler

import (
	"encoding/json"
	"insider-assignment/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	Service *service.SenderService
}

func NewHandler(svc *service.SenderService) *Handler {
	return &Handler{Service: svc}
}

// Health godoc
// @Summary Check service health
// @Description Returns whether the sender service is currently running.
// @Tags health
// @Produce plain
// @Success 200 {string} string "Sender Service is running set to: true/false"
// @Router /health [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	status := h.Service.IsRunning()
	w.Write([]byte("Sender Service is running set to: " + strconv.FormatBool(status)))
}

// Control godoc
// @Summary Start or stop the sender service
// @Description Starts or stops the background sender process.
// @Tags control
// @Param action path string true "Action (start|stop)"
// @Produce plain
// @Success 200 {string} string "Sender Service Started or Stopped"
// @Failure 400 {string} string "invalid action"
// @Failure 405 {string} string "method not allowed"
// @Router /control/{action} [post]
func (h *Handler) Control(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	action := vars["action"]

	switch action {
	case "start":
		h.Service.Start()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Sender Service Started"))
	case "stop":
		h.Service.Stop()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Sender Service Stopped"))
	default:
		http.Error(w, "invalid action", http.StatusBadRequest)
	}
}

// GetSent godoc
// @Summary Get all sent messages
// @Description Retrieves a paginated list of messages that have already been sent.
// @Tags messages
// @Param limit query int false "Maximum number of messages to return (default 50)"
// @Param offset query int false "Offset for pagination"
// @Produce json
// @Success 200 {array} models.Message
// @Failure 500 {string} string "failed to get sent messages"
// @Router /sent-messages [get]
func (h *Handler) GetSent(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 50
	}
	msgs, err := h.Service.GetSent(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "failed to get sent messages", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}
