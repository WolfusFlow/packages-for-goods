package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"pfg/internal/pack"
)

type Handler struct {
	service *pack.Service
}

func NewHandler(service *pack.Service) *Handler {
	return &Handler{service: service}
}

type orderRequest struct {
	Quantity int `json:"quantity"`
}

type packResponse struct {
	TotalItems int         `json:"totalItems"`
	TotalPacks int         `json:"totalPacks"`
	Packs      []packEntry `json:"packs"`
}

type packEntry struct {
	Size  int `json:"size"`
	Count int `json:"count"`
}

func (h *Handler) CalculatePacks(w http.ResponseWriter, r *http.Request) {
	var req orderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Quantity <= 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	result, err := h.service.Calculate(r.Context(), req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := packResponse{
		TotalItems: result.TotalItems,
		TotalPacks: result.TotalPacks,
	}
	for size, count := range result.Packs {
		resp.Packs = append(resp.Packs, packEntry{Size: size, Count: count})
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) ListPackSizes(w http.ResponseWriter, r *http.Request) {
	sizes, err := h.service.ListPacks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(sizes)
}

func (h *Handler) AddPackSize(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Size int `json:"size"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil || data.Size <= 0 {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		return
	}
	if err := h.service.AddPack(r.Context(), data.Size); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeletePackSize(w http.ResponseWriter, r *http.Request) {
	sizeStr := r.URL.Query().Get("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		return
	}
	if err := h.service.RemovePack(context.Background(), size); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
