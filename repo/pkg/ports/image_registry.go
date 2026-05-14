package ports

import "context"

type ImageRef struct {
	Registry   string
	Repository string
	Tag        string
	Digest     string
}

type ImageTag struct {
	Name   string
	Digest string
}

type ImageScanStatus struct {
	Status     string
	Critical   int
	High       int
	Medium     int
	Low        int
	ReportURL  string
	ProviderID string
}

type ImageRegistry interface {
	EnsureProject(ctx context.Context, tenantID string) error
	ListTags(ctx context.Context, repository string) ([]ImageTag, error)
	GetScanStatus(ctx context.Context, ref ImageRef) (ImageScanStatus, error)
}
