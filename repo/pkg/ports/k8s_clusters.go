package ports

import "context"

type K8sClusterState string

const (
	K8sClusterStateProvisioning K8sClusterState = "provisioning"
	K8sClusterStateRunning      K8sClusterState = "running"
	K8sClusterStateDeleting     K8sClusterState = "deleting"
)

type K8sClusterCreateRequest struct {
	TenantID       string
	IdempotencyKey string
	Name           string
	Version        string
}

type K8sClusterGetRequest struct {
	TenantID  string
	ClusterID string
}

type K8sClusterListRequest struct {
	TenantID string
}

type K8sClusterKubeconfigRequest struct {
	TenantID  string
	ClusterID string
}

type K8sClusterProxyRequest struct {
	TenantID       string
	ClusterID      string
	IdempotencyKey string
	Method         string
	Path           string
	Query          map[string]string
	Body           map[string]any
}

type K8sClusterRecord struct {
	ClusterID string
	TenantID  string
	Name      string
	Version   string
	State     K8sClusterState
	Reason    string
	CreatedAt int64
	UpdatedAt int64
}

type K8sClusterKubeconfigRecord struct {
	ClusterID  string
	TenantID   string
	Server     string
	Namespace  string
	CAData     string
	Token      string
	Kubeconfig string
	ExpiresAt  int64
	CreatedAt  int64
}

type K8sClusterProxyRecord struct {
	ClusterID  string
	TenantID   string
	Method     string
	Path       string
	Query      map[string]string
	StatusCode int
	Headers    map[string]string
	Body       map[string]any
	ProxiedAt  int64
}

type K8sClusterService interface {
	CreateCluster(ctx context.Context, req K8sClusterCreateRequest) (K8sClusterRecord, error)
	GetCluster(ctx context.Context, req K8sClusterGetRequest) (K8sClusterRecord, error)
	ListClusters(ctx context.Context, req K8sClusterListRequest) ([]K8sClusterRecord, error)
	DeleteCluster(ctx context.Context, req K8sClusterGetRequest) (K8sClusterRecord, error)
	GetKubeconfig(ctx context.Context, req K8sClusterKubeconfigRequest) (K8sClusterKubeconfigRecord, error)
	Proxy(ctx context.Context, req K8sClusterProxyRequest) (K8sClusterProxyRecord, error)
}
