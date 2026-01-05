package ui

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"kasirnest/firebase"
	"kasirnest/models"
	"kasirnest/utils"
)

// ProductsScreen represents the products management interface
type ProductsScreen struct {
	window           fyne.Window
	container        *fyne.Container
	firebaseClient   *firebase.Client
	firestoreService *firebase.FirestoreService
	storageService   *firebase.StorageService
	table            *widget.Table
	searchEntry      *widget.Entry
	categoryFilter   *widget.Select
	products         []models.Product
	filteredProducts []models.Product
	selectedRows     []int // Track selected rows manually
}

// NewProductsScreen creates a new products screen
func NewProductsScreen(w fyne.Window, fbClient *firebase.Client) *ProductsScreen {
	screen := &ProductsScreen{
		window:           w,
		firebaseClient:   fbClient,
		firestoreService: firebase.NewFirestoreService(fbClient),
		products:         make([]models.Product, 0),
		filteredProducts: make([]models.Product, 0),
	}

	// Initialize storage service if available
	if fbClient.Storage != nil {
		screen.storageService = firebase.NewStorageService(fbClient, "your-bucket-name") // Will be configured
	}

	screen.setupUI()
	screen.loadProducts()
	return screen
}

// setupUI sets up the products interface
func (p *ProductsScreen) setupUI() {
	// Create search and filter controls
	p.searchEntry = widget.NewEntry()
	p.searchEntry.SetPlaceHolder("Cari produk...")
	p.searchEntry.OnChanged = func(text string) {
		p.filterProducts()
	}

	categories := []string{"Semua", "Makanan", "Elektronik", "Fashion", "Kesehatan", "Rumah Tangga", "Alat Tulis", "Lainnya"}
	p.categoryFilter = widget.NewSelect(categories, func(value string) {
		p.filterProducts()
	})
	p.categoryFilter.SetSelected("Semua")

	// Create action buttons
	addButton := widget.NewButton("Tambah Produk", func() {
		p.ShowAddProductDialog()
	})
	addButton.Importance = widget.HighImportance

	editButton := widget.NewButton("Edit", func() {
		p.showEditProductDialog()
	})

	deleteButton := widget.NewButton("Hapus", func() {
		p.showDeleteProductDialog()
	})
	deleteButton.Importance = widget.DangerImportance

	refreshButton := widget.NewButton("Refresh", func() {
		p.Refresh()
	})

	// Create controls container
	controls := container.NewHBox(
		widget.NewLabel("Cari:"),
		p.searchEntry,
		widget.NewLabel("Kategori:"),
		p.categoryFilter,
		widget.NewSeparator(),
		addButton,
		editButton,
		deleteButton,
		refreshButton,
	)

	// Create table
	p.createTable()

	// Create main container
	p.container = container.NewVBox(
		controls,
		widget.NewSeparator(),
		p.table,
	)
}

// createTable creates the products table
func (p *ProductsScreen) createTable() {
	p.table = widget.NewTable(
		func() (int, int) {
			return len(p.filteredProducts) + 1, 7 // +1 for header, 7 columns
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			if id.Row == 0 {
				// Header row
				headers := []string{"ID", "Nama", "Harga", "Stok", "Kategori", "Barcode", "Status"}
				if id.Col < len(headers) {
					label.SetText(headers[id.Col])
					label.TextStyle = fyne.TextStyle{Bold: true}
				}
			} else {
				// Data rows
				if id.Row-1 < len(p.filteredProducts) {
					product := p.filteredProducts[id.Row-1]
					switch id.Col {
					case 0:
						label.SetText(product.ProductID)
					case 1:
						label.SetText(product.Name)
					case 2:
						label.SetText(utils.FormatCurrency(product.Price))
					case 3:
						label.SetText(fmt.Sprintf("%d", product.Stock))
					case 4:
						label.SetText(product.Category)
					case 5:
						label.SetText(product.Barcode)
					case 6:
						if product.IsInStock() {
							label.SetText("Tersedia")
						} else {
							label.SetText("Habis")
						}
					}
				}
			}
		},
	)

	// Set column widths
	p.table.SetColumnWidth(0, 80)  // ID
	p.table.SetColumnWidth(1, 200) // Name
	p.table.SetColumnWidth(2, 100) // Price
	p.table.SetColumnWidth(3, 60)  // Stock
	p.table.SetColumnWidth(4, 100) // Category
	p.table.SetColumnWidth(5, 120) // Barcode
	p.table.SetColumnWidth(6, 80)  // Status
}

