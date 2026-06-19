package runtime

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kubercloud/ani/pkg/ports"
)

type LocalNetworkService struct {
	mu                sync.RWMutex
	now               func() time.Time
	store             ports.NetworkResourceStore
	routeRenderer     ports.NetworkProviderRenderer
	routeDryRun       ports.NetworkProviderDryRun
	routeApply        ports.NetworkProviderApply
	routeStatus       ports.NetworkProviderStatusReader
	routeExecution    NetworkProviderExecutionConfig
	vpcs              map[string]ports.NetworkVPCRecord
	subnets           map[string]ports.NetworkSubnetRecord
	securityGroup     map[string]ports.NetworkSecurityGroupRecord
	loadBalancers     map[string]ports.NetworkLoadBalancerRecord
	routes            map[string]ports.NetworkRouteRecord
	vpcIdempotency    map[string]string
	subnetIdempotency map[string]string
	securityGroupIdem map[string]string
	loadBalancerIdem  map[string]string
	routeIdempotency  map[string]string
}

type NetworkServiceOption func(*LocalNetworkService)

type NetworkProviderExecutionConfig struct {
	UserID          string
	PermissionProof string
}

func WithNetworkServiceClock(now func() time.Time) NetworkServiceOption {
	return func(service *LocalNetworkService) {
		if now != nil {
			service.now = now
		}
	}
}

func WithNetworkResourceStore(store ports.NetworkResourceStore) NetworkServiceOption {
	return func(service *LocalNetworkService) {
		service.store = store
	}
}

func WithNetworkRouteProvider(
	renderer ports.NetworkProviderRenderer,
	dryRun ports.NetworkProviderDryRun,
	apply ports.NetworkProviderApply,
	status ports.NetworkProviderStatusReader,
	execution NetworkProviderExecutionConfig,
) NetworkServiceOption {
	return func(service *LocalNetworkService) {
		service.routeRenderer = renderer
		service.routeDryRun = dryRun
		service.routeApply = apply
		service.routeStatus = status
		service.routeExecution = execution
	}
}

func NewLocalNetworkService(options ...NetworkServiceOption) *LocalNetworkService {
	service := &LocalNetworkService{
		now:               func() time.Time { return time.Now().UTC() },
		vpcs:              map[string]ports.NetworkVPCRecord{},
		subnets:           map[string]ports.NetworkSubnetRecord{},
		securityGroup:     map[string]ports.NetworkSecurityGroupRecord{},
		loadBalancers:     map[string]ports.NetworkLoadBalancerRecord{},
		routes:            map[string]ports.NetworkRouteRecord{},
		vpcIdempotency:    map[string]string{},
		subnetIdempotency: map[string]string{},
		securityGroupIdem: map[string]string{},
		loadBalancerIdem:  map[string]string{},
		routeIdempotency:  map[string]string{},
	}
	for _, option := range options {
		option(service)
	}
	return service
}

func (s *LocalNetworkService) CreateVPC(ctx context.Context, request ports.NetworkVPCCreateRequest) (ports.NetworkVPCRecord, error) {
	if err := requireNetworkTenantAndName(request.TenantID, request.Name); err != nil {
		return ports.NetworkVPCRecord{}, err
	}
	idemKey, err := requireIdempotencyKey(request.TenantID, request.IdempotencyKey)
	if err != nil {
		return ports.NetworkVPCRecord{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if id, ok := s.vpcIdempotency[idemKey]; ok {
		if record, exists := s.vpcs[id]; exists {
			return record, nil
		}
	}
	now := s.now().UTC()
	record := ports.NetworkVPCRecord{
		TenantID:  request.TenantID,
		VPCID:     "vpc_" + uuid.NewString(),
		Name:      strings.TrimSpace(request.Name),
		CIDR:      firstNetworkNonEmpty(request.CIDR, "10.0.0.0/16"),
		State:     ports.NetworkResourceAvailable,
		Reason:    "created by local network profile",
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.vpcs[record.VPCID] = record
	s.vpcIdempotency[idemKey] = record.VPCID
	if err := s.upsertVPC(ctx, record); err != nil {
		return ports.NetworkVPCRecord{}, err
	}
	return record, nil
}

func (s *LocalNetworkService) ListVPCs(_ context.Context, request ports.NetworkResourceListRequest) ([]ports.NetworkVPCRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]ports.NetworkVPCRecord, 0, len(s.vpcs))
	for _, record := range s.vpcs {
		if record.TenantID == request.TenantID && record.State != ports.NetworkResourceDeleted {
			items = append(items, record)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].UpdatedAt.After(items[j].UpdatedAt) })
	return items, nil
}

