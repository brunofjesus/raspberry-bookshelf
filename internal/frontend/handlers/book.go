package handlers

import (
	"context"
	"net/http"

	"github.com/brunofjesus/raspberry-bookshelf/internal/entities"
	"github.com/brunofjesus/raspberry-bookshelf/internal/frontend/templates/modules"
	"github.com/go-chi/chi/v5"
)

type (
	GetBookFn   = func(ctx context.Context, bookId string) (*entities.Book, error)
	BookHandler struct {
		getBookFn GetBookFn
	}
)

func NewBookHandler(getBook GetBookFn) *BookHandler {
	return &BookHandler{
		getBookFn: getBook,
	}
}

func (h *BookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookID")
	book, err := h.getBookFn(r.Context(), bookID)
	if err != nil {
		http.Error(w, "Error fetching books", http.StatusInternalServerError)
		return
	}

	if book == nil {
		// TODO: maybe we should return some error component?
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	c := modules.BookInfo(book)
	err = c.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
