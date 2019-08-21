package sharedbehaviors

import (
	"fmt"
	"github.com/openshift/hive/test/functional/utils"
	"os"
	"strings"

	. "github.com/onsi/gomega"
)

func GetClusterBaiscInfromation() (clusterName, clusterID, infraID string) {
	infraID, _ = utils.RunShellCmd(fmt.Sprintf("cat %s | ", utils.GetMetadataFile()) + `jq -r '.infraID'`)
	infraID = strings.Trim(infraID, "\n")

	clusterID, _ = utils.RunShellCmd(fmt.Sprintf("cat %s | ", utils.GetMetadataFile()) + `jq -r '.clusterID'`)
	clusterID = strings.Trim(clusterID, "\n")

	clusterName, _ = utils.RunShellCmd(fmt.Sprintf("cat %s | ", utils.GetMetadataFile()) + `jq -r '.clusterName'`)
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
	values, _ = utils.RunShellCmd(awsCmd + ` resourcegroupstaggingapi get-tag-values --key openshiftClusterID | grep ` + clusterID)
	keys, _ = utils.RunShellCmd(awsCmd + ` resourcegroupstaggingapi get-tag-keys | grep "kubernetes.io/cluster/` + infraID + `"`)
	return
}

func CheckClusterResouceAreDestroyed(awsRegion, clusterID, infraID string) {
	workingDir := utils.GetAssetsFlag()

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
	cmdOut, _ := utils.RunShellCmd(fmt.Sprintf(script, workingDir, awsRegion, infraID))
	checkCmd := os.Getenv("QE_PROJECT_HOME") + "/bin/check_aws_resource.sh"
	awsFile := workingDir + "/aw.mine.txt"
	cmdOut, err := utils.RunShellCmd(fmt.Sprintf("sh %s %s %s", checkCmd, awsFile, awsRegion))
	utils.RunShellCmd(fmt.Sprintf(`rm -f %s/aw*.txt`, workingDir))
	fmt.Println("Run check output", cmdOut)
	Ω(err).NotTo(HaveOccurred())
}
