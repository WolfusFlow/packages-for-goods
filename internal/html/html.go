package html

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"pfg/internal/pack"

	"github.com/CloudyKit/jet/v6"
)

type HTMLHandler struct {
	service *pack.Service
	views   *jet.Set
}

func NewHTMLHandler(service *pack.Service, views *jet.Set) *HTMLHandler {
	return &HTMLHandler{service: service, views: views}
}

func (h *HTMLHandler) RenderPackList(w http.ResponseWriter, r *http.Request) {
	sizes, err := h.service.ListPacks(r.Context())
	if err != nil {
		http.Error(w, "Failed to load packs", http.StatusInternalServerError)
		return
	}
	view, err := h.views.GetTemplate("packs.jet")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	vars := make(jet.VarMap)
	vars.Set("packs", sizes)
	view.Execute(w, vars, nil)
}

func (h *HTMLHandler) HandleAddPack(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}
	sizeStr := r.FormValue("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		return
	}
	err = h.service.AddPack(r.Context(), size)
	if err != nil {
		http.Error(w, "Failed to add pack", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/packs", http.StatusSeeOther)
}

func (h *HTMLHandler) HandleDeletePack(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}
	sizeStr := r.FormValue("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		return
	}
	err = h.service.RemovePack(context.Background(), size)
	if err != nil {
		http.Error(w, "Failed to delete pack", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/packs", http.StatusSeeOther)
}

func (h *HTMLHandler) RenderCalculateForm(w http.ResponseWriter, r *http.Request) {
	t, err := h.views.GetTemplate("calculate.jet")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	var result *pack.PackResult
	if r.Method == http.MethodPost {
		quantityStr := r.FormValue("quantity")
		quantity, err := strconv.Atoi(quantityStr)
		if err == nil && quantity > 0 {
			resultVal, err := h.service.Calculate(r.Context(), quantity)
			if err != nil {
				http.Error(w, "Calculation failed: "+err.Error(), http.StatusInternalServerError)
				log.Println("Calculation error:", err)
				return
			}
			result = &resultVal
		}
	}

	vars := make(jet.VarMap)
	vars.Set("result", result)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(w, vars, nil)
	if err != nil {
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
		log.Println("Execution error:", err)
	}
}
