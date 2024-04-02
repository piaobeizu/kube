/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/03/22 17:30:25
 Desc     :
*/

package kube

type Interface interface {
	Pod() PodInterface
}

type FlatteItem struct {
	Name string
	Val  any
	Kind string
}
