package contract

// CreatePasteCommand is the service-layer input for creating a paste.
type CreatePasteCommand struct {
	Title      string
	Content    string
	Language   string
	Visibility string
}

// UpdatePasteCommand is the service-layer input for updating a paste.
type UpdatePasteCommand struct {
	Title      string
	Content    string
	Language   string
	Visibility string
}

// PasteResult is the service-layer output for paste operations.
type PasteResult struct {
	ID         int64
	OwnerID    int64
	Title      string
	ShortLink  string
	Content    string
	Language   string
	Visibility string
	CreatedAt  string
	UpdatedAt  string
}
