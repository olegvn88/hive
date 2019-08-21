package utils

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const (
	MetadataFile       = "/metadata.json"
	DefAssetsDirectory = "/tmp/artifacts/installer/"
)

var (
	assetsDirectory = ""
	kubeConfig      = "KUBECONFIG"
	binScheduler    = "scheduler"
	defKubeFile     = "kubeconfig"
	sharedDir       = "clusteroperator-e2e/shared"

	assetsFlag string
)

func init() {
	dir, err := os.Getwd()
	if err == nil {
		dir += "/test1/"
	} else {
		dir = "./test1/"
	}

	flag.StringVar(&assetsFlag, "ad", dir, "Let the installer uses the different assets directory")
}

func GetAssetsFlag() string {
	return assetsFlag
}

func sliceToMap(env []string) map[string]string {
	envMap := make(map[string]string)

	for _, e := range env {
		pair := strings.Split(e, "=")
		// invalid format of env
		if len(pair) == 1 {
			fmt.Println("Invalid env", e)
			continue
		}

		envMap[pair[0]] = pair[1]
	}

	return envMap
}

func oldOverrideEnv(env []string) (newEnv []string) {
	envMap := sliceToMap(env)
	matchEnv := make(map[string]bool)

	// The Existing Env
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		val, ok := envMap[pair[0]]
		if ok {
			// using my env
			newEnv = append(newEnv, pair[0]+"="+val)
			matchEnv[pair[0]] = true
		} else {
			// using os env
			newEnv = append(newEnv, e)
		}
	}

	// The New Env
	for key, val := range envMap {
		_, ok := matchEnv[key]
		if !ok {
			newEnv = append(newEnv, key+"="+val)
		}
	}

	return
}

func OverrideEnv(env []string) (newEnv []string) {
	newEnv = os.Environ()
	newEnv = append(newEnv, env...)

	return
}

func GetSharedMountPath() string {
	path, err := exec.LookPath(binScheduler)
	if err != nil {
		fmt.Println("Failed to get path of", binScheduler, err)
		return ""
	}

	// e.g: /tmp/shared/bin/scheduler
	// filepath.Dir(filepath.Dir()) => /tmp/shared
	return filepath.Dir(filepath.Dir(path))
}

func GetSharedDirectory() string {
	shared := GetSharedMountPath()
	if shared == "" {
		return ""
	}

	return GetSharedMountPath() + "/" + sharedDir
}

func InsenstiveContains(a, b string) bool {
	return strings.Contains(strings.ToLower(a), strings.ToLower(b))
}

func SetAssetsDirectory(dir string) string {
	if dir != "" {
		assetsDirectory = dir
	}

	return assetsDirectory
}

// Return "" if we can't find the metadata.json
// kubeconfig reside in the assetsDirectory
func GetAssetsDirectory() string {
	if assetsDirectory == "" {
		SetAssetsDirectory(assetsFlag)
	}

	// os.Setenv("XXXXX", "yyyyy")
	// os.Getenv("XXXXX") return yyyyy
	kubePath := os.Getenv(kubeConfig)
	if kubePath != "" {
		// export KUBECONFIG=/root/4.0.0-0.nightly-2019-03-28-030453/test1/auth/kubeconfig
		if !InsenstiveContains(kubePath, defKubeFile) {
			return assetsDirectory
		}

		for tmpDir := path.Dir(kubePath); tmpDir != "/"; tmpDir = path.Dir(tmpDir) {
			if _, err := os.Stat(tmpDir + MetadataFile); err == nil {
				assetsDirectory = tmpDir
				break
			}
		}
	}

	return assetsDirectory
}

func GetMetadataFile() string {
	return GetAssetsDirectory() + MetadataFile
}

func GetSshPublicKey(fileName string) string {
	file := os.Getenv("SSH_PUB_KEY_PATH")
	if file != "" {
		_, err := os.Stat(file)
		if err == nil {
			dat, err := ioutil.ReadFile(file)
			if err == nil {
				return string(dat)
			}
		}
	}

	sshKey := os.Getenv("SSH_PUB_KEY")
	if sshKey != "" {
		return sshKey
	}

	home := os.Getenv("HOME")
	if home != "" {
		file = filepath.Join(home, ".ssh", fileName)
		_, err := os.Stat(file)
		if err == nil {
			dat, err := ioutil.ReadFile(file)
			if err == nil {
				return string(dat)
			}
		}
	}

	return ""
}

//func GetRegistryPath() string {
//	registry_path := os.Getenv("QE_OPENSHIFT_REGISTRY_PATH")
//	if registry_path != "" {
//		return registry_path
//	}
//
//	return DefRegistryServer
//}

func GetAwsCredentials() (keyId, accessKey string) {
	fileName := os.Getenv("HOME") + "/.aws/credentials"
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "reading ", fileName, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "aws_access_key_id") {
			slice := strings.Split(line, "=")
			if len(slice) != 1 {
				keyId = slice[1]
			}
		}

		if strings.Contains(line, "aws_secret_access_key") {
			slice := strings.Split(line, "=")
			if len(slice) != 1 {
				accessKey = slice[1]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Scan file:", fileName, err)
		return
	}

	return
}
