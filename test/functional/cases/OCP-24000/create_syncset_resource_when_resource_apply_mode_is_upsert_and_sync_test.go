package config_test

import (
	"github.com/openshift/hive/test/functional/utils"

	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"strings"
	"testing"
)

func TestInstallConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite for Hive")
}

var (
	assetDir = utils.GetAssetsDirectory()
	hiveDir  = assetDir + "/hive"
	newEnv   = append(os.Environ(), fmt.Sprintf("KUBECONFIG=%s", hiveDir+"/kubeconfig"))
)

var _ = Describe("[200_000009_000004] OCP-24000 - [ipi-on-aws] Using Syncset when resourceApplyMode is Upsert and Sync", func() {

	Describe("Using Syncset when resourceApplyMode is Upsert", func() {
		Context("Using Syncset when resourceApplyMode is Upsert", func() {
			It("Create configmap object by Syncset when resourceApplyMode is Upsert", func() {
				utils.MakeDirectory(hiveDir)
				checkSyncsetIsCreated := createConfigMapOnRemoteClusterBySyncSetResource("upsert")
				Ω(checkSyncsetIsCreated).Should(Equal(true))

				checkSyncSetResourseIsDeleted := deleteSyncsetResource("upsert-test-syncset")
				Ω(checkSyncSetResourseIsDeleted).Should(Equal(true))

				configMapObectStillExists := checkConfigMapObectStillExists("upsert-foo")
				Ω(configMapObectStillExists).Should(Equal(true))
			})
		})
	})
	Describe("Using Syncset when resourceApplyMode is Sync", func() {
		Context("Using Syncset when resourceApplyMode is Sync", func() {
			It("Create configmap object by Syncset when resourceApplyMode is Sync", func() {
				checkSyncsetIsCreated := createConfigMapOnRemoteClusterBySyncSetResource("sync")
				Ω(checkSyncsetIsCreated).Should(Equal(true))

				checkSyncSetResourseIsDeleted := deleteSyncsetResource("sync-test-syncset")
				Ω(checkSyncSetResourseIsDeleted).Should(Equal(true))

				configMapObectStillExists := checkConfigMapObectStillExists("sync-foo")
				Ω(configMapObectStillExists).ShouldNot(Equal(true))
			})
		})
	})
})

func checkConfigMapObectStillExists(configMapName string) bool {
	result, err := utils.RunShellCmdWithEnv(newEnv, hiveDir, fmt.Sprintf(`oc get configmap %s -o yaml`, configMapName))
	printResultError(result, err)
	if utils.InsenstiveContains(result, "foo: bar_syncset") {
		return true
	}
	if utils.InsenstiveContains(result, "configmaps \""+configMapName+"\" not found") {
		return false
	}
	return false
}

func deleteSyncsetResource(syncsetName string) bool {
	script := `oc delete syncset %s`
	result, _ := utils.RunShellCmdWithEnv(os.Environ(), hiveDir, fmt.Sprintf(script, syncsetName))
	if utils.InsenstiveContains(result, "\""+syncsetName+"\" deleted") {
		return true
	}
	return false
}

func createObjectInCluster(fileName string, newEnv []string) {
	script := `oc create -f %s/%s`
	result, err := utils.RunShellCmdWithEnv(newEnv, hiveDir, fmt.Sprintf(script, hiveDir, fileName))
	printResultError(result, err)
}

func createConfigMapOnRemoteClusterBySyncSetResource(applyMode string) bool {
	syncSetFileName := applyMode + "-syncset-patch.yaml"
	syncSet := `echo "apiVersion: hive.openshift.io/v1alpha1
kind: SyncSet
metadata:
  name: %s-test-syncset
spec:
  clusterDeploymentRefs:
  - name: %s
  resourceApplyMode: \"%s\"
  resources:
  - kind: ConfigMap
    apiVersion: v1
    metadata:
      name: %s-foo
      namespace: default
    data:
      foo: bar_syncset" > %s/%s`
	result, err := utils.RunShellCmdWithEnv(os.Environ(), hiveDir, fmt.Sprintf(syncSet, applyMode, utils.GetHiveClusterName(), strings.Title(applyMode), applyMode, hiveDir, syncSetFileName))
	printResultError(result, err)
	createObjectInCluster(syncSetFileName, os.Environ())
	result, err = utils.RunShellCmdWithEnv(newEnv, hiveDir, fmt.Sprintf(`oc get configmap %s-foo -o yaml`, applyMode))
	if utils.InsenstiveContains(result, "foo: bar_sync") {
		return true
	}
	fmt.Println("result:", result, "err: ", err)
	return false
}

func printResultError(result string, err error) {
	if err != nil {
		fmt.Println("result:\n", result, "err:\n", err)
	}
}
