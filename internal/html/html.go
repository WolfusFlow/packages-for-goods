package html

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"pfg/internal/auth"
	"pfg/internal/config"
	"pfg/internal/pack"

	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
)

type HTMLHandler struct {
	service   *pack.Service
	templates *template.Template
	config    *config.Config
	logger    *zap.Logger
}

func NewHTMLHandler(
	service *pack.Service,
	templates *template.Template,
	config *config.Config,
	logger *zap.Logger,
) *HTMLHandler {
	return &HTMLHandler{
		service:   service,
		templates: templates,
		config:    config,
		logger:    logger,
	}
}

func (h *HTMLHandler) RenderWelcomePage(w http.ResponseWriter, r *http.Request) {
	isAdmin, email := adminInfoFromCookie(r)
	err := h.templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"Path":       r.URL.Path,
		"IsLoggedIn": isAdmin,
		"UserEmail":  email,
	})
	if err != nil {
		h.logger.Error("Failed to render welcome page", zap.Error(err))
		http.Error(w, "Template rendering failed", http.StatusInternalServerError)
	}
}

func (h *HTMLHandler) RenderPackList(w http.ResponseWriter, r *http.Request) {
	sizes, err := h.service.ListPacks(r.Context())
	if err != nil {
		h.logger.Error("Failed to load packs", zap.Error(err))
		http.Error(w, "Failed to load packs", http.StatusInternalServerError)
		return
	}

	isAdmin, email := adminInfoFromCookie(r)
	err = h.templates.ExecuteTemplate(w, "packs.html", map[string]interface{}{
		"packs":      sizes,
		"Path":       r.URL.Path,
		"IsLoggedIn": isAdmin,
		"UserEmail":  email,
	})

	if err != nil {
		h.logger.Error("Failed to render packs page", zap.Error(err))
		http.Error(w, "Template rendering failed", http.StatusInternalServerError)
	}
}

func (h *HTMLHandler) HandleAddPack(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.logger.Warn("Invalid form on AddPack", zap.Error(err))
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	sizeStr := r.FormValue("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		h.logger.Warn("Invalid pack size value", zap.String("input", sizeStr), zap.Error(err))
		http.Error(w, "Invalid size", http.StatusBadRequest)
		return
	}

	err = h.service.AddPack(r.Context(), size)
	if err != nil {
		h.logger.Warn("Duplicate or failed pack add", zap.Int("size", size), zap.Error(err))
		sizes, _ := h.service.ListPacks(r.Context())
		isAdmin, email := adminInfoFromCookie(r)
		h.templates.ExecuteTemplate(w, "packs.html", map[string]interface{}{
			"packs":      sizes,
			"Path":       r.URL.Path,
			"error":      err.Error(),
			"IsLoggedIn": isAdmin,
			"UserEmail":  email,
		})
		return
	}

	h.logger.Info("Pack added", zap.Int("size", size))
	http.Redirect(w, r, "/packs", http.StatusSeeOther)
}

func (h *HTMLHandler) HandleDeletePack(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.logger.Warn("Invalid form on DeletePack", zap.Error(err))
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	sizeStr := r.FormValue("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		h.logger.Warn("Invalid pack size for deletion", zap.String("input", sizeStr), zap.Error(err))
		http.Error(w, "Invalid size", http.StatusBadRequest)
		return
	}

	err = h.service.RemovePack(r.Context(), size)
	if err != nil {
		h.logger.Error("Failed to delete pack", zap.Int("size", size), zap.Error(err))
		http.Error(w, "Failed to delete pack", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Pack deleted", zap.Int("size", size))
	http.Redirect(w, r, "/packs", http.StatusSeeOther)
}

func (h *HTMLHandler) RenderCalculateForm(w http.ResponseWriter, r *http.Request) {
	var result *pack.PackResult

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			h.logger.Warn("Failed to parse form in calculate", zap.Error(err))
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		qtyStr := r.FormValue("quantity")
		qty, err := strconv.Atoi(qtyStr)
		if err != nil || qty <= 0 {
			h.logger.Warn("Invalid quantity input", zap.String("input", qtyStr), zap.Error(err))
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}

		val, err := h.service.Calculate(r.Context(), qty)
		if err != nil {
			h.logger.Error("Failed to calculate", zap.Int("qty", qty), zap.Error(err))
			http.Error(w, "Failed to calculate packs", http.StatusInternalServerError)
			return
		}

		result = &val
		h.logger.Info("HTML pack calculation completed", zap.Int("quantity", qty), zap.Any("result", val))
	}

	isAdmin, email := adminInfoFromCookie(r)
	err := h.templates.ExecuteTemplate(w, "calculate.html", map[string]interface{}{
		"result":     result,
		"Path":       r.URL.Path,
		"IsLoggedIn": isAdmin,
		"UserEmail":  email,
	})
	if err != nil {
		h.logger.Error("Failed to render calculate page", zap.Error(err))
		http.Error(w, "Template rendering failed", http.StatusInternalServerError)
	}
}

func (h *HTMLHandler) RenderUnauthorized(w http.ResponseWriter, r *http.Request) {
	h.logger.Warn("Unauthorized access attempt", zap.String("path", r.URL.Path))
	_ = h.templates.ExecuteTemplate(w, "unauthorized.html", map[string]any{
		"Path":       r.URL.Path,
		"IsLoggedIn": false,
		"UserEmail":  "",
	})
}

func (h *HTMLHandler) RenderLoginForm(w http.ResponseWriter, r *http.Request) {
	isAdmin, email := adminInfoFromCookie(r)
	h.templates.ExecuteTemplate(w, "login.html", map[string]any{
		"Path":       r.URL.Path,
		"IsLoggedIn": isAdmin,
		"UserEmail":  email,
	})
}

func (h *HTMLHandler) HandleLoginPost(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pass := r.FormValue("password")

	if email != h.config.AdminEmail || pass != h.config.AdminPassword {
		h.logger.Warn("Login failed", zap.String("email", email))
		isAdmin, _ := adminInfoFromCookie(r)
		h.templates.ExecuteTemplate(w, "login.html", map[string]any{
			"Error":      "Invalid credentials",
			"Path":       r.URL.Path,
			"IsLoggedIn": isAdmin,
			"UserEmail":  email,
		})
		return
	}

	_, token, _ := auth.TokenAuth.Encode(map[string]any{
		"email":   email,
		"isAdmin": true,
		"exp":     jwtauth.ExpireIn(30 * time.Minute),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "admin_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	h.logger.Info("Login successful", zap.String("email", email))
	http.Redirect(w, r, "/packs", http.StatusSeeOther)
}

func adminInfoFromCookie(r *http.Request) (isAdmin bool, email string) {
	cookie, err := r.Cookie("admin_token")
	if err != nil {
		return false, ""
	}

	token, err := auth.TokenAuth.Decode(cookie.Value)
	if err != nil {
		return false, ""
	}

	claims := token.PrivateClaims()
	admin, _ := claims["isAdmin"].(bool)
	emailStr, _ := claims["email"].(string)
	return admin, emailStr
}
