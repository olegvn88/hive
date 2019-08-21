package config_test

import (
	"strings"
	"testing"

	"github.com/openshift/hive/test/functional/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInstallConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite for Cloud-Credential Operator")
}

var _ = Describe("[200_000007_000003] OCP-24349 - [ipi-on-aws] Cloud-credential operator contains ClusterStatus.Condition.Type Upgradeable", func() {

	Describe("Check cloud credential status.condition.type contains Upgradeable", func() {
		Context("Scenario: Check cloud credential status.condition.type contains Upgradeable", func() {
			It("Check cloud credential status.condition.type contains Upgradeable", func() {

				script := `oc get co cloud-credential -o json | jq -r '.status.conditions | .[] .type'`
				cmdOut, _ := utils.RunShellCmd(script)

				Ω(strings.Contains(cmdOut, "Upgradeable")).Should(Equal(true))
			})
		})
	})
})
