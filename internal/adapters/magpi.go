package adapters

import (
	"context"
	"encoding/xml"
	"log/slog"
	"net/http"
	"time"

	"github.com/brunofjesus/raspberry-bookshelf/internal/entities"
)

const magPiBookShelfURL = "https://magpi.raspberrypi.com/bookshelf.xml"

// MagPiAPI is an adapter for fetching MagPi books and magazines.
type MagPiAPI struct {
	httpClient http.Client
}

// BookshelfXML represents the structure of the MagPi bookshelf XML response.
type BookshelfXML struct {
	MagPi []BookshelfItem `xml:"MAGPI>ITEM"`
	Books []BookshelfItem `xml:"BOOKS>ITEM"`
}

// BookshelfItem represents a single item (book or magazine) in the MagPi bookshelf.
// The Category field is set manually after unmarshalling.
type BookshelfItem struct {
	Title       string `xml:"TITLE"`
	Description string `xml:"DESC"`
	Cover       string `xml:"COVER"`
	File        string `xml:"FILE"`
	PDF         string `xml:"PDF"`
	Category    string
}

// IsLocked checks if the item is locked (i.e., does not have a PDF link).
func (i *BookshelfItem) IsLocked() bool {
	return i.PDF == ""
}

// ToBookEntity converts a BookshelfItem to an entities.Book.
func (i *BookshelfItem) ToBookEntity() entities.Book {
	return entities.Book{
		Title:       i.Title,
		Description: i.Description,
		Cover:       i.Cover,
		Link:        i.PDF,
		Category:    i.Category,
	}
}

// NewMagPiAPI creates a new instance of MagPiAPI with a configured HTTP client.
func NewMagPiAPI() *MagPiAPI {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	return &MagPiAPI{
		httpClient: client,
	}
}

// GetBooks fetches the list of MagPi books and magazines from the MagPi API.
// It returns a slice of Book entities or an error if the operation fails.
// The function uses concurrency to process magazines and books in parallel.
func (m *MagPiAPI) GetBooks(ctx context.Context) ([]entities.Book, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, magPiBookShelfURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("Cannot close response body", slog.Any("error", err))
		}
	}()

	var magPiXML BookshelfXML
	err = xml.NewDecoder(resp.Body).Decode(&magPiXML)
	if err != nil {
		return nil, err
	}

	result := make([]entities.Book, 0, len(magPiXML.Books)+len(magPiXML.MagPi))
	magzCh := make(chan entities.Book, 1)
	bookCh := make(chan entities.Book, 1)

	go func() {
		for _, item := range magPiXML.MagPi {
			item.Category = "MagPI"
			magzCh <- item.ToBookEntity()
		}
		close(magzCh)
	}()
	go func() {
		for _, item := range magPiXML.Books {
			item.Category = "Book"
			bookCh <- item.ToBookEntity()
		}
		close(bookCh)
	}()

	for count := 0; count < 2; {
		select {
		case v, ok := <-magzCh:
			if !ok {
				magzCh = nil
				count++
				continue
			}
			result = append(result, v)
		case v, ok := <-bookCh:
			if !ok {
				bookCh = nil
				count++
				continue
			}
			result = append(result, v)
		}
	}

	return result, nil
}
