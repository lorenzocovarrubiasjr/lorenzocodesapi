package models

type WorkHistoryItem struct {
	ID          string   `json:"id"`
	CompanyName string   `json:"companyName"`
	Role        string   `json:"role"`
	Skills      []string `json:"skills"`
	Url         string   `json:"url"`
	LogoUrl     string   `json:"logoUrl"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}
