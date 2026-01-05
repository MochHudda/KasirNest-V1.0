package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"kasirnest/firebase"
	"kasirnest/models"
	"kasirnest/utils"
)

// TransactionsScreen represents the transactions interface
type TransactionsScreen struct {
	window           fyne.Window
	container        *fyne.Container
	firebaseClient   *firebase.Client
	firestoreService *firebase.FirestoreService
	tabs             *container.DocTabs

	// POS (New Transaction) tab
	posContainer       *fyne.Container
	productSearch      *widget.Entry
	cartTable          *widget.Table
	totalLabel         *widget.Label
	currentTransaction *models.Transaction

	// Transaction History tab
	historyContainer *fyne.Container
	historyTable     *widget.Table
	transactions     []models.Transaction
}

// NewTransactionsScreen creates a new transactions screen
func NewTransactionsScreen(w fyne.Window, fbClient *firebase.Client) *TransactionsScreen {
	screen := &TransactionsScreen{
		window:           w,
		firebaseClient:   fbClient,
		firestoreService: firebase.NewFirestoreService(fbClient),
		transactions:     make([]models.Transaction, 0),
	}

	screen.setupUI()
	screen.loadTransactions()
	return screen
}

// setupUI sets up the transactions interface
func (t *TransactionsScreen) setupUI() {
	// Create tabs
	t.tabs = container.NewDocTabs()

	// Create POS tab
	t.setupPOSTab()

	// Create history tab
	t.setupHistoryTab()

	// Create simple container instead of tabs
	t.container = container.NewBorder(nil, nil, nil, nil, t.posContainer)
}

// setupPOSTab sets up the POS (Point of Sale) tab
func (t *TransactionsScreen) setupPOSTab() {
	// Initialize new transaction
	t.currentTransaction = &models.Transaction{
		TransID:       fmt.Sprintf("trans_%d", time.Now().Unix()),
		UserID:        "",
		Date:          time.Now(),
		Total:         0,
		PaymentMethod: models.PaymentCash,
		Items:         make([]models.TransactionItem, 0),
	}

	// Get current user
	userID, _, _ := GetCurrentUser()
	t.currentTransaction.UserID = userID

	// Create product search
	t.productSearch = widget.NewEntry()
	t.productSearch.SetPlaceHolder("Cari produk atau scan barcode...")
	t.productSearch.OnSubmitted = func(text string) {
		t.searchAndAddProduct(text)
	}

	searchButton := widget.NewButton("Cari", func() {
		t.searchAndAddProduct(t.productSearch.Text)
	})

	// Create cart table
	t.createCartTable()

	// Create total display
	t.totalLabel = widget.NewLabel("Total: Rp 0")
	t.totalLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Create action buttons
	clearButton := widget.NewButton("Clear", func() {
		t.clearTransaction()
	})

	paymentMethods := []string{models.PaymentCash, models.PaymentCard, models.PaymentDigital}
	paymentSelect := widget.NewSelect(paymentMethods, func(method string) {
		t.currentTransaction.PaymentMethod = method
	})
	paymentSelect.SetSelected(models.PaymentCash)

	processButton := widget.NewButton("Proses Pembayaran", func() {
		t.processPayment()
	})
	processButton.Importance = widget.HighImportance

	// Create product search container
	searchContainer := container.NewHBox(
		widget.NewLabel("Produk:"),
		t.productSearch,
		searchButton,
	)

	// Create bottom actions
	actionsContainer := container.NewHBox(
		clearButton,
		widget.NewLabel("Metode Bayar:"),
		paymentSelect,
		widget.NewSeparator(),
		t.totalLabel,
		processButton,
	)

	// Create POS container
	t.posContainer = container.NewVBox(
		searchContainer,
		widget.NewSeparator(),
		widget.NewLabel("Keranjang Belanja:"),
		t.cartTable,
		widget.NewSeparator(),
		actionsContainer,
	)
}

