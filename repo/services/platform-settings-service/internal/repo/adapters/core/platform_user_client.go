package core

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	anisdk "github.com/kubercloud/ani-sdks/core-go/anisdk"
	"github.com/kubercloud/ani/services/platform-settings-service/internal/repo/ports"
)

// CorePlatformUserClient 基于 Core Go SDK 实现 ports.CorePlatformUserClient。
type CorePlatformUserClient struct {
	sdk anisdk.Client
}

var _ ports.CorePlatformUserClient = (*CorePlatformUserClient)(nil)

// NewCorePlatformUserClient 从环境变量构造 Core 平台用户 API 客户端。
func NewCorePlatformUserClient() ports.CorePlatformUserClient {
	return &CorePlatformUserClient{sdk: newCoreSDKClient()}
}

// Create 调用 Core POST /admin/platform-users，成功后 GET 详情补齐 DTO。
func (c *CorePlatformUserClient) Create(ctx context.Context, in ports.PlatformUserCreateInput) (ports.PlatformUserDTO, error) {
	_ = ctx
	// 步骤 1：调用 Core 创建平台账号
	raw, err := c.sdk.Request("POST", "/admin/platform-users", anisdk.RequestOptions{
		Body: map[string]any{
			"email":        in.Email,
			"username":     in.Username,
			"display_name": in.DisplayName,
			"role":         in.Role,
			"password":     in.Password,
		},
	})
	if err != nil {
		return ports.PlatformUserDTO{}, mapSDKError(err)
	}
	// 步骤 2：从创建响应取出 user id
	obj, err := asObject(raw)
	if err != nil {
		return ports.PlatformUserDTO{}, err
	}
	id := stringField(obj, "id")
	if id == "" {
		return ports.PlatformUserDTO{}, fmt.Errorf("%w: create response missing id", ports.ErrCoreUnavailable)
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return ports.PlatformUserDTO{}, fmt.Errorf("%w: invalid create id: %v", ports.ErrCoreUnavailable, err)
	}
	// 步骤 3：回查详情补齐完整 DTO（创建接口仅返回 id/message）
	return c.Get(ctx, uid)
}

// List 调用 Core GET /admin/platform-users（游标 + 过滤参数透传）。
func (c *CorePlatformUserClient) List(ctx context.Context, filter ports.PlatformUserListFilter) (ports.PlatformUserListDTO, error) {
	_ = ctx
	// 步骤 1：组装游标与过滤 query
	params := anisdk.CursorParams(filter.Limit, filter.Cursor)
	if filter.Role != "" {
		params["role"] = filter.Role
	}
	if filter.Status != "" {
		params["status"] = filter.Status
	}
	if filter.Source != "" {
		params["source"] = filter.Source
	}
	if filter.Search != "" {
		params["search"] = filter.Search
	}
	// 步骤 2：调用 Core 列表接口
	raw, err := c.sdk.Request("GET", "/admin/platform-users", anisdk.RequestOptions{Params: params})
	if err != nil {
		return ports.PlatformUserListDTO{}, mapSDKError(err)
	}
	// 步骤 3：解码 items / next_cursor
	obj, err := asObject(raw)
	if err != nil {
		return ports.PlatformUserListDTO{}, err
	}
	itemsRaw, err := asObjectSlice(obj["items"])
	if err != nil {
		return ports.PlatformUserListDTO{}, err
	}
	items := make([]ports.PlatformUserDTO, 0, len(itemsRaw))
	for _, it := range itemsRaw {
		dto, err := decodePlatformUser(it)
		if err != nil {
			return ports.PlatformUserListDTO{}, err
		}
		items = append(items, dto)
	}
	return ports.PlatformUserListDTO{
		Items:      items,
		NextCursor: stringField(obj, "next_cursor"),
	}, nil
}

// Get 调用 Core GET /admin/platform-users/{userId}。
func (c *CorePlatformUserClient) Get(ctx context.Context, userID uuid.UUID) (ports.PlatformUserDTO, error) {
	_ = ctx
	// 步骤 1：拼路径并调用 Core 详情接口
	path := fmt.Sprintf("/admin/platform-users/%s", userID.String())
	raw, err := c.sdk.Request("GET", path, anisdk.RequestOptions{})
	if err != nil {
		return ports.PlatformUserDTO{}, mapSDKError(err)
	}
	// 步骤 2：解码为 PlatformUserDTO
	return decodePlatformUser(raw)
}

// ChangeRole 调用 Core PUT /admin/platform-users/{userId}/role。
func (c *CorePlatformUserClient) ChangeRole(ctx context.Context, userID uuid.UUID, role string) error {
	_ = ctx
	// 步骤 1：拼路径并提交新角色
	path := fmt.Sprintf("/admin/platform-users/%s/role", userID.String())
	_, err := c.sdk.Request("PUT", path, anisdk.RequestOptions{
		Body: map[string]any{"role": role},
	})
	// 步骤 2：映射 Core 错误码为领域哨兵
	return mapSDKError(err)
}

// ResetPassword 调用 Core POST /admin/platform-users/{userId}/reset-password。
func (c *CorePlatformUserClient) ResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	_ = ctx
	// 步骤 1：拼路径并提交新密码（明文仅本次透传）
	path := fmt.Sprintf("/admin/platform-users/%s/reset-password", userID.String())
	_, err := c.sdk.Request("POST", path, anisdk.RequestOptions{
		Body: map[string]any{"new_password": newPassword},
	})
	// 步骤 2：映射 Core 错误码为领域哨兵
	return mapSDKError(err)
}

// SetStatus 调用 Core disable/enable 端点（status=active → enable，否则 disable）。
func (c *CorePlatformUserClient) SetStatus(ctx context.Context, userID uuid.UUID, status string) error {
	_ = ctx
	// 步骤 1：按目标状态选择 disable / enable
	action := "disable"
	if status == "active" {
		action = "enable"
	}
	// 步骤 2：调用对应 Core 写接口
	path := fmt.Sprintf("/admin/platform-users/%s/%s", userID.String(), action)
	_, err := c.sdk.Request("POST", path, anisdk.RequestOptions{})
	// 步骤 3：映射 Core 错误码为领域哨兵
	return mapSDKError(err)
}

// SoftDelete 调用 Core DELETE /admin/platform-users/{userId}。
func (c *CorePlatformUserClient) SoftDelete(ctx context.Context, userID uuid.UUID) error {
	_ = ctx
	// 步骤 1：拼路径并调用 Core 软删除
	path := fmt.Sprintf("/admin/platform-users/%s", userID.String())
	_, err := c.sdk.Request("DELETE", path, anisdk.RequestOptions{})
	// 步骤 2：映射 Core 错误码为领域哨兵
	return mapSDKError(err)
}

// ListPlatformRoles 待 #7 增补 Core roles API / SDK operation 后实现。
func (c *CorePlatformUserClient) ListPlatformRoles(ctx context.Context) ([]ports.PlatformRoleDTO, error) {
	_ = ctx
	// TODO(issue-007): wire GET /admin/platform-users/roles once Core SDK operation exists.
	return nil, ports.ErrNotImplemented
}
