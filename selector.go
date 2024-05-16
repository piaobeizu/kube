/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wwangxiaoakng@modelbest.cn'
 @Time    : 2024/05/16 20:00:18
 Desc     :
*/

package kube

import (
	v1 "k8s.io/api/core/v1"
)

type NodeSelector struct {
	*v1.NodeSelector
}

func NewNodeSelector() *NodeSelector {
	return &NodeSelector{
		NodeSelector: &v1.NodeSelector{
			NodeSelectorTerms: []v1.NodeSelectorTerm{},
		},
	}
}

func (n *NodeSelector) NodeSelectorTerm(term NodeSelectorTerm) *NodeSelector {
	n.NodeSelectorTerms = append(n.NodeSelectorTerms, term.NodeSelectorTerm)
	return n
}

type NodeSelectorTerm struct {
	v1.NodeSelectorTerm
}

func NewNodeSelectorTerm() *NodeSelectorTerm {
	return &NodeSelectorTerm{
		NodeSelectorTerm: v1.NodeSelectorTerm{
			MatchExpressions: []v1.NodeSelectorRequirement{},
			MatchFields:      []v1.NodeSelectorRequirement{},
		},
	}
}

func (n *NodeSelectorTerm) MatchExpression(key string, operator v1.NodeSelectorOperator, values []string) *NodeSelectorTerm {
	n.MatchExpressions = append(n.MatchExpressions, v1.NodeSelectorRequirement{
		Key:      key,
		Operator: operator,
		Values:   values,
	})
	return n
}
