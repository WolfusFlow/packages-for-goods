package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"pfg/internal/pack"
)

type Handler struct {
	mu        sync.RWMutex
	packSizes []int
}

func NewHandler() *Handler {
	return &Handler{
		packSizes: []int{250, 500, 1000, 2000, 5000},
	}
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

	h.mu.RLock()
	sizes := append([]int{}, h.packSizes...)
	h.mu.RUnlock()

	result := pack.CalculateOptimalPacks(req.Quantity, sizes)

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
	h.mu.RLock()
	defer h.mu.RUnlock()
	json.NewEncoder(w).Encode(h.packSizes)
}

func (h *Handler) AddPackSize(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Size int `json:"size"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil || data.Size <= 0 {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, s := range h.packSizes {
		if s == data.Size {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	h.packSizes = append(h.packSizes, data.Size)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeletePackSize(w http.ResponseWriter, r *http.Request) {
	sizeStr := r.URL.Query().Get("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	filtered := h.packSizes[:0]
	for _, s := range h.packSizes {
		if s != size {
			filtered = append(filtered, s)
		}
	}
	h.packSizes = filtered
	w.WriteHeader(http.StatusNoContent)
}
