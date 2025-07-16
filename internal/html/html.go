package html

import (
	"html/template"
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

func (h *HTMLHandler) RenderWelcomePage(w http.ResponseWriter, r *http.Request) {
	err := h.templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Template rendering failed", http.StatusInternalServerError)
	}
}

func (h *HTMLHandler) RenderPackList(w http.ResponseWriter, r *http.Request) {
	sizes, err := h.service.ListPacks(r.Context())
	if err != nil {
		http.Error(w, "Failed to load packs", http.StatusInternalServerError)
		return
	}

	err = h.templates.ExecuteTemplate(w, "packs.html", map[string]interface{}{
		"packs": sizes,
	})
	if err != nil {
		http.Error(w, "Template rendering failed", http.StatusInternalServerError)
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
		qtyStr := r.FormValue("quantity")
		qty, err := strconv.Atoi(qtyStr)
		if err == nil && qty > 0 {
			val, err := h.service.Calculate(r.Context(), qty)
			if err == nil {
				result = &val
			}
		}
	}

	err := h.templates.ExecuteTemplate(w, "calculate.html", map[string]interface{}{
		"result": result,
	})
	if err != nil {
		http.Error(w, "Template rendering failed", http.StatusInternalServerError)
	}
}