func (s *LocalNetworkService) GetVPC(_ context.Context, request ports.NetworkResourceGetRequest) (ports.NetworkVPCRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, ok := s.vpcs[request.ResourceID]
	if !ok || record.TenantID != request.TenantID || record.State == ports.NetworkResourceDeleted {
		return ports.NetworkVPCRecord{}, ports.ErrNotFound
	}
	return record, nil
}

func (s *LocalNetworkService) DeleteVPC(ctx context.Context, request ports.NetworkResourceGetRequest) (ports.NetworkVPCRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.vpcs[request.ResourceID]
	if !ok || record.TenantID != request.TenantID || record.State == ports.NetworkResourceDeleted {
		return ports.NetworkVPCRecord{}, ports.ErrNotFound
	}
	now := s.now().UTC()
	record.State = ports.NetworkResourceDeleted
	record.Reason = "deleted by local network profile"
	record.UpdatedAt = now
	s.vpcs[record.VPCID] = record
	if err := s.upsertVPC(ctx, record); err != nil {
		return ports.NetworkVPCRecord{}, err
	}
	return record, nil
}

func (s *LocalNetworkService) CreateSubnet(ctx context.Context, request ports.NetworkSubnetCreateRequest) (ports.NetworkSubnetRecord, error) {
	if err := requireNetworkTenantAndName(request.TenantID, request.Name); err != nil {
		return ports.NetworkSubnetRecord{}, err
	}
	idemKey, err := requireIdempotencyKey(request.TenantID, request.IdempotencyKey)
	if err != nil {
		return ports.NetworkSubnetRecord{}, err
	}
	if strings.TrimSpace(request.VPCID) == "" {
		return ports.NetworkSubnetRecord{}, fmt.Errorf("%w: vpc_id is required", ports.ErrInvalid)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if id, ok := s.subnetIdempotency[idemKey]; ok {
		if record, exists := s.subnets[id]; exists {
			return record, nil
		}
	}
	now := s.now().UTC()
	record := ports.NetworkSubnetRecord{
		TenantID:  request.TenantID,
		SubnetID:  "subnet_" + uuid.NewString(),
		VPCID:     strings.TrimSpace(request.VPCID),
		Name:      strings.TrimSpace(request.Name),
		CIDR:      firstNetworkNonEmpty(request.CIDR, "10.0.1.0/24"),
		Gateway:   strings.TrimSpace(request.Gateway),
		State:     ports.NetworkResourceAvailable,
		Reason:    "created by local network profile",
		CreatedAt: now,
		UpdatedAt: now,
	}
	vpc, ok := s.vpcs[record.VPCID]
	if !ok || vpc.TenantID != request.TenantID || vpc.State == ports.NetworkResourceDeleted {
		return ports.NetworkSubnetRecord{}, fmt.Errorf("%w: vpc not found", ports.ErrNotFound)
	}
	s.subnets[record.SubnetID] = record
	s.subnetIdempotency[idemKey] = record.SubnetID
	if err := s.upsertSubnet(ctx, record); err != nil {
		return ports.NetworkSubnetRecord{}, err
	}
	return record, nil
}

func (s *LocalNetworkService) ListSubnets(_ context.Context, request ports.NetworkResourceListRequest) ([]ports.NetworkSubnetRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]ports.NetworkSubnetRecord, 0, len(s.subnets))
	for _, record := range s.subnets {
		if record.TenantID == request.TenantID && record.State != ports.NetworkResourceDeleted {
			items = append(items, record)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].UpdatedAt.After(items[j].UpdatedAt) })
	return items, nil
}

