package content

import "time"

type TechnicalItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type ProjectItem struct {
	ID          string   `json:"id"`
	Tag         string   `json:"tag"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Images      []string `json:"images"`
}

type PortfolioInfo struct {
	OwnerName    string `json:"ownerName"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	About        string `json:"about"`
	ContactEmail string `json:"contactEmail"`
	ContactPhone string `json:"contactPhone"`
	Location     string `json:"location"`
}

type ContentBody struct {
	Technical     []TechnicalItem `json:"technical"`
	Projects      []ProjectItem   `json:"projects"`
	PortfolioInfo PortfolioInfo   `json:"portfolioInfo"`
}

type GetAdminContentResponse struct {
	Locale    string      `json:"locale"`
	Version   int         `json:"version"`
	UpdatedAt time.Time   `json:"updated_at"`
	Content   ContentBody `json:"content"`
}

type PutAdminContentRequest struct {
	Version int         `json:"version"`
	Content ContentBody `json:"content"`
}

type PutAdminContentResponse struct {
	OK        bool      `json:"ok"`
	Version   int       `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PublishResponse struct {
	OK               bool      `json:"ok"`
	PublishedVersion int       `json:"published_version"`
	PublishedAt      time.Time `json:"published_at"`
}

type HistoryItem struct {
	Locale    string      `json:"locale"`
	Version   int         `json:"version"`
	Content   ContentBody `json:"content"`
	UpdatedBy string      `json:"updated_by"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type HistoryResponse struct {
	Locale  string        `json:"locale"`
	History []HistoryItem `json:"history"`
}

type PublicContentResponse struct {
	Locale  string      `json:"locale"`
	Content ContentBody `json:"content"`
}

type TechnicalListResponse struct {
	Locale    string          `json:"locale"`
	Version   int             `json:"version"`
	UpdatedAt time.Time       `json:"updated_at"`
	Items     []TechnicalItem `json:"items"`
}

type TechnicalMutationResponse struct {
	OK        bool          `json:"ok"`
	Version   int           `json:"version"`
	UpdatedAt time.Time     `json:"updated_at"`
	Item      TechnicalItem `json:"item"`
}

type TechnicalDeleteResponse struct {
	OK        bool      `json:"ok"`
	Version   int       `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedID string    `json:"deleted_id"`
}
