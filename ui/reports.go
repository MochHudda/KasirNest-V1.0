package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"kasirnest/firebase"
	"kasirnest/models"
	"kasirnest/utils"
)

// ReportsScreen represents the reports interface
type ReportsScreen struct {
	window           fyne.Window
	container        *fyne.Container
	firebaseClient   *firebase.Client
	firestoreService *firebase.FirestoreService

	// Report filters
	dateFromEntry    *widget.Entry
	dateToEntry      *widget.Entry
	reportTypeSelect *widget.Select

	// Report display
	summaryCard    *widget.Card
	chartContainer *fyne.Container
	detailsTable   *widget.Table

	// Data
	currentReport *models.Report
	transactions  []models.Transaction
}

// NewReportsScreen creates a new reports screen
func NewReportsScreen(w fyne.Window, fbClient *firebase.Client) *ReportsScreen {
	screen := &ReportsScreen{
		window:           w,
		firebaseClient:   fbClient,
		firestoreService: firebase.NewFirestoreService(fbClient),
		transactions:     make([]models.Transaction, 0),
	}

	screen.setupUI()
	screen.generateTodayReport()
	return screen
}

// setupUI sets up the reports interface
func (r *ReportsScreen) setupUI() {
	// Create filter controls
	r.setupFilterControls()

	// Create summary card
	r.createSummaryCard()

	// Create chart container (placeholder)
	r.chartContainer = container.NewVBox(
		widget.NewCard("Grafik Penjualan", "",
			widget.NewLabel("Grafik akan ditampilkan di sini")),
	)

	// Create details table
	r.createDetailsTable()

	// Create main layout
	filtersContainer := container.NewHBox(
		widget.NewLabel("Dari:"),
		r.dateFromEntry,
		widget.NewLabel("Sampai:"),
		r.dateToEntry,
		widget.NewLabel("Jenis:"),
		r.reportTypeSelect,
		widget.NewButton("Generate", func() {
			r.generateReport()
		}),
		widget.NewButton("Export", func() {
			r.exportReport()
		}),
	)

	// Create top section with summary and chart
	topSection := container.NewHBox(
		r.summaryCard,
		r.chartContainer,
	)

	r.container = container.NewVBox(
		filtersContainer,
		widget.NewSeparator(),
		topSection,
		widget.NewSeparator(),
		widget.NewLabel("Detail Produk Terlaris:"),
		r.detailsTable,
	)
}

// setupFilterControls sets up the filter controls
func (r *ReportsScreen) setupFilterControls() {
	// Date from entry
	r.dateFromEntry = widget.NewEntry()
	r.dateFromEntry.SetPlaceHolder("DD/MM/YYYY")
	today := time.Now()
	r.dateFromEntry.SetText(utils.FormatDateShort(today))

	// Date to entry
	r.dateToEntry = widget.NewEntry()
	r.dateToEntry.SetPlaceHolder("DD/MM/YYYY")
	r.dateToEntry.SetText(utils.FormatDateShort(today))

	// Report type select
	reportTypes := []string{"Harian", "Mingguan", "Bulanan", "Kustom"}
	r.reportTypeSelect = widget.NewSelect(reportTypes, func(reportType string) {
		r.updateDateRange(reportType)
	})
	r.reportTypeSelect.SetSelected("Harian")
}

// updateDateRange updates date range based on report type
func (r *ReportsScreen) updateDateRange(reportType string) {
	now := time.Now()

	switch reportType {
	case "Harian":
		r.dateFromEntry.SetText(utils.FormatDateShort(now))
		r.dateToEntry.SetText(utils.FormatDateShort(now))
	case "Mingguan":
		weekStart := now.AddDate(0, 0, -int(now.Weekday()))
		r.dateFromEntry.SetText(utils.FormatDateShort(weekStart))
		r.dateToEntry.SetText(utils.FormatDateShort(now))
	case "Bulanan":
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		r.dateFromEntry.SetText(utils.FormatDateShort(monthStart))
		r.dateToEntry.SetText(utils.FormatDateShort(now))
	}
}

// createSummaryCard creates the summary statistics card
func (r *ReportsScreen) createSummaryCard() {
	r.summaryCard = widget.NewCard("Ringkasan Penjualan", "",
		container.NewVBox(
			widget.NewLabel("Total Penjualan: Rp 0"),
			widget.NewLabel("Jumlah Transaksi: 0"),
			widget.NewLabel("Rata-rata per Transaksi: Rp 0"),
			widget.NewLabel("Produk Terjual: 0"),
		))
}

