package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"kasirnest/config"
	"kasirnest/firebase"
	"kasirnest/ui"
)

// Application represents the main application
type Application struct {
	fyneApp        fyne.App
	window         fyne.Window
	config         *config.Config
	firebaseClient *firebase.Client

	// Screens
	loginScreen     *ui.LoginScreen
	dashboardScreen *ui.DashboardScreen

	// Current state
	isLoggedIn bool
}

func main() {
	// Create the application
	app := &Application{}

	// Initialize the application
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Show the window and run
	app.window.ShowAndRun()
}

// Initialize initializes the application
func (a *Application) Initialize() error {
	// Create Fyne app
	a.fyneApp = app.NewWithID("com.kasirnest.pos")

	// Set app metadata (if method exists)
	// Note: SetMetadata may not exist in all Fyne versions
	// a.fyneApp.SetMetadata(&fyne.AppMetadata{
	//	ID:      "com.kasirnest.pos",
	//	Name:    "KasirNest",
	//	Version: "1.0.0",
	//	Icon:    nil,
	// })

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		// If config file doesn't exist, show setup dialog
		return a.showSetupDialog(err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return a.showConfigurationError(err)
	}

	a.config = cfg

	// Initialize Firebase
	if err := a.initializeFirebase(); err != nil {
		return a.showFirebaseError(err)
	}

	// Create main window
	a.createMainWindow()

	// Check if user is already logged in
	if ui.ValidateSession() {
		a.showDashboard()
	} else {
		a.showLogin()
	}

	return nil
}

// initializeFirebase initializes Firebase services
func (a *Application) initializeFirebase() error {
	client, err := firebase.Initialize(a.config.Firebase)
	if err != nil {
		return err
	}

	a.firebaseClient = client
	log.Println("Firebase initialized successfully")
	return nil
}

// createMainWindow creates the main application window
func (a *Application) createMainWindow() {
	width, height := a.config.GetWindowSize()

	a.window = a.fyneApp.NewWindow(a.config.App.Name)
	a.window.Resize(fyne.NewSize(float32(width), float32(height)))
	a.window.CenterOnScreen()

	// Set window icon
	if logoResource, err := fyne.LoadResourceFromPath("assets/logo.png"); err == nil {
		a.window.SetIcon(logoResource)
	}

	// Handle window close
	a.window.SetCloseIntercept(func() {
		if a.isLoggedIn {
			dialog.ShowConfirm("Keluar Aplikasi", "Apakah Anda yakin ingin keluar dari aplikasi?",
				func(confirm bool) {
					if confirm {
						a.cleanup()
						a.window.Close()
					}
				}, a.window)
		} else {
			a.cleanup()
			a.window.Close()
		}
	})
}

// showLogin shows the login screen
func (a *Application) showLogin() {
	a.isLoggedIn = false

	// Create auth service
	authService := firebase.NewAuthService(a.firebaseClient)

	// Create login screen
	a.loginScreen = ui.NewLoginScreen(a.window, authService, func() {
		a.onLoginSuccess()
	})

	// Set window content
	a.window.SetContent(a.loginScreen.GetContainer())
	a.window.SetTitle(a.config.App.Name + " - Login")

	log.Println("Login screen displayed")
}

// showDashboard shows the dashboard screen
func (a *Application) showDashboard() {
	a.isLoggedIn = true

	// Create dashboard screen
	a.dashboardScreen = ui.NewDashboardScreen(a.window, a.firebaseClient)

	// Set logout callback
	a.dashboardScreen.SetLogoutCallback(func() {
		a.onLogout()
	})

	// Set window content
	a.window.SetContent(a.dashboardScreen.GetContainer())

	// Update window title with user info
	_, _, userName := ui.GetCurrentUser()
	if userName != "" {
		a.window.SetTitle(a.config.App.Name + " - " + userName)
	} else {
		a.window.SetTitle(a.config.App.Name + " - Dashboard")
	}

	log.Println("Dashboard screen displayed")
}

// onLoginSuccess handles successful login
func (a *Application) onLoginSuccess() {
	log.Println("Login successful")
	a.showDashboard()
}

// onLogout handles user logout
func (a *Application) onLogout() {
	log.Println("User logged out")
	a.showLogin()
}

