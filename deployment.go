/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/03/22 19:48:22
 Desc     : deployment
*/

package kube

/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/03/22 22:53:57
 Desc     :
*/

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Deployment struct {
	ctx context.Context
	*LinkInfo
	*appsv1.Deployment
	client *KubeClient
}

func NewDeployment(ctx context.Context) *Deployment {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &Deployment{
		Deployment: &appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Deployment",
				APIVersion: "apps/v1",
			},
			ObjectMeta: metav1.ObjectMeta{},
			Spec:       appsv1.DeploymentSpec{},
		},
		LinkInfo: &LinkInfo{},
		client:   nil,
		ctx:      ctx,
	}
}

func (d *Deployment) Link(region, config string) *Deployment {
	d.Region = region
	d.Config = config
	d.client = NewKubeClient(region, config)
	return d
}

func (d *Deployment) Metadata(name, namespace string) *Deployment {
	d.Name, d.Namespace = name, namespace
	return d
}

func (d *Deployment) Labels(labels map[string]string) *Deployment {
	if d.Deployment.Labels == nil {
		d.Deployment.Labels = make(map[string]string)
	}
	for k, v := range labels {
		d.Deployment.Labels[k] = v
	}
	return d
}

func (d *Deployment) Replicas(replicas int32) *Deployment {
	if d.Deployment.Spec.Replicas == nil {
		d.Deployment.Spec.Replicas = new(int32)
	}
	d.Deployment.Spec.Replicas = &replicas
	return d
}

func (d *Deployment) Selector(selector map[string]string) *Deployment {
	if d.Deployment.Spec.Selector == nil {
		d.Deployment.Spec.Selector = &metav1.LabelSelector{}
	}
	d.Deployment.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: selector,
	}
	return d
}

func (d *Deployment) Template(pod *PodTemplate) *Deployment {
	if pod == nil {
		return d
	}
	d.Deployment.Spec.Template = pod.Template
	return d
}

func (d *Deployment) Strategy(strategy appsv1.DeploymentStrategy) *Deployment {
	d.Deployment.Spec.Strategy = strategy
	return d
}

func (d *Deployment) Create() error {
	_, err := d.client.AppsV1().Deployments(d.Namespace).
		Create(d.ctx, d.Deployment, metav1.CreateOptions{})
	return err
}

func (d *Deployment) Delete() error {
	return d.client.AppsV1().Deployments(d.Namespace).
		Delete(d.ctx, d.Name, metav1.DeleteOptions{})
}

func (d *Deployment) Update() error {
	_, err := d.client.AppsV1().Deployments(d.Namespace).Update(d.ctx, d.Deployment, metav1.UpdateOptions{})
	return err
}

func (d *Deployment) Get() (*appsv1.Deployment, error) {
	return d.client.AppsV1().Deployments(d.Namespace).Get(d.ctx, d.Name, metav1.GetOptions{})
}

func (d *Deployment) Empty() bool {
	_, err := d.client.AppsV1().Deployments(d.Namespace).
		Get(d.ctx, d.Name, metav1.GetOptions{})
	return errors.IsNotFound(err)
}

func (d *Deployment) CreateOrUpdate() error {
	_, err := d.client.AppsV1().Deployments(d.Namespace).Get(d.ctx, d.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return d.Create()
		}
		return err
	}
	return d.Update()
}

// implement the rollout method
func (d *Deployment) Rollout() error {
	if d.Deployment.Annotations == nil {
		d.Deployment.Annotations = make(map[string]string)
	}
	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().String())
	_, err := d.client.AppsV1().Deployments(d.Namespace).Patch(d.ctx, d.Deployment.Name,
		types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
	return err
}
