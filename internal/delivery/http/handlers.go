package delivery

import (
	"encoding/json"
	"github.com/maxzhovtyj/adtelligent-hw/internal/models"
	"github.com/maxzhovtyj/adtelligent-hw/internal/services"
	"log"
	"net/http"
	"strconv"
	"sync"
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
	req := acquire()
	defer release(req)

	var err error

	req.ID, err = strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.responseError(w, http.StatusBadRequest, "invalid source id")
		return
	}

	campaigns, err := h.sourceCampaigns(req)
	if err != nil {
		h.responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.responseJSON(w, campaigns)
}

var p sync.Pool

func acquire() *Request {
	r := p.Get()
	if r == nil {
		return new(Request)
	}

	return r.(*Request)
}

func release(r *Request) {
	r.Reset()
	p.Put(r)
}

type Request struct {
	ID int
}

func (r *Request) Reset() {
	r.ID = 0
}

func (h *Handler) sourceCampaigns(req *Request) ([]models.Campaign, error) {
	campaigns, err := h.services.GetSourceCampaigns(req.ID)
	if err != nil {
		return nil, err
	}

	return campaigns, nil
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
