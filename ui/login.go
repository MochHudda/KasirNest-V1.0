package ui

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"kasirnest/firebase"
	"kasirnest/utils"
)

// LoginScreen represents the login interface
type LoginScreen struct {
	window        fyne.Window
	container     *fyne.Container
	emailEntry    *widget.Entry
	passwordEntry *widget.Entry
	loginButton   *widget.Button
	loadingLabel  *widget.Label
	authService   *firebase.AuthService
	onLogin       func()
}

// NewLoginScreen creates a new login screen
func NewLoginScreen(w fyne.Window, authService *firebase.AuthService, onLoginSuccess func()) *LoginScreen {
	login := &LoginScreen{
		window:      w,
		authService: authService,
		onLogin:     onLoginSuccess,
	}

	login.setupUI()
	return login
}

// setupUI sets up the login interface
func (l *LoginScreen) setupUI() {
	// Create logo
	logoResource, err := fyne.LoadResourceFromPath("assets/logo.png")
	var logo *widget.Icon
	if err == nil {
		logo = widget.NewIcon(logoResource)
	} else {
		// Fallback to default icon if logo not found
		logo = widget.NewIcon(theme.LoginIcon())
	}

	// Create title
	title := widget.NewLabel("KasirNest")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := widget.NewLabel("Sistem Kasir Modern")
	subtitle.Alignment = fyne.TextAlignCenter

	// Create form fields
	l.emailEntry = widget.NewEntry()
	l.emailEntry.SetPlaceHolder("Email")
	l.emailEntry.Validator = func(text string) error {
		if !utils.ValidateEmail(text) {
			return fmt.Errorf("format email tidak valid")
		}
		return nil
	}

	l.passwordEntry = widget.NewPasswordEntry()
	l.passwordEntry.SetPlaceHolder("Password")
	l.passwordEntry.Validator = func(text string) error {
		if !utils.ValidateMinLength(text, 6) {
			return fmt.Errorf("password minimal 6 karakter")
		}
		return nil
	}

	// Create login button
	l.loginButton = widget.NewButton("Masuk", l.handleLogin)
	l.loginButton.Importance = widget.HighImportance

	// Create loading label
	l.loadingLabel = widget.NewLabel("")
	l.loadingLabel.Hide()

	// Handle Enter key press
	l.emailEntry.OnSubmitted = func(text string) {
		l.passwordEntry.FocusGained()
	}
	l.passwordEntry.OnSubmitted = func(text string) {
		l.handleLogin()
	}

	// Create form
	form := container.NewVBox(
		widget.NewCard("", "", container.NewVBox(
			container.NewCenter(logo),
			container.NewCenter(title),
			container.NewCenter(subtitle),
			widget.NewSeparator(),
			widget.NewLabel("Email:"),
			l.emailEntry,
			widget.NewLabel("Password:"),
			l.passwordEntry,
			container.NewPadded(l.loginButton),
			l.loadingLabel,
		)),
	)

	// Center the form
	l.container = container.NewCenter(form)
}

// handleLogin handles the login process
func (l *LoginScreen) handleLogin() {
	// Validate inputs
	email := l.emailEntry.Text
	password := l.passwordEntry.Text

	// Validate form
	if err := l.emailEntry.Validate(); err != nil {
		l.showError("Email: " + err.Error())
		return
	}

	if err := l.passwordEntry.Validate(); err != nil {
		l.showError("Password: " + err.Error())
		return
	}

	// Use password for authentication (placeholder for actual implementation)
	_ = password // This will be used for actual Firebase client-side auth

	// Show loading state
	l.setLoading(true)

	// Perform login (this is a simplified version)
	// In a real application, you would use Firebase Auth SDK on client side
	go func() {
		// Simulate authentication
		// Note: Server-side password verification is not supported by Firebase Admin SDK
		// This should be done using Firebase Auth SDK on the client side

		user, err := l.authService.GetUserByEmail(email)
		if err != nil {
			l.setLoading(false)
			l.showError("Login gagal: " + err.Error())
			return
		}

		if user == nil {
			l.setLoading(false)
			l.showError("User tidak ditemukan")
			return
		}

		// Store user session (simplified)
		l.storeUserSession(user.UID, user.Email, user.DisplayName)

		l.setLoading(false)
		l.onLogin()
	}()
}

// setLoading sets the loading state
func (l *LoginScreen) setLoading(loading bool) {
	if loading {
		l.loginButton.Disable()
		l.loadingLabel.SetText("Memproses login...")
		l.loadingLabel.Show()
	} else {
		l.loginButton.Enable()
		l.loadingLabel.Hide()
	}
}

// showError displays an error message
func (l *LoginScreen) showError(message string) {
	dialog.ShowError(fmt.Errorf(message), l.window)
}

// showInfo displays an info message
func (l *LoginScreen) showInfo(title, message string) {
	dialog.ShowInformation(title, message, l.window)
}

// GetContainer returns the login container
func (l *LoginScreen) GetContainer() *fyne.Container {
	return l.container
}

// storeUserSession stores user session data
func (l *LoginScreen) storeUserSession(uid, email, name string) {
	// Store session in preferences
	prefs := fyne.CurrentApp().Preferences()
	prefs.SetString("user_id", uid)
	prefs.SetString("user_email", email)
	prefs.SetString("user_name", name)
	prefs.SetString("login_time", fmt.Sprintf("%d", time.Now().Unix()))

	log.Printf("User logged in: %s (%s)", name, email)
}

// ValidateSession validates if user session is still valid
func ValidateSession() bool {
	prefs := fyne.CurrentApp().Preferences()
	userID := prefs.String("user_id")
	loginTimeStr := prefs.String("login_time")

	if userID == "" || loginTimeStr == "" {
		return false
	}

	// Check if session is expired (24 hours)
	loginTime, err := strconv.ParseInt(loginTimeStr, 10, 64)
	if err != nil {
		return false
	}

	sessionDuration := time.Since(time.Unix(loginTime, 0))
	maxSessionDuration := 24 * time.Hour

	return sessionDuration < maxSessionDuration
}

// GetCurrentUser returns current logged-in user info
func GetCurrentUser() (uid, email, name string) {
	prefs := fyne.CurrentApp().Preferences()
	return prefs.String("user_id"), prefs.String("user_email"), prefs.String("user_name")
}

// ClearSession clears user session
func ClearSession() {
	prefs := fyne.CurrentApp().Preferences()
	prefs.RemoveValue("user_id")
	prefs.RemoveValue("user_email")
	prefs.RemoveValue("user_name")
	prefs.RemoveValue("login_time")
}
