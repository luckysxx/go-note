package dto

// CreateSnippetRequest is the HTTP DTO for creating a snippet.
type CreateSnippetRequest struct {
	Type       string `json:"type"`
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content"`
	FileURL    string `json:"file_url"`
	FileSize   int64  `json:"file_size"`
	MimeType   string `json:"mime_type"`
	Language   string `json:"language"`
	Visibility string `json:"visibility"`
	GroupID    *int64 `json:"group_id"`
}

// UpdateSnippetRequest is the HTTP DTO for updating a snippet.
type UpdateSnippetRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content"`
	Language   string `json:"language"`
	Visibility string `json:"visibility"`
	GroupID    *int64 `json:"group_id"`
}

// SnippetResponse is the HTTP DTO for snippet responses.
type SnippetResponse struct {
	ID         int64  `json:"id"`
	OwnerID    int64  `json:"owner_id"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	FileURL    string `json:"file_url,omitempty"`
	FileSize   int64  `json:"file_size,omitempty"`
	MimeType   string `json:"mime_type,omitempty"`
	Language   string `json:"language"`
	Visibility string `json:"visibility"`
	GroupID    *int64 `json:"group_id,omitempty"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
