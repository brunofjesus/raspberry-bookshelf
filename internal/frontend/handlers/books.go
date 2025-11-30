package handlers

import (
	"context"
	"net/http"

	"github.com/brunofjesus/raspberry-bookshelf/internal/entities"
	"github.com/brunofjesus/raspberry-bookshelf/internal/frontend/templates/modules"
)

type (
	GetBooksFn   = func(ctx context.Context, category string) ([]entities.Book, error)
	BooksHandler struct {
		getBooksFn GetBooksFn
	}
)

// NewBooksHandler creates a new BooksHandler with the provided GetBooksFn.
// This handler is responsible for serving a list of books, optionally filtered by category.
// It returns a component that can be displayed on a page.
func NewBooksHandler(getBooks GetBooksFn) *BooksHandler {
	return &BooksHandler{
		getBooksFn: getBooks,
	}
}

func (h *BooksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	currentCategory := r.URL.Query().Get("cat")
	books, err := h.getBooksFn(r.Context(), currentCategory)
	if err != nil {
		http.Error(w, "Error fetching books", http.StatusInternalServerError)
		return
	}

	c := modules.Books(books)
	err = c.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
