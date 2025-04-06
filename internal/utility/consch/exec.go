package consch

import (
	"bytes"
	"fmt"
	"os/exec"
)

func execPython(node *conNode, containerNumber int) (string, error) {
	defer PyDoneExec(containerNumber)

	containerName := fmt.Sprintf("rcc-python_runner-%v", containerNumber)

	code := node.Code

	cmd := exec.Command("docker", "exec", containerName, "python3", "-c", code)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return stderr.String(), fmt.Errorf("error: %v", err)
	}

	return stdout.String(), nil
}
