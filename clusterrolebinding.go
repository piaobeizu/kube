/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/04/01 16:17:47
 Desc     :
*/

package kube

import (
	"context"
	"reflect"

	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterRoleBinding struct {
	*LinkInfo
	*v1.ClusterRoleBinding
	ctx    context.Context
	client *KubeClient
}

func NewClusterRoleBinding(ctx context.Context) *ClusterRoleBinding {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &ClusterRoleBinding{
		ClusterRoleBinding: &v1.ClusterRoleBinding{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ClusterRoleBinding",
				APIVersion: "rbac.authorization.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{},
			Subjects:   []v1.Subject{},
			RoleRef:    v1.RoleRef{},
		},
		ctx:      ctx,
		LinkInfo: &LinkInfo{},
		client:   nil,
	}
}

func (crb *ClusterRoleBinding) Metadata(name, namespace string) *ClusterRoleBinding {
	crb.Name = name
	crb.Namespace = namespace
	return crb
}

func (crb *ClusterRoleBinding) Link(region, config string) *ClusterRoleBinding {
	crb.Region = region
	crb.Config = config
	crb.client = NewKubeClient(region, config)
	return crb
}
func (crb *ClusterRoleBinding) Labels(labels map[string]string) *ClusterRoleBinding {
	if crb.ClusterRoleBinding.Labels == nil {
		crb.ClusterRoleBinding.Labels = make(map[string]string)
	}
	crb.ClusterRoleBinding.Labels = labels
	return crb
}

func (crb *ClusterRoleBinding) Annotations(annotations map[string]string) *ClusterRoleBinding {
	if crb.ClusterRoleBinding.Annotations == nil {
		crb.ClusterRoleBinding.Annotations = make(map[string]string)
	}
	crb.ClusterRoleBinding.Annotations = annotations
	return crb
}

func (crb *ClusterRoleBinding) Subject(kind, apigroup, name, namespace string) *ClusterRoleBinding {
	if crb.Subjects == nil {
		crb.Subjects = make([]v1.Subject, 0)
	}
	crb.Subjects = append(crb.Subjects, v1.Subject{
		Kind:      kind,
		APIGroup:  apigroup,
		Name:      name,
		Namespace: namespace,
	})
	return crb
}

func (crb *ClusterRoleBinding) RoleRef(kind, apigroup, name string) *ClusterRoleBinding {
	crb.ClusterRoleBinding.RoleRef.Kind = kind
	crb.ClusterRoleBinding.RoleRef.APIGroup = apigroup
	crb.ClusterRoleBinding.RoleRef.Name = name
	return crb
}

func (crb *ClusterRoleBinding) Create() error {
	_, err := crb.client.Clientset.RbacV1().ClusterRoleBindings().Create(crb.ctx, crb.ClusterRoleBinding, metav1.CreateOptions{})
	return err
}

func (crb *ClusterRoleBinding) Get() (*v1.ClusterRoleBinding, error) {
	return crb.client.Clientset.RbacV1().ClusterRoleBindings().Get(crb.ctx, crb.Name, metav1.GetOptions{})
}

func (crb *ClusterRoleBinding) Delete() error {
	return crb.client.Clientset.RbacV1().ClusterRoleBindings().Delete(crb.ctx, crb.Name, metav1.DeleteOptions{})
}

func (crb *ClusterRoleBinding) Update() error {
	_, err := crb.client.Clientset.RbacV1().ClusterRoleBindings().Update(crb.ctx, crb.ClusterRoleBinding, metav1.UpdateOptions{})
	return err
}

func (crb *ClusterRoleBinding) Empty() bool {
	_, err := crb.client.Clientset.RbacV1().ClusterRoleBindings().Get(crb.ctx, crb.Name, metav1.GetOptions{})
	return errors.IsNotFound(err)
}

func (crb *ClusterRoleBinding) List() (*v1.ClusterRoleBindingList, error) {
	return crb.client.Clientset.RbacV1().ClusterRoleBindings().List(crb.ctx, metav1.ListOptions{})
}

func (crb *ClusterRoleBinding) CreateOrUpdate() error {
	_, err := crb.client.Clientset.RbacV1().ClusterRoleBindings().Get(crb.ctx, crb.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return crb.Create()
		}
		return err
	}
	return crb.Update()
}

func (crb *ClusterRoleBinding) Equal() bool {
	clusterRoleBinding, err := crb.Get()
	if err != nil && !errors.IsNotFound(err) {
		panic(err)
	}
	return reflect.DeepEqual(clusterRoleBinding.Labels, crb.ClusterRoleBinding.Labels) &&
		reflect.DeepEqual(clusterRoleBinding.Annotations, crb.ClusterRoleBinding.Annotations) &&
		reflect.DeepEqual(clusterRoleBinding.Subjects, crb.ClusterRoleBinding.Subjects) &&
		reflect.DeepEqual(clusterRoleBinding.RoleRef, crb.ClusterRoleBinding.RoleRef)
}
