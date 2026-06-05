package runtime

import (
	"context"
	"testing"
	"time"

	"github.com/kubercloud/ani/pkg/ports"
)

func TestLocalMeteringServiceQueriesUsageWithDevProfile(t *testing.T) {
	service := NewLocalMeteringService()

	result, err := service.QueryUsage(context.Background(), ports.MeteringUsageQueryRequest{
		TenantID:     "tenant-a",
		ResourceType: ports.MeteringResourceTokenTotal,
	})
	if err != nil {
		t.Fatalf("QueryUsage() error = %v", err)
	}
	if result.DevProfile.Mode != "local" || result.DevProfile.Provider != "local-metering-service" || result.DevProfile.RealProvider {
		t.Fatalf("dev profile = %+v, want local non-real marker", result.DevProfile)
	}
	if len(result.Items) != 0 {
		t.Fatalf("items = %d, want empty local profile usage until events are reported", len(result.Items))
	}
}

func TestLocalMeteringServiceReportsTokenUsageIdempotently(t *testing.T) {
	service := NewLocalMeteringService(WithMeteringClock(func() time.Time {
		return time.Unix(2300, 0).UTC()
	}))
	request := ports.TokenUsageReportRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "token-usage-a",
		Source:         "model-service",
		Model:          "qwen2.5",
		InputTokens:    10,
		OutputTokens:   15,
		RequestID:      "req-a",
	}

	first, err := service.ReportTokenUsage(context.Background(), request)
	if err != nil {
		t.Fatalf("ReportTokenUsage(first) error = %v", err)
	}
	second, err := service.ReportTokenUsage(context.Background(), request)
	if err != nil {
		t.Fatalf("ReportTokenUsage(second) error = %v", err)
	}
	if second.ReportID != first.ReportID || second.State != ports.TokenUsageReportDuplicate {
		t.Fatalf("second = %+v, want duplicate replay of %q", second, first.ReportID)
	}
	if first.TotalTokens != 25 {
		t.Fatalf("total tokens = %d, want 25", first.TotalTokens)
	}

	usage, err := service.QueryUsage(context.Background(), ports.MeteringUsageQueryRequest{
		TenantID:     "tenant-a",
		ResourceType: ports.MeteringResourceTokenTotal,
	})
	if err != nil {
		t.Fatalf("QueryUsage after report error = %v", err)
	}
	if len(usage.Items) != 1 || usage.Items[0].TotalQuantity != 25 {
		t.Fatalf("usage = %+v, want one token_total item with 25", usage.Items)
	}
}
