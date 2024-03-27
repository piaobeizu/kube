package base

import (
	"fmt"
	"strconv"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	DefaultSchedulerName     = v1.DefaultSchedulerName
	GPUSchedulerName         = "gpu-scheduler"
	VGPUSchedulerName        = "volcano" // TODO change to vocano when in production
	CPUPriorityClass         = "jeeves-cpu-compute"
	GPUPriorityClass         = "jeeves-gpu-compute"
	HPCPriorityClass         = "jeeves-hpc-compute"
	VGPUPriorityClass        = "jeeves-gpu-compute" // TODO change to jeeves-vgpu-compute when in production
	ResourceNvidiaGPU        = "nvidia.com/gpu"
	ResourceVendorVGPU       = "baidu.com/cgpu"
	ResourceVendorGPUPercent = "baidu.com/cgpu_core"
	ResourceVendorGPUMem     = "baidu.com/cgpu_memory"
	ResourceRdma             = "rdma/hca"
	ResourceInfinityBand     = "mellanox.com/InfiniBand"
	A800PriorityClass        = "jeeves-a800-compute"

	GPUSeriesA100     = "a100"
	GPUSeriesA800     = "a800"
	A100TolerationKey = "gpu-series"
)

type Resource struct {
	CPUNum           uint   `json:"cpu_num" gorm:"not null;default:1"`
	GPUNum           uint   `json:"gpu_num" gorm:"not null;default:0"`
	GPUPercent       uint   `json:"gpu_percent" gorm:"not null;default:100"`
	GPUMem           uint   `json:"gpu_mem" gorm:"not null;default:12"`
	MemorySize       uint   `json:"memory_size" gorm:"not null;default:4"`
	EphemeralStorage uint   `json:"-" gorm:"not null;default:800"`
	GPUSeries        string `json:"gpu_series" gorm:"type:varchar(32);default:'pascal'"`
}

func (r Resource) ResourceList(region string) v1.ResourceList {
	cpuResource := fmt.Sprintf("%d", r.CPUNum)
	memoryResource := fmt.Sprintf("%dGi", r.MemorySize)
	ephemeralStorageResource := fmt.Sprintf("%dGi", r.EphemeralStorage)
	resources := v1.ResourceList{
		v1.ResourceCPU:              resource.MustParse(cpuResource),
		v1.ResourceMemory:           resource.MustParse(memoryResource),
		v1.ResourceEphemeralStorage: resource.MustParse(ephemeralStorageResource),
	}
	if r.GPUNum > 0 {
		gpuResource := fmt.Sprintf("%d", r.GPUNum)
		if r.GPUSeries == "vGPU" {
			gpuPercentResource := fmt.Sprintf("%d", r.GPUPercent)
			gpuMemResource := fmt.Sprintf("%d", r.GPUMem)
			resources[ResourceVendorVGPU] = resource.MustParse(gpuResource)
			resources[ResourceVendorGPUPercent] = resource.MustParse(gpuPercentResource)
			resources[ResourceVendorGPUMem] = resource.MustParse(gpuMemResource)
		} else {
			resources[ResourceNvidiaGPU] = resource.MustParse(gpuResource)
		}

		if (r.GPUSeries == GPUSeriesA100 || r.GPUSeries == GPUSeriesA800) && r.GPUNum >= 8 {
			if region == "klara-2-pek02" {
				resources[ResourceInfinityBand] = resource.MustParse(strconv.Itoa(1000))
			} else {
				resources[ResourceRdma] = resource.MustParse(strconv.Itoa(1000))
			}
		}
	}
	return resources
}

type SchedulingStrategy interface {
	Labels() map[string]string
	Annotations() map[string]string
	SchedulerName() string
	PriorityClassName() string
	Tolerations() []v1.Toleration
	PreferredSchedulingTerms() []v1.PreferredSchedulingTerm
	NodeSelectorTerms() []v1.NodeSelectorTerm
}
type ResourceStrategy interface {
	Requests() Resource
	Limits() Resource
	Resource() Resource
	Merge(Resource)
	SchedulingStrategy() SchedulingStrategy
	GetRegion() string
}

var _ SchedulingStrategy = CPUSchedulingStrategy{}
var _ SchedulingStrategy = GPUSchedulingStrategy{}
var _ SchedulingStrategy = VGPUSchedulingStrategy{}

type GPUSchedulingStrategy struct {
	Raw       Resource
	CanBorrow bool
}

func (r GPUSchedulingStrategy) Labels() map[string]string {
	return make(map[string]string)
}

func (r GPUSchedulingStrategy) Annotations() map[string]string {
	return make(map[string]string)
}

