/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2023/09/08 17:54:29
 Desc     :
*/

package kube

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	typev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type LinkInfo struct {
	Region string `json:"region"`
	Config string `json:"config"`
}

type KubeClient struct {
	*kubernetes.Clientset
}

func NewKubeClient(region, kubeconfig string) *KubeClient {
	return &KubeClient{getClient(region, kubeconfig)}
}

func (cs *KubeClient) Namespaces() typev1.NamespaceInterface {
	return cs.CoreV1().Namespaces()
}

func (cs *KubeClient) ResourceQuotas(namespace string) typev1.ResourceQuotaInterface {
	return cs.CoreV1().ResourceQuotas(namespace)
}

// func NewInClusterClientSet() *ClientSet {
// 	config, err := rest.InClusterConfig()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	// TODO: set QPS and Burst
// 	// config.QPS = 30?
// 	// config.Burst = 60
// 	// creates the clientset
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return &ClientSet{clientset}
// }

func BuildClusterConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig == "" {
		return rest.InClusterConfig()
	}
	var cfg rest.Config
	if err := yaml.Unmarshal([]byte(kubeconfig), &cfg); err != nil {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	// kubeconfig := parsePath(region)
	configPath := fmt.Sprintf("/tmp/%d", time.Now().Nanosecond())
	if err := os.WriteFile(configPath, []byte(kubeconfig), 0644); err != nil {
		return nil, err
	}
	rest, err := clientcmd.BuildConfigFromFlags("", configPath)
	os.Remove(configPath)
	return rest, err
}

var cs = clientSet{}

type clientSet struct {
	lock       sync.Mutex
	clientsets map[string]*kubernetes.Clientset
}

func getClient(region, kubeconfig string) *kubernetes.Clientset {
	if cs.clientsets == nil {
		cs.clientsets = make(map[string]*kubernetes.Clientset)
	}
	// TODO: set QPS and Burst
	// qps := viper.GetFloat64("api-server-qps")
	// burst := viper.GetInt("api-server-burst")
	if client, ok := cs.clientsets[region]; ok {
		return client
	} else {
		cfg, err := BuildClusterConfig(kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		// cfg.QPS = float32(qps)
		// cfg.Burst = burst
		client, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			panic(err.Error())
		}
		// map concurrency write
		cs.lock.Lock()
		cs.clientsets[region] = client
		cs.lock.Unlock()
		return client
	}
}

func parsePath(region string) string {
	path := os.Getenv("KUBECONFIG")
	if strings.HasPrefix(path, "") {
		path = homeDir()
	}
	return fmt.Sprintf("%s/.kube/%s.kubeconfig", path, region)
}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}
