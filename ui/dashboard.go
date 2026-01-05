package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"kasirnest/firebase"
)

// DashboardScreen represents the main dashboard
type DashboardScreen struct {
	window             fyne.Window
	container          *fyne.Container
	content            *container.DocTabs
	firebaseClient     *firebase.Client
	firestoreService   *firebase.FirestoreService
	productsScreen     *ProductsScreen
	transactionsScreen *TransactionsScreen
	reportsScreen      *ReportsScreen
}

// NewDashboardScreen creates a new dashboard screen
func NewDashboardScreen(w fyne.Window, fbClient *firebase.Client) *DashboardScreen {
	dashboard := &DashboardScreen{
		window:           w,
		firebaseClient:   fbClient,
		firestoreService: firebase.NewFirestoreService(fbClient),
	}

	dashboard.setupUI()
	return dashboard
}

// setupUI sets up the dashboard interface
func (d *DashboardScreen) setupUI() {
	// Get current user
	_, email, name := GetCurrentUser()

	// Create welcome card
	welcomeCard := d.createWelcomeCard(name, email)

	// Create quick stats card
	statsCard := d.createStatsCard()

	// Create top container with welcome and stats
	topContainer := container.NewHBox(
		welcomeCard,
		statsCard,
	)

	// Create tabs for different modules
	d.content = container.NewDocTabs()

	// Create and add screens
	d.productsScreen = NewProductsScreen(d.window, d.firebaseClient)
	d.transactionsScreen = NewTransactionsScreen(d.window, d.firebaseClient)
	d.reportsScreen = NewReportsScreen(d.window, d.firebaseClient)

	// Create simple container instead of tabs for now
	dashboardContent := d.createDashboardContent()
	d.container = container.NewBorder(
		topContainer,
		nil,
		nil,
		nil,
		dashboardContent,
	)

	// Create toolbar
	toolbar := d.createToolbar()

	// Create main container
	d.container = container.NewBorder(
		toolbar,   // top
		nil,       // bottom
		nil,       // left
		nil,       // right
		d.content, // center
	)
}

// createWelcomeCard creates the welcome card
func (d *DashboardScreen) createWelcomeCard(name, email string) *widget.Card {
	if name == "" {
		name = "Pengguna"
	}

	welcomeLabel := widget.NewLabel(fmt.Sprintf("Selamat datang, %s", name))
	welcomeLabel.TextStyle = fyne.TextStyle{Bold: true}

	emailLabel := widget.NewLabel(email)

	content := container.NewVBox(
		welcomeLabel,
		emailLabel,
		widget.NewLabel("Sistem siap digunakan"),
	)

	return widget.NewCard("Selamat Datang", "", content)
}

// createStatsCard creates the quick stats card
func (d *DashboardScreen) createStatsCard() *widget.Card {
	// Create stats labels
	totalProductsLabel := widget.NewLabel("Loading...")
	totalTransactionsLabel := widget.NewLabel("Loading...")
	todaysSalesLabel := widget.NewLabel("Loading...")

	// Load stats asynchronously
	go d.loadStats(totalProductsLabel, totalTransactionsLabel, todaysSalesLabel)

	content := container.NewVBox(
		container.NewHBox(widget.NewLabel("Total Produk:"), totalProductsLabel),
		container.NewHBox(widget.NewLabel("Total Transaksi:"), totalTransactionsLabel),
		container.NewHBox(widget.NewLabel("Penjualan Hari Ini:"), todaysSalesLabel),
	)

	return widget.NewCard("Statistik Cepat", "", content)
}

// createToolbar creates the toolbar
func (d *DashboardScreen) createToolbar() *widget.Toolbar {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			d.showSettings()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.LogoutIcon(), func() {
			d.handleLogout()
		}),
	)

	return toolbar
}

