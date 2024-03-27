/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/03/23 10:39:16
 Desc     :
*/

package kube

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Service struct {
	ctx context.Context
	*LinkInfo
	*v1.Service
	client *KubeClient
}

func NewService(ctx context.Context) *Service {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &Service{
		ctx: ctx,
		Service: &v1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{},
			Spec:       v1.ServiceSpec{},
		},
		LinkInfo: &LinkInfo{},
		client:   nil,
	}
}

func (s *Service) Link(region, config string) *Service {
	s.Region = region
	s.Config = config
	s.client = NewKubeClient(region, config)
	return s
}

func (s *Service) Metadata(name, namespace string) *Service {
	s.Name = name
	s.Namespace = namespace
	return s
}

func (s *Service) Port(name string, port, NodePort int32, protocal v1.Protocol) *Service {
	if s.Spec.Ports == nil {
		s.Spec.Ports = make([]v1.ServicePort, 0)
	}
	if port <= 0 {
		return s
	}
	s.Spec.Ports = append(s.Spec.Ports, v1.ServicePort{
		Port:       port,
		Name:       name,
		Protocol:   protocal,
		TargetPort: intstr.FromInt(int(port)),
	})
	if NodePort > 0 {
		s.Spec.Ports[len(s.Spec.Ports)-1].NodePort = NodePort
		s.Spec.Type = v1.ServiceTypeNodePort
	}
	return s
}

func (s *Service) Labels(labels map[string]string) *Service {
	if s.Service.Labels == nil {
		s.Service.Labels = make(map[string]string)
	}
	for k, v := range labels {
		s.Service.Labels[k] = v
	}
	return s
}

func (s *Service) Selector(selector map[string]string) *Service {
	if s.Spec.Selector == nil {
		s.Spec.Selector = make(map[string]string)
	}
	for k, v := range selector {
		s.Spec.Selector[k] = v
	}
	return s
}

func (s *Service) Type(serviceType v1.ServiceType) *Service {
	s.Spec.Type = serviceType
	return s
}

func (s *Service) Create() error {
	_, err := s.client.CoreV1().Services(s.Namespace).
		Create(s.ctx, s.Service, metav1.CreateOptions{})
	return err
}

func (s *Service) Update() error {
	_, err := s.client.CoreV1().Services(s.Namespace).
		Update(s.ctx, s.Service, metav1.UpdateOptions{})
	return err
}

func (s *Service) Delete() error {
	return s.client.CoreV1().Services(s.Namespace).
		Delete(s.ctx, s.Name, metav1.DeleteOptions{})
}

func (s *Service) Get(replace bool) (*v1.Service, error) {
	return s.client.CoreV1().Services(s.Namespace).Get(s.ctx, s.Name, metav1.GetOptions{})
}

func (s *Service) Empty() bool {
	_, err := s.client.CoreV1().Services(s.Namespace).Get(s.ctx, s.Name, metav1.GetOptions{})
	return errors.IsNotFound(err)
}

func (s *Service) CreateOrUpdate() error {
	_, err := s.client.CoreV1().Services(s.Namespace).Get(s.ctx, s.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return s.Create()
		}
		return err
	}
	return s.Update()
}
