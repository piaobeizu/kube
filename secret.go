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
	"reflect"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Secret struct {
	*LinkInfo
	*v1.Secret
	client *KubeClient
	ctx    context.Context
}

func NewSecret(ctx context.Context) *Secret {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &Secret{
		Secret: &v1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{},
			Data:       nil,
		},
		LinkInfo: &LinkInfo{},
		client:   nil,
		ctx:      ctx,
	}
}

func (s *Secret) Immutable(immutable bool) *Secret {
	s.Secret.Immutable = &immutable
	return s
}

func (s *Secret) Type(t v1.SecretType) *Secret {
	s.Secret.Type = t
	return s
}
func (s *Secret) Link(region, config string) *Secret {
	s.Region = region
	s.Config = config
	s.client = NewKubeClient(region, config)
	return s
}

func (s *Secret) Metadata(name, namespace string) *Secret {
	s.Name, s.Namespace = name, namespace
	return s
}

func (s *Secret) Data(data map[string][]byte) *Secret {
	if s.Secret.Data == nil {
		s.Secret.Data = make(map[string][]byte)
	}
	for k, v := range data {
		s.Secret.Data[k] = v
	}
	return s
}

func (s *Secret) StringData(data map[string]string) *Secret {
	if s.Secret.StringData == nil {
		s.Secret.StringData = make(map[string]string)
	}
	for k, v := range data {
		s.Secret.StringData[k] = v
	}
	return s
}

func (s *Secret) Labels(labels map[string]string) *Secret {
	if s.Secret.Labels == nil {
		s.Secret.Labels = make(map[string]string)
	}
	for k, v := range labels {
		s.Secret.Labels[k] = v
	}
	return s
}

func (s *Secret) Annotations(annotations map[string]string) *Secret {
	if s.Secret.Annotations == nil {
		s.Secret.Annotations = make(map[string]string)
	}
	for k, v := range annotations {
		s.Secret.Annotations[k] = v
	}
	return s
}

func (s *Secret) Create() error {
	Secrets := s.client.CoreV1().Secrets(s.Namespace)
	_, err := Secrets.Create(s.ctx, s.Secret, metav1.CreateOptions{})
	return err
}

func (s *Secret) Delete() error {
	Secrets := s.client.CoreV1().Secrets(s.Namespace)
	return Secrets.Delete(s.ctx, s.Name, metav1.DeleteOptions{})
}

func (s *Secret) Update() error {
	Secrets := s.client.CoreV1().Secrets(s.Namespace)
	_, err := Secrets.Update(s.ctx, s.Secret, metav1.UpdateOptions{})
	return err
}

func (s *Secret) Get() (*v1.Secret, error) {
	return s.client.CoreV1().Secrets(s.Namespace).Get(s.ctx, s.Name, metav1.GetOptions{})

}

func (s *Secret) CreateOrUpdate() error {
	_, err := s.client.CoreV1().Secrets(s.Namespace).Get(s.ctx, s.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return s.Create()
		}
		return err
	}
	return s.Update()
}

func (s *Secret) Empty() bool {
	_, err := s.client.CoreV1().Secrets(s.Namespace).Get(s.ctx, s.Name, metav1.GetOptions{})
	return errors.IsNotFound(err)
}

func (s *Secret) StringDataEqual() bool {
	secret, _ := s.Get()
	if secret != nil {
		return reflect.DeepEqual(secret.StringData, s.Secret.StringData)
	} else if s.Secret.StringData == nil {
		return true
	}
	return false
}

func (s *Secret) DataEqual() bool {
	secret, _ := s.Get()
	if secret != nil {
		return reflect.DeepEqual(secret.Data, s.Secret.Data)
	} else if s.Secret.Data == nil || len(s.Secret.Data) == 0 {
		return true
	}
	return false
}
