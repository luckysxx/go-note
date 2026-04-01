package contract

// CreateSnippetCommand 创建 Snippet 的服务层输入参数。
type CreateSnippetCommand struct {
	Type       string // code / note / file
	Title      string
	Content    string // code/note 时使用
	FileURL    string // file 时使用
	FileSize   int64  // file 时使用
	MimeType   string // file 时使用
	Language   string
	Visibility string
	GroupID    *int64 // 可选的分组 ID
}

// UpdateSnippetCommand 更新 Snippet 的服务层输入参数。
type UpdateSnippetCommand struct {
	Title      string
	Content    string
	Language   string
	Visibility string
	GroupID    *int64
}

// SnippetResult 服务层输出的 Snippet 数据结构。
type SnippetResult struct {
	ID         int64
	OwnerID    int64
	Type       string
	Title      string
	Content    string
	FileURL    string
	FileSize   int64
	MimeType   string
	Language   string
	Visibility string
	GroupID    *int64
	CreatedAt  string
	UpdatedAt  string
}
