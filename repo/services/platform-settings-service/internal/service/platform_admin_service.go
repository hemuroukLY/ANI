package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
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
	// 只读：不写审计。
	if s.coreClient == nil {
		return nil, businessError(codes.Unavailable, ports.ErrCoreUnavailable, "core client not configured")
	}
	// 步骤 1：分页默认 limit=20，上限 100；校验 status/source 枚举
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
	if req != nil {
		if detail := validateListPlatformAdminFilters(req.GetStatus(), req.GetSource()); detail != "" {
			return nil, businessError(codes.InvalidArgument, ports.ErrValidationFailed, detail)
		}
	}
	// 步骤 2：组装过滤条件（source 过滤枚举 local|third_party，透传 Core）
	filter := ports.PlatformUserListFilter{Limit: limit, Cursor: cursor}
	if req != nil {
		filter.Role = req.GetRole()
		filter.Status = req.GetStatus()
		filter.Source = req.GetSource()
		filter.Search = req.GetSearch()
	}
	// 步骤 3：调 Core List
	res, err := s.coreClient.List(ctx, filter)
	if err != nil {
		return nil, mapDomainError(err)
	}
	// 步骤 4：映射列表项（不含 email；剥 username 前缀；source 推断）
	items := make([]*platformsettingsv1.PlatformAdminListItem, 0, len(res.Items))
	for _, it := range res.Items {
		items = append(items, toPlatformAdminListItem(it))
	}
	return &platformsettingsv1.ListPlatformAdminsResponse{
		Items:      items,
		NextCursor: res.NextCursor,
	}, nil
}

func (s *PlatformAdminService) ListPlatformAdminRoles(context.Context, *platformsettingsv1.ListPlatformAdminRolesRequest) (*platformsettingsv1.ListPlatformAdminRolesResponse, error) {
	return nil, unimplemented()
}

func (s *PlatformAdminService) GetPlatformAdmin(ctx context.Context, req *platformsettingsv1.GetPlatformAdminRequest) (*platformsettingsv1.PlatformAdminDetail, error) {
	// 只读：不写审计。
	if s.coreClient == nil {
		return nil, businessError(codes.Unavailable, ports.ErrCoreUnavailable, "core client not configured")
	}
	if req == nil || strings.TrimSpace(req.GetUserId()) == "" {
		return nil, businessError(codes.InvalidArgument, ports.ErrValidationFailed, "user_id required")
	}
	// 步骤 1：校验 userId UUID
	userID, err := uuid.Parse(strings.TrimSpace(req.GetUserId()))
	if err != nil {
		return nil, businessError(codes.InvalidArgument, ports.ErrValidationFailed, "user_id must be a uuid")
	}
	// 步骤 2：调 Core Get；不存在 → PLATFORM_USER_NOT_FOUND
	dto, err := s.coreClient.Get(ctx, userID)
	if err != nil {
		return nil, mapDomainError(err)
	}
	// 步骤 3：映射详情全字段（含 email/created_at；不含 password）
	return toPlatformAdminDetail(dto), nil
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

func toPlatformAdminListItem(it ports.PlatformUserDTO) *platformsettingsv1.PlatformAdminListItem {
	item := &platformsettingsv1.PlatformAdminListItem{
		Id:       it.ID,
		Username: stripPlatformUsernamePrefix(it.Username),
		Role:     it.Role,
		Status:   it.Status,
		Source:   normalizePlatformSource(it.Username, it.Source),
	}
	if it.DisplayName != nil {
		item.DisplayName = *it.DisplayName
	}
	if it.LastLoginAt != nil {
		item.LastLoginAt = timestamppb.New(*it.LastLoginAt)
	}
	return item
}

func toPlatformAdminDetail(it ports.PlatformUserDTO) *platformsettingsv1.PlatformAdminDetail {
	detail := &platformsettingsv1.PlatformAdminDetail{
		Id:        it.ID,
		Email:     it.Email,
		Username:  stripPlatformUsernamePrefix(it.Username),
		Role:      it.Role,
		Status:    it.Status,
		Source:    normalizePlatformSource(it.Username, it.Source),
		CreatedAt: timestamppb.New(it.CreatedAt),
	}
	if it.DisplayName != nil {
		detail.DisplayName = *it.DisplayName
	}
	if it.LastLoginAt != nil {
		detail.LastLoginAt = timestamppb.New(*it.LastLoginAt)
	}
	return detail
}

// stripPlatformUsernamePrefix 对外响应剥除 local:/oidc: 存储前缀。
func stripPlatformUsernamePrefix(username string) string {
	switch {
	case strings.HasPrefix(username, "local:"):
		return strings.TrimPrefix(username, "local:")
	case strings.HasPrefix(username, "oidc:"):
		return strings.TrimPrefix(username, "oidc:")
	default:
		return username
	}
}

// normalizePlatformSource：oidc: → third_party，local: → local；已有合法 source 优先保留。
func normalizePlatformSource(username, source string) string {
	switch strings.TrimSpace(source) {
	case "local", "third_party", "unknown":
		return strings.TrimSpace(source)
	}
	switch {
	case strings.HasPrefix(username, "oidc:"):
		return "third_party"
	case strings.HasPrefix(username, "local:"):
		return "local"
	default:
		return "unknown"
	}
}
