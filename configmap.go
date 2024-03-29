/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/03/22 23:08:47
 Desc     :
*/

package kube

import (
	"context"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMap struct {
	*LinkInfo
	*v1.ConfigMap
	client *KubeClient
	ctx    context.Context
}

func NewConfigMap(ctx context.Context) *ConfigMap {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &ConfigMap{
		ConfigMap: &v1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{},
			Data:       nil,
		},
		LinkInfo: &LinkInfo{},
		client:   nil,
		ctx:      ctx,
	}
}

func (c *ConfigMap) Link(region, config string) *ConfigMap {
	c.Region = region
	c.Config = config
	c.client = NewKubeClient(region, config)
	return c
}

func (c *ConfigMap) Metadata(name, namespace string) *ConfigMap {
	c.Name, c.Namespace = name, namespace
	return c
}

func (c *ConfigMap) Data(data map[string]string) *ConfigMap {
	if c.ConfigMap.Data == nil {
		c.ConfigMap.Data = make(map[string]string)
	}
	for k, v := range data {
		c.ConfigMap.Data[k] = v
	}
	return c
}

func (c *ConfigMap) Labels(labels map[string]string) *ConfigMap {
	if c.ConfigMap.Labels == nil {
		c.ConfigMap.Labels = make(map[string]string)
	}
	for k, v := range labels {
		c.ConfigMap.Labels[k] = v
	}
	return c
}

func (c *ConfigMap) Annotations(annotations map[string]string) *ConfigMap {
	if c.ConfigMap.Annotations == nil {
		c.ConfigMap.Annotations = make(map[string]string)
	}
	for k, v := range annotations {
		c.ConfigMap.Annotations[k] = v
	}
	return c
}

func (c *ConfigMap) Create() error {
	configmaps := c.client.CoreV1().ConfigMaps(c.Namespace)
	_, err := configmaps.Create(c.ctx, c.ConfigMap, metav1.CreateOptions{})
	return err
}

func (c *ConfigMap) Delete() error {
	configmaps := c.client.CoreV1().ConfigMaps(c.Namespace)
	return configmaps.Delete(c.ctx, c.Name, metav1.DeleteOptions{})
}

func (c *ConfigMap) Update() error {
	configmaps := c.client.CoreV1().ConfigMaps(c.Namespace)
	_, err := configmaps.Update(c.ctx, c.ConfigMap, metav1.UpdateOptions{})
	return err
}

func (c *ConfigMap) Get() (*v1.ConfigMap, error) {
	return c.client.CoreV1().ConfigMaps(c.Namespace).Get(c.ctx, c.Name, metav1.GetOptions{})
}

func (c *ConfigMap) CreateOrUpdate() error {
	_, err := c.client.CoreV1().ConfigMaps(c.Namespace).Get(c.ctx, c.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return c.Create()
		}
		return err
	}
	return c.Update()
}

func (c *ConfigMap) Empty() bool {
	_, err := c.Get()
	return errors.IsNotFound(err)
}

func (c *ConfigMap) DataEqual() bool {
	cm, err := c.Get()
	if err != nil && !errors.IsNotFound(err) {
		panic(err)
	}
	if len(cm.Data) == 0 && len(c.ConfigMap.Data) == 0 {
		return true
	}
	return reflect.DeepEqual(cm.Data, c.ConfigMap.Data)

}
