/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/04/02 23:05:53
 Desc     :
*/

package kube

import (
	"encoding/json"
	"fmt"
	"testing"
)

var (
	aFI = `{
  "minreadyseconds": {
    "Name": "minReadySeconds",
    "Val": 0,
    "Kind": "int32"
  },
  "template.metadata.generatename": {
    "Name": "template.metadata.generateName",
    "Val": "",
    "Kind": "string"
  },
  "template.metadata.generation": {
    "Name": "template.metadata.generation",
    "Val": 0,
    "Kind": "int64"
  },
  "template.metadata.name": {
    "Name": "template.metadata.name",
    "Val": "",
    "Kind": "string"
  },
  "template.metadata.namespace": {
    "Name": "template.metadata.namespace",
    "Val": "",
    "Kind": "string"
  },
  "template.metadata.resourceversion": {
    "Name": "template.metadata.resourceVersion",
    "Val": "",
    "Kind": "string"
  },
  "template.metadata.selflink": {
    "Name": "template.metadata.selfLink",
    "Val": "",
    "Kind": "string"
  },
  "template.metadata.uid": {
    "Name": "template.metadata.uid",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.containers.0.env.0.name": {
    "Name": "template.spec.containers.0.env.0.name",
    "Val": "APP_CONFIG_PATH",
    "Kind": "string"
  },
  "template.spec.containers.0.env.0.value": {
    "Name": "template.spec.containers.0.env.0.value",
    "Val": "/etc/app/config.json",
    "Kind": "string"
  },
  "template.spec.containers.0.env.1.name": {
    "Name": "template.spec.containers.0.env.1.name",
    "Val": "APOLLO_PODIP",
    "Kind": "string"
  },
  "template.spec.containers.0.env.1.value": {
    "Name": "template.spec.containers.0.env.1.value",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.containers.0.env.1.valuefrom.fieldref.apiversion": {
    "Name": "template.spec.containers.0.env.1.valueFrom.fieldRef.apiVersion",
    "Val": "v1",
    "Kind": "string"
  },
  "template.spec.containers.0.env.1.valuefrom.fieldref.fieldpath": {
    "Name": "template.spec.containers.0.env.1.valueFrom.fieldRef.fieldPath",
    "Val": "status.podIP",
    "Kind": "string"
  },
  "template.spec.containers.0.env.2.name": {
    "Name": "template.spec.containers.0.env.2.name",
    "Val": "APOLLO_HOSTIP",
    "Kind": "string"
  },
  "template.spec.containers.0.env.2.value": {
    "Name": "template.spec.containers.0.env.2.value",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.containers.0.env.2.valuefrom.fieldref.apiversion": {
    "Name": "template.spec.containers.0.env.2.valueFrom.fieldRef.apiVersion",
    "Val": "v1",
    "Kind": "string"
  },
  "template.spec.containers.0.env.2.valuefrom.fieldref.fieldpath": {
    "Name": "template.spec.containers.0.env.2.valueFrom.fieldRef.fieldPath",
    "Val": "status.hostIP",
    "Kind": "string"
  },
  "template.spec.containers.0.image": {
    "Name": "template.spec.containers.0.image",
    "Val": "modelbest-registry.cn-beijing.cr.aliyuncs.com/deploy-images/apollo-labor:68fb69ac",
    "Kind": "string"
  },
  "template.spec.containers.0.imagepullpolicy": {
    "Name": "template.spec.containers.0.imagePullPolicy",
    "Val": "Always",
    "Kind": "string"
  },
  "template.spec.containers.0.name": {
    "Name": "template.spec.containers.0.name",
    "Val": "labor",
    "Kind": "string"
  },
  "template.spec.containers.0.stdin": {
    "Name": "template.spec.containers.0.stdin",
    "Val": false,
    "Kind": "bool"
  },
  "template.spec.containers.0.stdinonce": {
    "Name": "template.spec.containers.0.stdinOnce",
    "Val": false,
    "Kind": "bool"
  },
  "template.spec.containers.0.terminationmessagepath": {
    "Name": "template.spec.containers.0.terminationMessagePath",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.containers.0.terminationmessagepolicy": {
    "Name": "template.spec.containers.0.terminationMessagePolicy",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.containers.0.tty": {
    "Name": "template.spec.containers.0.tty",
    "Val": false,
    "Kind": "bool"
  },
  "template.spec.containers.0.volumemounts.0.mountpath": {
    "Name": "template.spec.containers.0.volumeMounts.0.mountPath",
    "Val": "/etc/app",
    "Kind": "string"
  },
  "template.spec.containers.0.volumemounts.0.name": {
    "Name": "template.spec.containers.0.volumeMounts.0.name",
    "Val": "apollo-config-labor",
    "Kind": "string"
  },
  "template.spec.containers.0.volumemounts.0.readonly": {
    "Name": "template.spec.containers.0.volumeMounts.0.readOnly",
    "Val": true,
    "Kind": "bool"
  },
  "template.spec.containers.0.volumemounts.0.subpath": {
    "Name": "template.spec.containers.0.volumeMounts.0.subPath",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.containers.0.volumemounts.0.subpathexpr": {
    "Name": "template.spec.containers.0.volumeMounts.0.subPathExpr",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.containers.0.workingdir": {
    "Name": "template.spec.containers.0.workingDir",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.dnspolicy": {
    "Name": "template.spec.dnsPolicy",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.hostipc": {
    "Name": "template.spec.hostIPC",
    "Val": false,
    "Kind": "bool"
  },
  "template.spec.hostname": {
    "Name": "template.spec.hostname",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.hostnetwork": {
    "Name": "template.spec.hostNetwork",
    "Val": false,
    "Kind": "bool"
  },
  "template.spec.hostpid": {
    "Name": "template.spec.hostPID",
    "Val": false,
    "Kind": "bool"
  },
  "template.spec.imagepullsecrets.0.name": {
    "Name": "template.spec.imagePullSecrets.0.name",
    "Val": "ali-registry-secret",
    "Kind": "string"
  },
  "template.spec.nodename": {
    "Name": "template.spec.nodeName",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.priorityclassname": {
    "Name": "template.spec.priorityClassName",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.restartpolicy": {
    "Name": "template.spec.restartPolicy",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.schedulername": {
    "Name": "template.spec.schedulerName",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.securitycontext.fsgroup": {
    "Name": "template.spec.securityContext.fsGroup",
    "Val": 0,
    "Kind": "int64"
  },
  "template.spec.securitycontext.runasgroup": {
    "Name": "template.spec.securityContext.runAsGroup",
    "Val": 0,
    "Kind": "int64"
  },
  "template.spec.securitycontext.runasuser": {
    "Name": "template.spec.securityContext.runAsUser",
    "Val": 0,
    "Kind": "int64"
  },
  "template.spec.serviceaccount": {
    "Name": "template.spec.serviceAccount",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.serviceaccountname": {
    "Name": "template.spec.serviceAccountName",
    "Val": "apollo-sa-labor",
    "Kind": "string"
  },
  "template.spec.subdomain": {
    "Name": "template.spec.subdomain",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.volumes.0.configmap.defaultmode": {
    "Name": "template.spec.volumes.0.configMap.defaultMode",
    "Val": 420,
    "Kind": "int32"
  },
  "template.spec.volumes.0.configmap.name": {
    "Name": "template.spec.volumes.0.configMap.name",
    "Val": "apollo-config-labor",
    "Kind": "string"
  },
  "template.spec.volumes.0.name": {
    "Name": "template.spec.volumes.0.name",
    "Val": "apollo-config-labor",
    "Kind": "string"
  },
  "updatestrategy.type": {
    "Name": "updateStrategy.type",
    "Val": "",
    "Kind": "string"
  }
}`
	bFI = `{
  "minreadyseconds": {
    "Name": "minReadySeconds",
    "Val": 0,
    "Kind": "int32"
  },
  "revisionhistorylimit": {
    "Name": "revisionHistoryLimit",
    "Val": 10,
    "Kind": "int32"
  },
  "template.metadata.generatename": {
    "Name": "template.metadata.generateName",
    "Val": "",
    "Kind": "string"
  },
  "template.metadata.generation": {
    "Name": "template.metadata.generation",
    "Val": 0,
    "Kind": "int64"
  },
  "template.metadata.name": {
    "Name": "template.metadata.name",
    "Val": "",
    "Kind": "string"
  },
  "template.metadata.namespace": {
    "Name": "template.metadata.namespace",
    "Val": "",
    "Kind": "string"
  },
  "template.metadata.resourceversion": {
    "Name": "template.metadata.resourceVersion",
    "Val": "",
    "Kind": "string"
  },
  "template.metadata.selflink": {
    "Name": "template.metadata.selfLink",
    "Val": "",
    "Kind": "string"
  },
  "template.metadata.uid": {
    "Name": "template.metadata.uid",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.containers.0.env.0.name": {
    "Name": "template.spec.containers.0.env.0.name",
    "Val": "APP_CONFIG_PATH",
    "Kind": "string"
  },
  "template.spec.containers.0.env.0.value": {
    "Name": "template.spec.containers.0.env.0.value",
    "Val": "/etc/app/config.json",
    "Kind": "string"
  },
  "template.spec.containers.0.env.1.name": {
    "Name": "template.spec.containers.0.env.1.name",
    "Val": "APOLLO_HOSTIP",
    "Kind": "string"
  },
  "template.spec.containers.0.env.1.value": {
    "Name": "template.spec.containers.0.env.1.value",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.containers.0.env.1.valuefrom.fieldref.apiversion": {
    "Name": "template.spec.containers.0.env.1.valueFrom.fieldRef.apiVersion",
    "Val": "v1",
    "Kind": "string"
  },
  "template.spec.containers.0.env.1.valuefrom.fieldref.fieldpath": {
    "Name": "template.spec.containers.0.env.1.valueFrom.fieldRef.fieldPath",
    "Val": "status.hostIP",
    "Kind": "string"
  },
  "template.spec.containers.0.env.2.name": {
    "Name": "template.spec.containers.0.env.2.name",
    "Val": "APOLLO_PODIP",
    "Kind": "string"
  },
  "template.spec.containers.0.env.2.value": {
    "Name": "template.spec.containers.0.env.2.value",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.containers.0.env.2.valuefrom.fieldref.apiversion": {
    "Name": "template.spec.containers.0.env.2.valueFrom.fieldRef.apiVersion",
    "Val": "v1",
    "Kind": "string"
  },
  "template.spec.containers.0.env.2.valuefrom.fieldref.fieldpath": {
    "Name": "template.spec.containers.0.env.2.valueFrom.fieldRef.fieldPath",
    "Val": "status.podIP",
    "Kind": "string"
  },
  "template.spec.containers.0.image": {
    "Name": "template.spec.containers.0.image",
    "Val": "modelbest-registry.cn-beijing.cr.aliyuncs.com/deploy-images/apollo-labor:68fb69dc",
    "Kind": "string"
  },
  "template.spec.containers.0.imagepullpolicy": {
    "Name": "template.spec.containers.0.imagePullPolicy",
    "Val": "Always",
    "Kind": "string"
  },
  "template.spec.containers.0.name": {
    "Name": "template.spec.containers.0.name",
    "Val": "labor",
    "Kind": "string"
  },
  "template.spec.containers.0.stdin": {
    "Name": "template.spec.containers.0.stdin",
    "Val": false,
    "Kind": "bool"
  },
  "template.spec.containers.0.stdinonce": {
    "Name": "template.spec.containers.0.stdinOnce",
    "Val": false,
    "Kind": "bool"
  },
  "template.spec.containers.0.terminationmessagepath": {
    "Name": "template.spec.containers.0.terminationMessagePath",
    "Val": "/dev/termination-log",
    "Kind": "string"
  },
  "template.spec.containers.0.terminationmessagepolicy": {
    "Name": "template.spec.containers.0.terminationMessagePolicy",
    "Val": "File",
    "Kind": "string"
  },
  "template.spec.containers.0.tty": {
    "Name": "template.spec.containers.0.tty",
    "Val": false,
    "Kind": "bool"
  },
  "template.spec.containers.0.volumemounts.0.mountpath": {
    "Name": "template.spec.containers.0.volumeMounts.0.mountPath",
    "Val": "/etc/app",
    "Kind": "string"
  },
  "template.spec.containers.0.volumemounts.0.name": {
    "Name": "template.spec.containers.0.volumeMounts.0.name",
    "Val": "apollo-config-labor",
    "Kind": "string"
  },
  "template.spec.containers.0.volumemounts.0.readonly": {
    "Name": "template.spec.containers.0.volumeMounts.0.readOnly",
    "Val": true,
    "Kind": "bool"
  },
  "template.spec.containers.0.volumemounts.0.subpath": {
    "Name": "template.spec.containers.0.volumeMounts.0.subPath",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.containers.0.volumemounts.0.subpathexpr": {
    "Name": "template.spec.containers.0.volumeMounts.0.subPathExpr",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.containers.0.workingdir": {
    "Name": "template.spec.containers.0.workingDir",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.dnspolicy": {
    "Name": "template.spec.dnsPolicy",
    "Val": "ClusterFirst",
    "Kind": "string"
  },
  "template.spec.hostipc": {
    "Name": "template.spec.hostIPC",
    "Val": false,
    "Kind": "bool"
  },
  "template.spec.hostname": {
    "Name": "template.spec.hostname",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.hostnetwork": {
    "Name": "template.spec.hostNetwork",
    "Val": false,
    "Kind": "bool"
  },
  "template.spec.hostpid": {
    "Name": "template.spec.hostPID",
    "Val": false,
    "Kind": "bool"
  },
  "template.spec.imagepullsecrets.0.name": {
    "Name": "template.spec.imagePullSecrets.0.name",
    "Val": "ali-registry-secret",
    "Kind": "string"
  },
  "template.spec.nodename": {
    "Name": "template.spec.nodeName",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.priorityclassname": {
    "Name": "template.spec.priorityClassName",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.restartpolicy": {
    "Name": "template.spec.restartPolicy",
    "Val": "Always",
    "Kind": "string"
  },
  "template.spec.schedulername": {
    "Name": "template.spec.schedulerName",
    "Val": "default-scheduler",
    "Kind": "string"
  },
  "template.spec.securitycontext.fsgroup": {
    "Name": "template.spec.securityContext.fsGroup",
    "Val": 0,
    "Kind": "int64"
  },
  "template.spec.securitycontext.runasgroup": {
    "Name": "template.spec.securityContext.runAsGroup",
    "Val": 0,
    "Kind": "int64"
  },
  "template.spec.securitycontext.runasuser": {
    "Name": "template.spec.securityContext.runAsUser",
    "Val": 0,
    "Kind": "int64"
  },
  "template.spec.serviceaccount": {
    "Name": "template.spec.serviceAccount",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.serviceaccountname": {
    "Name": "template.spec.serviceAccountName",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.subdomain": {
    "Name": "template.spec.subdomain",
    "Val": "",
    "Kind": "string"
  },
  "template.spec.terminationgraceperiodseconds": {
    "Name": "template.spec.terminationGracePeriodSeconds",
    "Val": 30,
    "Kind": "int64"
  },
  "template.spec.volumes.0.configmap.defaultmode": {
    "Name": "template.spec.volumes.0.configMap.defaultMode",
    "Val": 420,
    "Kind": "int32"
  },
  "template.spec.volumes.0.configmap.name": {
    "Name": "template.spec.volumes.0.configMap.name",
    "Val": "apollo-config-labor",
    "Kind": "string"
  },
  "template.spec.volumes.0.name": {
    "Name": "template.spec.volumes.0.name",
    "Val": "apollo-config-labor",
    "Kind": "string"
  },
  "updatestrategy.rollingupdate.maxsurge": {
    "Name": "updateStrategy.rollingUpdate.maxSurge",
    "Val": "",
    "Kind": "string"
  },
  "updatestrategy.rollingupdate.maxunavailable": {
    "Name": "updateStrategy.rollingUpdate.maxUnavailable",
    "Val": "",
    "Kind": "string"
  },
  "updatestrategy.type": {
    "Name": "updateStrategy.type",
    "Val": "RollingUpdate",
    "Kind": "string"
  }
}
`
)

func TestResourceEqual(t *testing.T) {
	var (
		a = make(map[string]FlatteItem)
		b = make(map[string]FlatteItem)
	)
	err := json.Unmarshal([]byte(aFI), &a)
	if err != nil {
		fmt.Printf("json.Unmarshal failed, err: %v", err)
	}
	err = json.Unmarshal([]byte(bFI), &b)
	if err != nil {
		fmt.Printf("json.Unmarshal failed, err: %v", err)
	}
	keys := []string{
		"^template.spec.containers.*.image$",
	}
	if !resourceEqual1(a, b, keys) {
		fmt.Print("Resource not equal\n")
	} else {
		fmt.Print("Resource equal\n")
	}
}

func TestResourceEqual2(t *testing.T) {

}
