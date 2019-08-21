package utils

import (
	"fmt"
	"os"
	"strings"

	. "github.com/onsi/gomega"
)

func GetClusterBaiscInfromation() (clusterName, clusterID, infraID string) {
	infraID, _ = RunShellCmd(fmt.Sprintf("cat %s | ", GetMetadataFile()) + `jq -r '.infraID'`)
	infraID = strings.Trim(infraID, "\n")

	clusterID, _ = RunShellCmd(fmt.Sprintf("cat %s | ", GetMetadataFile()) + `jq -r '.clusterID'`)
	clusterID = strings.Trim(clusterID, "\n")

	clusterName, _ = RunShellCmd(fmt.Sprintf("cat %s | ", GetMetadataFile()) + `jq -r '.clusterName'`)
	clusterName = strings.Trim(clusterName, "\n")

	return
}

func GetAwsCmd(awsRegion string) string {
	awsCmd := `aws `
	if awsRegion != "" {
		awsCmd = `aws --region ` + awsRegion
	}

	return awsCmd
}

func CheckResourceGroupsTaggingApi(awsRegion, clusterID, infraID string) (keys, values string) {
	awsCmd := GetAwsCmd(awsRegion)
	values, _ = RunShellCmd(awsCmd + ` resourcegroupstaggingapi get-tag-values --key openshiftClusterID | grep ` + clusterID)
	keys, _ = RunShellCmd(awsCmd + ` resourcegroupstaggingapi get-tag-keys | grep "kubernetes.io/cluster/` + infraID + `"`)
	return
}

func CheckClusterResouceAreDestroyed(awsRegion, clusterID, infraID string) {
	workingDir := GetAssetsFlag()

	keys, values := CheckResourceGroupsTaggingApi(awsRegion, clusterID, infraID)
	fmt.Println("aws get clusterID", values)
	//Ω(values).Should(BeEmpty()
	fmt.Println("aws get infraID", keys)
	//Ω(keys).Should(BeEmpty())

	var script = `tmp_dir="%s"
region="%s"

AWSCmd="aws"
if [[ x"$region" != x ]]
then
AWSCmd="aws --region ${region}"
fi
echo "aws: ${AWSCmd}"

${AWSCmd}  resourcegroupstaggingapi get-resources > ${tmp_dir}/aw.resource.txt
cat ${tmp_dir}/aw.resource.txt | jq -r '.ResourceTagMappingList[]  | select(.Tags[].Value | contains ("%s")) | .ResourceARN' | sort -u > ${tmp_dir}/aw.mine.txt
`
	cmdOut, _ := RunShellCmd(fmt.Sprintf(script, workingDir, awsRegion, infraID))
	checkCmd := os.Getenv("QE_PROJECT_HOME") + "/bin/check_aws_resource.sh"
	awsFile := workingDir + "/aw.mine.txt"
	cmdOut, err := RunShellCmd(fmt.Sprintf("sh %s %s %s", checkCmd, awsFile, awsRegion))
	RunShellCmd(fmt.Sprintf(`rm -f %s/aw*.txt`, workingDir))
	fmt.Println("Run check output", cmdOut)
	Ω(err).NotTo(HaveOccurred())
}
