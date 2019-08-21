package config_test

import (
	"fmt"
	"github.com/openshift/hive/test/functional/utils"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInstallConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite for Hive")
}

var _ = Describe("[200_000009_000002] OCP-23970 - [ipi-on-aws] The cluster name is limited by 63 characters", func() {

	Describe("Create cluster with cluster name more then 63 characters", func() {
		Context("Try to create cluster with cluster name more then 63 characters", func() {
			It("Try to create cluster with cluster name more then 63 characters", func() {
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
   AWS_ACCESS_KEY_ID="%s" \
   AWS_SECRET_ACCESS_KEY="%s" \
   | oc apply -f -
`
				clusterName := "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz1234567890ab"
				cmdOut, _ := utils.RunShellCmd(fmt.Sprintf(script, clusterName,
					os.Getenv("QE_OPENSHIFT_PAYLOAD_IMAGE"), utils.GetSshPublicKey("id_rsa.pub"),
					keyId, accessKey))
				fmt.Println("Create cluster deployment with cluster name more then 63 characters,", cmdOut)

				Ω(strings.Contains(cmdOut, "Invalid cluster deployment name (.meta.name): must be no more than 63 characters")).Should(Equal(true))

				// release the resource by OCP-999999
			})
		})
	})
})
