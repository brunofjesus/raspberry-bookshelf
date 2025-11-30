package adapters

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/brunofjesus/raspberry-bookshelf/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBooks(t *testing.T) {
	xmlContent, err := os.ReadFile("testdata/bookshelf.xml")
	require.Nil(t, err, "failed to read test XML file: %v", err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(xmlContent)
	}))
	defer server.Close()

	subject := MagPiAPI{
		httpClient:        server.Client(),
		magPiBookShelfURL: server.URL,
	}

	result, err := subject.GetBooks(t.Context())
	require.Nil(t, err, "get books returned error: %v", err)
	require.NotEmpty(t, result, "result cannot be empty")
	require.Equal(t, 5, len(result), "should have 3 items")

	expectedItems := []entities.Book{
		{
			Title:       "MagPI Locked Mag 1",
			Description: "Description for MagPi Mag 1",
			Cover:       "http://localhost/covers/1",
			Link:        "",
			Category:    "MagPI",
		},
		{
			Title:       "MagPI Available Mag 2",
			Description: "Description for the MagPI Available Mag 2",
			Cover:       "http://localhost/covers/2",
			Link:        "http://localhost/magpi/2",
			Category:    "MagPI",
		},
		{
			Title:       "MagPI Available Mag 3",
			Description: "Description for the MagPI Available Mag 3",
			Cover:       "http://localhost/covers/3",
			Link:        "http://localhost/magpi/3",
			Category:    "MagPI",
		},
		{
			Title:       "Some non available book yet",
			Description: "This is a forbidden book",
			Cover:       "http://localhost/covers/4",
			Link:        "",
			Category:    "Book",
		},
		{
			Title:       "Available Book 1",
			Description: "Description for the Available Book 1",
			Cover:       "http://localhost/covers/book1",
			Link:        "http://localhost/book/1",
			Category:    "Book",
		},
	}

	assert.ElementsMatch(t, expectedItems, result, "result should match")
}
