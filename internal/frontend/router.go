package frontend

import (
	"embed"
	"net/http"

	"github.com/brunofjesus/raspberry-bookshelf/internal/frontend/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static
var staticFs embed.FS

// NewHTTPRouter creates a new HTTP router with the provided handler functions.
func NewHTTPRouter(
	getCategoriesFn handlers.GetCategoriesFn,
	getBookFn handlers.GetBookFn,
	getBooksFn handlers.GetBooksFn,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fileServer := http.FileServer(http.FS(staticFs))
	r.Handle("/static/*", fileServer)

	r.Get("/", handlers.NewIndexHandler(getCategoriesFn).ServeHTTP)
	r.Get("/module/books", handlers.NewBooksHandler(getBooksFn).ServeHTTP)
	r.Get("/module/book/{bookID}", handlers.NewBookHandler(getBookFn).ServeHTTP)

	return r
}
