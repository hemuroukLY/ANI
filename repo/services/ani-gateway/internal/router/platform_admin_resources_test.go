package router

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol"
	commonv1 "github.com/kubercloud/ani/pkg/generated/pb/common/v1"
	platformsettingsv1 "github.com/kubercloud/ani/pkg/generated/pb/platform_settings/v1"
	"github.com/kubercloud/ani/services/ani-gateway/internal/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakePlatformAdminClient struct {
	lastList   *platformsettingsv1.ListPlatformAdminsRequest
	lastGetID  string
	lastCreate *platformsettingsv1.CreatePlatformAdminRequest
	listResp   *platformsettingsv1.ListPlatformAdminsResponse
	rolesResp  *platformsettingsv1.ListPlatformAdminRolesResponse
	getResp    *platformsettingsv1.PlatformAdminDetail
	err        error
}

func (f *fakePlatformAdminClient) CreatePlatformAdmin(_ context.Context, in *platformsettingsv1.CreatePlatformAdminRequest, _ ...grpc.CallOption) (*commonv1.IdempotentResult, error) {
	f.lastCreate = in
	if f.err != nil {
		return nil, f.err
	}
	return &commonv1.IdempotentResult{Id: "11111111-1111-1111-1111-111111111111", Message: "platform admin created"}, nil
}
func (f *fakePlatformAdminClient) ListPlatformAdmins(_ context.Context, in *platformsettingsv1.ListPlatformAdminsRequest, _ ...grpc.CallOption) (*platformsettingsv1.ListPlatformAdminsResponse, error) {
	f.lastList = in
	if f.err != nil {
		return nil, f.err
	}
	return f.listResp, nil
}
func (f *fakePlatformAdminClient) ListPlatformAdminRoles(context.Context, *platformsettingsv1.ListPlatformAdminRolesRequest, ...grpc.CallOption) (*platformsettingsv1.ListPlatformAdminRolesResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.rolesResp != nil {
		return f.rolesResp, nil
	}
	return &platformsettingsv1.ListPlatformAdminRolesResponse{}, nil
}
func (f *fakePlatformAdminClient) GetPlatformAdmin(_ context.Context, in *platformsettingsv1.GetPlatformAdminRequest, _ ...grpc.CallOption) (*platformsettingsv1.PlatformAdminDetail, error) {
	f.lastGetID = in.GetUserId()
	if f.err != nil {
		return nil, f.err
	}
	if f.getResp != nil {
		return f.getResp, nil
	}
	return &platformsettingsv1.PlatformAdminDetail{Id: in.GetUserId(), Username: "local:ops"}, nil
}
func (f *fakePlatformAdminClient) UpdatePlatformAdminRole(context.Context, *platformsettingsv1.UpdatePlatformAdminRoleRequest, ...grpc.CallOption) (*commonv1.IdempotentResult, error) {
	return f.errOrOK()
}
func (f *fakePlatformAdminClient) ResetPlatformAdminPassword(context.Context, *platformsettingsv1.ResetPlatformAdminPasswordRequest, ...grpc.CallOption) (*commonv1.IdempotentResult, error) {
	return f.errOrOK()
}
func (f *fakePlatformAdminClient) DisablePlatformAdmin(context.Context, *platformsettingsv1.DisablePlatformAdminRequest, ...grpc.CallOption) (*commonv1.IdempotentResult, error) {
	return f.errOrOK()
}
func (f *fakePlatformAdminClient) EnablePlatformAdmin(context.Context, *platformsettingsv1.EnablePlatformAdminRequest, ...grpc.CallOption) (*commonv1.IdempotentResult, error) {
	return f.errOrOK()
}
func (f *fakePlatformAdminClient) DeletePlatformAdmin(context.Context, *platformsettingsv1.DeletePlatformAdminRequest, ...grpc.CallOption) (*commonv1.IdempotentResult, error) {
	return f.errOrOK()
}
func (f *fakePlatformAdminClient) ListPlatformAdminAuditLogs(context.Context, *platformsettingsv1.ListPlatformAdminAuditLogsRequest, ...grpc.CallOption) (*platformsettingsv1.ListPlatformAdminAuditLogsResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &platformsettingsv1.ListPlatformAdminAuditLogsResponse{}, nil
}
func (f *fakePlatformAdminClient) errOrOK() (*commonv1.IdempotentResult, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &commonv1.IdempotentResult{Id: "u1", Message: "ok"}, nil
}

