# KasirNest - Desktop Cashier Application

Building a secure, modular desktop cashier (POS) application using Go, Fyne UI framework, and Firebase backend services.

## User Review Required

> [!IMPORTANT]
> **Firebase Project Setup Required**
> Before running this application, you'll need to:
> 1. Create a Firebase project at [Firebase Console](https://console.firebase.google.com)
> 2. Enable Firebase Authentication (Email/Password provider)
> 3. Create a Firestore Database
> 4. Optionally enable Firebase Storage for product images
> 5. Set up Firestore security rules to restrict access to authenticated users only
> 6. Obtain your Firebase project credentials (API Key, Project ID, etc.)

> [!WARNING]
> **Security Considerations**
> - The app will use `config/app.ini` to store Firebase credentials (gitignored)
> - Production builds will use obfuscation via `garble` to protect binary inspection
> - You should provide the actual Firebase credentials after project setup

> [!NOTE]
> **Build Tools Required**
> - Go 1.21 or later
> - `garble` for code obfuscation: `go install mvdan.cc/garble@latest`
> - Fyne build tools: `go install fyne.io/fyne/v2/cmd/fyne@latest`

---

## Proposed Changes

### Project Structure & Configuration

Creating the complete project structure as specified in requirements.

#### [NEW] [go.mod](file:///c:/ProjectHudda/KasirNest%20V1/go.mod)
Initialize Go module with required dependencies:
- `fyne.io/fyne/v2` - Desktop UI framework
- `firebase.google.com/go/v4` - Firebase Admin SDK
- `cloud.google.com/go/firestore` - Firestore client
- `cloud.google.com/go/storage` - Storage client
- `gopkg.in/ini.v1` - INI file parser for configuration

#### [NEW] [config/app.ini.example](file:///c:/ProjectHudda/KasirNest%20V1/config/app.ini.example)
Template configuration file showing required Firebase credentials structure (actual `app.ini` will be gitignored).

#### [NEW] [.gitignore](file:///c:/ProjectHudda/KasirNest%20V1/.gitignore)
Ignore sensitive files: `config/app.ini`, build artifacts, and IDE files.

---

### Data Models

Define Go structs for all data entities matching the Firestore schema.

#### [NEW] [models/user.go](file:///c:/ProjectHudda/KasirNest%20V1/models/user.go)
User model with fields: `UserID`, `Email`, `Name`, `Role`, `CreatedAt`, `LastLogin`.

#### [NEW] [models/product.go](file:///c:/ProjectHudda/KasirNest%20V1/models/product.go)
Product model with fields: `ProductID`, `Name`, `Price`, `Stock`, `Category`, `Barcode`, `ImageURL`, timestamps.

#### [NEW] [models/transaction.go](file:///c:/ProjectHudda/KasirNest%20V1/models/transaction.go)
Transaction model with fields: `TransID`, `UserID`, `Date`, `Total`, `PaymentMethod`, `Items` (array of transaction items).

#### [NEW] [models/category.go](file:///c:/ProjectHudda/KasirNest%20V1/models/category.go)
Category model with fields: `CategoryID`, `Name`, `Description`.

#### [NEW] [models/report.go](file:///c:/ProjectHudda/KasirNest%20V1/models/report.go)
Report model for daily/weekly summaries.

---

### Utility Functions

Helper functions for common operations.

#### [NEW] [utils/validator.go](file:///c:/ProjectHudda/KasirNest%20V1/utils/validator.go)
Input validation functions:
- Email format validation
- Price validation (must be positive number)
- Stock validation (cannot be negative)
- Required field checks

#### [NEW] [utils/formatter.go](file:///c:/ProjectHudda/KasirNest%20V1/utils/formatter.go)
Formatting utilities:
- Currency formatting (Rupiah)
- Date/time formatting
- Number formatting with thousand separators

#### [NEW] [utils/crypto.go](file:///c:/ProjectHudda/KasirNest%20V1/utils/crypto.go)
Encryption/decryption utilities for sensitive configuration data (optional enhancement).

---

### Firebase Integration

Modules for connecting to Firebase services.

#### [NEW] [firebase/client.go](file:///c:/ProjectHudda/KasirNest%20V1/firebase/client.go)
Initialize Firebase app and clients from configuration file.

#### [NEW] [firebase/auth.go](file:///c:/ProjectHudda/KasirNest%20V1/firebase/auth.go)
Firebase Authentication integration:
- Login with email/password
- Get user token
- Verify authentication state
- Logout functionality

#### [NEW] [firebase/firestore.go](file:///c:/ProjectHudda/KasirNest%20V1/firebase/firestore.go)
Firestore CRUD operations:
- Generic document CRUD methods
- Collection queries with filters
- Batch operations for transactions
- Real-time listeners (optional)

#### [NEW] [firebase/storage.go](file:///c:/ProjectHudda/KasirNest%20V1/firebase/storage.go)
Firebase Storage integration:
- Upload product images
- Generate download URLs
- Delete images

---

### Configuration Management

#### [NEW] [config/config.go](file:///c:/ProjectHudda/KasirNest%20V1/config/config.go)
Configuration loader that reads from `app.ini`:
- Firebase project credentials
- Application settings
- Environment-specific configurations

---

### UI Components (Fyne)

Desktop interface components using Fyne framework.

#### [NEW] [ui/login.go](file:///c:/ProjectHudda/KasirNest%20V1/ui/login.go)
Login screen:
- Email and password entry fields
- Login button with Firebase Auth integration
- Error message display
- Loading state during authentication

#### [NEW] [ui/dashboard.go](file:///c:/ProjectHudda/KasirNest%20V1/ui/dashboard.go)
Main dashboard after login:
- Navigation menu (Products, Transactions, Reports, Logout)
- Welcome message with user info
- Quick stats display
- Tab-based or split container layout

#### [NEW] [ui/products.go](file:///c:/ProjectHudda/KasirNest%20V1/ui/products.go)
Products management interface:
- Product list table with search/filter
- Add new product form
- Edit product dialog
- Delete confirmation
- Image upload support
- Stock and price validation

#### [NEW] [ui/transactions.go](file:///c:/ProjectHudda/KasirNest%20V1/ui/transactions.go)
Transactions interface:
- Create new transaction (POS screen)
- Product selection and quantity input
- Real-time total calculation
- Payment method selection
- Transaction history list
- Transaction details view

#### [NEW] [ui/reports.go](file:///c:/ProjectHudda/KasirNest%20V1/ui/reports.go)
Reports screen:
- Date range filter
- Sales summary (total, count)
- Top products list
- Export functionality (optional)

---

### Main Application

#### [NEW] [main.go](file:///c:/ProjectHudda/KasirNest%20V1/main.go)
Application entry point:
- Initialize configuration
- Initialize Firebase clients
- Create Fyne application
- Show login window
- Handle window lifecycle

---

### Security & Build

#### [NEW] [internal/secure.go](file:///c:/ProjectHudda/KasirNest%20V1/internal/secure.go)
Security utilities and anti-inspection hooks (placeholder for advanced protections).

#### [NEW] [build/build.sh](file:///c:/ProjectHudda/KasirNest%20V1/build/build.sh)
Build script with:
- `go build -ldflags="-s -w"` for size optimization
- `garble build` for code obfuscation
- Platform-specific builds (Windows, macOS, Linux)

#### [NEW] [build/build.bat](file:///c:/ProjectHudda/KasirNest%20V1/build/build.bat)
Windows batch version of build script.

---

### Assets & Documentation

#### [NEW] [assets/logo.png](file:///c:/ProjectHudda/KasirNest%20V1/assets/logo.png)
Application logo (placeholder or generated image).

#### [NEW] [README.md](file:///c:/ProjectHudda/KasirNest%20V1/README.md)
Comprehensive documentation:
- Project overview
- Prerequisites and setup
- Firebase configuration steps
- Build and run instructions
- Security best practices
- Distribution guidelines

#### [NEW] [FIREBASE_SETUP.md](file:///c:/ProjectHudda/KasirNest%20V1/FIREBASE_SETUP.md)
Detailed Firebase setup guide:
- Creating Firebase project
- Enabling services
- Setting up Firestore security rules
- Obtaining credentials

---

## Verification Plan

### Automated Tests

1. **Build Verification**
   ```bash
   cd "c:\ProjectHudda\KasirNest V1"
   go mod download
   go build -v
   ```

2. **Module Tests** (if time permits)
   ```bash
   go test ./utils/... -v
   go test ./firebase/... -v
   ```

### Manual Verification

1. **Configuration Setup**
   - Create `config/app.ini` from template
   - Fill in Firebase credentials
   - Verify app initializes without errors

2. **UI Testing**
   - Run application: `go run main.go`
   - Test login flow with Firebase Auth
   - Navigate through dashboard tabs
   - Test product CRUD operations
   - Create a test transaction
   - View reports with sample data

3. **Build Testing**
   - Run build script
   - Verify obfuscated binary works
   - Test `fyne package` for installer creation

4. **Security Verification**
   - Confirm no hardcoded credentials in source
   - Check that `config/app.ini` is gitignored
   - Verify obfuscated binary cannot be easily reverse-engineered
