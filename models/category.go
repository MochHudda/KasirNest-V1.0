package models

// Category represents a product category
type Category struct {
	CategoryID  string `json:"category_id" firestore:"category_id"`
	Name        string `json:"name" firestore:"name"`
	Description string `json:"description" firestore:"description"`
}

// DefaultCategories returns a slice of default categories
func DefaultCategories() []Category {
	return []Category{
		{CategoryID: "food", Name: "Makanan", Description: "Produk makanan dan minuman"},
		{CategoryID: "electronic", Name: "Elektronik", Description: "Perangkat elektronik"},
		{CategoryID: "fashion", Name: "Fashion", Description: "Pakaian dan aksesoris"},
		{CategoryID: "health", Name: "Kesehatan", Description: "Produk kesehatan dan kecantikan"},
		{CategoryID: "household", Name: "Rumah Tangga", Description: "Keperluan rumah tangga"},
		{CategoryID: "stationery", Name: "Alat Tulis", Description: "Alat tulis dan kantor"},
		{CategoryID: "other", Name: "Lainnya", Description: "Kategori lainnya"},
	}
}
