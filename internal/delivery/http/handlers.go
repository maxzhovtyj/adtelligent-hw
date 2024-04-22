package delivery

import (
	"encoding/json"
	"github.com/maxzhovtyj/adtelligent-hw/internal/services"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	services services.Services
}

func New(services services.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /sources/{id}/campaigns", h.sourceCampaignsHandler)

	return mux
}

func (h *Handler) sourceCampaignsHandler(w http.ResponseWriter, r *http.Request) {
	req := services.Acquire()
	defer services.Release(req)

	var err error

	req.ID, err = strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.responseError(w, http.StatusBadRequest, "invalid source id")
		return
	}

	for _, d := range strings.Split(r.URL.Query().Get("domains"), ",") {
		req.Domains = append(req.Domains, d)
	}

	campaigns, err := h.services.GetSourceCampaigns(req)
	if err != nil {
		h.responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.responseJSON(w, campaigns)
}

func (h *Handler) responseJSON(w http.ResponseWriter, data any) {
	raw, err := json.Marshal(data)
	if err != nil {
		h.responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(raw)
	if err != nil {
		h.responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) responseError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)

	_, err := w.Write([]byte(message))
	if err != nil {
		log.Println(err)
		return
	}
}
