/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/04/01 13:35:25
 Desc     :
*/

package kube

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceAccount struct {
	*LinkInfo
	*v1.ServiceAccount
	ctx    context.Context
	client *KubeClient
}

func NewServiceAccount(ctx context.Context) *ServiceAccount {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &ServiceAccount{
		ServiceAccount: &v1.ServiceAccount{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ServiceAccount",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{},
			Secrets:    []v1.ObjectReference{},
		},
		ctx:      ctx,
		LinkInfo: &LinkInfo{},
		client:   nil,
	}
}

func (sa *ServiceAccount) Metadata(name, namespace string) *ServiceAccount {
	sa.Name = name
	sa.Namespace = namespace
	return sa
}

func (sa *ServiceAccount) Labels(labels map[string]string) *ServiceAccount {
	if sa.ServiceAccount.Labels == nil {
		sa.ServiceAccount.Labels = make(map[string]string)
	}
	sa.ServiceAccount.Labels = labels
	return sa
}

func (sa *ServiceAccount) Annotations(annotations map[string]string) *ServiceAccount {
	if sa.ServiceAccount.Annotations == nil {
		sa.ServiceAccount.Annotations = make(map[string]string)
	}
	sa.ServiceAccount.Annotations = annotations
	return sa
}

func (sa *ServiceAccount) Link(region, config string) *ServiceAccount {
	sa.Region = region
	sa.Config = config
	sa.client = NewKubeClient(region, config)
	return sa
}

func (sa *ServiceAccount) AutomountServiceAccountToken(automount bool) *ServiceAccount {
	sa.ServiceAccount.AutomountServiceAccountToken = &automount
	return sa
}

func (sa *ServiceAccount) Secret(kind, secretName, namespace string) *ServiceAccount {
	if sa.Secrets == nil {
		sa.Secrets = make([]v1.ObjectReference, 0)
	}
	sa.Secrets = append(sa.Secrets, v1.ObjectReference{
		Kind:      kind,
		Name:      secretName,
		Namespace: namespace,
	})
	return sa
}

func (sa *ServiceAccount) Create() error {
	_, err := sa.client.CoreV1().ServiceAccounts(sa.Namespace).Create(sa.ctx, sa.ServiceAccount, metav1.CreateOptions{})
	return err
}

func (sa *ServiceAccount) Get() (*v1.ServiceAccount, error) {
	return sa.client.CoreV1().ServiceAccounts(sa.Namespace).Get(sa.ctx, sa.Name, metav1.GetOptions{})
}

func (sa *ServiceAccount) Update() error {
	_, err := sa.client.CoreV1().ServiceAccounts(sa.Namespace).Update(sa.ctx, sa.ServiceAccount, metav1.UpdateOptions{})
	return err
}

func (sa *ServiceAccount) Delete() error {
	return sa.client.CoreV1().ServiceAccounts(sa.Namespace).Delete(sa.ctx, sa.Name, metav1.DeleteOptions{})
}

func (sa *ServiceAccount) List() (*v1.ServiceAccountList, error) {
	return sa.client.CoreV1().ServiceAccounts(sa.Namespace).List(sa.ctx, metav1.ListOptions{})
}

func (sa *ServiceAccount) Empty() bool {
	_, err := sa.client.CoreV1().ServiceAccounts(sa.Namespace).Get(sa.ctx, sa.Name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		panic(err)
	}
	return errors.IsNotFound(err)
}

func (sa *ServiceAccount) CreateOrUpdate() error {
	_, err := sa.client.CoreV1().ServiceAccounts(sa.Namespace).Get(sa.ctx, sa.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return sa.Create()
		}
		return err
	}
	return sa.Update()
}

func (sa *ServiceAccount) Equal(keys []string) bool {
	serviceAccount, err := sa.Get()
	if err != nil && !errors.IsNotFound(err) {
		panic(err)
	}
	if len(keys) == 0 {
		keys = []string{}
	}
	return ResourceEqual(sa.ServiceAccount, serviceAccount, keys)
}
