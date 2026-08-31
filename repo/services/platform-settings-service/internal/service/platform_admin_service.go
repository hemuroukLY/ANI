package service

import (
	"context"

	commonv1 "github.com/kubercloud/ani/pkg/generated/pb/common/v1"
	platformsettingsv1 "github.com/kubercloud/ani/pkg/generated/pb/platform_settings/v1"
	"github.com/kubercloud/ani/services/platform-settings-service/internal/repo/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PlatformAdminService is the gRPC server skeleton for platform admin management.
type PlatformAdminService struct {
	platformsettingsv1.UnimplementedPlatformAdminServiceServer
}

var _ platformsettingsv1.PlatformAdminServiceServer = (*PlatformAdminService)(nil)

// NewPlatformAdminService returns a registerable gRPC server with UNIMPLEMENTED RPC bodies.
func NewPlatformAdminService() *PlatformAdminService {
	return &PlatformAdminService{}
}

// Register attaches this service to a gRPC server.
func (s *PlatformAdminService) Register(server *grpc.Server) {
	platformsettingsv1.RegisterPlatformAdminServiceServer(server, s)
}

func unimplemented() error {
	return status.Error(codes.Unimplemented, ports.ErrNotImplemented.Error())
}

func (s *PlatformAdminService) CreatePlatformAdmin(context.Context, *platformsettingsv1.CreatePlatformAdminRequest) (*commonv1.IdempotentResult, error) {
	return nil, unimplemented()
}

func (s *PlatformAdminService) ListPlatformAdmins(context.Context, *platformsettingsv1.ListPlatformAdminsRequest) (*platformsettingsv1.ListPlatformAdminsResponse, error) {
	return nil, unimplemented()
}

func (s *PlatformAdminService) ListPlatformAdminRoles(context.Context, *platformsettingsv1.ListPlatformAdminRolesRequest) (*platformsettingsv1.ListPlatformAdminRolesResponse, error) {
	return nil, unimplemented()
}

func (s *PlatformAdminService) GetPlatformAdmin(context.Context, *platformsettingsv1.GetPlatformAdminRequest) (*platformsettingsv1.PlatformAdminDetail, error) {
	return nil, unimplemented()
}

func (s *PlatformAdminService) UpdatePlatformAdminRole(context.Context, *platformsettingsv1.UpdatePlatformAdminRoleRequest) (*commonv1.IdempotentResult, error) {
	return nil, unimplemented()
}

func (s *PlatformAdminService) ResetPlatformAdminPassword(context.Context, *platformsettingsv1.ResetPlatformAdminPasswordRequest) (*commonv1.IdempotentResult, error) {
	return nil, unimplemented()
}

func (s *PlatformAdminService) DisablePlatformAdmin(context.Context, *platformsettingsv1.DisablePlatformAdminRequest) (*commonv1.IdempotentResult, error) {
	return nil, unimplemented()
}

func (s *PlatformAdminService) EnablePlatformAdmin(context.Context, *platformsettingsv1.EnablePlatformAdminRequest) (*commonv1.IdempotentResult, error) {
	return nil, unimplemented()
}

func (s *PlatformAdminService) DeletePlatformAdmin(context.Context, *platformsettingsv1.DeletePlatformAdminRequest) (*commonv1.IdempotentResult, error) {
	return nil, unimplemented()
}

func (s *PlatformAdminService) ListPlatformAdminAuditLogs(context.Context, *platformsettingsv1.ListPlatformAdminAuditLogsRequest) (*platformsettingsv1.ListPlatformAdminAuditLogsResponse, error) {
	return nil, unimplemented()
}