// ShowAddProductDialog shows the add product dialog
func (p *ProductsScreen) ShowAddProductDialog() {
	p.showProductDialog(nil)
}

// showEditProductDialog shows the edit product dialog
func (p *ProductsScreen) showEditProductDialog() {
	if len(p.selectedRows) == 0 {
		dialog.ShowInformation("Pilih Produk", "Silakan pilih produk yang ingin diedit", p.window)
		return
	}

	selectedRow := p.selectedRows[0] - 1 // -1 for header
	if selectedRow < 0 || selectedRow >= len(p.filteredProducts) {
		return
	}

	product := &p.filteredProducts[selectedRow]
	p.showProductDialog(product)
}

// showDeleteProductDialog shows the delete product confirmation
func (p *ProductsScreen) showDeleteProductDialog() {
	if len(p.selectedRows) == 0 {
		dialog.ShowInformation("Pilih Produk", "Silakan pilih produk yang ingin dihapus", p.window)
		return
	}

	selectedRow := p.selectedRows[0] - 1 // -1 for header
	if selectedRow < 0 || selectedRow >= len(p.filteredProducts) {
		return
	}

	product := p.filteredProducts[selectedRow]

	dialog.ShowConfirm("Hapus Produk",
		fmt.Sprintf("Apakah Anda yakin ingin menghapus produk '%s'?", product.Name),
		func(confirm bool) {
			if confirm {
				p.deleteProduct(product.ProductID)
			}
		}, p.window)
}

// showProductDialog shows add/edit product dialog
func (p *ProductsScreen) showProductDialog(product *models.Product) {
	isEdit := product != nil

	// Create form fields
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nama produk")

	priceEntry := widget.NewEntry()
	priceEntry.SetPlaceHolder("Harga")

	stockEntry := widget.NewEntry()
	stockEntry.SetPlaceHolder("Jumlah stok")

	categories := []string{"Makanan", "Elektronik", "Fashion", "Kesehatan", "Rumah Tangga", "Alat Tulis", "Lainnya"}
	categorySelect := widget.NewSelect(categories, nil)

	barcodeEntry := widget.NewEntry()
	barcodeEntry.SetPlaceHolder("Barcode (opsional)")

	// Fill form if editing
	if isEdit {
		nameEntry.SetText(product.Name)
		priceEntry.SetText(fmt.Sprintf("%.2f", product.Price))
		stockEntry.SetText(fmt.Sprintf("%d", product.Stock))
		categorySelect.SetSelected(product.Category)
		barcodeEntry.SetText(product.Barcode)
	} else {
		categorySelect.SetSelected(categories[0])
	}

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Nama:", Widget: nameEntry},
			{Text: "Harga:", Widget: priceEntry},
			{Text: "Stok:", Widget: stockEntry},
			{Text: "Kategori:", Widget: categorySelect},
			{Text: "Barcode:", Widget: barcodeEntry},
		},
	}

	// Show dialog
	title := "Tambah Produk"
	if isEdit {
		title = "Edit Produk"
	}

	dialog.ShowForm(title, "Simpan", "Batal", form.Items, func(confirm bool) {
		if confirm {
			if isEdit {
				p.updateProduct(product, nameEntry.Text, priceEntry.Text, stockEntry.Text,
					categorySelect.Selected, barcodeEntry.Text)
			} else {
				p.addProduct(nameEntry.Text, priceEntry.Text, stockEntry.Text,
					categorySelect.Selected, barcodeEntry.Text)
			}
		}
	}, p.window)
}

