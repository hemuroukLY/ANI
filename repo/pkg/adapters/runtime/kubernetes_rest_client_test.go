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

func TestKubernetesRESTClientServerSideDryRunUsesDryRunAll(t *testing.T) {
	var gotPath string
	var gotAuth string
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		gotPath = r.URL.String()
		gotAuth = r.Header.Get("Authorization")
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Query().Get("dryRun") != "All" {
			t.Fatalf("dryRun = %q, want All", r.URL.Query().Get("dryRun"))
		}
		return jsonResponse(http.StatusCreated, `{"kind":"Deployment"}`), nil
	})

	client := newTestKubernetesRESTClient(t, transport)
	client.bearerToken = "token-a"
	result, err := client.ServerSideDryRun(context.Background(), renderedDeployment(t))
	if err != nil {
		t.Fatalf("ServerSideDryRun() error = %v", err)
	}
	if !result.Accepted || result.Provider != "kubernetes" {
		t.Fatalf("result = %#v, want accepted kubernetes dry-run", result)
	}
	if !strings.Contains(gotPath, "/apis/apps/v1/namespaces/ani-tenant-tenant-a/deployments") {
		t.Fatalf("path = %q, want Deployment collection path", gotPath)
	}
	if gotAuth != "Bearer token-a" {
		t.Fatalf("Authorization = %q, want bearer token", gotAuth)
	}
}

func TestKubernetesRESTClientApplyUsesServerSideApply(t *testing.T) {
	var gotPath string
	var gotContentType string
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		gotPath = r.URL.String()
		gotContentType = r.Header.Get("Content-Type")
		if r.Method != http.MethodPatch {
			t.Fatalf("method = %s, want PATCH", r.Method)
		}
		if r.URL.Query().Get("fieldManager") != "ani-test" {
			t.Fatalf("fieldManager = %q, want ani-test", r.URL.Query().Get("fieldManager"))
		}
		if r.URL.Query().Get("force") != "true" {
			t.Fatalf("force = %q, want true", r.URL.Query().Get("force"))
		}
		return jsonResponse(http.StatusOK, `{"kind":"Deployment"}`), nil
	})

	client := newTestKubernetesRESTClient(t, transport)
	result, err := client.Apply(context.Background(), validProviderApplyRequest(t))
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if !result.Applied {
		t.Fatalf("Applied = false, reason = %s", result.Reason)
	}
	if len(result.ResourceRefs) != 1 || result.ResourceRefs[0] != "kubernetes/Deployment/app-01" {
		t.Fatalf("ResourceRefs = %#v, want deployment ref", result.ResourceRefs)
	}
	if !strings.Contains(gotPath, "/apis/apps/v1/namespaces/ani-tenant-tenant-a/deployments/app-01") {
		t.Fatalf("path = %q, want Deployment resource path", gotPath)
	}
	if gotContentType != kubernetesApplyPatchContentType {
		t.Fatalf("Content-Type = %q, want apply patch", gotContentType)
	}
}

func TestKubernetesRESTClientObserveDeploymentStatus(t *testing.T) {
	var gotPath string
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		gotPath = r.URL.String()
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		body, err := json.Marshal(map[string]any{
			"status": map[string]any{
				"availableReplicas": 1,
			},
		})
		if err != nil {
			t.Fatalf("marshal response: %v", err)
		}
		return jsonResponse(http.StatusOK, string(body)), nil
	})

	client := newTestKubernetesRESTClient(t, transport)
	observation, err := client.Observe(context.Background(), ports.WorkloadProviderStatusRequest{
		TenantID:   "tenant-a",
		InstanceID: "instance-a",
		Kind:       ports.WorkloadKindContainer,
		ApplyResult: ports.WorkloadProviderApplyResult{
			Applied:      true,
			Provider:     "kubernetes",
			ResourceRefs: []string{"kubernetes/Deployment/app-01"},
		},
	})
	if err != nil {
		t.Fatalf("Observe() error = %v", err)
	}
	if observation.Phase != "Running" {
		t.Fatalf("Phase = %q, want Running", observation.Phase)
	}
	if !strings.Contains(gotPath, "/apis/apps/v1/namespaces/ani-tenant-tenant-a/deployments/app-01") {
		t.Fatalf("path = %q, want Deployment resource path", gotPath)
	}
}

func TestKubernetesRESTClientSupportsKubeVirtVirtualMachine(t *testing.T) {
	var paths []string
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		paths = append(paths, r.URL.String())
		return jsonResponse(http.StatusOK, `{"kind":"VirtualMachine"}`), nil
	})

	client := newTestKubernetesRESTClient(t, transport)
	manifests, err := NewKubernetesDryRunRenderer(NewPlanningRuntime()).Render(context.Background(), ports.WorkloadSpec{
		TenantID: "tenant-a",
		Name:     "vm-01",
		Kind:     ports.WorkloadKindVM,
		VM: &ports.VMInstanceSpec{
			BootImage: "ubuntu.qcow2",
			RootDisk: ports.WorkloadStorageAttachment{
				Name:    "root",
				Kind:    ports.StorageAttachmentRootDisk,
				SizeGiB: 80,
			},
		},
	})
	if err != nil {
		t.Fatalf("Render(VM) error = %v", err)
	}
	if _, err := client.ServerSideDryRun(context.Background(), manifests); err != nil {
		t.Fatalf("ServerSideDryRun(VM) error = %v", err)
	}
	if len(paths) != 1 || !strings.Contains(paths[0], "/apis/kubevirt.io/v1/namespaces/ani-tenant-tenant-a/virtualmachines") {
		t.Fatalf("paths = %#v, want KubeVirt VirtualMachine collection", paths)
	}
}

func newTestKubernetesRESTClient(t *testing.T, transport http.RoundTripper) *KubernetesRESTClient {
	t.Helper()
	client, err := NewKubernetesRESTClient(KubernetesRESTClientConfig{
		Host:         "https://kubernetes.example.test",
		FieldManager: "ani-test",
		HTTPClient:   &http.Client{Transport: transport},
		Now:          func() time.Time { return time.Unix(900, 0) },
	})
	if err != nil {
		t.Fatalf("NewKubernetesRESTClient() error = %v", err)
	}
	return client
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