func (s *LocalNetworkService) GetSubnet(_ context.Context, request ports.NetworkResourceGetRequest) (ports.NetworkSubnetRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, ok := s.subnets[request.ResourceID]
	if !ok || record.TenantID != request.TenantID || record.State == ports.NetworkResourceDeleted {
		return ports.NetworkSubnetRecord{}, ports.ErrNotFound
	}
	return record, nil
}

func (s *LocalNetworkService) DeleteSubnet(ctx context.Context, request ports.NetworkResourceGetRequest) (ports.NetworkSubnetRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.subnets[request.ResourceID]
	if !ok || record.TenantID != request.TenantID || record.State == ports.NetworkResourceDeleted {
		return ports.NetworkSubnetRecord{}, ports.ErrNotFound
	}
	record.State = ports.NetworkResourceDeleted
	record.Reason = "deleted by local network profile"
	record.UpdatedAt = s.now().UTC()
	s.subnets[record.SubnetID] = record
	if err := s.upsertSubnet(ctx, record); err != nil {
		return ports.NetworkSubnetRecord{}, err
	}
	return record, nil
}

func (s *LocalNetworkService) CreateSecurityGroup(ctx context.Context, request ports.NetworkSecurityGroupCreateRequest) (ports.NetworkSecurityGroupRecord, error) {
	if err := requireNetworkTenantAndName(request.TenantID, request.Name); err != nil {
		return ports.NetworkSecurityGroupRecord{}, err
	}
	idemKey, err := requireIdempotencyKey(request.TenantID, request.IdempotencyKey)
	if err != nil {
		return ports.NetworkSecurityGroupRecord{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if id, ok := s.securityGroupIdem[idemKey]; ok {
		if record, exists := s.securityGroup[id]; exists {
			return record, nil
		}
	}
	now := s.now().UTC()
	record := ports.NetworkSecurityGroupRecord{
		TenantID:        request.TenantID,
		SecurityGroupID: "sg_" + uuid.NewString(),
		Name:            strings.TrimSpace(request.Name),
		Description:     strings.TrimSpace(request.Description),
		Rules:           append([]ports.NetworkSecurityGroupRule(nil), request.Rules...),
		State:           ports.NetworkResourceAvailable,
		Reason:          "created by local network profile",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	s.securityGroup[record.SecurityGroupID] = record
	s.securityGroupIdem[idemKey] = record.SecurityGroupID
	if err := s.upsertSecurityGroup(ctx, record); err != nil {
		return ports.NetworkSecurityGroupRecord{}, err
	}
	return record, nil
}

func (s *LocalNetworkService) ListSecurityGroups(_ context.Context, request ports.NetworkResourceListRequest) ([]ports.NetworkSecurityGroupRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]ports.NetworkSecurityGroupRecord, 0, len(s.securityGroup))
	for _, record := range s.securityGroup {
		if record.TenantID == request.TenantID && record.State != ports.NetworkResourceDeleted {
			items = append(items, record)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].UpdatedAt.After(items[j].UpdatedAt) })
	return items, nil
}

func (s *LocalNetworkService) GetSecurityGroup(_ context.Context, request ports.NetworkResourceGetRequest) (ports.NetworkSecurityGroupRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, ok := s.securityGroup[request.ResourceID]
	if !ok || record.TenantID != request.TenantID || record.State == ports.NetworkResourceDeleted {
		return ports.NetworkSecurityGroupRecord{}, ports.ErrNotFound
	}
	return record, nil
}

func (s *LocalNetworkService) DeleteSecurityGroup(ctx context.Context, request ports.NetworkResourceGetRequest) (ports.NetworkSecurityGroupRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.securityGroup[request.ResourceID]
	if !ok || record.TenantID != request.TenantID || record.State == ports.NetworkResourceDeleted {
		return ports.NetworkSecurityGroupRecord{}, ports.ErrNotFound
	}
	record.State = ports.NetworkResourceDeleted
	record.Reason = "deleted by local network profile"
	record.UpdatedAt = s.now().UTC()
	s.securityGroup[record.SecurityGroupID] = record
	if err := s.upsertSecurityGroup(ctx, record); err != nil {
		return ports.NetworkSecurityGroupRecord{}, err
	}
	return record, nil
}

