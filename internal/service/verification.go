package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"moveshare/internal/repository"
)

// VerificationService defines the interface for verification business logic
type VerificationService interface {
	UploadDocument(ctx context.Context, userID int64, docType repository.DocumentType, fileName, fileURL string) error
	GetUserDocuments(ctx context.Context, userID int64) ([]repository.VerificationDocument, error)
	GetDocumentByType(ctx context.Context, userID int64, docType repository.DocumentType) (*repository.VerificationDocument, error)
	DeleteDocument(ctx context.Context, userID, docID int64) error
	ReviewDocument(ctx context.Context, docID int64, status repository.DocumentStatus, rejectionReason *string, reviewedBy int64) error
	GetPendingDocuments(ctx context.Context, limit, offset int) ([]repository.VerificationDocument, error)
	GetVerificationStatus(ctx context.Context, userID int64) (*VerificationStatus, error)
}

// VerificationStatus represents the overall verification status for a user
type VerificationStatus struct {
	IsFullyVerified bool                                  `json:"is_fully_verified"`
	Documents       map[string]DocumentVerificationStatus `json:"documents"`
}

// DocumentVerificationStatus represents the status of a specific document type
type DocumentVerificationStatus struct {
	Status          repository.DocumentStatus `json:"status"`
	FileName        string                    `json:"file_name,omitempty"`
	FileURL         string                    `json:"file_url,omitempty"`
	RejectionReason *string                   `json:"rejection_reason,omitempty"`
	UploadedAt      *time.Time                `json:"uploaded_at,omitempty"`
	ReviewedAt      *time.Time                `json:"reviewed_at,omitempty"`
}

// verificationService implements VerificationService
type verificationService struct {
	verificationRepo repository.VerificationRepository
}

// NewVerificationService creates a new VerificationService
func NewVerificationService(verificationRepo repository.VerificationRepository) VerificationService {
	return &verificationService{
		verificationRepo: verificationRepo,
	}
}

// UploadDocument uploads a new verification document
func (s *verificationService) UploadDocument(ctx context.Context, userID int64, docType repository.DocumentType, fileName, fileURL string) error {
	if fileName == "" {
		return errors.New("file name is required")
	}
	if fileURL == "" {
		return errors.New("file URL is required")
	}

	// Validate document type
	if !s.isValidDocumentType(docType) {
		return errors.New("invalid document type")
	}

	// Validate file extension
	if !s.isValidFileType(fileName) {
		return errors.New("invalid file type. Only PDF, JPG, JPEG, PNG files are allowed")
	}

	// Check if document already exists and replace it
	existingDoc, err := s.verificationRepo.GetDocumentByType(ctx, userID, docType)
	if err == nil && existingDoc != nil {
		// Delete existing document
		err = s.verificationRepo.DeleteDocument(ctx, userID, existingDoc.ID)
		if err != nil {
			return fmt.Errorf("failed to replace existing document: %w", err)
		}
	}

	doc := &repository.VerificationDocument{
		UserID:       userID,
		DocumentType: docType,
		FileName:     fileName,
		FileURL:      fileURL,
		Status:       repository.DocumentStatusPending,
		UploadedAt:   time.Now(),
	}

	return s.verificationRepo.CreateDocument(ctx, doc)
}

// GetUserDocuments retrieves all documents for a user
func (s *verificationService) GetUserDocuments(ctx context.Context, userID int64) ([]repository.VerificationDocument, error) {
	return s.verificationRepo.GetUserDocuments(ctx, userID)
}

// GetDocumentByType retrieves a specific document type for a user
func (s *verificationService) GetDocumentByType(ctx context.Context, userID int64, docType repository.DocumentType) (*repository.VerificationDocument, error) {
	return s.verificationRepo.GetDocumentByType(ctx, userID, docType)
}

// DeleteDocument deletes a verification document
func (s *verificationService) DeleteDocument(ctx context.Context, userID, docID int64) error {
	return s.verificationRepo.DeleteDocument(ctx, userID, docID)
}