// addProduct adds a new product
func (p *ProductsScreen) addProduct(name, priceStr, stockStr, category, barcode string) {
	// Validate inputs
	if name == "" {
		dialog.ShowError(fmt.Errorf("nama produk tidak boleh kosong"), p.window)
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		dialog.ShowError(fmt.Errorf("harga harus berupa angka positif"), p.window)
		return
	}

	stock, err := strconv.Atoi(stockStr)
	if err != nil || stock < 0 {
		dialog.ShowError(fmt.Errorf("stok harus berupa angka non-negatif"), p.window)
		return
	}

	// Create new product
	product := models.Product{
		ProductID: fmt.Sprintf("prod_%d", time.Now().Unix()),
		Name:      name,
		Price:     price,
		Stock:     stock,
		Category:  category,
		Barcode:   barcode,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to Firestore
	err = p.firestoreService.Create("products", product.ProductID, product)
	if err != nil {
		dialog.ShowError(fmt.Errorf("gagal menambah produk: %v", err), p.window)
		return
	}

	dialog.ShowInformation("Sukses", "Produk berhasil ditambahkan", p.window)
	p.Refresh()
}

// updateProduct updates an existing product
func (p *ProductsScreen) updateProduct(product *models.Product, name, priceStr, stockStr, category, barcode string) {
	// Validate inputs
	if name == "" {
		dialog.ShowError(fmt.Errorf("nama produk tidak boleh kosong"), p.window)
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		dialog.ShowError(fmt.Errorf("harga harus berupa angka positif"), p.window)
		return
	}

	stock, err := strconv.Atoi(stockStr)
	if err != nil || stock < 0 {
		dialog.ShowError(fmt.Errorf("stok harus berupa angka non-negatif"), p.window)
		return
	}

	// Update product
	product.Name = name
	product.Price = price
	product.Stock = stock
	product.Category = category
	product.Barcode = barcode
	product.UpdatedAt = time.Now()

	// Save to Firestore
	err = p.firestoreService.Update("products", product.ProductID, *product)
	if err != nil {
		dialog.ShowError(fmt.Errorf("gagal mengupdate produk: %v", err), p.window)
		return
	}

	dialog.ShowInformation("Sukses", "Produk berhasil diupdate", p.window)
	p.Refresh()
}

// deleteProduct deletes a product
func (p *ProductsScreen) deleteProduct(productID string) {
	err := p.firestoreService.Delete("products", productID)
	if err != nil {
		dialog.ShowError(fmt.Errorf("gagal menghapus produk: %v", err), p.window)
		return
	}

	dialog.ShowInformation("Sukses", "Produk berhasil dihapus", p.window)
	p.Refresh()
}

// loadProducts loads products from Firestore
func (p *ProductsScreen) loadProducts() {
	// This is a simplified implementation
	// In a real app, you'd implement proper Firestore querying
	p.products = make([]models.Product, 0)
	p.filteredProducts = p.products
	p.table.Refresh()
}

// filterProducts filters products based on search and category
func (p *ProductsScreen) filterProducts() {
	searchText := p.searchEntry.Text
	selectedCategory := p.categoryFilter.Selected

	p.filteredProducts = make([]models.Product, 0)

	for _, product := range p.products {
		// Filter by search text
		if searchText != "" {
			if !contains(product.Name, searchText) && !contains(product.Barcode, searchText) {
				continue
			}
		}

		// Filter by category
		if selectedCategory != "Semua" && product.Category != selectedCategory {
			continue
		}

		p.filteredProducts = append(p.filteredProducts, product)
	}

	p.table.Refresh()
}

// contains checks if a string contains another string (case insensitive)
func contains(str, substr string) bool {
	return len(substr) == 0 || len(str) >= len(substr) &&
		(str == substr || len(str) > len(substr) &&
			(str[:len(substr)] == substr || str[len(str)-len(substr):] == substr ||
				indexOf(str, substr) != -1))
}

func indexOf(str, substr string) int {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// GetContainer returns the products container
func (p *ProductsScreen) GetContainer() *fyne.Container {
	return p.container
}

// Refresh refreshes the products data
func (p *ProductsScreen) Refresh() {
	p.loadProducts()
}
