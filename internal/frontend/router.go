package frontend

import (
	"net/http"

	"github.com/brunofjesus/raspberry-bookshelf/internal/frontend/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHTTPRouter(
	getCategoriesFn handlers.GetCategoriesFn,
	getBookFn handlers.GetBookFn,
	getBooksFn handlers.GetBooksFn,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	r.Get("/", handlers.NewIndexHandler(getCategoriesFn).ServeHTTP)
	r.Get("/module/books", handlers.NewBooksHandler(getBooksFn).ServeHTTP)
	r.Get("/module/book/{bookID}", handlers.NewBookHandler(getBookFn).ServeHTTP)
	//	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	//		w.Write([]byte("hello world"))
	//	})

	return r
}
