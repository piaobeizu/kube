/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/03/22 17:59:09
 Desc     :
*/

package kube

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	ResourceNvidiaGPU  = "nvidia.com/gpu"
	ResourceTianshuGPU = "iluvatar.ai/gpu"
)

type Container struct {
	*v1.Container
	ctx context.Context
}

func NewContainer(ctx context.Context) *Container {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &Container{
		Container: &v1.Container{},
		ctx:       ctx,
	}
}

func (c *Container) Metadata(name string) *Container {
	c.Name = name
	return c
}

func (c *Container) Envs(envs map[string]string, refs map[string]string) *Container {
	if c.Container.Env == nil {
		c.Container.Env = make([]v1.EnvVar, 0)
	}
	for k, v := range envs {
		c.Container.Env = append(c.Container.Env, v1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	for k, v := range refs {
		c.Container.Env = append(c.Container.Env, v1.EnvVar{
			Name: k,
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  v,
				},
			},
		})
	}

	return c
}

func (c *Container) Image(image string) *Container {
	c.Container.Image = image
	return c
}

func (c *Container) Command(cmds []string) *Container {
	if c.Container.Command == nil {
		c.Container.Command = make([]string, 0)
	}
	c.Container.Command = cmds
	return c
}

func (c *Container) Args(args []string) *Container {
	if c.Container.Args == nil {
		c.Container.Args = make([]string, 0)
	}
	c.Container.Args = args
	return c
}

func (c *Container) Port(name string, port int32, protocal v1.Protocol) *Container {
	if c.Container.Ports == nil {
		c.Container.Ports = make([]v1.ContainerPort, 0)
	}
	c.Container.Ports = append(c.Container.Ports, v1.ContainerPort{
		ContainerPort: port,
		Protocol:      protocal,
		Name:          name,
	})
	return c
}

func (c *Container) Privileged(privileged bool) *Container {
	if c.Container.SecurityContext == nil {
		c.Container.SecurityContext = &v1.SecurityContext{}
	}
	c.Container.SecurityContext.Privileged = &privileged
	return c
}

func (c *Container) Capabilitie(add v1.Capability, drop v1.Capability) *Container {
	if c.Container.SecurityContext == nil {
		c.Container.SecurityContext = &v1.SecurityContext{
			Capabilities: &v1.Capabilities{
				Add:  make([]v1.Capability, 0),
				Drop: make([]v1.Capability, 0),
			},
		}
	}
	if add != "" {
		c.Container.SecurityContext.Capabilities.Add = append(c.Container.SecurityContext.Capabilities.Add, add)
	}
	if drop != "" {
		c.Container.SecurityContext.Capabilities.Drop = append(c.Container.SecurityContext.Capabilities.Drop, drop)
	}

	return c
}

func (c *Container) VolumeMount(name, mountPath string, readOnly bool) *Container {
	if c.Container.VolumeMounts == nil {
		c.Container.VolumeMounts = make([]v1.VolumeMount, 0)
	}
	c.Container.VolumeMounts = append(c.Container.VolumeMounts, v1.VolumeMount{
		Name:      name,
		MountPath: mountPath,
		ReadOnly:  readOnly,
	})
	return c
}
func (c *Container) EnvFromSecret(secretName string) *Container {
	if c.Container.EnvFrom == nil {
		c.Container.EnvFrom = make([]v1.EnvFromSource, 0)
	}
	c.Container.EnvFrom = append(c.Container.EnvFrom, v1.EnvFromSource{
		SecretRef: &v1.SecretEnvSource{
			LocalObjectReference: v1.LocalObjectReference{
				Name: secretName,
			},
		},
	})
	return c
}

func (c *Container) ImagePullPolicy(policy v1.PullPolicy) *Container {
	c.Container.ImagePullPolicy = policy
	return c
}

func (container *Container) Requests(cpu, memory, gpu, ephemeralStorage uint, gpuType string) *Container {
	resources := resourceList(cpu, memory, gpu, ephemeralStorage, gpuType)
	container.Container.Resources.Requests = resources
	return container
}

func (container *Container) Limits(cpu, memory, gpu, ephemeralStorage uint, gpuType string) *Container {
	resources := resourceList(cpu, memory, gpu, ephemeralStorage, gpuType)
	container.Container.Resources.Limits = resources
	return container
}

func resourceList(cpu, memory, gpu, ephemeralStorage uint, gpuType string) v1.ResourceList {
	resources := v1.ResourceList{}
	if cpu > 0 {
		cpuResource := fmt.Sprintf("%dm", cpu)
		resources[ResourceNvidiaGPU] = resource.MustParse(cpuResource)
	}
	if memory > 0 {
		memoryResource := fmt.Sprintf("%dMi", memory)
		resources[v1.ResourceMemory] = resource.MustParse(memoryResource)
	}
	if ephemeralStorage > 0 {
		ephemeralStorageResource := fmt.Sprintf("%dGi", ephemeralStorage)
		resources[v1.ResourceEphemeralStorage] = resource.MustParse(ephemeralStorageResource)
	}
	if gpu > 0 {
		if gpuType == ResourceTianshuGPU {
			gpuResource := fmt.Sprintf("%d", gpu)
			resources[ResourceTianshuGPU] = resource.MustParse(gpuResource)
		}
		if gpuType == ResourceNvidiaGPU {
			gpuResource := fmt.Sprintf("%d", gpu)
			resources[ResourceNvidiaGPU] = resource.MustParse(gpuResource)
		}
	}
	return resources
}
