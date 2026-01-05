package firebase

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// Client holds Firebase service clients
type Client struct {
	App       *firebase.App
	Auth      *auth.Client
	Firestore *firestore.Client
	Storage   *storage.Client
	ctx       context.Context
}

// FirebaseConfig holds Firebase configuration
type FirebaseConfig struct {
	ProjectID        string
	PrivateKeyID     string
	PrivateKey       string
	ClientEmail      string
	ClientID         string
	AuthURI          string
	TokenURI         string
	AuthProviderX509 string
	ClientX509       string
	StorageBucket    string
}

// Initialize creates and initializes Firebase clients
func Initialize(config *FirebaseConfig) (*Client, error) {
	ctx := context.Background()

	// Create credentials JSON
	credentialsJSON := createCredentialsJSON(config)

	// Initialize Firebase app
	opt := option.WithCredentialsJSON(credentialsJSON)
	firebaseConfig := &firebase.Config{
		ProjectID:     config.ProjectID,
		StorageBucket: config.StorageBucket,
	}

	app, err := firebase.NewApp(ctx, firebaseConfig, opt)
	if err != nil {
		return nil, err
	}

	// Initialize Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Printf("Failed to initialize Auth client: %v", err)
		// Continue without Auth for now
	}

	// Initialize Firestore client
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	// Get raw storage client for our storage operations
	rawStorageClient, err := storage.NewClient(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create raw storage client: %v", err)
	}

	return &Client{
		App:       app,
		Auth:      authClient,
		Firestore: firestoreClient,
		Storage:   rawStorageClient,
		ctx:       ctx,
	}, nil
}

// GetContext returns the context used by the client
func (c *Client) GetContext() context.Context {
	return c.ctx
}

// Close closes all Firebase clients
func (c *Client) Close() error {
	if c.Firestore != nil {
		return c.Firestore.Close()
	}
	return nil
}

// createCredentialsJSON creates credentials JSON from config
func createCredentialsJSON(config *FirebaseConfig) []byte {
	credentialsTemplate := `{
		"type": "service_account",
		"project_id": "%s",
		"private_key_id": "%s",
		"private_key": "%s",
		"client_email": "%s",
		"client_id": "%s",
		"auth_uri": "%s",
		"token_uri": "%s",
		"auth_provider_x509_cert_url": "%s",
		"client_x509_cert_url": "%s"
	}`

	// Format the credentials JSON
	credentialsJSON := fmt.Sprintf(credentialsTemplate,
		config.ProjectID,
		config.PrivateKeyID,
		config.PrivateKey,
		config.ClientEmail,
		config.ClientID,
		config.AuthURI,
		config.TokenURI,
		config.AuthProviderX509,
		config.ClientX509,
	)

	return []byte(credentialsJSON)
}
