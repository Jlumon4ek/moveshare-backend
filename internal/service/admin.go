package service

import (
	"context"
	"moveshare/internal/models"
	"moveshare/internal/repository/admin"
)

type AdminService interface {
	GetUserCount(ctx context.Context) (int, error)
	GetChatConversationCount(ctx context.Context) (int, error)
	GetUsersList(ctx context.Context, limit, offset int) ([]models.UserCompanyInfo, error)
	GetActiveJobs(ctx context.Context, limit, offset int) ([]models.Job, error)
	ChangeUserStatus(ctx context.Context, userID int, newStatus string) error
	GetUserRole(ctx context.Context, userID int64) (string, error)
	ChangeVerificationFileStatus(ctx context.Context, fileID int, newStatus string) error
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

func (s *adminService) GetChatConversationCount(ctx context.Context) (int, error) {
	return s.adminRepo.GetChatConversationCount(ctx)
}

func (s *adminService) GetUsersList(ctx context.Context, limit, offset int) ([]models.UserCompanyInfo, error) {
	return s.adminRepo.GetUsersList(ctx, limit, offset)
}

func (s *adminService) GetActiveJobs(ctx context.Context, limit, offset int) ([]models.Job, error) {
	return s.adminRepo.GetAllJobs(ctx, limit, offset)
}

func (s *adminService) ChangeUserStatus(ctx context.Context, userID int, newStatus string) error {
	return s.adminRepo.ChangeUserStatus(ctx, userID, newStatus)
}

func (s *adminService) GetUserRole(ctx context.Context, userID int64) (string, error) {
	return s.adminRepo.GetUserRole(ctx, userID)
}

func (s *adminService) ChangeVerificationFileStatus(ctx context.Context, fileID int, newStatus string) error {
	return s.adminRepo.ChangeVerificationFileStatus(ctx, fileID, newStatus)
}
