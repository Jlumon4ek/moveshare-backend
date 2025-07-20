package handlers

import (
	"encoding/json"
	"moveshare/internal/repository"
	"moveshare/internal/service"
	"net/http"
)

type CompanyHandler struct {
	companyService service.CompanyService
}

func NewCompanyHandler(companyService service.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: companyService}
}

type CompanyResponse struct {
	Company *repository.Company `json:"company"`
}

type CompanyUpdateRequest struct {
	CompanyName        *string `json:"company_name"`
	EmailAddress       *string `json:"email_address"`
	Address            *string `json:"address"`
	State              *string `json:"state"`
	MCLicenseNumber    *string `json:"mc_license_number"`
	CompanyDescription *string `json:"company_description"`
	ContactPerson      *string `json:"contact_person"`
	PhoneNumber        *string `json:"phone_number"`
	City               *string `json:"city"`
	ZipCode            *string `json:"zip_code"`
	DotNumber          *string `json:"dot_number"`
}

// GetCompany handles fetching company data for the current user
// @Summary Get company data
// @Description Fetches company data for the current user
// @Tags Company Information
// @Produce json
// @Success 200 {object} CompanyResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /company [get]
func (h *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	company, err := h.companyService.GetCompany(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	resp := CompanyResponse{Company: company}
	w.Header().Set("Content-Type", "application/json")
	if company == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "company not found"})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// PatchCompany handles partial update of company data for the current user
// @Summary Update company data
// @Description Partially updates company data for the current user
// @Tags Company Information
// @Accept json
// @Produce json
// @Param body body CompanyUpdateRequest true "Company update data"
// @Success 200 {object} CompanyResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /company [patch]
func (h *CompanyHandler) PatchCompany(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req CompanyUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
		return
	}

	company := &repository.Company{UserID: userID}
	if req.CompanyName != nil {
		company.CompanyName = *req.CompanyName
	}
	if req.EmailAddress != nil {
		company.EmailAddress = *req.EmailAddress
	}
	if req.Address != nil {
		company.Address = *req.Address
	}
	if req.State != nil {
		company.State = *req.State
	}
	if req.MCLicenseNumber != nil {
		company.MCLicenseNumber = *req.MCLicenseNumber
	}
	if req.CompanyDescription != nil {
		company.CompanyDescription = *req.CompanyDescription
	}
	if req.ContactPerson != nil {
		company.ContactPerson = *req.ContactPerson
	}
	if req.PhoneNumber != nil {
		company.PhoneNumber = *req.PhoneNumber
	}
	if req.City != nil {
		company.City = *req.City
	}
	if req.ZipCode != nil {
		company.ZipCode = *req.ZipCode
	}
	if req.DotNumber != nil {
		company.DotNumber = *req.DotNumber
	}

	err := h.companyService.UpdateCompany(r.Context(), userID, company)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	updatedCompany, err := h.companyService.GetCompany(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	resp := CompanyResponse{Company: updatedCompany}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
