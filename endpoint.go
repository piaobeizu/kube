/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/03/26 15:46:30
 Desc     :
*/

package kube

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Endpoint struct {
	*LinkInfo
	*v1.Endpoints
	client *KubeClient
	ctx    context.Context
}

func NewEndpoint(ctx context.Context) *Endpoint {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &Endpoint{
		Endpoints: &v1.Endpoints{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Endpoints",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{},
			Subsets:    nil,
		},
		LinkInfo: &LinkInfo{},
		client:   nil,
		ctx:      ctx,
	}
}

func (e *Endpoint) Link(region, config string) *Endpoint {
	e.Region = region
	e.Config = config
	e.client = NewKubeClient(region, config)
	return e
}

func (e *Endpoint) Metadata(name, namespace string) *Endpoint {
	e.Name, e.Namespace = name, namespace
	return e
}

func (e *Endpoint) Subset(subsets []string) *Endpoint {
	if e.Subsets == nil {
		e.Endpoints.Subsets = make([]v1.EndpointSubset, 0)
	} else {
		e.Endpoints.Subsets = append(e.Endpoints.Subsets, v1.EndpointSubset{})
	}
	return e
}

func (e *Endpoint) Address(ip, nodeName, hostName string) *Endpoint {
	if e.Subsets == nil {
		e.Endpoints.Subsets = make([]v1.EndpointSubset, 0)
	}
	e.Endpoints.Subsets[len(e.Endpoints.Subsets)-1].Addresses = append(
		e.Endpoints.Subsets[len(e.Endpoints.Subsets)-1].Addresses, v1.EndpointAddress{
			IP:       ip,
			NodeName: &nodeName,
			Hostname: hostName,
		})
	return e
}

func (e *Endpoint) Port(name string, port int32, protocol v1.Protocol) *Endpoint {
	if e.Subsets == nil {
		e.Endpoints.Subsets = make([]v1.EndpointSubset, 0)
	}
	e.Endpoints.Subsets[len(e.Endpoints.Subsets)-1].Ports = append(
		e.Endpoints.Subsets[len(e.Endpoints.Subsets)-1].Ports, v1.EndpointPort{
			Name:     name,
			Port:     port,
			Protocol: protocol,
		})
	return e
}

func (e *Endpoint) Get() (*v1.Endpoints, error) {
	return e.client.Clientset.CoreV1().Endpoints(e.Namespace).Get(e.ctx, e.Name, metav1.GetOptions{})
}

func (e *Endpoint) Create() error {
	_, err := e.client.Clientset.CoreV1().Endpoints(e.Namespace).Create(e.ctx, e.Endpoints, metav1.CreateOptions{})
	return err
}

func (e *Endpoint) Delete() error {
	return e.client.Clientset.CoreV1().Endpoints(e.Namespace).Delete(e.ctx, e.Name, metav1.DeleteOptions{})
}

func (e *Endpoint) Update() error {
	_, err := e.client.Clientset.CoreV1().Endpoints(e.Namespace).Update(e.ctx, e.Endpoints, metav1.UpdateOptions{})
	return err
}

func (e *Endpoint) CreateOrUpdate() error {
	_, err := e.client.Clientset.CoreV1().Endpoints(e.Namespace).Get(e.ctx, e.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return e.Create()
		}
		return err
	}
	return e.Update()
}

func (e *Endpoint) Addresses() ([]v1.EndpointAddress, error) {
	endPoints, err := e.Get()
	if err != nil {
		return nil, err
	}
	addresses := make([]v1.EndpointAddress, 0)
	for _, subset := range endPoints.Subsets {
		addresses = append(addresses, subset.Addresses...)
	}
	return addresses, nil
}

func (e *Endpoint) Ports() ([]v1.EndpointPort, error) {
	endPoints, err := e.Get()
	if err != nil {
		return nil, err
	}
	ports := make([]v1.EndpointPort, 0)
	for _, subset := range endPoints.Subsets {
		ports = append(ports, subset.Ports...)
	}
	return ports, nil
}
