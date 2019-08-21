package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	disableSSHKey = "Oh, Don't use the ssh key"
	DefSSHKey     = "id_rsa.pub"
	defPlatform   = "aws"
	defRegion     = "ap-northeast-1"
	defBaseDomain = "qe1.devcluster.openshift.com"
	defPullSecret = `{"auths":{"www.openshift.com":{"auth":"b3Bl"}}}`
)

// DataToValidate as expected base domain and cluster name
type DataToValidate struct {
	SSH           string
	Platform      string
	Region        string
	BaseDomain    string
	ClusterName   string
	PullSecret    string
	OS_IMAGE      string
	RELEASE_IMAGE string
}

func NewcommandLineArgs(testData *DataToValidate) [][]string {
	commandLineArgs := make([][]string, 0)

	if testData.SSH != disableSSHKey {
		commandLineArgs = append(commandLineArgs, []string{"SSH", testData.SSH})
	}

	necessaryArgs := [][]string{
		{"Platform", testData.Platform},
		{"Region", testData.Region},
		{"Base Domain", testData.BaseDomain},
		{"Cluster Name", testData.ClusterName},
		// vt10x can't handle this: https://github.com/AlecAivazis/survey/issues/183
		{"Pull Secret", testData.PullSecret},
	}

	return append(commandLineArgs, necessaryArgs...)
}

func NewTestData(clusterName string) *DataToValidate {
	if clusterName == "" {
		return nil
	}

	// chekc: the environment has the sshkey
	enableSSHKey := false
	home := os.Getenv("HOME")
	if home != "" {
		_, err := filepath.Glob(filepath.Join(home, ".ssh", "*.pub"))
		if err == nil {
			enableSSHKey = true
		}
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = defRegion
	}

	rand := GetRandomInt(1000, 9999)
	testData := DataToValidate{
		Platform:    defPlatform,
		Region:      region,
		BaseDomain:  defBaseDomain,
		ClusterName: fmt.Sprintf("%s-%d", clusterName, rand),
		PullSecret:  defPullSecret,
	}

	if !enableSSHKey {
		testData.SSH = disableSSHKey
	} else {
		testData.SSH = DefSSHKey
	}

	return &testData
}
