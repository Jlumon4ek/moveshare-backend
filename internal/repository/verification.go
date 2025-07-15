package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DocumentType represents the type of verification document
type DocumentType string

const (
	DocumentTypeMCLicense       DocumentType = "mc_license"
	DocumentTypeDOTCertificate  DocumentType = "dot_certificate"
	DocumentTypeInsuranceCert   DocumentType = "insurance_certificate"
	DocumentTypeBusinessLicense DocumentType = "business_license"
)

// DocumentStatus represents the verification status of a document
type DocumentStatus string

const (
	DocumentStatusPending  DocumentStatus = "pending"
	DocumentStatusApproved DocumentStatus = "approved"
	DocumentStatusRejected DocumentStatus = "rejected"
)

// VerificationDocument represents a user's verification document
type VerificationDocument struct {
	ID              int64          `json:"id"`
	UserID          int64          `json:"user_id"`
	DocumentType    DocumentType   `json:"document_type"`
	FileName        string         `json:"file_name"`
	FileURL         string         `json:"file_url"`
	Status          DocumentStatus `json:"status"`
	RejectionReason *string        `json:"rejection_reason,omitempty"`
	UploadedAt      time.Time      `json:"uploaded_at"`
	ReviewedAt      *time.Time     `json:"reviewed_at,omitempty"`
	ReviewedBy      *int64         `json:"reviewed_by,omitempty"`
}

// VerificationRepository defines the interface for verification document operations
type VerificationRepository interface {
	CreateDocument(ctx context.Context, doc *VerificationDocument) error
	GetUserDocuments(ctx context.Context, userID int64) ([]VerificationDocument, error)
	GetDocumentByType(ctx context.Context, userID int64, docType DocumentType) (*VerificationDocument, error)
	UpdateDocumentStatus(ctx context.Context, docID int64, status DocumentStatus, rejectionReason *string, reviewedBy int64) error
	DeleteDocument(ctx context.Context, userID, docID int64) error
	GetPendingDocuments(ctx context.Context, limit, offset int) ([]VerificationDocument, error)
}

// verificationRepository implements VerificationRepository
type verificationRepository struct {
	db *pgxpool.Pool
}

// NewVerificationRepository creates a new VerificationRepository
func NewVerificationRepository(db *pgxpool.Pool) VerificationRepository {
	return &verificationRepository{db: db}
}

// CreateDocument creates a new verification document
func (r *verificationRepository) CreateDocument(ctx context.Context, doc *VerificationDocument) error {
	query := `
		INSERT INTO verification_documents (user_id, document_type, file_name, file_url, status, uploaded_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err := r.db.QueryRow(ctx, query,
		doc.UserID,
		doc.DocumentType,
		doc.FileName,
		doc.FileURL,
		doc.Status,
		doc.UploadedAt,
	).Scan(&doc.ID)

	return err
}

// GetUserDocuments retrieves all documents for a user
func (r *verificationRepository) GetUserDocuments(ctx context.Context, userID int64) ([]VerificationDocument, error) {
	query := `
		SELECT id, user_id, document_type, file_name, file_url, status, 
		       rejection_reason, uploaded_at, reviewed_at, reviewed_by
		FROM verification_documents 
		WHERE user_id = $1
		ORDER BY uploaded_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []VerificationDocument
	for rows.Next() {
		var doc VerificationDocument
		err := rows.Scan(
			&doc.ID,
			&doc.UserID,
			&doc.DocumentType,
			&doc.FileName,
			&doc.FileURL,
			&doc.Status,
			&doc.RejectionReason,
			&doc.UploadedAt,
			&doc.ReviewedAt,
			&doc.ReviewedBy,
		)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, rows.Err()
}

// GetDocumentByType retrieves a specific document type for a user
func (r *verificationRepository) GetDocumentByType(ctx context.Context, userID int64, docType DocumentType) (*VerificationDocument, error) {
	query := `
		SELECT id, user_id, document_type, file_name, file_url, status, 
		       rejection_reason, uploaded_at, reviewed_at, reviewed_by
		FROM verification_documents 
		WHERE user_id = $1 AND document_type = $2
		ORDER BY uploaded_at DESC
		LIMIT 1`

	var doc VerificationDocument
	err := r.db.QueryRow(ctx, query, userID, docType).Scan(
		&doc.ID,
		&doc.UserID,
		&doc.DocumentType,
		&doc.FileName,
		&doc.FileURL,
		&doc.Status,
		&doc.RejectionReason,
		&doc.UploadedAt,
		&doc.ReviewedAt,
		&doc.ReviewedBy,
	)

	if err != nil {
		return nil, err
	}

	return &doc, nil
}

// UpdateDocumentStatus updates the status of a verification document
func (r *verificationRepository) UpdateDocumentStatus(ctx context.Context, docID int64, status DocumentStatus, rejectionReason *string, reviewedBy int64) error {
	query := `
		UPDATE verification_documents 
		SET status = $1, rejection_reason = $2, reviewed_at = NOW(), reviewed_by = $3
		WHERE id = $4`

	_, err := r.db.Exec(ctx, query, status, rejectionReason, reviewedBy, docID)
	return err
}

// DeleteDocument deletes a verification document
func (r *verificationRepository) DeleteDocument(ctx context.Context, userID, docID int64) error {
	query := `DELETE FROM verification_documents WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, docID, userID)
	return err
}

// GetPendingDocuments retrieves pending documents for admin review
func (r *verificationRepository) GetPendingDocuments(ctx context.Context, limit, offset int) ([]VerificationDocument, error) {
	query := `
		SELECT vd.id, vd.user_id, vd.document_type, vd.file_name, vd.file_url, 
		       vd.status, vd.rejection_reason, vd.uploaded_at, vd.reviewed_at, vd.reviewed_by
		FROM verification_documents vd
		WHERE vd.status = $1
		ORDER BY vd.uploaded_at ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, DocumentStatusPending, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []VerificationDocument
	for rows.Next() {
		var doc VerificationDocument
		err := rows.Scan(
			&doc.ID,
			&doc.UserID,
			&doc.DocumentType,
			&doc.FileName,
			&doc.FileURL,
			&doc.Status,
			&doc.RejectionReason,
			&doc.UploadedAt,
			&doc.ReviewedAt,
			&doc.ReviewedBy,
		)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, rows.Err()
}
