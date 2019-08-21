package config_test

import (
	"github.com/openshift/hive/test/functional/utils"

	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInstallConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite for Hive")
}

var (
	testData    = utils.NewTestData("olnester-23289")
	clusterName = testData.ClusterName
)
var _ = Describe("[200_000009_000003] OCP-23289 - [ipi-on-aws] Report current number of install job retries in cluster deployment status", func() {

	Describe("Report current number of install job retries in cluster deployment status", func() {
		Context("Report current number of install job retries in cluster deployment status", func() {

			AfterEach(func() {
				fmt.Printf("Deleting cluster deployment %s\n %s\n", clusterName, CurrentGinkgoTestDescription().TestText)
				script := `oc delete clusterdeployment %s`
				cmdOut, _ := utils.RunShellCmd(fmt.Sprintf(script, clusterName))
				fmt.Println(">>>>>>>>>>>>result<<<<<<<<<<<<<\n", cmdOut)
			})

			It("Get current number of install job retries in cluster deployment status", func() {
				keyId, accessKey := utils.GetAwsCredentials()
				// Create clusterdeployment:
				script := `export CLUSTER_NAME="%s"
export RELEASE_IMAGE="%s"
export SSH_PUB_KEY="%s"
//export PULL_SECRET=$(cat ${PULL_SECRET_PATH})
export PULL_SECRET="$(cat ${HOME}/config.json)"

//sed -i 's/region: us-east-1/region: %s/g' $GOPATH/src/github.com/openshift/hive/config/templates/cluster-deployment.yaml
oc process -f $GOPATH/src/github.com/openshift/hive/config/templates/cluster-deployment.yaml \
   RELEASE_IMAGE="${RELEASE_IMAGE}" \
   CLUSTER_NAME="${CLUSTER_NAME}" \
   SSH_KEY="${SSH_PUB_KEY}" \
   PULL_SECRET="${PULL_SECRET}" \
   BASE_DOMAIN="%s" \
   AWS_ACCESS_KEY_ID="%s" \
   AWS_SECRET_ACCESS_KEY="%s" \
   | oc apply -f -
`
				cmdOut, _ := utils.RunShellCmd(fmt.Sprintf(script, clusterName,
					os.Getenv("QE_OPENSHIFT_PAYLOAD_IMAGE"), utils.GetSshPublicKey("id_rsa.pub"),
					testData.BaseDomain+"1", keyId, accessKey))
				fmt.Println("Create cluster deployment with wrong base domain\n", cmdOut)
				err := waitUntilRetryIsAppeared(40)
				Ω(err).ShouldNot(HaveOccurred())

				//check that count of retries is displayed in the cluster deployment
				script = `oc get clusterdeployment %s -o yaml | grep installRestarts`
				cmdOut, _ = utils.RunShellCmd(fmt.Sprintf(script, clusterName))
				Ω(strings.Contains(cmdOut, "installRestarts:")).Should(Equal(true))
			})
		})
	})
})

func waitUntilRetryIsAppeared(timeOut int) error {
	time.Sleep(40 * time.Second) // wait for 40 seconds until install pod starts
	return utils.WaitTimeOut(timeOut, 2, 0, checkInstallationError)
}

func checkInstallationError() (bool, error) {
	script := `oc logs -f %s -c hive`
	podName := strings.TrimSuffix(GetInstallPodName(clusterName), "\n")
	cmdOut, _ := utils.RunShellCmd(fmt.Sprintf(script, podName))

	if utils.InsenstiveContains(cmdOut, "exit status 1") {
		time.Sleep(2 * time.Second)
		return true, nil
	}
	return false, nil
}

func GetInstallPodName(clustername string) string {
	podName := clustername + "-install"
	script := `oc get pod |grep %s | cut -d' ' -f1`
	installPodName, _ := utils.RunShellCmd(fmt.Sprintf(script, podName))
	if utils.InsenstiveContains(installPodName, podName) {
		return installPodName
	}
	return "pod is not found"
}
