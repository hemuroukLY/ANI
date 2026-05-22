package router

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
	runtimeadapter "github.com/kubercloud/ani/pkg/adapters/runtime"
	"github.com/kubercloud/ani/pkg/ports"
)

type k8sClusterAPI struct{ service ports.K8sClusterService }
type k8sClusterCreateRequest struct {
	IdempotencyKey string `json:"idempotency_key"`
	Name           string `json:"name"`
	Version        string `json:"version"`
}
type k8sClusterProxyRequest struct {
	IdempotencyKey string            `json:"idempotency_key"`
	Method         string            `json:"method"`
	Path           string            `json:"path"`
	Query          map[string]string `json:"query"`
	Body           map[string]any    `json:"body"`
}
type k8sClusterResponse struct {
	ID         string                 `json:"id"`
	TenantID   string                 `json:"tenant_id"`
	Name       string                 `json:"name"`
	Version    string                 `json:"version,omitempty"`
	State      string                 `json:"state"`
	Reason     string                 `json:"reason,omitempty"`
	DevProfile coreDevProfileResponse `json:"dev_profile"`
	CreatedAt  string                 `json:"created_at"`
	UpdatedAt  string                 `json:"updated_at"`
}
type k8sClusterKubeconfigResponse struct {
	ClusterID  string                 `json:"cluster_id"`
	TenantID   string                 `json:"tenant_id"`
	Server     string                 `json:"server"`
	Namespace  string                 `json:"namespace"`
	CAData     string                 `json:"ca_data"`
	Token      string                 `json:"token"`
	Kubeconfig string                 `json:"kubeconfig"`
	ExpiresAt  string                 `json:"expires_at"`
	DevProfile coreDevProfileResponse `json:"dev_profile"`
}
type k8sClusterProxyResponse struct {
	ClusterID  string                 `json:"cluster_id"`
	TenantID   string                 `json:"tenant_id"`
	Method     string                 `json:"method"`
	Path       string                 `json:"path"`
	Query      map[string]string      `json:"query"`
	StatusCode int                    `json:"status_code"`
	Headers    map[string]string      `json:"headers"`
	Body       map[string]any         `json:"body"`
	ProxiedAt  string                 `json:"proxied_at"`
	DevProfile coreDevProfileResponse `json:"dev_profile"`
}