// createDetailsTable creates the details table for top products
func (r *ReportsScreen) createDetailsTable() {
	r.detailsTable = widget.NewTable(
		func() (int, int) {
			if r.currentReport == nil {
				return 1, 4 // Header only
			}
			return len(r.currentReport.TopProducts) + 1, 4 // +1 for header
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			if id.Row == 0 {
				// Header row
				headers := []string{"Produk", "Terjual", "Pendapatan", "Persentase"}
				if id.Col < len(headers) {
					label.SetText(headers[id.Col])
					label.TextStyle = fyne.TextStyle{Bold: true}
				}
			} else if r.currentReport != nil {
				// Data rows
				if id.Row-1 < len(r.currentReport.TopProducts) {
					product := r.currentReport.TopProducts[id.Row-1]
					switch id.Col {
					case 0:
						label.SetText(product.Name)
					case 1:
						label.SetText(fmt.Sprintf("%d", product.TotalSold))
					case 2:
						label.SetText(utils.FormatCurrency(product.TotalRevenue))
					case 3:
						percentage := 0.0
						if r.currentReport.TotalSales > 0 {
							percentage = (product.TotalRevenue / r.currentReport.TotalSales) * 100
						}
						label.SetText(utils.FormatPercentage(percentage))
					}
				}
			}
		},
	)

	// Set column widths
	r.detailsTable.SetColumnWidth(0, 200) // Product
	r.detailsTable.SetColumnWidth(1, 80)  // Sold
	r.detailsTable.SetColumnWidth(2, 100) // Revenue
	r.detailsTable.SetColumnWidth(3, 80)  // Percentage
}

// generateReport generates report based on selected criteria
func (r *ReportsScreen) generateReport() {
	// Parse date range
	dateFrom, err := r.parseDate(r.dateFromEntry.Text)
	if err != nil {
		widget.NewLabel("Format tanggal tidak valid").Show()
		return
	}

	dateTo, err := r.parseDate(r.dateToEntry.Text)
	if err != nil {
		widget.NewLabel("Format tanggal tidak valid").Show()
		return
	}

	// Load transactions for the date range
	r.loadTransactionsForDateRange(dateFrom, dateTo)

	// Generate report
	r.currentReport = models.NewDailyReport(dateFrom, r.transactions)

	// Update UI
	r.updateSummaryCard()
	r.detailsTable.Refresh()
}

// generateTodayReport generates today's report
func (r *ReportsScreen) generateTodayReport() {
	today := time.Now()
	r.loadTransactionsForDateRange(today, today)
	r.currentReport = models.NewDailyReport(today, r.transactions)
	r.updateSummaryCard()
	r.detailsTable.Refresh()
}

// updateSummaryCard updates the summary card with current report data
func (r *ReportsScreen) updateSummaryCard() {
	if r.currentReport == nil {
		return
	}

	avgPerTransaction := 0.0
	if r.currentReport.TotalTransactions > 0 {
		avgPerTransaction = r.currentReport.TotalSales / float64(r.currentReport.TotalTransactions)
	}

	totalProductsSold := 0
	for _, product := range r.currentReport.TopProducts {
		totalProductsSold += product.TotalSold
	}

	content := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Total Penjualan: %s", utils.FormatCurrency(r.currentReport.TotalSales))),
		widget.NewLabel(fmt.Sprintf("Jumlah Transaksi: %d", r.currentReport.TotalTransactions)),
		widget.NewLabel(fmt.Sprintf("Rata-rata per Transaksi: %s", utils.FormatCurrency(avgPerTransaction))),
		widget.NewLabel(fmt.Sprintf("Produk Terjual: %d", totalProductsSold)),
	)

	r.summaryCard = widget.NewCard("Ringkasan Penjualan", "", content)
}

// loadTransactionsForDateRange loads transactions for specified date range
func (r *ReportsScreen) loadTransactionsForDateRange(dateFrom, dateTo time.Time) {
	// This is a simplified implementation
	// In a real app, you'd query Firestore with date filters
	r.transactions = make([]models.Transaction, 0)

	// Create some dummy data for demonstration
	if dateFrom.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		// Create sample transaction for today
		sampleTransaction := models.Transaction{
			TransID:       "sample_001",
			Date:          time.Now(),
			Total:         50000,
			PaymentMethod: models.PaymentCash,
			Items: []models.TransactionItem{
				{
					ProductID: "prod_001",
					Name:      "Produk Demo A",
					Quantity:  2,
					Price:     15000,
					Subtotal:  30000,
				},
				{
					ProductID: "prod_002",
					Name:      "Produk Demo B",
					Quantity:  1,
					Price:     20000,
					Subtotal:  20000,
				},
			},
		}
		r.transactions = append(r.transactions, sampleTransaction)
	}
}

// parseDate parses date string in DD/MM/YYYY format
func (r *ReportsScreen) parseDate(dateStr string) (time.Time, error) {
	return time.Parse("02/01/2006", dateStr)
}

// exportReport exports the current report
func (r *ReportsScreen) exportReport() {
	if r.currentReport == nil {
		widget.NewLabel("Tidak ada data untuk diekspor").Show()
		return
	}

	// This is a placeholder for export functionality
	// In a real app, you'd implement CSV/PDF export
	widget.NewLabel("Fitur export akan diimplementasikan").Show()
}

// GetContainer returns the reports container
func (r *ReportsScreen) GetContainer() *fyne.Container {
	return r.container
}

// Refresh refreshes the reports data
func (r *ReportsScreen) Refresh() {
	r.generateReport()
}