// createDashboardContent creates the main dashboard content
func (d *DashboardScreen) createDashboardContent() *fyne.Container {
	// Create quick action buttons
	newTransactionBtn := widget.NewButton("Transaksi Baru", func() {
		d.content.SelectIndex(2) // Select transactions tab
		d.transactionsScreen.StartNewTransaction()
	})
	newTransactionBtn.Importance = widget.HighImportance

	addProductBtn := widget.NewButton("Tambah Produk", func() {
		d.content.SelectIndex(1) // Select products tab
		d.productsScreen.ShowAddProductDialog()
	})

	viewReportsBtn := widget.NewButton("Lihat Laporan", func() {
		d.content.SelectIndex(3) // Select reports tab
	})

	// Create quick actions card
	quickActionsCard := widget.NewCard("Aksi Cepat", "",
		container.NewVBox(
			newTransactionBtn,
			addProductBtn,
			viewReportsBtn,
		))

	// Create recent activities (placeholder)
	recentActivities := widget.NewCard("Aktivitas Terbaru", "",
		widget.NewLabel("Belum ada aktivitas terbaru"))

	return container.NewHBox(
		quickActionsCard,
		recentActivities,
	)
}

// loadStats loads statistics data
func (d *DashboardScreen) loadStats(productsLabel, transactionsLabel, salesLabel *widget.Label) {
	// Load product count
	productCount, err := d.firestoreService.GetCollectionSize("products")
	if err != nil {
		productsLabel.SetText("Error")
	} else {
		productsLabel.SetText(fmt.Sprintf("%d", productCount))
	}

	// Load transaction count
	transactionCount, err := d.firestoreService.GetCollectionSize("transactions")
	if err != nil {
		transactionsLabel.SetText("Error")
	} else {
		transactionsLabel.SetText(fmt.Sprintf("%d", transactionCount))
	}

	// Load today's sales (placeholder)
	salesLabel.SetText("Rp 0")
}

// showSettings shows the settings dialog
func (d *DashboardScreen) showSettings() {
	// Create settings form
	themeSelect := widget.NewSelect([]string{"Light", "Dark"}, nil)
	themeSelect.SetSelected("Light")

	windowSizeEntry := widget.NewEntry()
	windowSizeEntry.SetText("1200x800")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Theme:", Widget: themeSelect},
			{Text: "Window Size:", Widget: windowSizeEntry},
		},
	}

	// Show dialog
	dialog.ShowForm("Pengaturan", "Simpan", "Batal", form.Items, func(confirm bool) {
		if confirm {
			// Apply settings
			d.applySettings(themeSelect.Selected, windowSizeEntry.Text)
		}
	}, d.window)
}

// applySettings applies the settings
func (d *DashboardScreen) applySettings(theme, windowSize string) {
	// Apply theme
	if theme == "Dark" {
		fyne.CurrentApp().Settings().SetTheme(&darkTheme{})
	} else {
		fyne.CurrentApp().Settings().SetTheme(&lightTheme{})
	}

	// Apply window size (simplified)
	if windowSize != "" {
		// Parse and apply window size
		// This is a simplified implementation
	}
}

// handleLogout handles user logout
func (d *DashboardScreen) handleLogout() {
	dialog.ShowConfirm("Keluar", "Apakah Anda yakin ingin keluar?", func(confirm bool) {
		if confirm {
			ClearSession()

			// Reset window content to login screen
			// This would be handled by the main application
			d.onLogout()
		}
	}, d.window)
}

// onLogout callback for logout (to be set by main app)
var logoutCallback func()

func (d *DashboardScreen) onLogout() {
	if logoutCallback != nil {
		logoutCallback()
	}
}

// SetLogoutCallback sets the logout callback
func (d *DashboardScreen) SetLogoutCallback(callback func()) {
	logoutCallback = callback
}

// GetContainer returns the dashboard container
func (d *DashboardScreen) GetContainer() *fyne.Container {
	return d.container
}

// Refresh refreshes the dashboard data
func (d *DashboardScreen) Refresh() {
	// Refresh all child screens
	if d.productsScreen != nil {
		d.productsScreen.Refresh()
	}
	if d.transactionsScreen != nil {
		d.transactionsScreen.Refresh()
	}
	if d.reportsScreen != nil {
		d.reportsScreen.Refresh()
	}
}

// Simple theme implementations
type lightTheme struct{}

func (t *lightTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}
func (t *lightTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}
func (t *lightTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (t *lightTheme) Size(name fyne.ThemeSizeName) float32 { return theme.DefaultTheme().Size(name) }

type darkTheme struct{}

func (t *darkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, theme.VariantDark)
}
func (t *darkTheme) Font(style fyne.TextStyle) fyne.Resource { return theme.DefaultTheme().Font(style) }
func (t *darkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (t *darkTheme) Size(name fyne.ThemeSizeName) float32 { return theme.DefaultTheme().Size(name) }