func (r GPUSchedulingStrategy) SchedulerName() string {
	return GPUSchedulerName
}

func (r GPUSchedulingStrategy) PriorityClassName() string {
	if r.Raw.GPUSeries == GPUSeriesA100 {
		return HPCPriorityClass
	} else if r.Raw.GPUSeries == GPUSeriesA800 {
		return A800PriorityClass
	}
	return GPUPriorityClass
}

func (r GPUSchedulingStrategy) PreferredSchedulingTerms() []v1.PreferredSchedulingTerm {
	return []v1.PreferredSchedulingTerm{}
}

func (r GPUSchedulingStrategy) NodeSelectorTerms() []v1.NodeSelectorTerm {
	var poolKey string
	if r.Raw.GPUNum >= 4 {
		poolKey = "jeeves-graphics-pack/4-8"
	} else {
		poolKey = "jeeves-graphics-pack/1-2"
	}

	matchExpressions := []v1.NodeSelectorRequirement{
		{
			Key:      "node-role.kubernetes.io/jeeves-gpu",
			Operator: v1.NodeSelectorOpExists,
		},
		{
			Key:      poolKey,
			Operator: v1.NodeSelectorOpExists,
		},
		{
			Key:      "node.kubernetes.io/instance-series",
			Operator: v1.NodeSelectorOpIn,
			Values:   []string{r.Raw.GPUSeries},
		},
	}

	return []v1.NodeSelectorTerm{{MatchExpressions: matchExpressions}}

}

func (r GPUSchedulingStrategy) Tolerations() []v1.Toleration {
	var key string
	if r.Raw.GPUNum >= 4 {
		key = "jeeves-graphics-pack/4-8"
	} else {
		key = "jeeves-graphics-pack/1-2"
	}
	toleration := []v1.Toleration{
		{
			Key:    "nvidia.com/gpu",
			Effect: v1.TaintEffectNoSchedule,
		},
		{
			Key:    key,
			Effect: v1.TaintEffectNoSchedule,
		},
	}

	if r.Raw.GPUSeries == GPUSeriesA100 {
		toleration = append(toleration, v1.Toleration{
			Key:      A100TolerationKey,
			Operator: v1.TolerationOpEqual,
			Value:    GPUSeriesA100,
			Effect:   v1.TaintEffectNoSchedule,
		})
	}

	return toleration
}

type VGPUSchedulingStrategy struct {
	Raw       Resource
	CanBorrow bool
}

func (r VGPUSchedulingStrategy) Labels() map[string]string {
	return map[string]string{
		"cce.baidubce.com/baidu-cgpu.affinity": "offline",
	}
}

func (r VGPUSchedulingStrategy) Annotations() map[string]string {
	return map[string]string{
		"scheduling.k8s.io/job-enable-oversell": "false",
		"scheduling.volcano.sh/queue-name":      "default",
	}
}

func (r VGPUSchedulingStrategy) SchedulerName() string {
	return VGPUSchedulerName
}

func (r VGPUSchedulingStrategy) PriorityClassName() string {
	return VGPUPriorityClass
}

func (r VGPUSchedulingStrategy) PreferredSchedulingTerms() []v1.PreferredSchedulingTerm {
	return []v1.PreferredSchedulingTerm{}
}

func (r VGPUSchedulingStrategy) NodeSelectorTerms() []v1.NodeSelectorTerm {
	return []v1.NodeSelectorTerm{
		{
			MatchExpressions: []v1.NodeSelectorRequirement{
				{
					Key:      "node-role.kubernetes.io/jeeves-gpu",
					Operator: v1.NodeSelectorOpExists,
				},
				{
					Key:      "jeeves-graphics-pack/vgpu",
					Operator: v1.NodeSelectorOpExists,
				},
			},
		},
	}

}

func (r VGPUSchedulingStrategy) Tolerations() []v1.Toleration {
	return []v1.Toleration{
		{
			Key:    "nvidia.com/gpu",
			Effect: v1.TaintEffectNoSchedule,
		},
		{
			Key:    "jeeves-graphics-pack/vgpu",
			Effect: v1.TaintEffectNoSchedule,
		},
	}
}

type CPUSchedulingStrategy struct {
	Raw Resource
}

func (r CPUSchedulingStrategy) Labels() map[string]string {
	return make(map[string]string)
}

func (r CPUSchedulingStrategy) Annotations() map[string]string {
	return make(map[string]string)
}

func (r CPUSchedulingStrategy) SchedulerName() string {
	return DefaultSchedulerName
}

func (r CPUSchedulingStrategy) PriorityClassName() string {
	return CPUPriorityClass
}

