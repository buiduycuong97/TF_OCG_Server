package response

type CategoryResponse struct {
	CategoryID int32  `json:"categoryId"`
	Name       string `json:"name"`
	Handle     string `json:"handle"`
	Image      string `json:"image"`
}