// ReviewDocument updates the status of a verification document (admin only)
func (s *verificationService) ReviewDocument(ctx context.Context, docID int64, status repository.DocumentStatus, rejectionReason *string, reviewedBy int64) error {
	if status != repository.DocumentStatusApproved && status != repository.DocumentStatusRejected {
		return errors.New("invalid status for review")
	}

	if status == repository.DocumentStatusRejected && (rejectionReason == nil || *rejectionReason == "") {
		return errors.New("rejection reason is required when rejecting a document")
	}

	return s.verificationRepo.UpdateDocumentStatus(ctx, docID, status, rejectionReason, reviewedBy)
}

// GetPendingDocuments retrieves pending documents for admin review
func (s *verificationService) GetPendingDocuments(ctx context.Context, limit, offset int) ([]repository.VerificationDocument, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.verificationRepo.GetPendingDocuments(ctx, limit, offset)
}

// GetVerificationStatus returns the overall verification status for a user
func (s *verificationService) GetVerificationStatus(ctx context.Context, userID int64) (*VerificationStatus, error) {
	documents, err := s.verificationRepo.GetUserDocuments(ctx, userID)
	if err != nil {
		return nil, err
	}

	status := &VerificationStatus{
		IsFullyVerified: true,
		Documents: map[string]DocumentVerificationStatus{
			string(repository.DocumentTypeMCLicense):       {Status: repository.DocumentStatusPending},
			string(repository.DocumentTypeDOTCertificate):  {Status: repository.DocumentStatusPending},
			string(repository.DocumentTypeInsuranceCert):   {Status: repository.DocumentStatusPending},
			string(repository.DocumentTypeBusinessLicense): {Status: repository.DocumentStatusPending},
		},
	}

	// Create a map for quick lookup of latest documents by type
	latestDocs := make(map[repository.DocumentType]*repository.VerificationDocument)
	for i := range documents {
		doc := &documents[i]
		if existing, exists := latestDocs[doc.DocumentType]; !exists || doc.UploadedAt.After(existing.UploadedAt) {
			latestDocs[doc.DocumentType] = doc
		}
	}

	// Update status based on latest documents
	requiredTypes := []repository.DocumentType{
		repository.DocumentTypeMCLicense,
		repository.DocumentTypeDOTCertificate,
		repository.DocumentTypeInsuranceCert,
		repository.DocumentTypeBusinessLicense,
	}

	for _, docType := range requiredTypes {
		docStatus := DocumentVerificationStatus{Status: repository.DocumentStatusPending}

		if doc, exists := latestDocs[docType]; exists {
			docStatus = DocumentVerificationStatus{
				Status:          doc.Status,
				FileName:        doc.FileName,
				FileURL:         doc.FileURL,
				RejectionReason: doc.RejectionReason,
				UploadedAt:      &doc.UploadedAt,
				ReviewedAt:      doc.ReviewedAt,
			}
		}

		status.Documents[string(docType)] = docStatus

		// Check if all documents are approved
		if docStatus.Status != repository.DocumentStatusApproved {
			status.IsFullyVerified = false
		}
	}

	return status, nil
}

// isValidDocumentType checks if the document type is valid
func (s *verificationService) isValidDocumentType(docType repository.DocumentType) bool {
	validTypes := map[repository.DocumentType]bool{
		repository.DocumentTypeMCLicense:       true,
		repository.DocumentTypeDOTCertificate:  true,
		repository.DocumentTypeInsuranceCert:   true,
		repository.DocumentTypeBusinessLicense: true,
	}
	return validTypes[docType]
}

// isValidFileType checks if the file type is allowed
func (s *verificationService) isValidFileType(fileName string) bool {
	fileName = strings.ToLower(fileName)
	allowedExtensions := []string{".pdf", ".jpg", ".jpeg", ".png"}

	for _, ext := range allowedExtensions {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}
	return false
}
