package bookshelf

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"sort"

	"github.com/brunofjesus/raspberry-bookshelf/internal/entities"
)

type Storage struct {
	books           []entities.Book
	bookIDMap       map[string]*entities.Book
	bookCategoryMap map[string][]entities.Book
}

func NewStorage() *Storage {
	return &Storage{
		books:           []entities.Book{},
		bookIDMap:       map[string]*entities.Book{},
		bookCategoryMap: map[string][]entities.Book{},
	}
}

func (s *Storage) GetByID(ctx context.Context, id string) (*entities.Book, error) {
	book := s.bookIDMap[id]
	return book, nil
}

func (s *Storage) Get(ctx context.Context, category string) ([]entities.Book, error) {
	if len(category) > 0 {
		slog.Debug("getting books in category", slog.String("category", category))
		return s.bookCategoryMap[category], nil
	}
	return s.books, nil
}

func (s *Storage) GetCategories(ctx context.Context) ([]string, error) {
	result := []string{}
	categories := maps.Keys(s.bookCategoryMap)
	for c := range categories {
		result = append(result, c)
	}
	sort.StringSlice(result).Sort()
	return result, nil
}

func (s *Storage) ReplaceAll(ctx context.Context, books []entities.Book) error {
	slog.Debug("replacing books", slog.Int("size", len(books)))
	bookSlice := make([]entities.Book, 0, len(books))
	bookIDMap := make(map[string]*entities.Book)
	bookCategoryMap := make(map[string][]entities.Book)

	for _, book := range books {
		if err := s.genBookID(ctx, &book); err != nil {
			return fmt.Errorf("error generating id: %w", err)
		}

		bookSlice = append(bookSlice, book)
		bookIDMap[book.ID] = &book
		bookCategoryMap[book.Category] = append(bookCategoryMap[book.Category], book)
	}

	s.books = bookSlice
	s.bookCategoryMap = bookCategoryMap
	s.bookIDMap = bookIDMap
	return nil
}

func (s *Storage) genBookID(_ context.Context, book *entities.Book) error {
	if book == nil {
		return errors.New("book is required")
	}

	if book.ID != "" {
		slog.Info("book already had an id", slog.String("bookID", book.ID))
		return nil
	}

	hash := sha1.New()
	_, err := io.WriteString(hash, fmt.Sprintf("%s:%s", book.Cover, book.Title))
	if err != nil {
		return err
	}
	book.ID = hex.EncodeToString(hash.Sum(nil))
	return nil
}
