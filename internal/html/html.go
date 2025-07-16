package html

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"pfg/internal/pack"
)

type HTMLHandler struct {
	service   *pack.Service
	templates *template.Template
}

func NewHTMLHandler(service *pack.Service, templates *template.Template) *HTMLHandler {
	return &HTMLHandler{
		service:   service,
		templates: templates,
	}
}

func (h *HTMLHandler) RenderPackList(w http.ResponseWriter, r *http.Request) {
	sizes, err := h.service.ListPacks(r.Context())
	if err != nil {
		http.Error(w, "Failed to load packs", http.StatusInternalServerError)
		return
	}

	data := struct {
		Packs []int
	}{
		Packs: sizes,
	}

	err = h.templates.ExecuteTemplate(w, "packs.html", data)
	if err != nil {
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
		log.Println("packs.html render error:", err)
	}
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

	err = h.service.RemovePack(r.Context(), size)
	if err != nil {
		http.Error(w, "Failed to delete pack", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/packs", http.StatusSeeOther)
}

func (h *HTMLHandler) RenderCalculateForm(w http.ResponseWriter, r *http.Request) {
	var result *pack.PackResult

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err == nil {
			quantityStr := r.FormValue("quantity")
			quantity, err := strconv.Atoi(quantityStr)
			if err == nil && quantity > 0 {
				res, err := h.service.Calculate(r.Context(), quantity)
				if err == nil {
					result = &res
				}
			}
		}
	}

	data := struct {
		Result *pack.PackResult
	}{
		Result: result,
	}

	err := h.templates.ExecuteTemplate(w, "calculate.html", data)
	if err != nil {
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
		log.Println("calculate.html render error:", err)
	}
}
