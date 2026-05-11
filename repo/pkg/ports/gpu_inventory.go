package ports

import "context"

type GPUVendor string

const (
	GPUVendorNVIDIA  GPUVendor = "nvidia"
	GPUVendorHuawei  GPUVendor = "huawei"
	GPUVendorHygon   GPUVendor = "hygon"
	GPUVendorUnknown GPUVendor = "unknown"
)

type GPUVirtualizationMode string

const (
	GPUVirtualizationNone GPUVirtualizationMode = "none"
	GPUVirtualizationMIG  GPUVirtualizationMode = "mig"
	GPUVirtualizationVGPU GPUVirtualizationMode = "vgpu"
)

type GPUDeviceClass struct {
	Vendor             GPUVendor
	Model              string
	MemoryMiB          int64
	ResourceName       string
	VirtualizationMode GPUVirtualizationMode
	DriverVersion      string
	RuntimeVersion     string
	Capabilities       []string
}

type GPUNodeClass struct {
	NodeName      string
	Vendor        GPUVendor
	Model         string
	KernelVersion string
	OSImage       string
	Pool          string
	Labels        map[string]string
	Taints        []string
	Devices       []GPUDeviceClass
	Ready         bool
	Reason        string
}

type GPUDiscoveryFilter struct {
	Vendors []GPUVendor
	Pool    string
	Labels  map[string]string
}

type GPUSchedulingRequest struct {
	TenantID             string
	WorkloadID           string
	PreferredVendors     []GPUVendor
	PreferredModels      []string
	RequiredMemoryMiB    int64
	RequiredCount        int
	VirtualizationModes  []GPUVirtualizationMode
	RequiredCapabilities []string
	Pool                 string
}

type GPUSchedulingDecision struct {
	NodeSelector     map[string]string
	Tolerations      []string
	ResourceName     string
	ResourceQuantity string
	RuntimeClassName string
	SchedulerName    string
	QueueName        string
	Reasons          []string
}

// GPUInventory discovers heterogeneous GPU capacity and maps workload intent to
// scheduling constraints. Implementations may use Kubernetes labels, GPU Feature
// Discovery, vendor device plugins, or customer inventory systems.
type GPUInventory interface {
	ListNodeClasses(ctx context.Context, filter GPUDiscoveryFilter) ([]GPUNodeClass, error)
	GetNodeClass(ctx context.Context, nodeName string) (GPUNodeClass, error)
	PlanScheduling(ctx context.Context, request GPUSchedulingRequest) (GPUSchedulingDecision, error)
}
