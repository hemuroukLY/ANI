package runtime

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kubercloud/ani/pkg/ports"
)

func TestK8sClusterProxyForwardingServiceForwardsToResolvedAPIServer(t *testing.T) {
	base := NewLocalK8sClusterService()
	cluster, err := base.CreateCluster(context.Background(), ports.K8sClusterCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "create-vc-a",
		Name:           "vc-a",
		Version:        "v1.30.0",
	})
	if err != nil {
		t.Fatal(err)
	}

	transport := &capturingK8sProxyRoundTripper{
		statusCode: http.StatusCreated,
		headers: http.Header{
			"Content-Type":          []string{"application/json"},
			"X-Kubernetes-Audit-ID": []string{"audit-1"},
		},
		body: `{"kind":"Pod","metadata":{"name":"demo-pod"}}`,
	}

	resolver := staticK8sProxyTargetResolver{target: ports.K8sClusterProxyTarget{
		TenantID:    "tenant-a",
		ClusterID:   cluster.ClusterID,
		Server:      "https://tenant-a-vcluster.example",
		BearerToken: "tenant-token",
	}}
	service := NewK8sClusterProxyForwardingService(
		base,
		resolver,
		WithK8sClusterProxyForwardingHTTPClient(&http.Client{Transport: transport}),
		WithK8sClusterProxyForwardingClock(func() time.Time { return time.Unix(700, 0) }),
	)

	result, err := service.Proxy(context.Background(), ports.K8sClusterProxyRequest{
		TenantID:       "tenant-a",
		ClusterID:      cluster.ClusterID,
		IdempotencyKey: "proxy-1",
		Method:         "post",
		Path:           "api/v1/namespaces/default/pods",
		Query:          map[string]string{"limit": "20"},
		Body:           map[string]any{"metadata": map[string]any{"name": "demo-pod"}},
	})
	if err != nil {
		t.Fatalf("Proxy() error = %v", err)
	}

	if transport.method != http.MethodPost {
		t.Fatalf("upstream method = %s, want POST", transport.method)
	}
	if transport.path != "/api/v1/namespaces/default/pods" {
		t.Fatalf("upstream path = %s", transport.path)
	}
	if transport.query != "limit=20" {
		t.Fatalf("upstream query = %s, want limit=20", transport.query)
	}
	if transport.authorization != "Bearer tenant-token" {
		t.Fatalf("upstream authorization = %q", transport.authorization)
	}
	if metadata, _ := transport.decodedBody["metadata"].(map[string]any); metadata["name"] != "demo-pod" {
		t.Fatalf("upstream body = %+v", transport.decodedBody)
	}
	if result.StatusCode != http.StatusCreated || result.Body["kind"] != "Pod" {
		t.Fatalf("proxy result = %+v", result)
	}
	if result.Headers["x-kubernetes-audit-id"] != "audit-1" {
		t.Fatalf("proxy headers = %+v", result.Headers)
	}
	if result.ProxiedAt != 700 {
		t.Fatalf("ProxiedAt = %d, want 700", result.ProxiedAt)
	}
}

func TestK8sClusterProxyForwardingServiceRejectsMismatchedResolvedTarget(t *testing.T) {
	base := NewLocalK8sClusterService()
	cluster, err := base.CreateCluster(context.Background(), ports.K8sClusterCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "create-vc-a",
		Name:           "vc-a",
	})
	if err != nil {
		t.Fatal(err)
	}
	service := NewK8sClusterProxyForwardingService(
		base,
		staticK8sProxyTargetResolver{target: ports.K8sClusterProxyTarget{
			TenantID:  "tenant-b",
			ClusterID: cluster.ClusterID,
			Server:    "https://tenant-b.invalid",
		}},
	)

	if _, err := service.Proxy(context.Background(), ports.K8sClusterProxyRequest{
		TenantID:       "tenant-a",
		ClusterID:      cluster.ClusterID,
		IdempotencyKey: "proxy-1",
		Method:         "GET",
		Path:           "/version",
	}); err == nil {
		t.Fatalf("want mismatched resolved target error")
	}
}

