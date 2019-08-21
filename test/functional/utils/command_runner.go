package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	defTmpDir = "/tmp"
)

func CaseInsenstiveContains(a, b string) bool {
	return strings.Contains(strings.ToUpper(a), strings.ToUpper(b))
}

func RunShellCmd(script string) (string, error) {
	return RunShellCmdWithEnv(nil, defTmpDir, script)
}

func RunShellEnvCmd(env []string, script string) (string, error) {
	return RunShellCmdWithEnv(env, defTmpDir, script)
}

func RunCurrentShell(script string) (string, error) {
	tmpfile, err := ioutil.TempFile("./", "tmpShell")
	if err != nil {
		fmt.Println("Create Temp file error:", err)
		return "", err
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(script)); err != nil {
		fmt.Println("Write Temp file error:", tmpfile.Name(), err)
		return "", err
	}

	cmdOut, err := RunOSCommandWithEnv(nil, "./", "/bin/bash", tmpfile.Name())
	if err := tmpfile.Close(); err != nil {
		fmt.Println("Close Temp file error:", tmpfile.Name(), err)
		return "", err
	}

	return cmdOut, err
}

func RunShellCmdWithEnv(env []string, workingDir, script string) (string, error) {
	tmpDir, err := ioutil.TempDir(workingDir, "ginkgo")
	if err != nil {
		fmt.Println("Create Tmp Dir error:", err)
		return "", err
	}

	defer os.RemoveAll(tmpDir) // clean up

	scriptFile := filepath.Join(tmpDir, "scriptFileile")
	if err := ioutil.WriteFile(scriptFile, []byte(script), 0666); err != nil {
		fmt.Println("Write Temp file error:", scriptFile, err)
		return "", err
	}

	if workingDir == defTmpDir {
		// run in the tmpDir
		workingDir = tmpDir
	}

	return RunOSCommandWithEnv(env, workingDir, "/bin/bash", scriptFile)
}

func RunOSCommand(command string, arguments ...string) (string, error) {
	cmdOut, err := RunOSCommandWithArgs(command, arguments...)
	if err != nil {
		fmt.Println("Error run cmd:", cmdOut, err)
	}

	return cmdOut, err
}

// copy from https://github.com/openshift/installer/blob/master/tests/bdd-smoke/utils/command_runner.go
// RunOSCommandWithArgs executes a command in the operating system with arguments
func RunOSCommandWithArgs(command string, arguments ...string) (string, error) {
	return RunOSCommandWithEnv(nil, "", command, arguments...)
}

// If Env is not set, the process inherits environment of the calling process.
func RunOSCommandWithEnv(env []string, workingDir string, command string, arguments ...string) (string, error) {
	cmd := exec.Command(command, arguments...)
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	if env != nil {
		cmd.Env = env
	}
	cmd.Stdin = strings.NewReader("")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	/*
		if err != nil {
			fmt.Println("Error: failed running the command \"" + command + "\" with the arguments:")
		}
	*/
	return outb.String() + errb.String(), err
}