package service

import (
	"context"
	"moveshare/internal/models"
	"moveshare/internal/repository/admin"
)

type AdminService interface {
	GetUserCount(ctx context.Context) (int, error)
	GetPendingUsersCount(ctx context.Context) (int, error)
	GetChatConversationCount(ctx context.Context) (int, error)
	GetActiveJobsCount(ctx context.Context) (int, error)
	GetUsersList(ctx context.Context, limit, offset int) ([]models.UserCompanyInfo, error)
	GetUsersListTotal(ctx context.Context) (int, error)
	GetJobsList(ctx context.Context, limit, offset int, statuses []string) ([]models.JobManagementInfo, error)
	GetJobsListTotal(ctx context.Context, statuses []string) (int, error)
	// GetActiveJobs(ctx context.Context, limit, offset int) ([]models.Job, error)
	ChangeUserStatus(ctx context.Context, userID int, newStatus string) error
	GetUserRole(ctx context.Context, userID int64) (string, error)
	ChangeVerificationFileStatus(ctx context.Context, fileID int, newStatus string) error
	GetUserFullInfo(ctx context.Context, userID int64) (*models.UserFullInfo, error)
	GetPlatformAnalytics(ctx context.Context, days int) (*models.PlatformAnalytics, error)
	GetSystemSettings(ctx context.Context) (*models.SystemSettings, error)
	UpdateSystemSettings(ctx context.Context, settings *models.SystemSettings) error
}

type adminService struct {
	adminRepo admin.AdminRepository
}

func NewAdminService(adminRepo admin.AdminRepository) AdminService {
	return &adminService{
		adminRepo: adminRepo,
	}
}

func (s *adminService) GetUserCount(ctx context.Context) (int, error) {
	return s.adminRepo.GetUserCount(ctx)
}

func (s *adminService) GetPendingUsersCount(ctx context.Context) (int, error) {
	return s.adminRepo.GetPendingUsersCount(ctx)
}

func (s *adminService) GetChatConversationCount(ctx context.Context) (int, error) {
	return s.adminRepo.GetChatConversationCount(ctx)
}

func (s *adminService) GetActiveJobsCount(ctx context.Context) (int, error) {
	return s.adminRepo.GetActiveJobsCount(ctx)
}

func (s *adminService) GetUsersList(ctx context.Context, limit, offset int) ([]models.UserCompanyInfo, error) {
	return s.adminRepo.GetUsersList(ctx, limit, offset)
}

func (s *adminService) GetJobsList(ctx context.Context, limit, offset int, statuses []string) ([]models.JobManagementInfo, error) {
	return s.adminRepo.GetJobsList(ctx, limit, offset, statuses)
}

func (s *adminService) GetUsersListTotal(ctx context.Context) (int, error) {
	return s.adminRepo.GetUsersListTotal(ctx)
}

func (s *adminService) GetJobsListTotal(ctx context.Context, statuses []string) (int, error) {
	return s.adminRepo.GetJobsListTotal(ctx, statuses)
}

// func (s *adminService) GetActiveJobs(ctx context.Context, limit, offset int) ([]models.Job, error) {
// 	return s.adminRepo.GetAllJobs(ctx, limit, offset)
// }


func (s *adminService) ChangeUserStatus(ctx context.Context, userID int, newStatus string) error {
	return s.adminRepo.ChangeUserStatus(ctx, userID, newStatus)
}

func (s *adminService) GetUserRole(ctx context.Context, userID int64) (string, error) {
	return s.adminRepo.GetUserRole(ctx, userID)
}

func (s *adminService) ChangeVerificationFileStatus(ctx context.Context, fileID int, newStatus string) error {
	return s.adminRepo.ChangeVerificationFileStatus(ctx, fileID, newStatus)
}

func (s *adminService) GetUserFullInfo(ctx context.Context, userID int64) (*models.UserFullInfo, error) {
	return s.adminRepo.GetUserFullInfo(ctx, userID)
}

func (s *adminService) GetPlatformAnalytics(ctx context.Context, days int) (*models.PlatformAnalytics, error) {
	// Get top companies (limit to 5)
	topCompanies, err := s.adminRepo.GetTopCompanies(ctx, days, 5)
	if err != nil {
		return nil, err
	}

	// Get busiest routes (limit to 5)
	busiestRoutes, err := s.adminRepo.GetBusiestRoutes(ctx, days, 5)
	if err != nil {
		return nil, err
	}

	return &models.PlatformAnalytics{
		TopCompanies:  topCompanies,
		BusiestRoutes: busiestRoutes,
	}, nil
}

func (s *adminService) GetSystemSettings(ctx context.Context) (*models.SystemSettings, error) {
	return s.adminRepo.GetSystemSettings(ctx)
}

func (s *adminService) UpdateSystemSettings(ctx context.Context, settings *models.SystemSettings) error {
	return s.adminRepo.UpdateSystemSettings(ctx, settings)
}