func (s *LocalNetworkService) CreateLoadBalancer(ctx context.Context, request ports.NetworkLoadBalancerCreateRequest) (ports.NetworkLoadBalancerRecord, error) {
	if err := requireNetworkTenantAndName(request.TenantID, request.Name); err != nil {
		return ports.NetworkLoadBalancerRecord{}, err
	}
	idemKey, err := requireIdempotencyKey(request.TenantID, request.IdempotencyKey)
	if err != nil {
		return ports.NetworkLoadBalancerRecord{}, err
	}
	if strings.TrimSpace(request.VPCID) == "" {
		return ports.NetworkLoadBalancerRecord{}, fmt.Errorf("%w: vpc_id is required", ports.ErrInvalid)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if id, ok := s.loadBalancerIdem[idemKey]; ok {
		if record, exists := s.loadBalancers[id]; exists {
			return record, nil
		}
	}
	now := s.now().UTC()
	record := ports.NetworkLoadBalancerRecord{
		TenantID:       request.TenantID,
		LoadBalancerID: "lb_" + uuid.NewString(),
		Name:           strings.TrimSpace(request.Name),
		VPCID:          strings.TrimSpace(request.VPCID),
		SubnetID:       strings.TrimSpace(request.SubnetID),
		Scheme:         firstNetworkNonEmpty(request.Scheme, "internal"),
		VIP:            "local-dev",
		Listeners:      append([]ports.NetworkLoadBalancerListener(nil), request.Listeners...),
		State:          ports.NetworkResourceAvailable,
		Reason:         "created by local network profile",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	vpc, ok := s.vpcs[record.VPCID]
	if !ok || vpc.TenantID != request.TenantID || vpc.State == ports.NetworkResourceDeleted {
		return ports.NetworkLoadBalancerRecord{}, fmt.Errorf("%w: vpc not found", ports.ErrNotFound)
	}
	s.loadBalancers[record.LoadBalancerID] = record
	s.loadBalancerIdem[idemKey] = record.LoadBalancerID
	if err := s.upsertLoadBalancer(ctx, record); err != nil {
		return ports.NetworkLoadBalancerRecord{}, err
	}
	return record, nil
}

func (s *LocalNetworkService) ListLoadBalancers(_ context.Context, request ports.NetworkResourceListRequest) ([]ports.NetworkLoadBalancerRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]ports.NetworkLoadBalancerRecord, 0, len(s.loadBalancers))
	for _, record := range s.loadBalancers {
		if record.TenantID == request.TenantID && record.State != ports.NetworkResourceDeleted {
			items = append(items, record)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].UpdatedAt.After(items[j].UpdatedAt) })
	return items, nil
}

func (s *LocalNetworkService) GetLoadBalancer(_ context.Context, request ports.NetworkResourceGetRequest) (ports.NetworkLoadBalancerRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, ok := s.loadBalancers[request.ResourceID]
	if !ok || record.TenantID != request.TenantID || record.State == ports.NetworkResourceDeleted {
		return ports.NetworkLoadBalancerRecord{}, ports.ErrNotFound
	}
	return record, nil
}

func (s *LocalNetworkService) DeleteLoadBalancer(ctx context.Context, request ports.NetworkResourceGetRequest) (ports.NetworkLoadBalancerRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.loadBalancers[request.ResourceID]
	if !ok || record.TenantID != request.TenantID || record.State == ports.NetworkResourceDeleted {
		return ports.NetworkLoadBalancerRecord{}, ports.ErrNotFound
	}
	record.State = ports.NetworkResourceDeleted
	record.Reason = "deleted by local network profile"
	record.UpdatedAt = s.now().UTC()
	s.loadBalancers[record.LoadBalancerID] = record
	if err := s.upsertLoadBalancer(ctx, record); err != nil {
		return ports.NetworkLoadBalancerRecord{}, err
	}
	return record, nil
}