func (r CPUSchedulingStrategy) PreferredSchedulingTerms() []v1.PreferredSchedulingTerm {
	return []v1.PreferredSchedulingTerm{}
}

func (r CPUSchedulingStrategy) NodeSelectorTerms() []v1.NodeSelectorTerm {
	return []v1.NodeSelectorTerm{
		{
			MatchExpressions: []v1.NodeSelectorRequirement{
				{
					Key:      "node-role.kubernetes.io/jeeves-cpu",
					Operator: v1.NodeSelectorOpExists,
				},
			},
		},
	}

}

func (r CPUSchedulingStrategy) Tolerations() []v1.Toleration {
	return []v1.Toleration{
		{
			Key:    "node-role.kubernetes.io/jeeves-cpu",
			Effect: v1.TaintEffectNoSchedule,
		},
	}
}

var _ ResourceStrategy = HalfMemoryResourceStrategy{}
var _ ResourceStrategy = HETStrategy{}

type HalfMemoryResourceStrategy struct {
	Raw             Resource
	Region          string
	UseGPUScheduler bool
	CanBorrow       bool
}

func (r HalfMemoryResourceStrategy) Requests() Resource {
	resource := r.Resource()
	resource.CPUNum = 1
	resource.MemorySize = (resource.MemorySize + 1) / 2
	resource.EphemeralStorage = (resource.EphemeralStorage + 1) / 2
	return resource
}

func (r HalfMemoryResourceStrategy) Limits() Resource {
	return r.Resource()
}

func (r HalfMemoryResourceStrategy) Resource() Resource {
	return r.Raw
}

func (r HalfMemoryResourceStrategy) Merge(resource Resource) {
	r.Raw.CPUNum = maxUInt(r.Raw.CPUNum, resource.CPUNum)
	r.Raw.MemorySize = maxUInt(r.Raw.MemorySize, resource.MemorySize)
	r.Raw.EphemeralStorage = maxUInt(r.Raw.EphemeralStorage, resource.EphemeralStorage)
	r.Raw.GPUNum = maxUInt(r.Raw.GPUNum, resource.GPUNum)
}

func (r HalfMemoryResourceStrategy) SchedulingStrategy() SchedulingStrategy {
	if r.Raw.GPUSeries == "vGPU" {
		return &VGPUSchedulingStrategy{r.Raw, r.CanBorrow}
	} else if r.Raw.GPUNum > 0 {
		return &GPUSchedulingStrategy{r.Raw, r.CanBorrow}
	}
	return &CPUSchedulingStrategy{r.Raw}
}

func (r HalfMemoryResourceStrategy) GetRegion() string {
	return r.Region
}

// Heterogeneous Architecture Strategy
// Gives guaranteed resource in pure-CPU computation, buf halven the host resource requests when GPU is used.
// when in pure-CPU computation, node selection is done in a tiling manner, while in GPU compuation,
// scheduler will select used nodes with available resources first, leaving bulks of free resource for large
// job like multi-GPU training.
type HETStrategy struct {
	Raw    Resource
	Region string
	// UseGPUScheduler bool
	// CanBorrow       bool
}

func (r HETStrategy) Requests() Resource {
	resource := r.Resource()
	if resource.GPUNum > 0 {
		resource.CPUNum = (resource.CPUNum + 1) / 2
		resource.MemorySize = (resource.MemorySize + 1) / 2
		resource.EphemeralStorage = (resource.EphemeralStorage + 1) / 2
	}
	return resource
}

func (r HETStrategy) Limits() Resource {
	return r.Resource()
}

func (r HETStrategy) Resource() Resource {
	return r.Raw
}

func (r HETStrategy) Merge(resource Resource) {
	r.Raw.CPUNum = maxUInt(r.Raw.CPUNum, resource.CPUNum)
	r.Raw.MemorySize = maxUInt(r.Raw.MemorySize, resource.MemorySize)
	r.Raw.EphemeralStorage = maxUInt(r.Raw.EphemeralStorage, resource.EphemeralStorage)
	r.Raw.GPUNum = maxUInt(r.Raw.GPUNum, resource.GPUNum)
}

func (r HETStrategy) SchedulingStrategy() SchedulingStrategy {
	if r.Raw.GPUSeries == "vGPU" {
		return &VGPUSchedulingStrategy{Raw: r.Raw}
	} else if r.Raw.GPUNum > 0 {
		return &GPUSchedulingStrategy{Raw: r.Raw}
	}
	return &CPUSchedulingStrategy{r.Raw}
}

func (r HETStrategy) GetRegion() string {
	return r.Region
}

func maxUInt(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}
