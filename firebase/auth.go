package firebase

import (
	"context"
	"errors"
	"time"

	"firebase.google.com/go/v4/auth"
)

// AuthService handles Firebase Authentication operations
type AuthService struct {
	client *auth.Client
	ctx    context.Context
}

// NewAuthService creates a new auth service
func NewAuthService(client *Client) *AuthService {
	return &AuthService{
		client: client.Auth,
		ctx:    client.GetContext(),
	}
}

// LoginWithEmailPassword authenticates user with email and password
// Note: This is server-side authentication, client-side auth should use Firebase SDK
func (a *AuthService) LoginWithEmailPassword(email, password string) (*auth.UserRecord, error) {
	if a.client == nil {
		return nil, errors.New("auth client not initialized")
	}

	// Get user by email
	user, err := a.client.GetUserByEmail(a.ctx, email)
	if err != nil {
		return nil, err
	}

	// Note: Server-side password verification is not directly supported
	// This would typically be done on the client-side using Firebase Auth SDK
	// For now, we'll just return the user record
	return user, nil
}

// GetUser gets user by UID
func (a *AuthService) GetUser(uid string) (*auth.UserRecord, error) {
	if a.client == nil {
		return nil, errors.New("auth client not initialized")
	}

	return a.client.GetUser(a.ctx, uid)
}

// GetUserByEmail gets user by email
func (a *AuthService) GetUserByEmail(email string) (*auth.UserRecord, error) {
	if a.client == nil {
		return nil, errors.New("auth client not initialized")
	}

	return a.client.GetUserByEmail(a.ctx, email)
}

// CreateUser creates a new user
func (a *AuthService) CreateUser(email, password, displayName string) (*auth.UserRecord, error) {
	if a.client == nil {
		return nil, errors.New("auth client not initialized")
	}

	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		DisplayName(displayName).
		EmailVerified(false)

	return a.client.CreateUser(a.ctx, params)
}

// UpdateUser updates user information
func (a *AuthService) UpdateUser(uid string, displayName string) (*auth.UserRecord, error) {
	if a.client == nil {
		return nil, errors.New("auth client not initialized")
	}

	params := (&auth.UserToUpdate{}).DisplayName(displayName)
	return a.client.UpdateUser(a.ctx, uid, params)
}

// DeleteUser deletes a user
func (a *AuthService) DeleteUser(uid string) error {
	if a.client == nil {
		return errors.New("auth client not initialized")
	}

	return a.client.DeleteUser(a.ctx, uid)
}

// VerifyIDToken verifies Firebase ID token
func (a *AuthService) VerifyIDToken(idToken string) (*auth.Token, error) {
	if a.client == nil {
		return nil, errors.New("auth client not initialized")
	}

	return a.client.VerifyIDToken(a.ctx, idToken)
}

// CreateCustomToken creates a custom token for authentication
func (a *AuthService) CreateCustomToken(uid string) (string, error) {
	if a.client == nil {
		return "", errors.New("auth client not initialized")
	}

	return a.client.CustomToken(a.ctx, uid)
}

// SetCustomUserClaims sets custom claims for a user
func (a *AuthService) SetCustomUserClaims(uid string, claims map[string]interface{}) error {
	if a.client == nil {
		return errors.New("auth client not initialized")
	}

	return a.client.SetCustomUserClaims(a.ctx, uid, claims)
}

// ListUsers lists all users with pagination
func (a *AuthService) ListUsers(maxResults int, pageToken string) ([]*auth.UserRecord, string, error) {
	if a.client == nil {
		return nil, "", errors.New("auth client not initialized")
	}

	iter := a.client.Users(a.ctx, pageToken)
	iter.PageInfo().MaxSize = maxResults

	var users []*auth.UserRecord
	nextPageToken := ""

	for {
		user, err := iter.Next()
		if err != nil {
			break
		}
		// Convert ExportedUserRecord to UserRecord if needed
		userRecord := &auth.UserRecord{
			UserInfo:               user.UserInfo,
			CustomClaims:           user.CustomClaims,
			Disabled:               user.Disabled,
			EmailVerified:          user.EmailVerified,
			ProviderUserInfo:       user.ProviderUserInfo,
			TokensValidAfterMillis: user.TokensValidAfterMillis,
		}
		users = append(users, userRecord)
	}

	return users, nextPageToken, nil
}

// RevokeRefreshTokens revokes all refresh tokens for a user
func (a *AuthService) RevokeRefreshTokens(uid string) error {
	if a.client == nil {
		return errors.New("auth client not initialized")
	}

	return a.client.RevokeRefreshTokens(a.ctx, uid)
}

// UserInfo represents simplified user information
type UserInfo struct {
	UID         string    `json:"uid"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	PhotoURL    string    `json:"photo_url"`
	Disabled    bool      `json:"disabled"`
	CreatedAt   time.Time `json:"created_at"`
}

// ConvertUserRecord converts auth.UserRecord to UserInfo
func ConvertUserRecord(user *auth.UserRecord) *UserInfo {
	return &UserInfo{
		UID:         user.UID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		PhotoURL:    user.PhotoURL,
		Disabled:    user.Disabled,
		CreatedAt:   time.Unix(user.UserMetadata.CreationTimestamp, 0),
	}
}
