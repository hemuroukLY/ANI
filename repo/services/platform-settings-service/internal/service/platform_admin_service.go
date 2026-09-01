package service

import (
	"context"
	"strings"

	commonv1 "github.com/kubercloud/ani/pkg/generated/pb/common/v1"
	platformsettingsv1 "github.com/kubercloud/ani/pkg/generated/pb/platform_settings/v1"
	"github.com/kubercloud/ani/services/platform-settings-service/internal/repo/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PlatformAdminService is the gRPC server for platform admin management.
type PlatformAdminService struct {
	platformsettingsv1.UnimplementedPlatformAdminServiceServer
	coreClient ports.CorePlatformUserClient
	auditStore ports.PlatformAdminAuditStore
}

var _ platformsettingsv1.PlatformAdminServiceServer = (*PlatformAdminService)(nil)

// NewPlatformAdminService constructs the service with Core client + audit store injected.
func NewPlatformAdminService(coreClient ports.CorePlatformUserClient, auditStore ports.PlatformAdminAuditStore) *PlatformAdminService {
	return &PlatformAdminService{
		coreClient: coreClient,
		auditStore: auditStore,
	}
}

// Register attaches this service to a gRPC server.
func (s *PlatformAdminService) Register(server *grpc.Server) {
	platformsettingsv1.RegisterPlatformAdminServiceServer(server, s)
}

func (s *PlatformAdminService) CreatePlatformAdmin(ctx context.Context, req *platformsettingsv1.CreatePlatformAdminRequest) (*commonv1.IdempotentResult, error) {
	const action = "platform_admin.create"
	if req == nil {
		err := businessError(codes.InvalidArgument, ports.ErrValidationFailed, "request required")
		writeAuditFailure(ctx, s.auditStore, action, nil, err)
		return nil, err
	}

	// 步骤 1：入参校验；失败审计写入 role（不含 password）
	// 幂等由网关层处理，本服务不消费 idempotency_key。
	if detail := validateCreatePlatformAdmin(req.GetEmail(), req.GetUsername(), req.GetDisplayName(), req.GetRole(), req.GetPassword()); detail != "" {
		err := businessError(codes.InvalidArgument, ports.ErrValidationFailed, detail)
		writeAuditFailure(ctx, s.auditStore, action, map[string]any{"role": strings.TrimSpace(req.GetRole())}, err)
		return nil, err
	}
	if s.coreClient == nil {
		err := businessError(codes.Unavailable, ports.ErrCoreUnavailable, "core client not configured")
		writeAuditFailure(ctx, s.auditStore, action, nil, err)
		return nil, err
	}

	// 步骤 2：调 Core 创建（冲突透传；失败审计不含 role）
	createdID, err := s.coreClient.Create(ctx, ports.PlatformUserCreateInput{
		Email:       strings.TrimSpace(req.GetEmail()),
		Username:    strings.TrimSpace(req.GetUsername()),
		DisplayName: strings.TrimSpace(req.GetDisplayName()),
		Role:        strings.TrimSpace(req.GetRole()),
		Password:    req.GetPassword(),
	})
	if err != nil {
		mapped := mapDomainError(err)
		writeAuditFailure(ctx, s.auditStore, action, nil, mapped)
		return nil, mapped
	}

	// 步骤 3：成功审计仅 target_id（不含 password/role）
	writeAuditSuccess(ctx, s.auditStore, action, map[string]any{
		"target_id": createdID,
	})

	return &commonv1.IdempotentResult{
		Id:      createdID,
		Message: "platform admin created",
	}, nil
}

func (s *PlatformAdminService) ListPlatformAdmins(ctx context.Context, req *platformsettingsv1.ListPlatformAdminsRequest) (*platformsettingsv1.ListPlatformAdminsResponse, error) {
	if s.coreClient == nil {
		return nil, businessError(codes.Unavailable, ports.ErrCoreUnavailable, "core client not configured")
	}
	limit := 20
	cursor := ""
	if req != nil && req.GetPage() != nil {
		if l := int(req.GetPage().GetLimit()); l > 0 {
			limit = l
			if limit > 100 {
				limit = 100
			}
		}
		cursor = req.GetPage().GetCursor()
	}
	filter := ports.PlatformUserListFilter{Limit: limit, Cursor: cursor}
	if req != nil {
		filter.Role = req.GetRole()
		filter.Status = req.GetStatus()
		filter.Source = req.GetSource()
		filter.Search = req.GetSearch()
	}
	res, err := s.coreClient.List(ctx, filter)
	if err != nil {
		return nil, mapDomainError(err)
	}
	items := make([]*platformsettingsv1.PlatformAdminListItem, 0, len(res.Items))
	for _, it := range res.Items {
		item := &platformsettingsv1.PlatformAdminListItem{
			Id:       it.ID,
			Username: it.Username,
			Role:     it.Role,
			Status:   it.Status,
			Source:   it.Source,
		}
		if it.DisplayName != nil {
			item.DisplayName = *it.DisplayName
		}
		if it.LastLoginAt != nil {
			item.LastLoginAt = timestamppb.New(*it.LastLoginAt)
		}
		items = append(items, item)
	}
	return &platformsettingsv1.ListPlatformAdminsResponse{
		Items:      items,
		NextCursor: res.NextCursor,
	}, nil
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
