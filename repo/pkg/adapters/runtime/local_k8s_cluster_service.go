package runtime

import (
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kubercloud/ani/pkg/ports"
)

type localK8sClusterService struct {
	mu   sync.Mutex
	byID map[string]ports.K8sClusterRecord
	idem map[string]string
}

func NewLocalK8sClusterService() ports.K8sClusterService {
	return &localK8sClusterService{byID: map[string]ports.K8sClusterRecord{}, idem: map[string]string{}}
}

func (s *localK8sClusterService) CreateCluster(_ context.Context, req ports.K8sClusterCreateRequest) (ports.K8sClusterRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if req.TenantID == "" || req.Name == "" || req.IdempotencyKey == "" {
		return ports.K8sClusterRecord{}, fmt.Errorf("%w: tenant_id/name/idempotency_key required", ports.ErrInvalid)
	}
	key := req.TenantID + ":" + req.IdempotencyKey
	if id, ok := s.idem[key]; ok {
		return s.byID[id], nil
	}
	now := time.Now().Unix()
	rec := ports.K8sClusterRecord{ClusterID: "k8sclu-" + uuid.NewString(), TenantID: req.TenantID, Name: req.Name, Version: req.Version, State: ports.K8sClusterStateRunning, Reason: "local vcluster profile", CreatedAt: now, UpdatedAt: now}
	s.byID[rec.ClusterID] = rec
	s.idem[key] = rec.ClusterID
	return rec, nil
}

func (s *localK8sClusterService) GetCluster(_ context.Context, req ports.K8sClusterGetRequest) (ports.K8sClusterRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.byID[req.ClusterID]
	if !ok || rec.TenantID != req.TenantID {
		return ports.K8sClusterRecord{}, ports.ErrNotFound
	}
	return rec, nil
}
func (s *localK8sClusterService) ListClusters(_ context.Context, req ports.K8sClusterListRequest) ([]ports.K8sClusterRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []ports.K8sClusterRecord{}
	for _, r := range s.byID {
		if r.TenantID == req.TenantID {
			out = append(out, r)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt < out[j].CreatedAt })
	return out, nil
}
func (s *localK8sClusterService) DeleteCluster(_ context.Context, req ports.K8sClusterGetRequest) (ports.K8sClusterRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.byID[req.ClusterID]
	if !ok || rec.TenantID != req.TenantID {
		return ports.K8sClusterRecord{}, ports.ErrNotFound
	}
	rec.State = ports.K8sClusterStateDeleting
	rec.UpdatedAt = time.Now().Unix()
	s.byID[req.ClusterID] = rec
	return rec, nil
}

func (s *localK8sClusterService) GetKubeconfig(_ context.Context, req ports.K8sClusterKubeconfigRequest) (ports.K8sClusterKubeconfigRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.byID[req.ClusterID]
	if !ok || rec.TenantID != req.TenantID {
		return ports.K8sClusterKubeconfigRecord{}, ports.ErrNotFound
	}
	if rec.State != ports.K8sClusterStateRunning {
		return ports.K8sClusterKubeconfigRecord{}, fmt.Errorf("%w: kubeconfig requires a running k8s cluster", ports.ErrConflict)
	}
	now := time.Now().Unix()
	server := fmt.Sprintf("https://%s.local.ani.invalid", rec.ClusterID)
	namespace := "tenant-" + req.TenantID
	token := "local-kubeconfig-" + uuid.NewString()
	caData := base64.StdEncoding.EncodeToString([]byte("local-dev-profile-ca:" + rec.ClusterID))
	kubeconfig := fmt.Sprintf(`apiVersion: v1
kind: Config
clusters:
- name: %s
  cluster:
    server: %s
    certificate-authority-data: %s
contexts:
- name: %s
  context:
    cluster: %s
    namespace: %s
    user: %s-user
current-context: %s
users:
- name: %s-user
  user:
    token: %s
`, rec.Name, server, caData, rec.Name, rec.Name, namespace, rec.Name, rec.Name, rec.Name, token)
	return ports.K8sClusterKubeconfigRecord{
		ClusterID:  rec.ClusterID,
		TenantID:   rec.TenantID,
		Server:     server,
		Namespace:  namespace,
		CAData:     caData,
		Token:      token,
		Kubeconfig: kubeconfig,
		ExpiresAt:  now + 3600,
		CreatedAt:  now,
	}, nil
}

func (s *localK8sClusterService) Proxy(_ context.Context, req ports.K8sClusterProxyRequest) (ports.K8sClusterProxyRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.byID[req.ClusterID]
	if !ok || rec.TenantID != req.TenantID {
		return ports.K8sClusterProxyRecord{}, ports.ErrNotFound
	}
	if rec.State != ports.K8sClusterStateRunning {
		return ports.K8sClusterProxyRecord{}, fmt.Errorf("%w: proxy requires a running k8s cluster", ports.ErrConflict)
	}
	method := strings.ToUpper(strings.TrimSpace(req.Method))
	path := normalizeK8sProxyPath(req.Path)
	if method == "" || path == "" {
		return ports.K8sClusterProxyRecord{}, fmt.Errorf("%w: method/path required for k8s proxy", ports.ErrInvalid)
	}
	if !isAllowedK8sProxyPath(path) {
		return ports.K8sClusterProxyRecord{}, fmt.Errorf("%w: k8s proxy path must start with /api/, /apis/, /healthz, /livez, /readyz or /version", ports.ErrInvalid)
	}
	if req.IdempotencyKey == "" {
		return ports.K8sClusterProxyRecord{}, fmt.Errorf("%w: idempotency_key required for k8s proxy", ports.ErrInvalid)
	}
	now := time.Now().Unix()
	query := copyStringMap(req.Query)
	body := copyAnyMap(req.Body)
	return ports.K8sClusterProxyRecord{
		ClusterID:  rec.ClusterID,
		TenantID:   rec.TenantID,
		Method:     method,
		Path:       path,
		Query:      query,
		StatusCode: 200,
		Headers: map[string]string{
			"content-type":              "application/json",
			"x-ani-provider":            "local-k8s-cluster-service",
			"x-ani-k8s-cluster-version": rec.Version,
		},
		Body: map[string]any{
			"apiVersion": "v1",
			"kind":       "ANIProxyPreview",
			"metadata": map[string]any{
				"cluster_id": rec.ClusterID,
				"tenant_id":  rec.TenantID,
				"path":       path,
				"method":     method,
			},
			"request": map[string]any{
				"query": query,
				"body":  body,
			},
			"message": "local dev profile; request was not forwarded to a real vCluster API server",
		},
		ProxiedAt: now,
	}, nil
}

func normalizeK8sProxyPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

func isAllowedK8sProxyPath(path string) bool {
	switch {
	case path == "/healthz", path == "/livez", path == "/readyz", path == "/version":
		return true
	case strings.HasPrefix(path, "/api/"), strings.HasPrefix(path, "/apis/"):
		return true
	default:
		return false
	}
}

func copyStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return map[string]string{}
	}
	out := make(map[string]string, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func copyAnyMap(in map[string]any) map[string]any {
	if len(in) == 0 {
		return map[string]any{}
	}
	out := make(map[string]any, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}
