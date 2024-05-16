package kube

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodsGetter interface {
	Pods() PodInterface
}
type PodInterface interface {
	Create() error
	Delete() error
	Update() error
	Get() (*Pod, error)
	CreateOrUpdate() error
}

type Pod struct {
	*LinkInfo
	*v1.Pod
	client *KubeClient
	ctx    context.Context
}

func NewPod(ctx context.Context) *Pod {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &Pod{
		Pod: &v1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{},
			Spec:       v1.PodSpec{},
		},
		LinkInfo: &LinkInfo{},
		client:   nil,
		ctx:      ctx,
	}
}

func (p *Pod) Link(region, config string) *Pod {
	p.Region = region
	p.Config = config
	p.client = NewKubeClient(region, config)
	return p
}

func (p *Pod) Metadata(name, namespace string) *Pod {
	p.Name, p.Namespace = name, namespace
	return p
}

func (p *Pod) Labels(labels map[string]string) *Pod {
	if p.Pod.Labels == nil {
		p.Pod.Labels = make(map[string]string)
	}
	for k, v := range labels {
		p.Pod.Labels[k] = v
	}
	return p
}

func (p *Pod) Annotations(annotations map[string]string) *Pod {
	if p.Pod.Annotations == nil {
		p.Pod.Annotations = make(map[string]string)
	}
	for k, v := range annotations {
		p.Pod.Annotations[k] = v
	}
	return p
}

func (p *Pod) Container(container v1.Container) *Pod {
	if p.Pod.Spec.Containers == nil {
		p.Pod.Spec.Containers = make([]v1.Container, 0)
	}
	p.Pod.Spec.Containers = append(p.Pod.Spec.Containers, container)
	return p
}

func (p *Pod) Volume(volumes v1.Volume) *Pod {
	if p.Pod.Spec.Volumes == nil {
		p.Pod.Spec.Volumes = make([]v1.Volume, 0)
	}
	p.Pod.Spec.Volumes = append(p.Pod.Spec.Volumes, volumes)
	return p
}

func (p *Pod) RestartPolicy(policy v1.RestartPolicy) *Pod {
	p.Pod.Spec.RestartPolicy = policy
	return p
}

func (p *Pod) Toleration(toleration []v1.Toleration) *Pod {
	p.Pod.Spec.Tolerations = toleration
	return p
}

func (p *Pod) ImagePullSecrets(secrets []string) *Pod {
	for _, secret := range secrets {
		p.Pod.Spec.ImagePullSecrets = append(p.Pod.Spec.ImagePullSecrets, v1.LocalObjectReference{Name: secret})
	}
	return p
}

func (p *Pod) Create() error {
	pods := p.client.CoreV1().Pods(p.Namespace)
	_, err := pods.Create(p.ctx, p.Pod, metav1.CreateOptions{})
	return err
}

func (p Pod) Delete() error {
	pods := p.client.CoreV1().Pods(p.Namespace)
	return pods.Delete(p.ctx, p.Pod.Name, metav1.DeleteOptions{})
}

func (p *Pod) Update() error {
	pods := p.client.CoreV1().Pods(p.Namespace)
	_, err := pods.Update(p.ctx, p.Pod, metav1.UpdateOptions{})
	return err
}

func (p *Pod) Get(replace bool) (*Pod, error) {
	pods := p.client.CoreV1().Pods(p.Namespace)
	pod, err := pods.Get(p.ctx, p.Pod.Name, metav1.GetOptions{})
	if replace {
		p.Pod = pod
	}
	return p, err

}

func (p *Pod) CreateOrUpdate() error {
	pod, err := p.Get(false)
	if err != nil {
		return err
	}
	if pod.Pod == nil {
		return p.Create()
	}
	return p.Update()
}

type PodTemplate struct {
	*v1.PodTemplate
	ctx context.Context
}

func NewPodTemplate(ctx context.Context) *PodTemplate {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &PodTemplate{
		PodTemplate: &v1.PodTemplate{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Template:   v1.PodTemplateSpec{},
		},
		ctx: ctx,
	}
}

func (pt *PodTemplate) Metadata(name string) *PodTemplate {
	pt.PodTemplate.Name = name
	return pt
}

func (pt *PodTemplate) Labels(labels map[string]string) *PodTemplate {
	if pt.PodTemplate.Template.Labels == nil {
		pt.PodTemplate.Template.Labels = make(map[string]string)
	}
	for k, v := range labels {
		pt.PodTemplate.Template.Labels[k] = v
	}
	return pt
}

func (pt *PodTemplate) Annotations(annotations map[string]string) *PodTemplate {
	if pt.PodTemplate.Annotations == nil {
		pt.PodTemplate.Annotations = make(map[string]string)
	}
	for k, v := range annotations {
		pt.PodTemplate.Annotations[k] = v
	}
	return pt
}

func (pt *PodTemplate) Container(container Container) *PodTemplate {
	if pt.PodTemplate.Template.Spec.Containers == nil {
		pt.PodTemplate.Template.Spec.Containers = make([]v1.Container, 0)
	}
	pt.PodTemplate.Template.Spec.Containers = append(pt.PodTemplate.Template.Spec.Containers, *container.Container)
	return pt
}

