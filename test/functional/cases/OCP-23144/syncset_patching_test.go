package config_test

import (
	"github.com/openshift/hive/test/functional/utils"

	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"strings"
	"testing"
	"time"
)

func TestInstallConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite for Hive")
}

var (
	assetDir = utils.GetAssetsDirectory()
	hiveDir  = assetDir + "/hive"
	newEnv   = append(os.Environ(), fmt.Sprintf("KUBECONFIG=%s", hiveDir + "/kubeconfig"))
)

var _ = Describe("[200_000009_000004] OCP-23144 - [ipi-on-aws] SyncSet patching", func() {

	Describe("SyncSet patching", func() {
		Context("SyncSet patching", func() {
			It("SyncSet patchings", func() {
				configMapIsCreated := createConfigmapOnRemoteCluster()
				Ω(configMapIsCreated).Should(Equal(true))

				syncSetPatchIsApplied := applySyncSetPatch()
				Ω(syncSetPatchIsApplied).Should(Equal(true))
			})
		})
	})
})

func applySyncSetPatch() bool {
	syncSetFileName := "syncset-patch.yaml"
	syncSet := `echo "apiVersion: hive.openshift.io/v1alpha1
kind: SyncSet
metadata:
  name: test-patch
spec:
  clusterDeploymentRefs:
  - name: %s
  patches:
  - kind: ConfigMap
    apiVersion: v1
    name: foo
    namespace: default
    patch: |-
      { \"data\": { \"foo\": \"new-bar\" } }
    patchType: merge" > %s/%s
`
	result, err := utils.RunShellCmdWithEnv(os.Environ(), hiveDir, fmt.Sprintf(syncSet, utils.GetHiveClusterName(), hiveDir, syncSetFileName))
	printResultError(result, err)
	createObjectInCluster(syncSetFileName, os.Environ())
	time.Sleep(5 * time.Second)
	//loginToRemoteCluster()
	result, err = utils.RunShellCmdWithEnv(newEnv, hiveDir, fmt.Sprintf(`oc get configmap foo -o yaml`))
	if utils.InsenstiveContains(result, "foo: new-bar") {
		return true
	}
	fmt.Println("result:", result, "err: ", err)
	return false
}

func createObjectInCluster(fileName string, newEnv []string) {
	script := `oc create -f %s/%s`
	result, err := utils.RunShellCmdWithEnv(newEnv, hiveDir, fmt.Sprintf(script, hiveDir, fileName))
	printResultError(result, err)
}

func createConfigmapOnRemoteCluster() bool {
	//loginToRemoteCluster()
	configmapFileName := "configmap.yaml"
	configMapExample := `echo "apiVersion: v1
kind: ConfigMap
metadata:
  name: foo
  namespace: default
data:
  foo: bar" > %s/%s`
	result, err := utils.RunShellCmdWithEnv(newEnv, hiveDir, fmt.Sprintf(configMapExample, hiveDir, configmapFileName))
	printResultError(result, err)
	createObjectInCluster(configmapFileName, newEnv)

	result, err = utils.RunShellCmdWithEnv(newEnv, hiveDir, fmt.Sprintf(`oc get configmap foo -o yaml`))
	if utils.InsenstiveContains(result, "foo: bar") {
		return true
	}
	printResultError(result, err)
	return false
}

func printResultError(result string, err error) {
	if err != nil {
		fmt.Println("result:\n", result, "err:\n", err)
	}
}

func loginToRemoteCluster() error {
	utils.MakeDirectory(hiveDir)
	clusterName := strings.TrimSuffix(utils.GetHiveClusterName(), "\n")
	getClusterAdminPassword := `oc get secret %s-admin-password -o json | jq -r '.data.password' | base64 -d`

	clusterAdminPassword, _ := utils.RunShellCmd(fmt.Sprintf(getClusterAdminPassword, clusterName))
	loginCommand := `oc login -u kubeadmin -p %s`
	result, err := utils.RunShellEnvCmd(newEnv, fmt.Sprintf(loginCommand, clusterAdminPassword))
	printResultError(result, err)
	return err
}