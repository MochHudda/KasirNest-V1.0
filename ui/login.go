package ui

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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

// setupUI sets up the login interface with modern dark design
func (l *LoginScreen) setupUI() {
	// Create dark background
	darkBg := canvas.NewRectangle(color.NRGBA{R: 28, G: 42, B: 56, A: 255}) // Dark navy

	// Create KasirNest logo (stylized K with shopping cart)
	logoText := canvas.NewText("K", color.NRGBA{R: 56, G: 189, B: 190, A: 255}) // Teal color
	logoText.TextSize = 80
	logoText.TextStyle = fyne.TextStyle{Bold: true}
	logoText.Alignment = fyne.TextAlignCenter

	// Cart wheels (orange circles)
	wheel1 := canvas.NewCircle(color.NRGBA{R: 255, G: 138, B: 48, A: 255}) // Orange
	wheel1.Resize(fyne.NewSize(12, 12))
	wheel2 := canvas.NewCircle(color.NRGBA{R: 255, G: 138, B: 48, A: 255}) // Orange
	wheel2.Resize(fyne.NewSize(12, 12))

	// Position wheels relative to the K
	logoContainer := container.NewWithoutLayout(
		logoText,
		wheel1,
		wheel2,
	)
	logoContainer.Resize(fyne.NewSize(120, 100))
	wheel1.Move(fyne.NewPos(70, 85))
	wheel2.Move(fyne.NewPos(95, 85))

	// KasirNest title
	kasirText := canvas.NewText("KasirNest", color.NRGBA{R: 56, G: 189, B: 190, A: 255}) // Teal
	kasirText.TextSize = 36
	kasirText.TextStyle = fyne.TextStyle{Bold: true}
	kasirText.Alignment = fyne.TextAlignCenter

	// V1.0 version text
	versionText := canvas.NewText("V1.0", color.NRGBA{R: 255, G: 138, B: 48, A: 255}) // Orange
	versionText.TextSize = 20
	versionText.TextStyle = fyne.TextStyle{Bold: true}
	versionText.Alignment = fyne.TextAlignCenter

	// Login title
	loginTitle := canvas.NewText("Login", color.NRGBA{R: 255, G: 255, B: 255, A: 255}) // White
	loginTitle.TextSize = 28
	loginTitle.TextStyle = fyne.TextStyle{Bold: true}
	loginTitle.Alignment = fyne.TextAlignCenter

	// Create custom styled entries
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

	// Email field with envelope icon
	emailIcon := canvas.NewImageFromResource(theme.MailSendIcon())
	emailIcon.Resize(fyne.NewSize(20, 20))
	emailContainer := container.NewBorder(nil, nil,
		container.NewPadded(emailIcon), nil,
		l.emailEntry,
	)

	// Password field with lock icon
	passwordIcon := canvas.NewImageFromResource(theme.ViewRefreshIcon()) // Using available icon as substitute
	passwordIcon.Resize(fyne.NewSize(20, 20))
	passwordContainer := container.NewBorder(nil, nil,
		container.NewPadded(passwordIcon), nil,
		l.passwordEntry,
	)

	// Create modern login button
	l.loginButton = widget.NewButton("Log in", l.handleLogin)
	l.loginButton.Importance = widget.HighImportance

	// Create sign up link
	signUpLabel := widget.NewRichTextFromMarkdown("Don't have an account? **[Sign up](signup)**")
	signUpLabel.Wrapping = fyne.TextWrapOff

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

	// Create main login form with proper spacing
	formContent := container.NewVBox(
		// Logo and branding section
		container.NewCenter(logoContainer),
		container.NewVBox(
			container.NewCenter(kasirText),
			container.NewCenter(versionText),
		),
		widget.NewSeparator(),

		// Login title
		container.NewCenter(loginTitle),

		// Spacing
		canvas.NewRectangle(color.Transparent),

		// Input fields
		emailContainer,
		passwordContainer,

		// Spacing
		canvas.NewRectangle(color.Transparent),

		// Login button
		l.loginButton,

		// Loading indicator
		l.loadingLabel,

		// Sign up link
		container.NewCenter(signUpLabel),
	)

	// Add padding around the form
	paddedForm := container.NewPadded(formContent)

	// Create the main container with dark background
	l.container = container.NewMax(
		darkBg,
		container.NewCenter(
			container.NewVBox(
				canvas.NewRectangle(color.Transparent), // Top spacing
				paddedForm,
				canvas.NewRectangle(color.Transparent), // Bottom spacing
			),
		),
	)
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
