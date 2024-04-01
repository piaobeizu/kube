/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/04/01 11:50:34
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

type ClusterRole struct {
	*LinkInfo
	*v1.ClusterRole
	ctx    context.Context
	client *KubeClient
}

func NewClusterRole(ctx context.Context) *ClusterRole {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &ClusterRole{
		ClusterRole: &v1.ClusterRole{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ClusterRole",
				APIVersion: "rbac.authorization.k8s.io/v1",
			},
			Rules:           []v1.PolicyRule{},
			AggregationRule: &v1.AggregationRule{},
		},
		ctx:      ctx,
		LinkInfo: &LinkInfo{},
		client:   nil,
	}
}

func (cr *ClusterRole) Metadata(name string) *ClusterRole {
	cr.Name = name
	return cr
}
func (cr *ClusterRole) Labels(labels map[string]string) *ClusterRole {
	if cr.ClusterRole.Labels == nil {
		cr.ClusterRole.Labels = make(map[string]string)
	}
	cr.ClusterRole.Labels = labels
	return cr
}

func (cr *ClusterRole) Annotations(annotations map[string]string) *ClusterRole {
	if cr.ClusterRole.Annotations == nil {
		cr.ClusterRole.Annotations = make(map[string]string)
	}
	cr.ClusterRole.Annotations = annotations
	return cr
}

func (cr *ClusterRole) Link(region, config string) *ClusterRole {
	cr.Region = region
	cr.Config = config
	cr.client = NewKubeClient(region, config)
	return cr
}
func (cr *ClusterRole) Rule(verbs []string, apiGroups []string, resoueces []string) *ClusterRole {
	if cr.Rules == nil {
		cr.Rules = make([]v1.PolicyRule, 0)
	}
	cr.Rules = append(cr.Rules, v1.PolicyRule{
		Verbs:     verbs,
		APIGroups: apiGroups,
		Resources: resoueces,
	})
	return cr
}

func (cr *ClusterRole) AggregationRule(matchLabels map[string]string) *ClusterRole {
	if cr.ClusterRole.AggregationRule == nil {
		cr.ClusterRole.AggregationRule = &v1.AggregationRule{
			ClusterRoleSelectors: []metav1.LabelSelector{},
		}
	}
	cr.ClusterRole.AggregationRule.ClusterRoleSelectors = append(cr.ClusterRole.AggregationRule.ClusterRoleSelectors, metav1.LabelSelector{
		MatchLabels: matchLabels,
	})
	return cr
}

func (cr *ClusterRole) Create() (*v1.ClusterRole, error) {
	return cr.client.RbacV1().ClusterRoles().Create(cr.ctx,
		cr.ClusterRole, metav1.CreateOptions{})
}

func (cr *ClusterRole) Update() (*v1.ClusterRole, error) {
	return cr.client.RbacV1().ClusterRoles().Update(cr.ctx,
		cr.ClusterRole, metav1.UpdateOptions{})
}

func (cr *ClusterRole) Delete() error {
	return cr.client.RbacV1().ClusterRoles().Delete(cr.ctx, cr.Name, metav1.DeleteOptions{})
}

func (cr *ClusterRole) Get() (*v1.ClusterRole, error) {
	return cr.client.RbacV1().ClusterRoles().Get(cr.ctx, cr.Name, metav1.GetOptions{})
}

func (cr *ClusterRole) Empty() bool {
	_, err := cr.client.RbacV1().ClusterRoles().Get(cr.ctx, cr.Name, metav1.GetOptions{})
	return errors.IsNotFound(err)
}

func (cr *ClusterRole) CreateOrUpdate() (*v1.ClusterRole, error) {
	_, err := cr.client.RbacV1().ClusterRoles().Get(cr.ctx, cr.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return cr.Create()
		}
		return nil, err
	}
	return cr.Update()
}

func (cr *ClusterRole) List() (*v1.ClusterRoleList, error) {
	return cr.client.RbacV1().ClusterRoles().List(cr.ctx, metav1.ListOptions{})
}

func (cr *ClusterRole) Equal() bool {
	clusterRole, err := cr.Get()
	if !errors.IsNotFound(err) {
		panic(err)
	}
	if len(clusterRole.Rules) == 0 && len(cr.Rules) == 0 {
		return true
	}
	return reflect.DeepEqual(cr.ClusterRole.Rules, clusterRole.Rules)
}
