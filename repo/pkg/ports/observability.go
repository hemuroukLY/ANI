package ports

import (
	"context"
	"time"
)

type ObservabilityResultType string

const (
	ObservabilityResultVector ObservabilityResultType = "vector"
	ObservabilityResultMatrix ObservabilityResultType = "matrix"
	ObservabilityResultScalar ObservabilityResultType = "scalar"
	ObservabilityResultString ObservabilityResultType = "string"
)

type ObservabilityAlertSeverity string

const (
	ObservabilityAlertSeverityInfo     ObservabilityAlertSeverity = "info"
	ObservabilityAlertSeverityWarning  ObservabilityAlertSeverity = "warning"
	ObservabilityAlertSeverityCritical ObservabilityAlertSeverity = "critical"
)

type ObservabilityAlertRuleState string

const (
	ObservabilityAlertRuleActive   ObservabilityAlertRuleState = "active"
	ObservabilityAlertRuleDisabled ObservabilityAlertRuleState = "disabled"
	ObservabilityAlertRuleDeleted  ObservabilityAlertRuleState = "deleted"
)

type ObservabilityQueryRequest struct {
	TenantID string
	Query    string
	Time     time.Time
	Timeout  time.Duration
}

type ObservabilityQuerySample struct {
	Metric    map[string]string
	Value     float64
	Timestamp time.Time
}

type ObservabilityQueryResult struct {
	Query      string
	ResultType ObservabilityResultType
	Results    []ObservabilityQuerySample
	DevProfile DevProfileInfo
}

type ObservabilityAlertRuleRecord struct {
	TenantID    string
	RuleID      string
	Name        string
	PromQL      string
	Duration    time.Duration
	Severity    ObservabilityAlertSeverity
	Labels      map[string]string
	Annotations map[string]string
	Enabled     bool
	State       ObservabilityAlertRuleState
	DevProfile  DevProfileInfo
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ObservabilityAlertRuleCreateRequest struct {
	TenantID       string
	IdempotencyKey string
	Name           string
	PromQL         string
	Duration       time.Duration
	Severity       ObservabilityAlertSeverity
	Labels         map[string]string
	Annotations    map[string]string
	Enabled        bool
}

type ObservabilityAlertRuleUpdateRequest struct {
	TenantID       string
	RuleID         string
	IdempotencyKey string
	Name           string
	PromQL         string
	Duration       time.Duration
	Severity       ObservabilityAlertSeverity
	Labels         map[string]string
	Annotations    map[string]string
	Enabled        *bool
}

type ObservabilityAlertRuleGetRequest struct {
	TenantID string
	RuleID   string
}

type ObservabilityAlertRuleListRequest struct {
	TenantID string
	Limit    int
	Cursor   string
}

type ObservabilityService interface {
	Query(ctx context.Context, request ObservabilityQueryRequest) (ObservabilityQueryResult, error)
	CreateAlertRule(ctx context.Context, request ObservabilityAlertRuleCreateRequest) (ObservabilityAlertRuleRecord, error)
	ListAlertRules(ctx context.Context, request ObservabilityAlertRuleListRequest) ([]ObservabilityAlertRuleRecord, error)
	GetAlertRule(ctx context.Context, request ObservabilityAlertRuleGetRequest) (ObservabilityAlertRuleRecord, error)
	UpdateAlertRule(ctx context.Context, request ObservabilityAlertRuleUpdateRequest) (ObservabilityAlertRuleRecord, error)
	DeleteAlertRule(ctx context.Context, request ObservabilityAlertRuleGetRequest) (ObservabilityAlertRuleRecord, error)
}
