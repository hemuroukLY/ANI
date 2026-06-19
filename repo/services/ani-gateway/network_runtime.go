package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	runtimeadapter "github.com/kubercloud/ani/pkg/adapters/runtime"
	"github.com/kubercloud/ani/pkg/ports"
)

type gatewayNetworkRuntimeConfig struct {
	ProviderMode              string
	ProviderApply             bool
	ProviderUserID            string
	ProviderProof             string
	KubernetesAPIHost         string
	KubernetesBearerToken     string
	KubernetesProviderManager string
	KubernetesHTTPClient      *http.Client
}

func gatewayNetworkRuntimeConfigFromEnv() gatewayNetworkRuntimeConfig {
	return gatewayNetworkRuntimeConfig{
		ProviderMode:              os.Getenv("NETWORK_PROVIDER"),
		ProviderApply:             strings.EqualFold(strings.TrimSpace(os.Getenv("NETWORK_PROVIDER_APPLY_ENABLED")), "true"),
		ProviderUserID:            os.Getenv("NETWORK_PROVIDER_USER_ID"),
		ProviderProof:             os.Getenv("NETWORK_PROVIDER_PERMISSION_PROOF"),
		KubernetesAPIHost:         os.Getenv("KUBERNETES_API_HOST"),
		KubernetesBearerToken:     os.Getenv("KUBERNETES_BEARER_TOKEN"),
		KubernetesProviderManager: os.Getenv("KUBERNETES_PROVIDER_FIELD_MANAGER"),
	}
}

func newGatewayNetworkService(cfg gatewayNetworkRuntimeConfig) (ports.NetworkService, error) {
	switch mode := strings.TrimSpace(cfg.ProviderMode); mode {
	case "", "local", "not_configured":
		return nil, nil
	case "kubeovn_rest":
		if strings.TrimSpace(cfg.ProviderUserID) == "" || strings.TrimSpace(cfg.ProviderProof) == "" {
			return nil, fmt.Errorf("%w: network provider requires NETWORK_PROVIDER_USER_ID and NETWORK_PROVIDER_PERMISSION_PROOF", ports.ErrInvalid)
		}
		client, err := runtimeadapter.NewKubernetesRESTClient(runtimeadapter.KubernetesRESTClientConfig{
			Host:         cfg.KubernetesAPIHost,
			BearerToken:  cfg.KubernetesBearerToken,
			FieldManager: cfg.KubernetesProviderManager,
			HTTPClient:   cfg.KubernetesHTTPClient,
		})
		if err != nil {
			return nil, err
		}
		provider := runtimeadapter.NewKubeOVNNetworkProviderAdapter(
			client,
			runtimeadapter.WithKubeOVNNetworkProviderApplyEnabled(cfg.ProviderApply),
		)
		return runtimeadapter.NewLocalNetworkService(
			runtimeadapter.WithNetworkRouteProvider(
				runtimeadapter.NewKubeOVNNetworkRenderer(),
				provider,
				provider,
				provider,
				runtimeadapter.NetworkProviderExecutionConfig{
					UserID:          cfg.ProviderUserID,
					PermissionProof: cfg.ProviderProof,
				},
			),
		), nil
	default:
		return nil, fmt.Errorf("%w: unsupported NETWORK_PROVIDER %q", ports.ErrUnsupported, mode)
	}
}