// setupHistoryTab sets up the transaction history tab
func (t *TransactionsScreen) setupHistoryTab() {
	// Create filter controls
	dateFromEntry := widget.NewEntry()
	dateFromEntry.SetPlaceHolder("DD/MM/YYYY")

	dateToEntry := widget.NewEntry()
	dateToEntry.SetPlaceHolder("DD/MM/YYYY")

	filterButton := widget.NewButton("Filter", func() {
		// Implement date filtering
		t.filterTransactions(dateFromEntry.Text, dateToEntry.Text)
	})

	refreshButton := widget.NewButton("Refresh", func() {
		t.loadTransactions()
	})

	// Create filter container
	filterContainer := container.NewHBox(
		widget.NewLabel("Dari:"),
		dateFromEntry,
		widget.NewLabel("Sampai:"),
		dateToEntry,
		filterButton,
		refreshButton,
	)

	// Create history table
	t.createHistoryTable()

	// Create history container
	t.historyContainer = container.NewVBox(
		filterContainer,
		widget.NewSeparator(),
		t.historyTable,
	)
}

// createCartTable creates the shopping cart table
func (t *TransactionsScreen) createCartTable() {
	t.cartTable = widget.NewTable(
		func() (int, int) {
			return len(t.currentTransaction.Items) + 1, 6 // +1 for header, 6 columns
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			if id.Row == 0 {
				// Header row
				headers := []string{"Produk", "Harga", "Qty", "Subtotal", "", ""}
				if id.Col < len(headers) {
					label.SetText(headers[id.Col])
					label.TextStyle = fyne.TextStyle{Bold: true}
				}
			} else {
				// Data rows
				if id.Row-1 < len(t.currentTransaction.Items) {
					item := t.currentTransaction.Items[id.Row-1]
					switch id.Col {
					case 0:
						label.SetText(item.Name)
					case 1:
						label.SetText(utils.FormatCurrency(item.Price))
					case 2:
						label.SetText(fmt.Sprintf("%d", item.Quantity))
					case 3:
						label.SetText(utils.FormatCurrency(item.Subtotal))
					case 4:
						// Edit quantity button placeholder
						label.SetText("Edit")
					case 5:
						// Remove button placeholder
						label.SetText("Hapus")
					}
				}
			}
		},
	)

	// Set column widths
	t.cartTable.SetColumnWidth(0, 200) // Product
	t.cartTable.SetColumnWidth(1, 80)  // Price
	t.cartTable.SetColumnWidth(2, 50)  // Qty
	t.cartTable.SetColumnWidth(3, 100) // Subtotal
	t.cartTable.SetColumnWidth(4, 60)  // Edit
	t.cartTable.SetColumnWidth(5, 60)  // Remove
}

// createHistoryTable creates the transaction history table
func (t *TransactionsScreen) createHistoryTable() {
	t.historyTable = widget.NewTable(
		func() (int, int) {
			return len(t.transactions) + 1, 5 // +1 for header, 5 columns
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			if id.Row == 0 {
				// Header row
				headers := []string{"ID", "Tanggal", "Total", "Metode", "Item"}
				if id.Col < len(headers) {
					label.SetText(headers[id.Col])
					label.TextStyle = fyne.TextStyle{Bold: true}
				}
			} else {
				// Data rows
				if id.Row-1 < len(t.transactions) {
					trans := t.transactions[id.Row-1]
					switch id.Col {
					case 0:
						label.SetText(trans.TransID)
					case 1:
						label.SetText(utils.FormatDateTimeShort(trans.Date))
					case 2:
						label.SetText(utils.FormatCurrency(trans.Total))
					case 3:
						label.SetText(trans.PaymentMethod)
					case 4:
						label.SetText(fmt.Sprintf("%d item", len(trans.Items)))
					}
				}
			}
		},
	)

	// Set column widths
	t.historyTable.SetColumnWidth(0, 120) // ID
	t.historyTable.SetColumnWidth(1, 120) // Date
	t.historyTable.SetColumnWidth(2, 100) // Total
	t.historyTable.SetColumnWidth(3, 80)  // Method
	t.historyTable.SetColumnWidth(4, 60)  // Items
}

