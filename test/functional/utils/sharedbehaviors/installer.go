package sharedbehaviors

import (
	"github.com/openshift/hive/test/functional/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"os"
	"strings"
)

type SharedInstallerInputs struct {
	AssetsDir string
	CmdOut    string
	NewEnv    []string
}

func OverrideEnvironmentEnv(testData *utils.DataToValidate) []string {
	// In CI-OPERATOR By default:
	// ci-operator/templates/openshift/installer/cluster-launch-installer-e2e.yaml
	// - name: OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE
	//   value: ${RELEASE_IMAGE_LATEST}

	// By manual:
	// export OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE=registry.svc.ci.openshift.org/ocp/release:4.0.0-0.nightly-2019-03-28-030453

	// Get the environment values from the current Env
	env := os.Environ()
	// Values from testData have the higher priority
	if testData.OS_IMAGE != "" {
		env = append(env, fmt.Sprintf("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE=%s", testData.OS_IMAGE))
	}
	if testData.RELEASE_IMAGE != "" {
		env = append(env, fmt.Sprintf("OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE=%s", testData.RELEASE_IMAGE))
	}

	return env
}

func SharedInstallLogBehavior(inputs *SharedInstallerInputs) {
	It("When I execute openshift-install create install-config", func() {
		cmdOut := inputs.CmdOut
		// ToDo: get more logs
		// Ω(cmdOut).Should(ContainSubstring("SSH Public Key"))
		// Ω(cmdOut).Should(ContainSubstring("Platform aws"))
		// Ω(cmdOut).Should(ContainSubstring("Region"))
		Ω(cmdOut).Should(ContainSubstring("Base Domain "))
		Ω(cmdOut).Should(ContainSubstring("Cluster Name"))
		Ω(cmdOut).Should(ContainSubstring("Pull Secret"))
	})
}

func DeleteAssetsDir(assetsDir string) {
	var info, err = os.Stat(assetsDir)

	if err != nil { // file existed
		return
	}

	err = os.RemoveAll(assetsDir)
	if err != nil {
		fmt.Println("Remove dir error: ", err)
		os.Exit(-1)
	}

	err = os.MkdirAll(assetsDir, info.Mode())
	if err != nil {
		fmt.Println("Create dir error: ", err)
		os.Exit(-1)
	}
}

func checkClusterHealth(payload_version string) string {
	if payload_version == "" {
		payload_version = os.Getenv("QE_OPENSHIFT_VERSION")
	}

	Ω(payload_version).ShouldNot(BeEmpty())
	cmdOut, _ := utils.RunShellCmd(`check_cluster_health.sh ` + payload_version)
	return cmdOut
}

func SharedCheckClusterHealthPrint(payload_version string) {
	fmt.Println(checkClusterHealth(payload_version))
}

func SharedCheckClusterHealthFunc(payload_version string) {
	cmdOut := checkClusterHealth(payload_version)
	if strings.Contains(cmdOut, "Errrorr") {
		fmt.Println("Cluster is broken", cmdOut)
		os.Exit(-1)
	}
	Ω(cmdOut).ShouldNot(ContainSubstring("Errrorr"))
}

func SharedCheckClusterHealthBehavior(payload_version string) {
	It("After installation, Run check the cluster health", func() {
		SharedCheckClusterHealthFunc(payload_version)
	})
}
