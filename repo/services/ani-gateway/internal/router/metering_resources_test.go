package router

import (
	"context"
	"testing"

	"github.com/kubercloud/ani/pkg/ports"
)

func TestMeteringAPIUsageResponseMarksLocalProfile(t *testing.T) {
	api := newMeteringAPI()

	result, err := api.service.QueryUsage(context.Background(), ports.MeteringUsageQueryRequest{
		TenantID: "tenant-a",
	})
	if err != nil {
		t.Fatalf("QueryUsage error = %v", err)
	}
	response := meteringUsageFromResult(result)
	if response.Total != 0 {
		t.Fatalf("total = %d, want 0", response.Total)
	}
	requireLocalCoreDevProfile(t, response.DevProfile, "local-metering-service")
}

func TestMeteringAPITokenUsageReportResponse(t *testing.T) {
	api := newMeteringAPI()

	report, err := api.service.ReportTokenUsage(context.Background(), ports.TokenUsageReportRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "token-usage-router",
		Source:         "model-service",
		Model:          "qwen2.5",
		InputTokens:    7,
		OutputTokens:   11,
	})
	if err != nil {
		t.Fatalf("ReportTokenUsage error = %v", err)
	}
	response := tokenUsageReportFromRecord(report)
	if response.ID == "" || response.TotalTokens != 18 || response.State != "accepted" {
		t.Fatalf("response = %+v, want accepted total 18", response)
	}
	requireLocalCoreDevProfile(t, response.DevProfile, "local-metering-service")
}
