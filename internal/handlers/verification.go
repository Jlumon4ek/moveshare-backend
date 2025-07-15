package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"

	"moveshare/internal/repository"
	"moveshare/internal/service"
)

// VerificationHandler handles HTTP requests for verification operations
type VerificationHandler struct {
	verificationService service.VerificationService
	minioClient         *minio.Client
	minioBucket         string
}

// NewVerificationHandler creates a new VerificationHandler
func NewVerificationHandler(verificationService service.VerificationService, minioClient *minio.Client, minioBucket string) *VerificationHandler {
	return &VerificationHandler{
		verificationService: verificationService,
		minioClient:         minioClient,
		minioBucket:         minioBucket,
	}
}

// Request/Response models
type UploadDocumentRequest struct {
	DocumentType string `json:"document_type" example:"mc_license"`
}

type DocumentResponse struct {
	Document *repository.VerificationDocument `json:"document"`
}

type DocumentsResponse struct {
	Documents []repository.VerificationDocument `json:"documents"`
}

type VerificationStatusResponse struct {
	Status *service.VerificationStatus `json:"status"`
}

type ReviewDocumentRequest struct {
	Status          string  `json:"status" example:"approved"`
	RejectionReason *string `json:"rejection_reason,omitempty"`
}

// UploadDocument godoc
// @Summary      Upload verification document
// @Description  Upload a verification document (MC License, DOT Certificate, Insurance Certificate, Business License)
// @Tags         verification
// @Accept       multipart/form-data
// @Produce      json
// @Param        document_type formData string true "Document type" Enums(mc_license, dot_certificate, insurance_certificate, business_license)
// @Param        file formData file true "Document file (PDF, JPG, JPEG, PNG)"
// @Success      201 {object} DocumentResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /verification/documents [post]
func (h *VerificationHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, `{"error": "failed to parse form"}`, http.StatusBadRequest)
		return
	}

	documentType := r.FormValue("document_type")
	if documentType == "" {
		http.Error(w, `{"error": "document_type is required"}`, http.StatusBadRequest)
		return
	}

	// Validate document type
	docType := repository.DocumentType(documentType)
	validTypes := map[repository.DocumentType]bool{
		repository.DocumentTypeMCLicense:       true,
		repository.DocumentTypeDOTCertificate:  true,
		repository.DocumentTypeInsuranceCert:   true,
		repository.DocumentTypeBusinessLicense: true,
	}
	if !validTypes[docType] {
		http.Error(w, `{"error": "invalid document type"}`, http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error": "file is required"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	if !h.isValidDocumentFileType(fileHeader.Filename) {
		http.Error(w, `{"error": "invalid file type. Only PDF, JPG, JPEG, PNG files are allowed"}`, http.StatusBadRequest)
		return
	}

	// Upload to Minio
	fileURL, err := h.uploadDocumentToMinio(file, fileHeader, userID, docType)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "failed to upload file: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Save to database
	err = h.verificationService.UploadDocument(r.Context(), userID, docType, fileHeader.Filename, fileURL)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// Get the uploaded document
	document, err := h.verificationService.GetDocumentByType(r.Context(), userID, docType)
	if err != nil {
		http.Error(w, `{"error": "failed to retrieve uploaded document"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(DocumentResponse{Document: document})
}

// GetUserDocuments godoc
// @Summary      Get user verification documents
// @Description  Get all verification documents for the authenticated user
// @Tags         verification
// @Produce      json
// @Success      200 {object} DocumentsResponse
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /verification/documents [get]
func (h *VerificationHandler) GetUserDocuments(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	documents, err := h.verificationService.GetUserDocuments(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "failed to get documents"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DocumentsResponse{Documents: documents})
}

// GetVerificationStatus godoc
// @Summary      Get verification status
// @Description  Get the overall verification status for the authenticated user
// @Tags         verification
// @Produce      json
// @Success      200 {object} VerificationStatusResponse
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /verification/status [get]
func (h *VerificationHandler) GetVerificationStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	status, err := h.verificationService.GetVerificationStatus(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "failed to get verification status"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(VerificationStatusResponse{Status: status})
}

// DeleteDocument godoc
// @Summary      Delete verification document
// @Description  Delete a verification document
// @Tags         verification
// @Param        id path int true "Document ID"
// @Success      204
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /verification/documents/{id} [delete]
func (h *VerificationHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	docIDStr := chi.URLParam(r, "id")
	docID, err := strconv.ParseInt(docIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid document ID"}`, http.StatusBadRequest)
		return
	}

	err = h.verificationService.DeleteDocument(r.Context(), userID, docID)
	if err != nil {
		http.Error(w, `{"error": "failed to delete document"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Admin endpoints

// GetPendingDocuments godoc
// @Summary      Get pending verification documents (Admin)
// @Description  Get all pending verification documents for admin review
// @Tags         verification
// @Produce      json
// @Param        limit query int false "Limit" default(20)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} DocumentsResponse
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /admin/verification/documents/pending [get]
func (h *VerificationHandler) GetPendingDocuments(w http.ResponseWriter, r *http.Request) {
	// Note: You'll need to implement admin role checking
	// userID, ok := r.Context().Value(userIDKey).(int64)
	// if !ok || !isAdmin(userID) {
	//     http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
	//     return
	// }

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	documents, err := h.verificationService.GetPendingDocuments(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, `{"error": "failed to get pending documents"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DocumentsResponse{Documents: documents})
}

// ReviewDocument godoc
// @Summary      Review verification document (Admin)
// @Description  Approve or reject a verification document
// @Tags         verification
// @Accept       json
// @Produce      json
// @Param        id path int true "Document ID"
// @Param        request body ReviewDocumentRequest true "Review request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /admin/verification/documents/{id}/review [post]
func (h *VerificationHandler) ReviewDocument(w http.ResponseWriter, r *http.Request) {
	// Note: You'll need to implement admin role checking
	reviewerID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	docIDStr := chi.URLParam(r, "id")
	docID, err := strconv.ParseInt(docIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid document ID"}`, http.StatusBadRequest)
		return
	}

	var req ReviewDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	status := repository.DocumentStatus(req.Status)
	if status != repository.DocumentStatusApproved && status != repository.DocumentStatusRejected {
		http.Error(w, `{"error": "status must be 'approved' or 'rejected'"}`, http.StatusBadRequest)
		return
	}

	err = h.verificationService.ReviewDocument(r.Context(), docID, status, req.RejectionReason, reviewerID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Document reviewed successfully"})
}

// Helper methods

// uploadDocumentToMinio uploads a document file to Minio
func (h *VerificationHandler) uploadDocumentToMinio(file multipart.File, fileHeader *multipart.FileHeader, userID int64, docType repository.DocumentType) (string, error) {
	ctx := context.Background()

	// Generate unique filename
	filename := fmt.Sprintf("verification/%d/%s_%d_%s",
		userID,
		docType,
		time.Now().UnixNano(),
		fileHeader.Filename)

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := h.minioClient.PutObject(ctx, h.minioBucket, filename, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	// Generate presigned URL (valid for 24 hours)
	url, err := h.minioClient.PresignedGetObject(ctx, h.minioBucket, filename, 24*time.Hour, nil)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

// isValidDocumentFileType checks if the file type is valid for documents
func (h *VerificationHandler) isValidDocumentFileType(filename string) bool {
	filename = strings.ToLower(filename)
	ext := filepath.Ext(filename)

	allowedExtensions := map[string]bool{
		".pdf":  true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	return allowedExtensions[ext]
}
