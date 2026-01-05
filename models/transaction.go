package models

import (
	"time"
)

// Transaction represents a sale transaction
type Transaction struct {
	TransID       string            `json:"trans_id" firestore:"trans_id"`
	UserID        string            `json:"user_id" firestore:"user_id"`
	Date          time.Time         `json:"date" firestore:"date"`
	Total         float64           `json:"total" firestore:"total"`
	PaymentMethod string            `json:"payment_method" firestore:"payment_method"`
	Items         []TransactionItem `json:"items" firestore:"items"`
}

// TransactionItem represents an item in a transaction
type TransactionItem struct {
	ProductID string  `json:"product_id" firestore:"product_id"`
	Name      string  `json:"name" firestore:"name"`
	Quantity  int     `json:"qty" firestore:"qty"`
	Price     float64 `json:"price" firestore:"price"`
	Subtotal  float64 `json:"subtotal" firestore:"subtotal"`
}

// Payment method constants
const (
	PaymentCash    = "cash"
	PaymentCard    = "card"
	PaymentDigital = "digital"
)

// AddItem adds an item to the transaction
func (t *Transaction) AddItem(product *Product, quantity int) error {
	if !product.CanSell(quantity) {
		return ErrInsufficientStock
	}

	item := TransactionItem{
		ProductID: product.ProductID,
		Name:      product.Name,
		Quantity:  quantity,
		Price:     product.Price,
		Subtotal:  product.CalculateSubtotal(quantity),
	}

	t.Items = append(t.Items, item)
	t.calculateTotal()
	return nil
}

// RemoveItem removes an item from the transaction by product ID
func (t *Transaction) RemoveItem(productID string) {
	for i, item := range t.Items {
		if item.ProductID == productID {
			t.Items = append(t.Items[:i], t.Items[i+1:]...)
			break
		}
	}
	t.calculateTotal()
}

// UpdateItemQuantity updates the quantity of an item in the transaction
func (t *Transaction) UpdateItemQuantity(productID string, newQuantity int) {
	for i, item := range t.Items {
		if item.ProductID == productID {
			if newQuantity <= 0 {
				t.RemoveItem(productID)
				return
			}
			t.Items[i].Quantity = newQuantity
			t.Items[i].Subtotal = item.Price * float64(newQuantity)
			break
		}
	}
	t.calculateTotal()
}

// calculateTotal calculates the total amount of the transaction
func (t *Transaction) calculateTotal() {
	total := 0.0
	for _, item := range t.Items {
		total += item.Subtotal
	}
	t.Total = total
}

// GetItemCount returns the total number of items in the transaction
func (t *Transaction) GetItemCount() int {
	return len(t.Items)
}

// GetTotalQuantity returns the total quantity of all items
func (t *Transaction) GetTotalQuantity() int {
	totalQty := 0
	for _, item := range t.Items {
		totalQty += item.Quantity
	}
	return totalQty
}
