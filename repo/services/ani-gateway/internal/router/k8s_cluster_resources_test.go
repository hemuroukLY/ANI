package router

import (
	"context"
	"testing"

	"github.com/kubercloud/ani/pkg/ports"
)

func TestK8sClusterAPIDevProfileAndIdempotency(t *testing.T) {
	api := newK8sClusterAPI()
	a, err := api.service.CreateCluster(context.Background(), ports.K8sClusterCreateRequest{TenantID: "t1", IdempotencyKey: "idem-1", Name: "vc-a"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := api.service.CreateCluster(context.Background(), ports.K8sClusterCreateRequest{TenantID: "t1", IdempotencyKey: "idem-1", Name: "vc-a"})
	if err != nil {
		t.Fatal(err)
	}
	if a.ClusterID != b.ClusterID {
		t.Fatalf("want idempotent cluster id, got %s != %s", a.ClusterID, b.ClusterID)
	}
	resp := k8sClusterFromRecord(a)
	requireLocalCoreDevProfile(t, resp.DevProfile, "local-k8s-cluster-service")

	kubeconfig, err := api.service.GetKubeconfig(context.Background(), ports.K8sClusterKubeconfigRequest{TenantID: "t1", ClusterID: a.ClusterID})
	if err != nil {
		t.Fatal(err)
	}
	if kubeconfig.Kubeconfig == "" || kubeconfig.Token == "" || kubeconfig.Server == "" {
		t.Fatalf("want kubeconfig content, got %+v", kubeconfig)
	}
	kubeconfigResp := k8sClusterKubeconfigFromRecord(kubeconfig)
	requireLocalCoreDevProfile(t, kubeconfigResp.DevProfile, "local-k8s-cluster-service")

	proxy, err := api.service.Proxy(context.Background(), ports.K8sClusterProxyRequest{
		TenantID:       "t1",
		ClusterID:      a.ClusterID,
		IdempotencyKey: "idem-proxy-1",
		Method:         "get",
		Path:           "/api/v1/namespaces/default/pods",
		Query:          map[string]string{"limit": "20"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if proxy.Method != "GET" || proxy.StatusCode != 200 || proxy.Body["kind"] != "ANIProxyPreview" {
		t.Fatalf("unexpected proxy response: %+v", proxy)
	}
	proxyResp := k8sClusterProxyFromRecord(proxy)
	requireLocalCoreDevProfile(t, proxyResp.DevProfile, "local-k8s-cluster-service")

	if _, err := api.service.Proxy(context.Background(), ports.K8sClusterProxyRequest{
		TenantID:       "t1",
		ClusterID:      a.ClusterID,
		IdempotencyKey: "idem-proxy-2",
		Method:         "GET",
		Path:           "/forbidden",
	}); err == nil {
		t.Fatalf("want invalid path error")
	}
}
