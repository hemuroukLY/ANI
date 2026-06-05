package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestParseCoreCommandsRejectsServicesResources(t *testing.T) {
	_, err := parseCommand([]string{"model", "list"})
	if err == nil || !strings.Contains(err.Error(), "services resources are not supported") {
		t.Fatalf("parseCommand error = %v, want Services rejection", err)
	}
}

func TestParseCoreCommandsBuildsInstanceListRequest(t *testing.T) {
	command, err := parseCommand([]string{"instances", "list", "--limit", "10", "--cursor", "next"})
	if err != nil {
		t.Fatalf("parseCommand error = %v", err)
	}
	if command.Method != http.MethodGet || command.Path != "/instances" {
		t.Fatalf("command = %+v, want GET /instances", command)
	}
	if command.Query.Get("limit") != "10" || command.Query.Get("cursor") != "next" {
		t.Fatalf("query = %s, want limit and cursor", command.Query.Encode())
	}
}

func TestParseCoreCommandsBuildsExpandedCoreListRequests(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		wantPath string
	}{
		{name: "vpcs", args: []string{"network-vpcs", "list"}, wantPath: "/networks/vpcs"},
		{name: "subnets", args: []string{"network-subnets", "list"}, wantPath: "/networks/subnets"},
		{name: "security groups", args: []string{"network-security-groups", "list"}, wantPath: "/networks/security-groups"},
		{name: "load balancers", args: []string{"network-load-balancers", "list"}, wantPath: "/networks/load-balancers"},
		{name: "volumes", args: []string{"volumes", "list"}, wantPath: "/volumes"},
		{name: "filesystems", args: []string{"filesystems", "list"}, wantPath: "/filesystems"},
		{name: "objects", args: []string{"objects", "list"}, wantPath: "/objects"},
		{name: "vector stores", args: []string{"vector-stores", "list"}, wantPath: "/vector-stores"},
		{name: "encryption keys", args: []string{"encryption-keys", "list"}, wantPath: "/encryption/keys"},
		{name: "observability alert rules", args: []string{"observability-alert-rules", "list"}, wantPath: "/observability/alert-rules"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			command, err := parseCommand(append(tt.args, "--limit", "5"))
			if err != nil {
				t.Fatalf("parseCommand error = %v", err)
			}
			if command.Method != http.MethodGet || command.Path != tt.wantPath {
				t.Fatalf("command = %+v, want GET %s", command, tt.wantPath)
			}
			if command.Query.Get("limit") != "5" {
				t.Fatalf("query = %s, want limit=5", command.Query.Encode())
			}
		})
	}
}

func TestParseCoreCommandsBuildsObservabilityQueryRequest(t *testing.T) {
	command, err := parseCommand([]string{"observability-query", "get", "--query", "up", "--tenant-id", "tenant-a"})
	if err != nil {
		t.Fatalf("parseCommand error = %v", err)
	}
	if command.Method != http.MethodGet || command.Path != "/observability/query" {
		t.Fatalf("command = %+v, want GET /observability/query", command)
	}
	if command.Query.Get("query") != "up" || command.Query.Get("tenant_id") != "tenant-a" {
		t.Fatalf("query = %s, want query and tenant_id", command.Query.Encode())
	}
}

func TestRunSendsCoreRequestWithBearerToken(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotAuth string
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		gotAuth = req.Header.Get("Authorization")
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"items":[],"total":0}`)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}, nil
	})}

	body, err := execute(t.Context(), client, "http://127.0.0.1:4010/api/v1", "dev-token", command{
		Method: http.MethodGet,
		Path:   "/instances",
		Query:  nil,
	})

	if err != nil {
		t.Fatalf("execute error = %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/api/v1/instances" {
		t.Fatalf("request = %s %s, want GET /api/v1/instances", gotMethod, gotPath)
	}
	if gotAuth != "Bearer dev-token" {
		t.Fatalf("Authorization = %q, want bearer token", gotAuth)
	}
	if !strings.Contains(string(body), `"total":0`) {
		t.Fatalf("body = %s, want response body", string(body))
	}
}

func TestRunPrintsVersionWithoutCoreAPIRequest(t *testing.T) {
	oldVersion := Version
	oldBuildTime := BuildTime
	Version = "v0.1.0-test"
	BuildTime = "20260604T000000Z"
	t.Cleanup(func() {
		Version = oldVersion
		BuildTime = oldBuildTime
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{"--version"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("run exit = %d, stderr = %s", exitCode, stderr.String())
	}
	output := stdout.String()
	if !strings.Contains(output, "ani version v0.1.0-test") || !strings.Contains(output, "build 20260604T000000Z") {
		t.Fatalf("version output = %q, want version and build time", output)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunPrintsJSONVersionForReleaseEvidence(t *testing.T) {
	oldVersion := Version
	oldBuildTime := BuildTime
	Version = "v0.1.0-test"
	BuildTime = "20260604T000000Z"
	t.Cleanup(func() {
		Version = oldVersion
		BuildTime = oldBuildTime
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{"--version", "--version-format", "json"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("run exit = %d, stderr = %s", exitCode, stderr.String())
	}
	var metadata map[string]string
	if err := json.Unmarshal(stdout.Bytes(), &metadata); err != nil {
		t.Fatalf("version output is not JSON: %v; output = %q", err, stdout.String())
	}
	if metadata["name"] != "ani" || metadata["scope"] != "core" {
		t.Fatalf("metadata = %#v, want ani core metadata", metadata)
	}
	if metadata["version"] != "v0.1.0-test" || metadata["build_time"] != "20260604T000000Z" {
		t.Fatalf("metadata = %#v, want version and build time", metadata)
	}
}
