/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/03/22 18:09:53
 Desc     :
*/

package kube

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type Volume struct {
	*v1.Volume
	ctx context.Context
}

func NewVolume(ctx context.Context) *Volume {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &Volume{
		Volume: &v1.Volume{},
		ctx:    ctx,
	}
}

func (v *Volume) Metadata(name string) *Volume {
	v.Name = name
	return v
}

func (v *Volume) Secret(secretName string, items []v1.KeyToPath, mode int32) *Volume {
	v.VolumeSource = v1.VolumeSource{
		Secret: &v1.SecretVolumeSource{
			SecretName:  secretName,
			Items:       items,
			DefaultMode: &mode,
		},
	}
	return v
}

func (v *Volume) ConfigMap(configMapName string, items []v1.KeyToPath, mode int32) *Volume {
	v.VolumeSource = v1.VolumeSource{
		ConfigMap: &v1.ConfigMapVolumeSource{
			LocalObjectReference: v1.LocalObjectReference{
				Name: configMapName,
			},
			Items:       items,
			DefaultMode: &mode,
		},
	}
	return v
}

func (v *Volume) EmptyDir(medium string, size *resource.Quantity) *Volume {
	v.VolumeSource = v1.VolumeSource{
		EmptyDir: &v1.EmptyDirVolumeSource{
			Medium:    v1.StorageMedium(medium),
			SizeLimit: size,
		},
	}
	return v
}

func (v *Volume) HostPath(path string, t v1.HostPathType) *Volume {
	v.VolumeSource = v1.VolumeSource{
		HostPath: &v1.HostPathVolumeSource{
			Path: path,
			Type: &t,
		},
	}
	return v
}