func (pt *PodTemplate) Volume(volume Volume) *PodTemplate {
	if pt.PodTemplate.Template.Spec.Volumes == nil {
		pt.PodTemplate.Template.Spec.Volumes = make([]v1.Volume, 0)
	}
	pt.PodTemplate.Template.Spec.Volumes = append(pt.PodTemplate.Template.Spec.Volumes, *volume.Volume)
	return pt
}

func (pt *PodTemplate) RestartPolicy(policy v1.RestartPolicy) *PodTemplate {
	pt.PodTemplate.Template.Spec.RestartPolicy = policy
	return pt
}

func (pt *PodTemplate) SecurityContext(user, group, fsGroup int64) *PodTemplate {
	if pt.PodTemplate.Template.Spec.SecurityContext == nil {
		pt.PodTemplate.Template.Spec.SecurityContext = &v1.PodSecurityContext{}
	}
	pt.PodTemplate.Template.Spec.SecurityContext.RunAsUser = &user
	pt.PodTemplate.Template.Spec.SecurityContext.RunAsGroup = &group
	pt.PodTemplate.Template.Spec.SecurityContext.FSGroup = &fsGroup
	return pt
}

func (pt *PodTemplate) Toleration(key string, operator v1.TolerationOperator, value string, effect v1.TaintEffect, tolerationSeconds int64) *PodTemplate {
	if pt.PodTemplate.Template.Spec.Tolerations == nil {
		pt.PodTemplate.Template.Spec.Tolerations = make([]v1.Toleration, 0)
	}
	pt.PodTemplate.Template.Spec.Tolerations = append(pt.PodTemplate.Template.Spec.Tolerations, v1.Toleration{
		Key:               key,
		Operator:          operator,
		Value:             value,
		Effect:            effect,
		TolerationSeconds: &tolerationSeconds,
	})
	return pt
}

func (pt *PodTemplate) ImagePullSecrets(secrets []string) *PodTemplate {
	for _, secret := range secrets {
		pt.PodTemplate.Template.Spec.ImagePullSecrets = append(pt.PodTemplate.Template.Spec.ImagePullSecrets, v1.LocalObjectReference{Name: secret})
	}
	return pt
}

func (pt *PodTemplate) ServiceAccount(serviceAccount string) *PodTemplate {
	pt.Template.Spec.ServiceAccountName = serviceAccount
	return pt
}

func (pt *PodTemplate) PodSecurityContext(securityContext v1.PodSecurityContext) *PodTemplate {
	pt.Template.Spec.SecurityContext = &securityContext
	return pt
}

func (pt *PodTemplate) TerminationGracePeriodSeconds(seconds int64) *PodTemplate {
	pt.Template.Spec.TerminationGracePeriodSeconds = &seconds
	return pt
}

func (pt *PodTemplate) AutomountServiceAccountToken(auto bool) *PodTemplate {
	pt.Template.Spec.AutomountServiceAccountToken = &auto
	return pt
}

func (pt *PodTemplate) NodeSelector(selecotrs map[string]string) *PodTemplate {
	pt.Template.Spec.NodeSelector = selecotrs
	return pt
}

func (pt *PodTemplate) RequiredDuringSchedulingIgnoredDuringExecution(selector NodeSelector) *PodTemplate {
	if pt.Template.Spec.Affinity == nil && pt.Template.Spec.Affinity.NodeAffinity == nil {
		pt.Template.Spec.Affinity = &v1.Affinity{
			NodeAffinity: &v1.NodeAffinity{},
		}
	}
	pt.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = selector.NodeSelector
	return pt
}

func (pt *PodTemplate) PreferredDuringSchedulingIgnoredDuringExecution(weight int32, selectorTerm NodeSelectorTerm) *PodTemplate {
	if pt.Template.Spec.Affinity == nil && pt.Template.Spec.Affinity.NodeAffinity == nil {
		pt.Template.Spec.Affinity = &v1.Affinity{
			NodeAffinity: &v1.NodeAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []v1.PreferredSchedulingTerm{},
			},
		}
	}
	pt.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(pt.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution, v1.PreferredSchedulingTerm{
		Weight:     weight,
		Preference: selectorTerm.NodeSelectorTerm,
	})
	return pt
}

func (pt *PodTemplate) PodAffinity(affinity v1.PodAffinity) *PodTemplate {
	if pt.Template.Spec.Affinity == nil {
		pt.Template.Spec.Affinity = &v1.Affinity{}
	}
	pt.Template.Spec.Affinity.PodAffinity = &affinity
	return pt
}

func (pt *PodTemplate) PodAntiAffinity(affinity v1.PodAntiAffinity) *PodTemplate {
	if pt.Template.Spec.Affinity == nil {
		pt.Template.Spec.Affinity = &v1.Affinity{}
	}
	pt.Template.Spec.Affinity.PodAntiAffinity = &affinity
	return pt
}