func setupPlatformAdminTestServer(t *testing.T, client platformsettingsv1.PlatformAdminServiceClient) *server.Hertz {
	t.Helper()
	h := server.Default()
	h.Use(middleware.RequestID())
	registerPlatformAdminsAPI(h.Group("/api/v1/svc"), &platformAdminsAPI{client: client})
	return h
}

func performPlatformAdmin(h *server.Hertz, method, path, body string) *protocol.Response {
	var bodyArg *ut.Body
	if body != "" {
		bodyArg = &ut.Body{Body: strings.NewReader(body), Len: len(body)}
	}
	headers := []ut.Header{{Key: "Content-Type", Value: "application/json"}}
	return ut.PerformRequest(h.Engine, method, path, bodyArg, headers...).Result()
}

func TestMapPlatformAdminError_Table(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name   string
		err    error
		status int
		code   string
	}{
		{"validation", status.Error(codes.InvalidArgument, "VALIDATION_FAILED: bad"), http.StatusBadRequest, "VALIDATION_FAILED"},
		{"not found", status.Error(codes.NotFound, "PLATFORM_USER_NOT_FOUND: x"), http.StatusNotFound, "PLATFORM_USER_NOT_FOUND"},
		{"role not found", status.Error(codes.NotFound, "ROLE_NOT_FOUND"), http.StatusNotFound, "ROLE_NOT_FOUND"},
		{"username", status.Error(codes.AlreadyExists, "USERNAME_ALREADY_EXISTS"), http.StatusConflict, "USERNAME_ALREADY_EXISTS"},
		{"already exists fallback", status.Error(codes.AlreadyExists, "conflict"), http.StatusConflict, "USERNAME_ALREADY_EXISTS"},
		{"last admin", status.Error(codes.FailedPrecondition, "LAST_PLATFORM_ADMIN"), http.StatusUnprocessableEntity, "LAST_PLATFORM_ADMIN"},
		{"password", status.Error(codes.FailedPrecondition, "PASSWORD_SAME_AS_OLD"), http.StatusUnprocessableEntity, "PASSWORD_SAME_AS_OLD"},
		{"role change", status.Error(codes.FailedPrecondition, "ROLE_CHANGE_INVALID"), http.StatusUnprocessableEntity, "ROLE_CHANGE_INVALID"},
		{"core", status.Error(codes.Unavailable, "CORE_UNAVAILABLE: down"), http.StatusBadGateway, "CORE_UNAVAILABLE"},
		{"unimplemented", status.Error(codes.Unimplemented, "NOT_IMPLEMENTED"), http.StatusNotImplemented, "NOT_IMPLEMENTED"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := &app.RequestContext{}
			mapPlatformAdminError(c, tc.err)
			if c.Response.StatusCode() != tc.status {
				t.Fatalf("status=%d want %d", c.Response.StatusCode(), tc.status)
			}
			body := string(c.Response.Body())
			if !strings.Contains(body, tc.code) {
				t.Fatalf("body=%s", body)
			}
		})
	}
}

