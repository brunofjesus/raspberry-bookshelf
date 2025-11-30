package entities

// Book represents a book or magazine entity.
// It is part of the domain layer and used across the application.
type Book struct {
	ID          string
	Title       string
	Description string
	Cover       string
	Link        string
	Category    string
}
