package registry

import (
	"context"

	"github.com/kubercloud/ani/pkg/ports"
)

type NotConfigured struct{}

var _ ports.ImageRegistry = NotConfigured{}

func (NotConfigured) EnsureProject(context.Context, string) error {
	return ports.ErrNotConfigured
}

func (NotConfigured) ListTags(context.Context, string) ([]ports.ImageTag, error) {
	return nil, ports.ErrNotConfigured
}

func (NotConfigured) GetScanStatus(context.Context, ports.ImageRef) (ports.ImageScanStatus, error) {
	return ports.ImageScanStatus{}, ports.ErrNotConfigured
}
