package firebase

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// FirestoreService handles Firestore operations
type FirestoreService struct {
	client *firestore.Client
	ctx    context.Context
}

// NewFirestoreService creates a new Firestore service
func NewFirestoreService(client *Client) *FirestoreService {
	return &FirestoreService{
		client: client.Firestore,
		ctx:    client.GetContext(),
	}
}

// Create creates a new document in a collection
func (fs *FirestoreService) Create(collection, documentID string, data interface{}) error {
	if fs.client == nil {
		return errors.New("firestore client not initialized")
	}

	_, err := fs.client.Collection(collection).Doc(documentID).Set(fs.ctx, data)
	return err
}

// CreateWithAutoID creates a new document with auto-generated ID
func (fs *FirestoreService) CreateWithAutoID(collection string, data interface{}) (string, error) {
	if fs.client == nil {
		return "", errors.New("firestore client not initialized")
	}

	doc, _, err := fs.client.Collection(collection).Add(fs.ctx, data)
	if err != nil {
		return "", err
	}
	return doc.ID, nil
}

// Get retrieves a document by ID
func (fs *FirestoreService) Get(collection, documentID string, dest interface{}) error {
	if fs.client == nil {
		return errors.New("firestore client not initialized")
	}

	doc, err := fs.client.Collection(collection).Doc(documentID).Get(fs.ctx)
	if err != nil {
		return err
	}

	if !doc.Exists() {
		return errors.New("document does not exist")
	}

	return doc.DataTo(dest)
}

// Update updates an existing document
func (fs *FirestoreService) Update(collection, documentID string, data interface{}) error {
	if fs.client == nil {
		return errors.New("firestore client not initialized")
	}

	_, err := fs.client.Collection(collection).Doc(documentID).Set(fs.ctx, data, firestore.MergeAll)
	return err
}

// Delete deletes a document
func (fs *FirestoreService) Delete(collection, documentID string) error {
	if fs.client == nil {
		return errors.New("firestore client not initialized")
	}

	_, err := fs.client.Collection(collection).Doc(documentID).Delete(fs.ctx)
	return err
}

// List retrieves all documents from a collection
func (fs *FirestoreService) List(collection string, dest interface{}) error {
	if fs.client == nil {
		return errors.New("firestore client not initialized")
	}

	iter := fs.client.Collection(collection).Documents(fs.ctx)
	defer iter.Stop()

	var results []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		data := doc.Data()
		// Add document ID to the data
		data["id"] = doc.Ref.ID
		results = append(results, data)
	}

	// This is a simplified implementation
	// In a real application, you'd want to properly handle type conversion
	return nil
}

// Query executes a query with filters
func (fs *FirestoreService) Query(collection string, filters []QueryFilter, dest interface{}) error {
	if fs.client == nil {
		return errors.New("firestore client not initialized")
	}

	query := fs.client.Collection(collection).Query

	// Apply filters
	for _, filter := range filters {
		switch filter.Operator {
		case "==":
			query = query.Where(filter.Field, "==", filter.Value)
		case "!=":
			query = query.Where(filter.Field, "!=", filter.Value)
		case ">":
			query = query.Where(filter.Field, ">", filter.Value)
		case ">=":
			query = query.Where(filter.Field, ">=", filter.Value)
		case "<":
			query = query.Where(filter.Field, "<", filter.Value)
		case "<=":
			query = query.Where(filter.Field, "<=", filter.Value)
		case "array-contains":
			query = query.Where(filter.Field, "array-contains", filter.Value)
		case "in":
			query = query.Where(filter.Field, "in", filter.Value)
		default:
			return fmt.Errorf("unsupported operator: %s", filter.Operator)
		}
	}

	iter := query.Documents(fs.ctx)
	defer iter.Stop()

	var results []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		data := doc.Data()
		data["id"] = doc.Ref.ID
		results = append(results, data)
	}

	return nil
}

// BatchWrite performs multiple operations in a single batch
func (fs *FirestoreService) BatchWrite(operations []BatchOperation) error {
	if fs.client == nil {
		return errors.New("firestore client not initialized")
	}

	batch := fs.client.Batch()

	for _, op := range operations {
		docRef := fs.client.Collection(op.Collection).Doc(op.DocumentID)

		switch op.Operation {
		case "create", "set":
			batch.Set(docRef, op.Data)
		case "update":
			batch.Set(docRef, op.Data, firestore.MergeAll)
		case "delete":
			batch.Delete(docRef)
		default:
			return fmt.Errorf("unsupported batch operation: %s", op.Operation)
		}
	}

	_, err := batch.Commit(fs.ctx)
	return err
}

// GetCollectionSize returns the number of documents in a collection
func (fs *FirestoreService) GetCollectionSize(collection string) (int, error) {
	if fs.client == nil {
		return 0, errors.New("firestore client not initialized")
	}

	iter := fs.client.Collection(collection).Documents(fs.ctx)
	defer iter.Stop()

	count := 0
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}
		count++
	}

	return count, nil
}

// QueryFilter represents a query filter
type QueryFilter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// BatchOperation represents a batch operation
type BatchOperation struct {
	Operation  string      `json:"operation"` // "create", "update", "delete"
	Collection string      `json:"collection"`
	DocumentID string      `json:"document_id"`
	Data       interface{} `json:"data,omitempty"`
}

// Exists checks if a document exists
func (fs *FirestoreService) Exists(collection, documentID string) (bool, error) {
	if fs.client == nil {
		return false, errors.New("firestore client not initialized")
	}

	doc, err := fs.client.Collection(collection).Doc(documentID).Get(fs.ctx)
	if err != nil {
		return false, err
	}

	return doc.Exists(), nil
}
