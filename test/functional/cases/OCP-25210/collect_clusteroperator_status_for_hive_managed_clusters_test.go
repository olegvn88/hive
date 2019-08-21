package test_config

import (
	"encoding/json"
	"fmt"
	"github.com/openshift/hive/test/functional/utils"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInstallConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite for Hive")
}

var (
	assetDir = utils.GetAssetsDirectory()
	hiveDir  = assetDir + "/hive"
)

var _ = Describe("[200_000009_000010] OCP-25210 - [ipi-on-aws] Collect ClusterOperator Status for Hive Managed Clusters", func() {
	Describe("Collect ClusterOperator Status for Hive Managed Clusters", func() {
		Context("Collect ClusterOperator Status for Hive Managed Clusters", func() {
			It("Collect ClusterOperator Status for Hive Managed Clusters", func() {
				utils.MakeDirectory(hiveDir)
				result := getClusterOperatorStatuses()
				Ω(result).Should(Equal(true))
			})
		})
	})
})

func getClusterOperatorStatuses() bool {
	var clusterOperator ClusterOperator
	var result bool
	script := "oc get clusteroperator --config=%s/kubeconfig -o json"
	cmd, err := utils.RunShellCmd(fmt.Sprintf(script, hiveDir))
	utils.PrintError(err)
	err = json.Unmarshal([]byte(cmd), &clusterOperator)
	utils.PrintError(err)
	for i := 0; i < len(clusterOperator.Items); i++ {
		for _, condition := range clusterOperator.Items[i].Status.ConditionsOperator {
			if GetClusterState(clusterOperator.Items[i].Metadata.Name, condition.Type, condition.Status) {
				return true
			}
			fmt.Println(clusterOperator.Items[i].Metadata.Name, condition.Type, result)
		}
	}
	return result
}

func GetClusterState(operatorName, conditionType string, conditionStatus string) bool {
	var result bool
	var clusterState ClusterState
	script := "oc get clusterstate %s -o json"
	cmd, err := utils.RunShellCmd(fmt.Sprintf(script, strings.TrimSuffix(utils.GetHiveClusterName(""), "\n")))
	utils.PrintError(err)
	err = json.Unmarshal([]byte(cmd), &clusterState)
	utils.PrintError(err)
	for i := 0; i < len(clusterState.Status.ClusterOperators); i++ {
		for _, condition := range clusterState.Status.ClusterOperators[i].Conditions {
			if clusterState.Status.ClusterOperators[i].Name == operatorName && condition.Type == conditionType && condition.Status == conditionStatus {

				result = true
			}
		}
	}
	return result
}

type ClusterOperator struct {
	ApiVersion string `json:"apiVersion"`
	Items      []struct {
		ApiVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			CreationTimestamp string `json:"creationTimestamp"`
			Generation        int32  `json:"generation"`
			Name              string `json:"name"`
			ResourceVersion   string `json:"resourceVersion"`
			SelfLink          string `json:"selfLink"`
			Uid               string `json:"uid"`
		} `json:"metadata"`
		SpecOperator struct {
		} `json:"spec"`
		Status struct {
			ConditionsOperator []struct {
				LastTransitionTime string `json:"lastTransitionTime"`
				Reason             string `json:"reason"`
				Status             string `json:"status"`
				Type               string `json:"type"`
			} `json:"conditions"`
			RelatedObjects []struct {
				Group    string `json:"group"`
				Name     string `json:"name"`
				Resource string `json:"resource"`
			} `json:"relatedObjects"`
			Versions []struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"versions"`
		} `json:"status"`
	} `json:"items"`
}

type ClusterState struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		CreationTimestamp string `json:"creationTimestamp"`
		Generation        int32  `json:"generation"`
		Name              string `json:"name"`
		Namespace         string `json:"namespace"`
		OwnerReferences   []struct {
			ApiVersion         string `json:"apiVersion"`
			BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
			Controller         bool   `json:"controller"`
			Kind               string `json:"kind"`
			Name               string `json:"name"`
			Uid                string `json:"uid"`
		} `json:"ownerReferences"`
		ResourceVersion string `json:"resourceVersion"`
		SelfLink        string `json:"selfLink"`
		Uid             string `json:"uid"`
	} `json:"metadata"`
	Spec struct {
	} `json:"spec"`
	Status struct {
		ClusterOperators []struct {
			Conditions []struct {
				LastTransitionTime string `json:"lastTransitionTime"`
				Reason             string `json:"reason"`
				Status             string `json:"status"`
				Type               string `json:"type"`
			} `json:"conditions"`
			Name string `json:"name"`
		} `json:"clusterOperators"`
		lastUpdated string `json:"lastUpdated"`
	} `json:"status"`
}
