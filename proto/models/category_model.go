package models

type Category struct {
	CategoryID int32  `json:"categoryId"`
	Handle     string `json:"handle"`
	Name       string `json:"name"`
	Image      string `json:"image"`
}
