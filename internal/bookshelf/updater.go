package bookshelf

import (
	"context"
	"log/slog"
	"time"

	"github.com/brunofjesus/raspberry-bookshelf/internal/entities"
)

type BookClient interface {
	GetBooks(ctx context.Context) ([]entities.Book, error)
}

type BookReferenceStorage interface {
	ReplaceAll(ctx context.Context, books []entities.Book) error
}

type BookshelfUpdater struct {
	bookClient BookClient
	storage    BookReferenceStorage
	interval   time.Duration
}

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
