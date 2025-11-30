package bookshelf

import (
	"context"
	"log/slog"
	"time"

	"github.com/brunofjesus/raspberry-bookshelf/internal/entities"
)

// BookClient defines the interface for fetching books from an external source.
// It is used by the BookshelfUpdater to get the latest book data.
type BookClient interface {
	GetBooks(ctx context.Context) ([]entities.Book, error)
}

// BookReferenceStorage defines the interface for storing book references.
// It is used by the BookshelfUpdater to update the stored book data.
type BookReferenceStorage interface {
	ReplaceAll(ctx context.Context, books []entities.Book) error
}

// BookshelfUpdater is responsible for periodically updating the bookshelf
// by fetching new book data from a BookClient and storing it in a BookReferenceStorage.
type BookshelfUpdater struct {
	bookClient BookClient
	storage    BookReferenceStorage
	interval   time.Duration
}

// NewBookshelfUpdater creates a new instance of BookshelfUpdater.
// It takes a BookClient, a BookReferenceStorage, and an update interval as parameters.
func NewBookshelfUpdater(
	bookClient BookClient,
	storage BookReferenceStorage,
	interval time.Duration,
) *BookshelfUpdater {
	return &BookshelfUpdater{
		bookClient: bookClient,
		storage:    storage,
		interval:   interval,
	}
}

// Run starts the bookshelf updater, which periodically fetches new book data
// and updates the storage. It runs until the provided context is done.
// This operation is blocking, you might want to run it in a separate goroutine.
func (u *BookshelfUpdater) Run(ctx context.Context) error {
	slog.Debug("starting the updater")
	for {
		select {
		case <-ctx.Done():
			slog.Debug("context is done, exiting the bookshelf updater")
			return nil
		default:
			books, err := u.bookClient.GetBooks(ctx)
			if err != nil {
				slog.ErrorContext(ctx, "failed to get books", slog.Any("error", err))
			} else if err := u.storage.ReplaceAll(ctx, books); err != nil {
				slog.ErrorContext(ctx, "failed to update books", slog.Any("error", err))
			}
			slog.Debug("updater got new books, sleeping", slog.Any("interval", u.interval))
		}
		time.Sleep(u.interval)
	}
}
