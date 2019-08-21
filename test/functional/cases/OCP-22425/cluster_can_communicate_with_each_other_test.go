package config_test
/*
import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	s "github.com/openshift/hive/test/functional/utils/sharedbehaviors"

	"github.com/openshift/hive/test/functional/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInstallConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite for Hive")
}

var (
	sharedInputs = s.SharedInstallerInputs{}
	testData     = utils.NewTestData("olnester-22425")
)

var _ = Describe("[200_000009_000001] OCP-22425 - [ipi-on-aws] Clusters that use the same subdomain can communicate with each other", func() {
	BeforeEach(func() {
		sharedInputs.AssetsDir = utils.GetAssetsFlag()
	})

	Describe("Set up Hive", func() {
		// If we run in the environment of ci-operator, The initialization scripts will download the kubectl binary
		Context("Make sure that you have installed kubectl and kustomize tool", func() {
			It("install kustomize tool", func() {
				script := `curl -LO https://github.com/kubernetes-sigs/kustomize/releases/download/v2.0.2/kustomize_2.0.2_linux_amd64
chmod +x kustomize_2.0.2_linux_amd64
mv ./kustomize_2.0.2_linux_amd64 kustomize
mv -f kustomize $GOPATH/bin`
				utils.RunShellCmd(script)

				// Install mockgen:
				utils.RunShellCmd(`go get github.com/golang/mock/gomock; go install github.com/golang/mock/mockgen`)

				// Deploy hive
				script = `mkdir -p $GOPATH/src/github.com/openshift
cd $GOPATH/src/github.com/openshift 
rm -rf hive
git clone https://github.com/openshift/hive.git 
cd hive

export KUBECONFIG=%s/auth/kubeconfig
make deploy`
				cmdOut, _ := utils.RunShellCmd(fmt.Sprintf(script, sharedInputs.AssetsDir))
				fmt.Println("Deploy Hive,", cmdOut)
				time.Sleep(2 * time.Minute) //wait 2 minutes until hive is installed

				keyId, accessKey := utils.GetAwsCredentials()
				// Create clusterdeployment:
				script = `export CLUSTER_NAME="%s"
export RELEASE_IMAGE="%s"
export SSH_PUB_KEY="%s"
export PULL_SECRET=$(cat ${PULL_SECRET_PATH})
export KUBECONFIG=%s/auth/kubeconfig

sed -i 's/region: us-east-1/region: %s/g' $GOPATH/src/github.com/openshift/hive/config/templates/cluster-deployment.yaml
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
				// If create clusterdeployment failed, run: `oc delete clusterdeployment ${CLUSTER_NAME}`
				// fmt.Println(fmt.Sprintf(script, testData.ClusterName+"-hive", "ap-northeast-1", testData.BaseDomain))
				clusterName := testData.ClusterName + "-hive"
				cmdOut, _ = utils.RunShellCmd(fmt.Sprintf(script, clusterName,
					os.Getenv("QE_OPENSHIFT_PAYLOAD_IMAGE"), utils.GetSshPublicKey("id_rsa.pub"),
					sharedInputs.AssetsDir, "ap-northeast-1", testData.BaseDomain, keyId, accessKey))
				fmt.Println("Create clusterdeployment,", cmdOut)

				err := waitInstallerFinished(1 * 60 * 60) // wait 1 hour
				if err != nil {
					cmd := `oc log $(oc get pod --all-namespaces |grep hive |grep 22425 |grep install | awk '{print $2}') -c hive | tail -n 200`
					cmdOut, _ = utils.RunShellCmd(cmd)
					fmt.Println("Hive log:", cmdOut)
				}
				Ω(err).ShouldNot(HaveOccurred())

				// check the secrets and cluster health
				err = checkPeerCluster(clusterName)
				Ω(err).ShouldNot(HaveOccurred())

				// release the resource by OCP-999999
			})
		})
	})
})

func checkInstallerProgressing() (bool, error) {
	cmdOut, err := utils.RunShellCmd(`oc get pod  |grep shlao | grep hive-install |grep Completed |wc -l`)
	if err != nil {
		fmt.Println("Failed to get pod of hive-install", cmdOut, err)
		return false, err
	}

	if utils.InsenstiveContains(cmdOut, "1") {
		return true, nil
	}

	return false, nil
}

func waitInstallerFinished(timeOut int) error {
	time.Sleep(5 * time.Minute)
	return utils.WaitTimeOut(timeOut, 20, 0, checkInstallerProgressing)
}

func makeDirectory(dir string) string {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Println("Make dir error ", err)
		dir = ""
	}

	return dir
}

func checkPeerCluster(clusterName string) error {
	assetDir := utils.GetAssetsDirectory()
	hiveDir := assetDir + "/hive"
	makeDirectory(hiveDir)

	kubeConfig := hiveDir + "/kubeconfig"
	secretKube := strings.TrimSuffix(utils.GetKubeConfigName(), "\n")
	getKubeConfig := `oc get secret %s -o json | jq -r '.data.kubeconfig' | base64 -d > %s`
	utils.RunShellCmd(fmt.Sprintf(getKubeConfig, secretKube, kubeConfig))

	// check the New cluster
	newEnv := os.Environ()
	newEnv = append(newEnv, fmt.Sprintf("KUBECONFIG=%s", kubeConfig))
	utils.RunShellEnvCmd(newEnv, `check_cluster_health.sh `+os.Getenv("QE_OPENSHIFT_VERSION"))
	_, err := utils.RunShellEnvCmd(newEnv, "oc get node")

	return err
}
*/