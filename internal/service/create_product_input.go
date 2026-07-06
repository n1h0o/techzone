package service

type CreateProductInput struct {
	Name        string `json:"name" example:"iPhone 17"`
	Description string `json:"description" example:"Apple smartphone"`
	Price       int64  `json:"price" example:"99990"`
	Stock       int    `json:"stock" example:"10"`
	ImageURL    string `json:"image_url" example:"https://example.com/iphone.jpg"`
}
