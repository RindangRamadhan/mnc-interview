package toc

type ResponseToc struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Table       string `json:"table"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	TableID     string `json:"table_id"`
	Order       string `json:"order"`
}
