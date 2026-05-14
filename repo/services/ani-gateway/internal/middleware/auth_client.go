package middleware

import (
	"context"
	"os"
	"time"

	authv1 "github.com/kubercloud/ani/pkg/generated/pb/auth/v1"
	commonv1 "github.com/kubercloud/ani/pkg/generated/pb/common/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient interface {
	ValidateToken(ctx context.Context, token string) (*commonv1.TenantContext, error)
	CheckPermission(ctx context.Context, req *authv1.CheckPermissionRequest) (*authv1.CheckPermissionResponse, error)
}

type grpcAuthClient struct {
	client  authv1.AuthServiceClient
	timeout time.Duration
}

func NewAuthClientFromEnv() AuthClient {
	addr := os.Getenv("AUTH_SERVICE_ADDR")
	if addr == "" {
		addr = "127.0.0.1:9101"
	}
	timeout := 2 * time.Second
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}
	return &grpcAuthClient{
		client:  authv1.NewAuthServiceClient(conn),
		timeout: timeout,
	}
}

func (c *grpcAuthClient) ValidateToken(ctx context.Context, token string) (*commonv1.TenantContext, error) {
	callCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.ValidateToken(callCtx, &authv1.ValidateTokenRequest{Token: token})
}

func (c *grpcAuthClient) CheckPermission(ctx context.Context, req *authv1.CheckPermissionRequest) (*authv1.CheckPermissionResponse, error) {
	callCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.CheckPermission(callCtx, req)
}
