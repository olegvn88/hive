package sharedbehaviors

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/openshift/hive/test/functional/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func SharedLaunchUPIBehavior() {
	It("Launch cluster on UPI", func() {
		if needToCallJenkins() {
			_, err := callJenkins("launch_def_cluster")
			Ω(err).ShouldNot(HaveOccurred())
		}
	})
}

func SharedDestroyUPIBehavior() {
	It("Destroy cluster on UPI", func() {
		if needToCallJenkins() {
			_, err := callJenkins("destroy_def_cluster")
			Ω(err).ShouldNot(HaveOccurred())
		}
	})
}

func callJenkins(cmd string) (string, error) {
	dir := filepath.Dir(os.Getenv("KUBECONFIG"))
	if dir == "" {
		return "", errors.New("Failed get KUBECONFIG")
	}

	cmdOut, err := utils.RunShellCmd("mkdir -p " + dir)
	if err != nil {
		fmt.Printf("Failed to create dir:%s for:%s\n", dir, cmdOut)
		return cmdOut, err
	}

	cmdOut, err = utils.RunShellCmd(cmd)
	if err != nil {
		fmt.Printf("Failed to Run :%s cluster:%s\n", cmd, cmdOut)
	}
	return cmdOut, err
}

func needToCallJenkins() bool {
	platform := os.Getenv("QE_PLATFORM")
	if platform == "" {
		return false
	}
	return true
}
