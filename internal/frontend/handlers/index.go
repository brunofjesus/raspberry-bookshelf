package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/brunofjesus/raspberry-bookshelf/internal/frontend/templates"
)

type (
	GetCategoriesFn func(ctx context.Context) ([]string, error)
	IndexHandler    struct {
		getCategoriesFn GetCategoriesFn
	}
)

// NewIndexHandler creates a new IndexHandler with the provided GetCategoriesFn.
// This handler is responsible for serving the index page, which includes
// the list of categories and highlights the current category if provided on the
// NavBar.
func NewIndexHandler(getCategories GetCategoriesFn) *IndexHandler {
	return &IndexHandler{
		getCategoriesFn: getCategories,
	}
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	currentCategory := r.URL.Query().Get("cat")
	categories, err := h.getCategoriesFn(r.Context())
	if err != nil {
		slog.Error("cannot get list of categories", slog.Any("error", err))
	}
	c := templates.PageIndex(currentCategory)

	err = templates.Layout(c, "Bookshelf", currentCategory, categories).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
