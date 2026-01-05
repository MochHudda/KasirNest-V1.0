package models

import (
	"time"
)

// Product represents a product in the inventory
type Product struct {
	ProductID string    `json:"product_id" firestore:"product_id"`
	Name      string    `json:"name" firestore:"name"`
	Price     float64   `json:"price" firestore:"price"`
	Stock     int       `json:"stock" firestore:"stock"`
	Category  string    `json:"category" firestore:"category"`
	Barcode   string    `json:"barcode" firestore:"barcode"`
	ImageURL  string    `json:"image_url" firestore:"image_url"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}

// IsInStock checks if product has available stock
func (p *Product) IsInStock() bool {
	return p.Stock > 0
}

// CanSell checks if product can be sold with given quantity
func (p *Product) CanSell(quantity int) bool {
	return p.Stock >= quantity && quantity > 0
}

// CalculateSubtotal calculates subtotal for given quantity
func (p *Product) CalculateSubtotal(quantity int) float64 {
	return p.Price * float64(quantity)
}

// UpdateStock reduces stock by given quantity
func (p *Product) UpdateStock(quantity int) error {
	if !p.CanSell(quantity) {
		return ErrInsufficientStock
	}
	p.Stock -= quantity
	p.UpdatedAt = time.Now()
	return nil
}
