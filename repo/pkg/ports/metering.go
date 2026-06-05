package ports

import (
	"context"
	"time"
)

type MeteringResourceType string

const (
	MeteringResourceInstanceCPUSeconds    MeteringResourceType = "instance_cpu_seconds"
	MeteringResourceInstanceMemorySeconds MeteringResourceType = "instance_memory_gib_seconds"
	MeteringResourceInstanceGPUSeconds    MeteringResourceType = "instance_gpu_seconds"
	MeteringResourceTokenInput            MeteringResourceType = "token_input"
	MeteringResourceTokenOutput           MeteringResourceType = "token_output"
	MeteringResourceTokenTotal            MeteringResourceType = "token_total"
)

type TokenUsageReportState string

const (
	TokenUsageReportAccepted  TokenUsageReportState = "accepted"
	TokenUsageReportDuplicate TokenUsageReportState = "duplicate"
)

type MeteringUsageQueryRequest struct {
	TenantID     string
	StartTime    time.Time
	EndTime      time.Time
	ResourceType MeteringResourceType
	GroupBy      string
}

type MeteringUsageRecord struct {
	TenantID      string
	ResourceType  MeteringResourceType
	TotalQuantity float64
	Unit          string
	Period        string
}

type MeteringUsageResult struct {
	Items      []MeteringUsageRecord
	DevProfile DevProfileInfo
}

type TokenUsageReportRequest struct {
	TenantID       string
	IdempotencyKey string
	Source         string
	Model          string
	InputTokens    int64
	OutputTokens   int64
	RequestID      string
	InstanceID     string
	OccurredAt     time.Time
	Labels         map[string]string
}

type TokenUsageReportRecord struct {
	TenantID     string
	ReportID     string
	Source       string
	Model        string
	InputTokens  int64
	OutputTokens int64
	TotalTokens  int64
	RequestID    string
	InstanceID   string
	State        TokenUsageReportState
	DevProfile   DevProfileInfo
	CreatedAt    time.Time
}

type MeteringService interface {
	QueryUsage(ctx context.Context, request MeteringUsageQueryRequest) (MeteringUsageResult, error)
	ReportTokenUsage(ctx context.Context, request TokenUsageReportRequest) (TokenUsageReportRecord, error)
}