// searchAndAddProduct searches and adds product to cart
func (t *TransactionsScreen) searchAndAddProduct(query string) {
	if query == "" {
		return
	}

	// This is a simplified implementation
	// In a real app, you'd search products from Firestore
	// For now, we'll create a dummy product
	product := &models.Product{
		ProductID: "dummy",
		Name:      "Produk Demo - " + query,
		Price:     10000,
		Stock:     100,
		Category:  "Lainnya",
	}

	// Add to transaction
	err := t.currentTransaction.AddItem(product, 1)
	if err != nil {
		dialog.ShowError(err, t.window)
		return
	}

	// Update UI
	t.updateCartUI()
	t.productSearch.SetText("")
}

// updateCartUI updates the cart UI
func (t *TransactionsScreen) updateCartUI() {
	t.cartTable.Refresh()
	t.totalLabel.SetText("Total: " + utils.FormatCurrency(t.currentTransaction.Total))
}

// clearTransaction clears the current transaction
func (t *TransactionsScreen) clearTransaction() {
	dialog.ShowConfirm("Clear Transaksi", "Apakah Anda yakin ingin menghapus semua item?", func(confirm bool) {
		if confirm {
			t.currentTransaction.Items = make([]models.TransactionItem, 0)
			t.currentTransaction.Total = 0
			t.updateCartUI()
		}
	}, t.window)
}

// processPayment processes the payment
func (t *TransactionsScreen) processPayment() {
	if len(t.currentTransaction.Items) == 0 {
		dialog.ShowInformation("Keranjang Kosong", "Tidak ada item dalam keranjang", t.window)
		return
	}

	// Show payment confirmation
	message := fmt.Sprintf("Total: %s\nMetode: %s\nProses pembayaran?",
		utils.FormatCurrency(t.currentTransaction.Total),
		t.currentTransaction.PaymentMethod)

	dialog.ShowConfirm("Konfirmasi Pembayaran", message, func(confirm bool) {
		if confirm {
			t.saveTransaction()
		}
	}, t.window)
}

// saveTransaction saves the transaction
func (t *TransactionsScreen) saveTransaction() {
	// Save to Firestore
	err := t.firestoreService.Create("transactions", t.currentTransaction.TransID, *t.currentTransaction)
	if err != nil {
		dialog.ShowError(fmt.Errorf("gagal menyimpan transaksi: %v", err), t.window)
		return
	}

	dialog.ShowInformation("Sukses", "Transaksi berhasil diproses", t.window)

	// Reset transaction
	t.StartNewTransaction()

	// Refresh history
	t.loadTransactions()
}

// StartNewTransaction starts a new transaction
func (t *TransactionsScreen) StartNewTransaction() {
	userID, _, _ := GetCurrentUser()
	t.currentTransaction = &models.Transaction{
		TransID:       fmt.Sprintf("trans_%d", time.Now().Unix()),
		UserID:        userID,
		Date:          time.Now(),
		Total:         0,
		PaymentMethod: models.PaymentCash,
		Items:         make([]models.TransactionItem, 0),
	}
	t.updateCartUI()
}

// loadTransactions loads transaction history
func (t *TransactionsScreen) loadTransactions() {
	// This is a simplified implementation
	// In a real app, you'd load transactions from Firestore
	t.transactions = make([]models.Transaction, 0)

	if t.historyTable != nil {
		t.historyTable.Refresh()
	}
}

// filterTransactions filters transactions by date range
func (t *TransactionsScreen) filterTransactions(dateFrom, dateTo string) {
	// Implement date filtering logic
	// For now, just reload all transactions
	t.loadTransactions()
}

// GetContainer returns the transactions container
func (t *TransactionsScreen) GetContainer() *fyne.Container {
	return t.container
}

// Refresh refreshes the transactions data
func (t *TransactionsScreen) Refresh() {
	t.loadTransactions()
}
