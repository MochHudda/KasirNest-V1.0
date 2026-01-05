package models

import (
	"fmt"
	"time"
)

// Report represents a sales report
type Report struct {
	ReportID          string       `json:"report_id" firestore:"report_id"`
	Date              time.Time    `json:"date" firestore:"date"`
	TotalSales        float64      `json:"total_sales" firestore:"total_sales"`
	TotalTransactions int          `json:"total_transactions" firestore:"total_transactions"`
	TopProducts       []TopProduct `json:"top_products" firestore:"top_products"`
}

// TopProduct represents a top-selling product in reports
type TopProduct struct {
	ProductID    string  `json:"product_id" firestore:"product_id"`
	Name         string  `json:"name" firestore:"name"`
	TotalSold    int     `json:"total_sold" firestore:"total_sold"`
	TotalRevenue float64 `json:"total_revenue" firestore:"total_revenue"`
}

// NewDailyReport creates a daily report from transactions
func NewDailyReport(date time.Time, transactions []Transaction) *Report {
	report := &Report{
		ReportID:          generateReportID(date),
		Date:              date,
		TotalSales:        0,
		TotalTransactions: len(transactions),
		TopProducts:       make([]TopProduct, 0),
	}

	// Calculate total sales and analyze products
	productSales := make(map[string]*TopProduct)

	for _, transaction := range transactions {
		report.TotalSales += transaction.Total

		// Analyze items in transaction
		for _, item := range transaction.Items {
			if existing, exists := productSales[item.ProductID]; exists {
				existing.TotalSold += item.Quantity
				existing.TotalRevenue += item.Subtotal
			} else {
				productSales[item.ProductID] = &TopProduct{
					ProductID:    item.ProductID,
					Name:         item.Name,
					TotalSold:    item.Quantity,
					TotalRevenue: item.Subtotal,
				}
			}
		}
	}

	// Convert map to slice and sort by revenue
	for _, product := range productSales {
		report.TopProducts = append(report.TopProducts, *product)
	}

	// Sort by total revenue (descending)
	for i := 0; i < len(report.TopProducts)-1; i++ {
		for j := i + 1; j < len(report.TopProducts); j++ {
			if report.TopProducts[i].TotalRevenue < report.TopProducts[j].TotalRevenue {
				report.TopProducts[i], report.TopProducts[j] = report.TopProducts[j], report.TopProducts[i]
			}
		}
	}

	// Limit to top 10 products
	if len(report.TopProducts) > 10 {
		report.TopProducts = report.TopProducts[:10]
	}

	return report
}

// generateReportID generates a unique report ID based on date
func generateReportID(date time.Time) string {
	return fmt.Sprintf("report_%s", date.Format("20060102"))
}

// NewWeeklyReport creates a weekly report
func NewWeeklyReport(startDate time.Time, transactions []Transaction) *Report {
	report := NewDailyReport(startDate, transactions)
	report.ReportID = fmt.Sprintf("weekly_report_%s", startDate.Format("20060102"))
	return report
}

// NewMonthlyReport creates a monthly report
func NewMonthlyReport(month time.Time, transactions []Transaction) *Report {
	report := NewDailyReport(month, transactions)
	report.ReportID = fmt.Sprintf("monthly_report_%s", month.Format("200601"))
	return report
}

// calculateTopProducts calculates top-selling products from transactions
func calculateTopProducts(transactions []Transaction) []TopProduct {
	productStats := make(map[string]*TopProduct)

	// Calculate product statistics
	for _, trans := range transactions {
		for _, item := range trans.Items {
			productID := item.ProductID
			if stat, exists := productStats[productID]; exists {
				stat.TotalSold += item.Quantity
				stat.TotalRevenue += item.Subtotal
			} else {
				productStats[productID] = &TopProduct{
					ProductID:    productID,
					Name:         item.Name,
					TotalSold:    item.Quantity,
					TotalRevenue: item.Subtotal,
				}
			}
		}
	}

	// Convert to slice
	topProducts := make([]TopProduct, 0)
	for _, stat := range productStats {
		topProducts = append(topProducts, *stat)
	}

	return topProducts
}