func TestK8sClusterProxyForwardingServiceListsWorkloadsFromResolvedAPIServer(t *testing.T) {
	clusterProvider := &fakeK8sClusterForwardingProvider{
		result: ports.K8sClusterProviderApplyResult{
			Applied:      true,
			Provider:     "vcluster",
			ResourceRefs: []string{"vcluster/HelmRelease/vc-workloads"},
			Reason:       "vCluster Helm release applied",
		},
	}
	base := NewLocalK8sClusterService(WithK8sClusterProviderApply(clusterProvider))
	cluster, err := base.CreateCluster(context.Background(), ports.K8sClusterCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "create-vc-workloads",
		Name:           "vc-workloads",
		Version:        "v1.36.1",
	})
	if err != nil {
		t.Fatal(err)
	}

	transport := &capturingK8sProxyRoundTripper{
		statusCode: http.StatusOK,
		headers:    http.Header{"Content-Type": []string{"application/json"}},
		body: `{
			"kind":"DeploymentList",
			"items":[{
				"metadata":{"name":"web","namespace":"default","creationTimestamp":"2026-06-19T10:00:00Z"},
				"spec":{"replicas":3,"template":{"spec":{"containers":[{"image":"registry.example/web:v1"}]}}},
				"status":{"readyReplicas":2}
			}]
		}`,
	}
	resolver := staticK8sProxyTargetResolver{target: ports.K8sClusterProxyTarget{
		TenantID:    "tenant-a",
		ClusterID:   cluster.ClusterID,
		Server:      "https://tenant-a-vcluster.example",
		BearerToken: "tenant-token",
	}}
	service := NewK8sClusterProxyForwardingService(
		base,
		resolver,
		WithK8sClusterProxyForwardingHTTPClient(&http.Client{Transport: transport}),
	)

	workloads, err := service.ListWorkloads(context.Background(), ports.K8sClusterWorkloadListRequest{
		TenantID:  "tenant-a",
		ClusterID: cluster.ClusterID,
		Namespace: "default",
		Kind:      "Deployment",
	})
	if err != nil {
		t.Fatalf("ListWorkloads() error = %v", err)
	}
	if transport.method != http.MethodGet || transport.path != "/apis/apps/v1/namespaces/default/deployments" {
		t.Fatalf("upstream request = %s %s, want GET deployments", transport.method, transport.path)
	}
	if transport.authorization != "Bearer tenant-token" {
		t.Fatalf("upstream authorization = %q", transport.authorization)
	}
	if len(workloads) != 1 {
		t.Fatalf("workloads = %+v, want one Deployment", workloads)
	}
	got := workloads[0]
	if got.Name != "web" || got.Namespace != "default" || got.Kind != "Deployment" || got.Replicas != 3 || got.ReadyReplicas != 2 || got.Image != "registry.example/web:v1" || got.Status != ports.K8sWorkloadPending {
		t.Fatalf("workload = %+v, want parsed pending Deployment", got)
	}
}

type staticK8sProxyTargetResolver struct {
	target ports.K8sClusterProxyTarget
}

func (r staticK8sProxyTargetResolver) ResolveK8sClusterProxyTarget(context.Context, ports.K8sClusterGetRequest) (ports.K8sClusterProxyTarget, error) {
	return r.target, nil
}

type capturingK8sProxyRoundTripper struct {
	statusCode    int
	headers       http.Header
	body          string
	method        string
	path          string
	query         string
	authorization string
	decodedBody   map[string]any
}

type fakeK8sClusterForwardingProvider struct {
	result ports.K8sClusterProviderApplyResult
}

func (p *fakeK8sClusterForwardingProvider) ApplyK8sCluster(context.Context, ports.K8sClusterProviderApplyRequest) (ports.K8sClusterProviderApplyResult, error) {
	return p.result, nil
}

func (t *capturingK8sProxyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	t.method = req.Method
	t.path = req.URL.Path
	t.query = req.URL.RawQuery
	t.authorization = req.Header.Get("Authorization")
	if req.Body != nil {
		defer req.Body.Close()
		if err := json.NewDecoder(req.Body).Decode(&t.decodedBody); err != nil {
			return nil, err
		}
	}
	return &http.Response{
		StatusCode: t.statusCode,
		Header:     t.headers,
		Body:       io.NopCloser(strings.NewReader(t.body)),
		Request:    req,
	}, nil
}
