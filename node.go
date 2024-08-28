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

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Node struct {
	*LinkInfo
	*v1.Node
	client *KubeClient
	ctx    context.Context
}

func NewNode(ctx context.Context) *Node {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &Node{
		Node: &v1.Node{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Node",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{},
		},
		LinkInfo: &LinkInfo{},
		client:   nil,
		ctx:      ctx,
	}
}

func (n *Node) Link(region, config string) *Node {
	n.Region = region
	n.Config = config
	n.client = NewKubeClient(region, config)
	return n
}

func (n *Node) Metadata(name string) *Node {
	n.Name = name
	return n
}

func (n *Node) Labels(labels map[string]string) *Node {
	if n.Node.Labels == nil {
		n.Node.Labels = make(map[string]string)
	}
	for k, v := range labels {
		n.Node.Labels[k] = v
	}
	return n
}

func (n *Node) Annotations(annotations map[string]string) *Node {
	if n.Node.Annotations == nil {
		n.Node.Annotations = make(map[string]string)
	}
	for k, v := range annotations {
		n.Node.Annotations[k] = v
	}
	return n
}

func (n *Node) Create() error {
	Nodes := n.client.CoreV1().Nodes()
	_, err := Nodes.Create(n.ctx, n.Node, metav1.CreateOptions{})
	return err
}

func (n *Node) Delete() error {
	Nodes := n.client.CoreV1().Nodes()
	return Nodes.Delete(n.ctx, n.Name, metav1.DeleteOptions{})
}

func (n *Node) Update() error {
	Nodes := n.client.CoreV1().Nodes()
	_, err := Nodes.Update(n.ctx, n.Node, metav1.UpdateOptions{})
	return err
}

func (n *Node) Get() (*v1.Node, error) {
	return n.client.CoreV1().Nodes().Get(n.ctx, n.Name, metav1.GetOptions{})
}

func (n *Node) Fetch() (*Node, error) {
	node, err := n.client.CoreV1().Nodes().Get(n.ctx, n.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	n.Node = node
	return n, nil
}

func (n *Node) List() (*v1.NodeList, error) {
	return n.client.CoreV1().Nodes().List(n.ctx, metav1.ListOptions{})
}

func (n *Node) ListPods() ([]v1.Pod, error) {
	podList, err := n.client.CoreV1().Pods("").List(n.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var podsOnNode []v1.Pod
	for _, pod := range podList.Items {
		if pod.Spec.NodeName == n.Name {
			podsOnNode = append(podsOnNode, pod)
		}
	}
	return podsOnNode, nil
}

func (n *Node) CreateOrUpdate() error {
	_, err := n.client.CoreV1().Nodes().Get(n.ctx, n.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return n.Create()
		}
		return err
	}
	return n.Update()
}

func (n *Node) Empty() bool {
	_, err := n.client.CoreV1().Nodes().Get(n.ctx, n.Name, metav1.GetOptions{})
	return errors.IsNotFound(err)
}

func (n *Node) Equal(keys []string) bool {
	secret, err := n.Get()
	if err != nil && !errors.IsNotFound(err) {
		panic(err)
	}
	if len(keys) == 0 {
		keys = []string{"Data", "Type", "StringData", "Immutable"}
	}
	return ResourceEqual(n.Node, secret, keys)
}
