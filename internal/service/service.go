package service

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/brunofjesus/raspberry-bookshelf/internal/adapters"
	"github.com/brunofjesus/raspberry-bookshelf/internal/bookshelf"
	"github.com/brunofjesus/raspberry-bookshelf/internal/frontend"
	"golang.org/x/sync/errgroup"
)

// Runner defines an interface for components that can be run.
type Runner interface {
	Run(ctx context.Context) error
}

// Service represents the main application service.
// It holds the data needed for the application to run.
type Service struct {
	bookUpdater Runner
	bookStorage *bookshelf.Storage
}

// New creates a new instance of the Service.
// It initializes the necessary components such as the book client,
// book storage, and book updater.
func New() Service {
	bookClient := adapters.NewMagPiAPI()
	bookStorage := bookshelf.NewStorage()

	updater := bookshelf.NewBookshelfUpdater(
		bookClient,
		bookStorage,
		1*time.Hour,
	)

	return Service{
		bookUpdater: updater,
		bookStorage: bookStorage,
	}
}

// Run starts the service, including the book updater and the HTTP web server.
func (s Service) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		slog.Debug("Starting the book updater")
		return s.bookUpdater.Run(ctx)
	})

	g.Go(func() error {
		slog.Debug("Starting the HTTP Web Server")
		router := frontend.NewHTTPRouter(
			s.bookStorage.GetCategories,
			s.bookStorage.GetByID,
			s.bookStorage.Get,
		)
		return http.ListenAndServe("0.0.0.0:8080", router)
	})

	// Block until all goroutines finish
	return g.Wait()
}
