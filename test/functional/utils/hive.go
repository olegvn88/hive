package utils

import (
	"fmt"
	"os"
	"strings"
)

func MakeDirectory(dir string) string {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Println("Make dir error ", err)
		dir = ""
	}
	kubeConfig := dir + "/kubeconfig"
	secretKube := strings.TrimSuffix(GetKubeConfigName(), "\n")
	getKubeConfig := `oc get secret %s -o json | jq -r '.data.kubeconfig' | base64 -d > %s`
	RunShellCmd(fmt.Sprintf(getKubeConfig, secretKube, kubeConfig))

	return dir
}

func GetKubeConfigName() string {
	kubeConfigName, err := RunShellCmd(`oc get cd -o json | jq '.items[] .status .adminKubeconfigSecret .name'`)
	if err != nil {
		return "No resources found"
	}
	return kubeConfigName
}

//func GetHiveClusterName() string {
//	clusterName, err := RunShellCmd(`oc get cd | cut -d' ' -f1 | tail -1`)
//	if err != nil {
//		return "No resources found"
//	}
//	return clusterName
//}

func GetHiveClusterName(name string) string {
	var script string
	if name != "" {
		script = `oc get cd | cut -d' ' -f1 | grep %s`
		clusterName, err := RunShellCmd(fmt.Sprintf(script, name))
		if err != nil {
			fmt.Println("resource not found")
			return "resource not found"
		}
		return clusterName
	}
	script = `oc get cd | cut -d' ' -f1 `
	clusterName, err := RunShellCmd(script)
	PrintError(err)
	return clusterName
}

func PrintError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