func newK8sClusterAPI() *k8sClusterAPI {
	return &k8sClusterAPI{service: runtimeadapter.NewLocalK8sClusterService()}
}
func registerK8sClusterResources(v1 *route.RouterGroup) {
	api := newK8sClusterAPI()
	v1.GET("/k8s-clusters", api.listClusters)
	v1.POST("/k8s-clusters", api.createCluster)
	v1.GET("/k8s-clusters/:cluster_id", api.getCluster)
	v1.DELETE("/k8s-clusters/:cluster_id", api.deleteCluster)
	v1.GET("/k8s-clusters/:cluster_id/kubeconfig", api.getKubeconfig)
	v1.POST("/k8s-clusters/:cluster_id/proxy", api.proxy)
}
func (api *k8sClusterAPI) createCluster(ctx context.Context, c *app.RequestContext) {
	var req k8sClusterCreateRequest
	if err := c.BindJSON(&req); err != nil {
		writeDemoError(c, http.StatusBadRequest, "BAD_REQUEST", "invalid k8s cluster request")
		return
	}
	rec, err := api.service.CreateCluster(ctx, ports.K8sClusterCreateRequest{TenantID: demoTenantID(c), IdempotencyKey: req.IdempotencyKey, Name: req.Name, Version: req.Version})
	if err != nil {
		writeK8sClusterError(c, err)
		return
	}
	c.JSON(http.StatusCreated, k8sClusterFromRecord(rec))
}
func (api *k8sClusterAPI) listClusters(ctx context.Context, c *app.RequestContext) {
	recs, err := api.service.ListClusters(ctx, ports.K8sClusterListRequest{TenantID: demoTenantID(c)})
	if err != nil {
		writeK8sClusterError(c, err)
		return
	}
	items := make([]k8sClusterResponse, 0, len(recs))
	for _, r := range recs {
		items = append(items, k8sClusterFromRecord(r))
	}
	c.JSON(http.StatusOK, map[string]any{"items": items, "total": len(items), "next_cursor": nil})
}
func (api *k8sClusterAPI) getCluster(ctx context.Context, c *app.RequestContext) {
	rec, err := api.service.GetCluster(ctx, ports.K8sClusterGetRequest{TenantID: demoTenantID(c), ClusterID: c.Param("cluster_id")})
	if err != nil {
		writeK8sClusterError(c, err)
		return
	}
	c.JSON(http.StatusOK, k8sClusterFromRecord(rec))
}
func (api *k8sClusterAPI) deleteCluster(ctx context.Context, c *app.RequestContext) {
	rec, err := api.service.DeleteCluster(ctx, ports.K8sClusterGetRequest{TenantID: demoTenantID(c), ClusterID: c.Param("cluster_id")})
	if err != nil {
		writeK8sClusterError(c, err)
		return
	}
	c.JSON(http.StatusOK, k8sClusterFromRecord(rec))
}
func (api *k8sClusterAPI) getKubeconfig(ctx context.Context, c *app.RequestContext) {
	rec, err := api.service.GetKubeconfig(ctx, ports.K8sClusterKubeconfigRequest{TenantID: demoTenantID(c), ClusterID: c.Param("cluster_id")})
	if err != nil {
		writeK8sClusterError(c, err)
		return
	}
	c.JSON(http.StatusOK, k8sClusterKubeconfigFromRecord(rec))
}
func (api *k8sClusterAPI) proxy(ctx context.Context, c *app.RequestContext) {
	var req k8sClusterProxyRequest
	if err := c.BindJSON(&req); err != nil {
		writeDemoError(c, http.StatusBadRequest, "BAD_REQUEST", "invalid k8s cluster proxy request")
		return
	}
	rec, err := api.service.Proxy(ctx, ports.K8sClusterProxyRequest{
		TenantID:       demoTenantID(c),
		ClusterID:      c.Param("cluster_id"),
		IdempotencyKey: req.IdempotencyKey,
		Method:         req.Method,
		Path:           req.Path,
		Query:          req.Query,
		Body:           req.Body,
	})
	if err != nil {
		writeK8sClusterError(c, err)
		return
	}
	c.JSON(http.StatusOK, k8sClusterProxyFromRecord(rec))
}
func k8sClusterFromRecord(r ports.K8sClusterRecord) k8sClusterResponse {
	return k8sClusterResponse{ID: r.ClusterID, TenantID: r.TenantID, Name: r.Name, Version: r.Version, State: string(r.State), Reason: r.Reason, DevProfile: localCoreDevProfile("local-k8s-cluster-service", "Core dev/local profile; vCluster lifecycle is simulated"), CreatedAt: time.Unix(r.CreatedAt, 0).UTC().Format(time.RFC3339), UpdatedAt: time.Unix(r.UpdatedAt, 0).UTC().Format(time.RFC3339)}
}
func k8sClusterKubeconfigFromRecord(r ports.K8sClusterKubeconfigRecord) k8sClusterKubeconfigResponse {
	return k8sClusterKubeconfigResponse{ClusterID: r.ClusterID, TenantID: r.TenantID, Server: r.Server, Namespace: r.Namespace, CAData: r.CAData, Token: r.Token, Kubeconfig: r.Kubeconfig, ExpiresAt: time.Unix(r.ExpiresAt, 0).UTC().Format(time.RFC3339), DevProfile: localCoreDevProfile("local-k8s-cluster-service", "Core dev/local profile; kubeconfig targets a simulated vCluster endpoint")}
}
func k8sClusterProxyFromRecord(r ports.K8sClusterProxyRecord) k8sClusterProxyResponse {
	return k8sClusterProxyResponse{ClusterID: r.ClusterID, TenantID: r.TenantID, Method: r.Method, Path: r.Path, Query: r.Query, StatusCode: r.StatusCode, Headers: r.Headers, Body: r.Body, ProxiedAt: time.Unix(r.ProxiedAt, 0).UTC().Format(time.RFC3339), DevProfile: localCoreDevProfile("local-k8s-cluster-service", "Core dev/local profile; proxy response is simulated and not forwarded to a real vCluster API server")}
}
func writeK8sClusterError(c *app.RequestContext, err error) {
	switch {
	case errors.Is(err, ports.ErrNotFound):
		writeDemoError(c, http.StatusNotFound, "NOT_FOUND", err.Error())
	case errors.Is(err, ports.ErrConflict):
		writeDemoError(c, http.StatusConflict, "CONFLICT", err.Error())
	case errors.Is(err, ports.ErrInvalid):
		writeDemoError(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
	default:
		writeDemoError(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
	}
}