func (s *LocalNetworkService) CreateRoute(ctx context.Context, request ports.NetworkRouteCreateRequest) (ports.NetworkRouteRecord, error) {
	idemKey, err := requireIdempotencyKey(request.TenantID, request.IdempotencyKey)
	if err != nil {
		return ports.NetworkRouteRecord{}, err
	}
	if strings.TrimSpace(request.VPCID) == "" {
		return ports.NetworkRouteRecord{}, fmt.Errorf("%w: vpc_id is required", ports.ErrInvalid)
	}
	if strings.TrimSpace(request.DestinationCIDR) == "" || strings.TrimSpace(request.NextHopType) == "" || strings.TrimSpace(request.NextHopID) == "" {
		return ports.NetworkRouteRecord{}, fmt.Errorf("%w: destination_cidr/next_hop_type/next_hop_id are required", ports.ErrInvalid)
	}
	nextHopType := strings.ToLower(strings.TrimSpace(request.NextHopType))
	if nextHopType != "gateway" && nextHopType != "instance" && nextHopType != "nat" {
		return ports.NetworkRouteRecord{}, fmt.Errorf("%w: unsupported route next_hop_type %q", ports.ErrUnsupported, request.NextHopType)
	}
	s.mu.Lock()
	if id, ok := s.routeIdempotency[idemKey]; ok {
		if record, exists := s.routes[id]; exists {
			s.mu.Unlock()
			return record, nil
		}
	}
	vpc, ok := s.vpcs[strings.TrimSpace(request.VPCID)]
	if !ok || vpc.TenantID != request.TenantID || vpc.State == ports.NetworkResourceDeleted {
		s.mu.Unlock()
		return ports.NetworkRouteRecord{}, fmt.Errorf("%w: vpc not found", ports.ErrNotFound)
	}
	providerConfigured := s.routeProviderConfigured()
	record := ports.NetworkRouteRecord{
		TenantID:        request.TenantID,
		RouteID:         "rt_" + uuid.NewString(),
		VPCID:           strings.TrimSpace(request.VPCID),
		DestinationCIDR: strings.TrimSpace(request.DestinationCIDR),
		NextHopType:     nextHopType,
		NextHopID:       strings.TrimSpace(request.NextHopID),
		Description:     strings.TrimSpace(request.Description),
		State:           ports.NetworkResourceAvailable,
		CreatedAt:       s.now().UTC(),
	}
	if providerConfigured {
		record.State = ports.NetworkResourcePending
	}
	s.routes[record.RouteID] = record
	s.routeIdempotency[idemKey] = record.RouteID
	s.mu.Unlock()
	if !providerConfigured {
		return record, nil
	}
	applied, err := s.applyRouteProvider(ctx, record)
	s.mu.Lock()
	defer s.mu.Unlock()
	if err != nil {
		record.State = ports.NetworkResourceFailed
		s.routes[record.RouteID] = record
		return ports.NetworkRouteRecord{}, err
	}
	record = applied
	if _, exists := s.routes[record.RouteID]; exists {
		s.routes[record.RouteID] = record
	}
	return record, nil
}

func (s *LocalNetworkService) ListRoutes(_ context.Context, request ports.NetworkRouteListRequest) ([]ports.NetworkRouteRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]ports.NetworkRouteRecord, 0, len(s.routes))
	for _, record := range s.routes {
		if record.TenantID != request.TenantID {
			continue
		}
		if strings.TrimSpace(request.VPCID) != "" && record.VPCID != strings.TrimSpace(request.VPCID) {
			continue
		}
		items = append(items, record)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	return items, nil
}

func (s *LocalNetworkService) upsertVPC(ctx context.Context, record ports.NetworkVPCRecord) error {
	if s.store == nil {
		return nil
	}
	return s.store.UpsertVPC(ctx, record)
}

func (s *LocalNetworkService) upsertSubnet(ctx context.Context, record ports.NetworkSubnetRecord) error {
	if s.store == nil {
		return nil
	}
	return s.store.UpsertSubnet(ctx, record)
}

func (s *LocalNetworkService) upsertSecurityGroup(ctx context.Context, record ports.NetworkSecurityGroupRecord) error {
	if s.store == nil {
		return nil
	}
	return s.store.UpsertSecurityGroup(ctx, record)
}