func TestPlatformAdmins_ListForwardParams(t *testing.T) {
	fake := &fakePlatformAdminClient{
		listResp: &platformsettingsv1.ListPlatformAdminsResponse{
			Items: []*platformsettingsv1.PlatformAdminListItem{
				{Id: "u1", Username: "local:ops", Role: "platform-ops", Status: "active", Source: "local"},
			},
			NextCursor: "n1",
		},
	}
	h := setupPlatformAdminTestServer(t, fake)
	resp := performPlatformAdmin(h, http.MethodGet, "/api/v1/svc/platform-admins?limit=7&cursor=c1&role=platform-ops&status=active&source=local&search=ops", "")
	if resp.StatusCode() != http.StatusOK {
		t.Fatalf("status=%d body=%s", resp.StatusCode(), resp.Body())
	}
	if fake.lastList == nil || fake.lastList.GetRole() != "platform-ops" || fake.lastList.GetPage().GetLimit() != 7 ||
		fake.lastList.GetPage().GetCursor() != "c1" || fake.lastList.GetSearch() != "ops" {
		t.Fatalf("lastList=%+v", fake.lastList)
	}
	var payload map[string]any
	if err := json.Unmarshal(resp.Body(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["next_cursor"] != "n1" {
		t.Fatalf("payload=%v", payload)
	}
}

func TestPlatformAdmins_RolesDoesNotCollideWithUserID(t *testing.T) {
	fake := &fakePlatformAdminClient{
		rolesResp: &platformsettingsv1.ListPlatformAdminRolesResponse{
			Items: []*platformsettingsv1.PlatformRole{{Name: "platform-admin", Label: "Admin"}},
		},
	}
	h := setupPlatformAdminTestServer(t, fake)

	roles := performPlatformAdmin(h, http.MethodGet, "/api/v1/svc/platform-admins/roles", "")
	if roles.StatusCode() != http.StatusOK {
		t.Fatalf("roles status=%d body=%s", roles.StatusCode(), roles.Body())
	}
	if fake.lastGetID != "" {
		t.Fatalf("roles hit get with userId=%s", fake.lastGetID)
	}
	var rolesPayload map[string]any
	_ = json.Unmarshal(roles.Body(), &rolesPayload)
	items, _ := rolesPayload["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("roles payload=%v", rolesPayload)
	}

	get := performPlatformAdmin(h, http.MethodGet, "/api/v1/svc/platform-admins/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "")
	if get.StatusCode() != http.StatusOK {
		t.Fatalf("get status=%d body=%s", get.StatusCode(), get.Body())
	}
	if fake.lastGetID != "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa" {
		t.Fatalf("lastGetID=%s", fake.lastGetID)
	}
}

func TestPlatformAdmins_CreateForwardBody(t *testing.T) {
	fake := &fakePlatformAdminClient{}
	h := setupPlatformAdminTestServer(t, fake)
	body := `{"email":"a@x.com","username":"ops","display_name":"Ops","role":"platform-ops","password":"Abcd1234!","idempotency_key":"44444444-4444-4444-4444-444444444444"}`
	resp := performPlatformAdmin(h, http.MethodPost, "/api/v1/svc/platform-admins", body)
	if resp.StatusCode() != http.StatusOK {
		t.Fatalf("status=%d body=%s", resp.StatusCode(), resp.Body())
	}
	if fake.lastCreate == nil || fake.lastCreate.GetEmail() != "a@x.com" || fake.lastCreate.GetIdempotencyKey() == "" {
		t.Fatalf("lastCreate=%+v", fake.lastCreate)
	}
}

func TestHandler_CreateFlow(t *testing.T) {
	const id = "11111111-1111-1111-1111-111111111111"
	fake := &fakePlatformAdminClient{
		listResp: &platformsettingsv1.ListPlatformAdminsResponse{
			Items: []*platformsettingsv1.PlatformAdminListItem{
				{Id: id, Username: "ops", DisplayName: "Ops", Role: "platform-ops", Status: "active", Source: "local"},
			},
		},
		getResp: &platformsettingsv1.PlatformAdminDetail{
			Id: id, Email: "ops@ani.io", Username: "ops", DisplayName: "Ops",
			Role: "platform-ops", Status: "active", Source: "local",
		},
	}
	h := setupPlatformAdminTestServer(t, fake)

	createBody := `{"email":"ops@ani.io","username":"ops","display_name":"Ops","role":"platform-ops","password":"Abcd1234!","idempotency_key":"44444444-4444-4444-4444-444444444444"}`
	createResp := performPlatformAdmin(h, http.MethodPost, "/api/v1/svc/platform-admins", createBody)
	if createResp.StatusCode() != http.StatusOK {
		t.Fatalf("create status=%d body=%s", createResp.StatusCode(), createResp.Body())
	}
	var created map[string]any
	if err := json.Unmarshal(createResp.Body(), &created); err != nil {
		t.Fatal(err)
	}
	if created["id"] != id || created["message"] != "platform admin created" {
		t.Fatalf("created=%v", created)
	}
	if fake.lastCreate == nil || fake.lastCreate.GetPassword() == "" || fake.lastCreate.GetIdempotencyKey() == "" {
		t.Fatalf("lastCreate=%+v", fake.lastCreate)
	}

	listResp := performPlatformAdmin(h, http.MethodGet, "/api/v1/svc/platform-admins?role=platform-ops", "")
	if listResp.StatusCode() != http.StatusOK {
		t.Fatalf("list status=%d body=%s", listResp.StatusCode(), listResp.Body())
	}
	var listed map[string]any
	_ = json.Unmarshal(listResp.Body(), &listed)
	items, _ := listed["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("list=%v", listed)
	}
	first, _ := items[0].(map[string]any)
	if _, hasEmail := first["email"]; hasEmail {
		t.Fatalf("list item must not include email: %v", first)
	}
	for _, key := range []string{"id", "username", "display_name", "role", "status", "source"} {
		if first[key] == nil || first[key] == "" {
			t.Fatalf("missing list field %s: %v", key, first)
		}
	}
	if first["username"] != "ops" {
		t.Fatalf("username=%v", first["username"])
	}

	detailResp := performPlatformAdmin(h, http.MethodGet, "/api/v1/svc/platform-admins/"+id, "")
	if detailResp.StatusCode() != http.StatusOK {
		t.Fatalf("detail status=%d body=%s", detailResp.StatusCode(), detailResp.Body())
	}
	var detail map[string]any
	_ = json.Unmarshal(detailResp.Body(), &detail)
	for _, key := range []string{"id", "email", "username", "display_name", "role", "status", "source"} {
		if detail[key] == nil || detail[key] == "" {
			t.Fatalf("missing detail field %s: %v", key, detail)
		}
	}
	if detail["email"] != "ops@ani.io" || detail["role"] != "platform-ops" || detail["username"] != "ops" {
		t.Fatalf("detail=%v", detail)
	}
}

func TestPlatformAdmins_ListInvalidLimit(t *testing.T) {
	fake := &fakePlatformAdminClient{}
	h := setupPlatformAdminTestServer(t, fake)
	resp := performPlatformAdmin(h, http.MethodGet, "/api/v1/svc/platform-admins?limit=abc", "")
	if resp.StatusCode() != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", resp.StatusCode(), resp.Body())
	}
	if !strings.Contains(string(resp.Body()), "VALIDATION_FAILED") {
		t.Fatalf("body=%s", resp.Body())
	}
	if fake.lastList != nil {
		t.Fatalf("unexpected list call on invalid limit")
	}
}

func TestPlatformAdmins_ListEmpty(t *testing.T) {
	fake := &fakePlatformAdminClient{
		listResp: &platformsettingsv1.ListPlatformAdminsResponse{Items: nil, NextCursor: ""},
	}
	h := setupPlatformAdminTestServer(t, fake)
	resp := performPlatformAdmin(h, http.MethodGet, "/api/v1/svc/platform-admins", "")
	if resp.StatusCode() != http.StatusOK {
		t.Fatalf("status=%d body=%s", resp.StatusCode(), resp.Body())
	}
	var payload map[string]any
	if err := json.Unmarshal(resp.Body(), &payload); err != nil {
		t.Fatal(err)
	}
	items, _ := payload["items"].([]any)
	if len(items) != 0 {
		t.Fatalf("want empty items got %v", payload)
	}
	if payload["next_cursor"] != nil {
		t.Fatalf("want next_cursor null/empty got %v", payload["next_cursor"])
	}
}

func TestPlatformAdmins_GetNotFound(t *testing.T) {
	fake := &fakePlatformAdminClient{err: status.Error(codes.NotFound, "PLATFORM_USER_NOT_FOUND")}
	h := setupPlatformAdminTestServer(t, fake)
	resp := performPlatformAdmin(h, http.MethodGet, "/api/v1/svc/platform-admins/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "")
	if resp.StatusCode() != http.StatusNotFound {
		t.Fatalf("status=%d body=%s", resp.StatusCode(), resp.Body())
	}
	if !strings.Contains(string(resp.Body()), "PLATFORM_USER_NOT_FOUND") {
		t.Fatalf("body=%s", resp.Body())
	}
}

func TestPlatformAdmins_UnimplementedMaps501(t *testing.T) {
	fake := &fakePlatformAdminClient{err: status.Error(codes.Unimplemented, "NOT_IMPLEMENTED")}
	h := setupPlatformAdminTestServer(t, fake)
	resp := performPlatformAdmin(h, http.MethodPost, "/api/v1/svc/platform-admins",
		`{"email":"a@x.com","username":"ops","display_name":"Ops","role":"platform-ops","password":"Abcd1234!","idempotency_key":"44444444-4444-4444-4444-444444444444"}`)
	if resp.StatusCode() != http.StatusNotImplemented {
		t.Fatalf("status=%d body=%s", resp.StatusCode(), resp.Body())
	}
}
