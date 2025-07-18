package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"pfg/internal/pack"

	"go.uber.org/zap"
)

type Handler struct {
	service *pack.Service
	logger  *zap.Logger
}

func NewHandler(service *pack.Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

type orderRequest struct {
	Quantity int `json:"quantity"`
}

type packEntry struct {
	Size  int `json:"size"`
	Count int `json:"count"`
}

type packResponse struct {
	Requested   int         `json:"requested"`
	Fulfilled   int         `json:"fulfilled"`
	Overpacked  int         `json:"overpacked"`
	TotalPacks  int         `json:"totalPacks"`
	PackDetails []packEntry `json:"packs"`
}

func (h *Handler) CalculatePacks(w http.ResponseWriter, r *http.Request) {
	var req orderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Quantity <= 0 {
		h.logger.Warn("Invalid request for CalculatePacks", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result, err := h.service.Calculate(r.Context(), req.Quantity)
	if err != nil {
		h.logger.Error("Failed to calculate packs", zap.Error(err), zap.Int("quantity", req.Quantity))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := packResponse{
		Requested:   req.Quantity,
		Fulfilled:   result.TotalItems,
		Overpacked:  result.TotalItems - req.Quantity,
		TotalPacks:  result.TotalPacks,
		PackDetails: []packEntry{},
	}

	for size, count := range result.Packs {
		resp.PackDetails = append(resp.PackDetails, packEntry{Size: size, Count: count})
	}

	h.logger.Info("Pack calculation completed", zap.Int("quantity", req.Quantity), zap.Any("response", resp))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) ListPackSizes(w http.ResponseWriter, r *http.Request) {
	sizes, err := h.service.ListPacks(r.Context())
	if err != nil {
		h.logger.Error("Failed to list pack sizes", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Pack sizes listed", zap.Int("count", len(sizes)))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sizes)
}

func (h *Handler) AddPackSize(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Size int `json:"size"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil || data.Size <= 0 {
		h.logger.Warn("Invalid pack size input", zap.Error(err))
		http.Error(w, "Invalid size", http.StatusBadRequest)
		return
	}
	if err := h.service.AddPack(r.Context(), data.Size); err != nil {
		h.logger.Error("Failed to add pack size", zap.Int("size", data.Size), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Pack size added", zap.Int("size", data.Size))
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeletePackSize(w http.ResponseWriter, r *http.Request) {
	sizeStr := r.URL.Query().Get("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		h.logger.Warn("Invalid size in delete request", zap.String("raw", sizeStr), zap.Error(err))
		http.Error(w, "Invalid size", http.StatusBadRequest)
		return
	}
	if err := h.service.RemovePack(context.Background(), size); err != nil {
		h.logger.Error("Failed to delete pack size", zap.Int("size", size), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Pack size deleted", zap.Int("size", size))
	w.WriteHeader(http.StatusNoContent)
}