func (s *LocalNetworkService) upsertLoadBalancer(ctx context.Context, record ports.NetworkLoadBalancerRecord) error {
	if s.store == nil {
		return nil
	}
	return s.store.UpsertLoadBalancer(ctx, record)
}

func (s *LocalNetworkService) routeProviderConfigured() bool {
	return s.routeRenderer != nil || s.routeDryRun != nil || s.routeApply != nil || s.routeStatus != nil
}

func (s *LocalNetworkService) applyRouteProvider(ctx context.Context, record ports.NetworkRouteRecord) (ports.NetworkRouteRecord, error) {
	if s.routeRenderer == nil || s.routeDryRun == nil || s.routeApply == nil || s.routeStatus == nil {
		return ports.NetworkRouteRecord{}, fmt.Errorf("%w: network route provider requires renderer, dry-run, apply, and status adapters", ports.ErrNotConfigured)
	}
	userID := strings.TrimSpace(s.routeExecution.UserID)
	permissionProof := strings.TrimSpace(s.routeExecution.PermissionProof)
	if userID == "" || permissionProof == "" {
		return ports.NetworkRouteRecord{}, fmt.Errorf("%w: network route provider requires explicit user id and permission proof", ports.ErrInvalid)
	}
	manifests, err := s.routeRenderer.RenderRoute(ctx, record)
	if err != nil {
		return ports.NetworkRouteRecord{}, err
	}
	requestedAt := s.now().UTC()
	dryRun, err := s.routeDryRun.DryRun(ctx, ports.NetworkProviderDryRunRequest{
		TenantID:        record.TenantID,
		UserID:          userID,
		ResourceKind:    "route",
		ResourceID:      record.RouteID,
		Operation:       ports.NetworkProviderOperationCreate,
		Manifests:       manifests,
		PermissionProof: permissionProof,
		RequestedAt:     requestedAt,
	})
	if err != nil {
		return ports.NetworkRouteRecord{}, err
	}
	apply, err := s.routeApply.Apply(ctx, ports.NetworkProviderApplyRequest{
		TenantID:        record.TenantID,
		UserID:          userID,
		ResourceKind:    "route",
		ResourceID:      record.RouteID,
		Operation:       ports.NetworkProviderOperationCreate,
		Manifests:       manifests,
		PermissionProof: permissionProof,
		DryRunResult:    dryRun,
		RequestedAt:     requestedAt,
	})
	if err != nil {
		return ports.NetworkRouteRecord{}, err
	}
	observation, err := s.routeStatus.Observe(ctx, ports.NetworkProviderStatusRequest{
		TenantID:        record.TenantID,
		UserID:          userID,
		ResourceKind:    "route",
		ResourceID:      record.RouteID,
		ApplyResult:     apply,
		PermissionProof: permissionProof,
		RequestedAt:     requestedAt,
	})
	if err != nil {
		return ports.NetworkRouteRecord{}, err
	}
	if observation.State == "" {
		record.State = ports.NetworkResourceAvailable
	} else {
		record.State = observation.State
	}
	record.Provider = firstNetworkNonEmpty(observation.Provider, apply.Provider)
	record.RealProvider = apply.Applied
	return record, nil
}

func requireNetworkTenantAndName(tenantID string, name string) error {
	if strings.TrimSpace(tenantID) == "" {
		return fmt.Errorf("%w: tenant_id is required", ports.ErrInvalid)
	}
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: name is required", ports.ErrInvalid)
	}
	return nil
}

func firstNetworkNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func requireIdempotencyKey(tenantID string, key string) (string, error) {
	tenantID = strings.TrimSpace(tenantID)
	key = strings.TrimSpace(key)
	if tenantID == "" {
		return "", fmt.Errorf("%w: tenant_id is required", ports.ErrInvalid)
	}
	if key == "" {
		return "", fmt.Errorf("%w: idempotency_key is required", ports.ErrInvalid)
	}
	return tenantID + "\x00" + key, nil
}

var _ ports.NetworkService = (*LocalNetworkService)(nil)
