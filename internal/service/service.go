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

type Runner interface {
	Run(ctx context.Context) error
}

type Service struct {
	bookUpdater Runner
	bookStorage *bookshelf.Storage
}

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