// showSetupDialog shows setup dialog when config file is missing
func (a *Application) showSetupDialog(configError error) error {
	// Create a simple window for setup dialog
	setupApp := app.New()
	setupWindow := setupApp.NewWindow("KasirNest - Setup")
	setupWindow.Resize(fyne.NewSize(500, 300))
	setupWindow.CenterOnScreen()

	message := `File konfigurasi tidak ditemukan!

Silakan buat file 'config/app.ini' berdasarkan template 'config/app.ini.example' dan isi dengan konfigurasi Firebase Anda.

Langkah-langkah:
1. Copy file 'config/app.ini.example' ke 'config/app.ini'
2. Edit file 'config/app.ini' dengan kredensial Firebase Anda
3. Jalankan aplikasi kembali

Error: ` + configError.Error()

	content := container.NewVBox(
		widget.NewLabel("Setup Konfigurasi"),
		widget.NewSeparator(),
		widget.NewRichTextFromMarkdown(message),
		widget.NewSeparator(),
		container.NewHBox(
			widget.NewButton("Buka Folder Config", func() {
				// Open config folder in file manager
				if err := a.openConfigFolder(); err != nil {
					dialog.ShowError(err, setupWindow)
				}
			}),
			widget.NewButton("Keluar", func() {
				setupWindow.Close()
			}),
		),
	)

	setupWindow.SetContent(container.NewPadded(content))
	setupWindow.ShowAndRun()

	return configError
}

// showConfigurationError shows configuration validation error
func (a *Application) showConfigurationError(configError error) error {
	setupApp := app.New()
	setupWindow := setupApp.NewWindow("KasirNest - Configuration Error")
	setupWindow.Resize(fyne.NewSize(500, 250))
	setupWindow.CenterOnScreen()

	message := `Konfigurasi tidak valid!

Silakan periksa file 'config/app.ini' dan pastikan semua field Firebase telah diisi dengan benar.

Error: ` + configError.Error()

	content := container.NewVBox(
		widget.NewLabel("Error Konfigurasi"),
		widget.NewSeparator(),
		widget.NewLabel(message),
		widget.NewSeparator(),
		container.NewHBox(
			widget.NewButton("Buka File Config", func() {
				// Open config file
				if err := a.openConfigFile(); err != nil {
					dialog.ShowError(err, setupWindow)
				}
			}),
			widget.NewButton("Keluar", func() {
				setupWindow.Close()
			}),
		),
	)

	setupWindow.SetContent(container.NewPadded(content))
	setupWindow.ShowAndRun()

	return configError
}

// showFirebaseError shows Firebase initialization error
func (a *Application) showFirebaseError(firebaseError error) error {
	setupApp := app.New()
	setupWindow := setupApp.NewWindow("KasirNest - Firebase Error")
	setupWindow.Resize(fyne.NewSize(500, 250))
	setupWindow.CenterOnScreen()

	message := `Gagal menghubungkan ke Firebase!

Periksa konfigurasi Firebase di file 'config/app.ini':
- Project ID
- Service Account credentials
- Internet connection

Error: ` + firebaseError.Error()

	content := container.NewVBox(
		widget.NewLabel("Error Firebase"),
		widget.NewSeparator(),
		widget.NewLabel(message),
		widget.NewSeparator(),
		container.NewHBox(
			widget.NewButton("Coba Lagi", func() {
				setupWindow.Close()
				// Restart application
				os.Exit(1)
			}),
			widget.NewButton("Keluar", func() {
				setupWindow.Close()
			}),
		),
	)

	setupWindow.SetContent(container.NewPadded(content))
	setupWindow.ShowAndRun()

	return firebaseError
}

// openConfigFolder opens the config folder in file manager
func (a *Application) openConfigFolder() error {
	// This is platform-specific
	// For Windows: explorer config
	// For macOS: open config
	// For Linux: xdg-open config
	return nil
}

// openConfigFile opens the config file in default editor
func (a *Application) openConfigFile() error {
	// This is platform-specific
	return nil
}

// cleanup performs cleanup before application exit
func (a *Application) cleanup() {
	log.Println("Cleaning up application...")

	// Close Firebase connections
	if a.firebaseClient != nil {
		if err := a.firebaseClient.Close(); err != nil {
			log.Printf("Error closing Firebase client: %v", err)
		}
	}

	// Save any pending configuration changes
	if a.config != nil {
		// Save config if needed
	}

	log.Println("Application cleanup completed")
}

// Helper function to check if running in debug mode
func (a *Application) isDebugMode() bool {
	return a.config != nil && a.config.IsDebug()
}

// Helper function to log debug messages
func (a *Application) debugLog(message string) {
	if a.isDebugMode() {
		log.Println("[DEBUG] " + message)
	}
}
