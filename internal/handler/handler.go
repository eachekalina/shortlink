package handler

import (
	"encoding/json"
	"github.com/eachekalina/shortlink/internal/service"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	s *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{s: s}
}

type CreateLinkRequest struct {
	Link string `json:"link"`
}

type CreateLinkResponse struct {
	Path string `json:"path"`
}

func (h *Handler) HandleCreateLink(w http.ResponseWriter, r *http.Request) {
	var req CreateLinkRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	path, err := h.s.CreateShortLink(r.Context(), req.Link)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := CreateLinkResponse{Path: path}
	bytes, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(bytes)
}

func (h *Handler) HandleLink(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	link, err := h.s.GetLink(r.Context(), path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, link, http.StatusMovedPermanently)
}
