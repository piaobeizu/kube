/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/03/22 19:20:11
 Desc     :
*/

package kube

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DaemonSet struct {
	*LinkInfo
	*appsv1.DaemonSet
	client *KubeClient
	ctx    context.Context
}

func NewDaemonSet(ctx context.Context) *DaemonSet {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &DaemonSet{
		DaemonSet: &appsv1.DaemonSet{
			TypeMeta: metav1.TypeMeta{
				Kind:       "DaemonSet",
				APIVersion: "apps/v1",
			},
			ObjectMeta: metav1.ObjectMeta{},
			Spec:       appsv1.DaemonSetSpec{},
		},
		LinkInfo: &LinkInfo{},
		client:   nil,
		ctx:      ctx,
	}
}

func (d *DaemonSet) Link(region, config string) *DaemonSet {
	d.Region = region
	d.Config = config
	d.client = NewKubeClient(region, config)
	return d
}

func (d *DaemonSet) Metadata(name, namespace string) *DaemonSet {
	d.Name, d.Namespace = name, namespace
	return d
}

func (d *DaemonSet) Labels(labels map[string]string) *DaemonSet {
	if d.DaemonSet.Labels == nil {
		d.DaemonSet.Labels = make(map[string]string)
	}
	for k, v := range labels {
		d.DaemonSet.Labels[k] = v
	}
	return d
}

func (d *DaemonSet) Annotations(annotations map[string]string) *DaemonSet {
	if d.DaemonSet.Annotations == nil {
		d.DaemonSet.Annotations = make(map[string]string)
	}
	for k, v := range annotations {
		d.DaemonSet.Annotations[k] = v
	}
	return d
}

func (d *DaemonSet) Selector(labels map[string]string) *DaemonSet {
	if d.DaemonSet.Spec.Selector == nil {
		d.DaemonSet.Spec.Selector = &metav1.LabelSelector{}
	}
	d.DaemonSet.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: labels,
	}
	return d
}

func (d *DaemonSet) Template(pt *PodTemplate) *DaemonSet {
	if pt == nil {
		return d
	}
	d.DaemonSet.Spec.Template = pt.Template
	return d
}

func (d *DaemonSet) Create() error {
	daemonsets := d.client.AppsV1().DaemonSets(d.Namespace)
	_, err := daemonsets.Create(d.ctx, d.DaemonSet, metav1.CreateOptions{})
	return err
}

func (d *DaemonSet) Delete() error {
	daemonsets := d.client.AppsV1().DaemonSets(d.Namespace)
	return daemonsets.Delete(d.ctx, d.Name, metav1.DeleteOptions{})
}

func (d *DaemonSet) Update() error {
	daemonsets := d.client.AppsV1().DaemonSets(d.Namespace)
	_, err := daemonsets.Update(d.ctx, d.DaemonSet, metav1.UpdateOptions{})
	return err
}

func (d *DaemonSet) Get() (*v1.DaemonSet, error) {
	return d.client.AppsV1().DaemonSets(d.Namespace).
		Get(d.ctx, d.Name, metav1.GetOptions{})
}

func (d *DaemonSet) CreateOrUpdate() error {
	_, err := d.client.AppsV1().DaemonSets(d.Namespace).Get(d.ctx, d.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return d.Create()
		}
		return err
	}
	return d.Update()
}

func (d *DaemonSet) Empty() bool {
	_, err := d.client.AppsV1().DaemonSets(d.Namespace).Get(d.ctx, d.Name, metav1.GetOptions{})
	return errors.IsNotFound(err)
}

// implement the rollout method
func (d *DaemonSet) Rollout() error {
	if d.DaemonSet.Annotations == nil {
		d.DaemonSet.Annotations = make(map[string]string)
	}
	d.DaemonSet.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Local().Format(time.RFC3339)
	_, err := d.client.AppsV1().DaemonSets(d.Namespace).Update(d.ctx, d.DaemonSet, metav1.UpdateOptions{})
	return err
}
