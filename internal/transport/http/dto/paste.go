package dto

// CreatePasteRequest is the HTTP DTO for creating a paste.
type CreatePasteRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Language   string `json:"language" binding:"required"`
	Visibility string `json:"visibility"`
}

// UpdatePasteRequest is the HTTP DTO for updating a paste.
type UpdatePasteRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Language   string `json:"language" binding:"required"`
	Visibility string `json:"visibility"`
}

// PasteResponse is the HTTP DTO for paste responses.
type PasteResponse struct {
	ID         int64  `json:"id"`
	OwnerID    int64  `json:"owner_id"`
	Title      string `json:"title"`
	ShortLink  string `json:"short_link,omitempty"`
	Content    string `json:"content"`
	Language   string `json:"language"`
	Visibility string `json:"visibility"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
