package containers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CreateContainer(language string) (string, error) {
	image := getDockerImage(language)

	fmt.Printf("choose Image: %s\n", image)

	cmd := exec.Command("docker", "run", "--detach",
		"--network=none",                   // No network access
		"--cap-drop=ALL",                   // Drop all capabilities
		"--memory=100m",                    // Limit memory
		"--cpus=0.5",                       // Limit CPU usage
		"--pids-limit=50",                  // Limit the number of processes
		"--read-only",                      // Read-only filesystem
		"--tmpfs", "/tmp:exec,rw,size=64m", // Writable /tmp directory
		image,
		"tail", "-f", "/dev/null")

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to create container: %v, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

func RemoveContainer(containerID string) error {
	cmd := exec.Command("docker", "rm", "-f", containerID)

	return cmd.Run()
}

func getDockerImage(language string) string {
	images := map[string]string{
		"javascript": "node:alpine",
	}

	image, exists := images[strings.ToLower(language)]

	if !exists {
		return "alpine:latest"
	}

	return image
}

func CreateFileInContainer(containerId string, filename string) error {
	content, err := os.ReadFile(filename)

	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	baseFileName := filepath.Base(filename)

	// create the file directly
	echoCmd := exec.Command(
		"docker", "exec", containerId,
		"sh", "-c",
		fmt.Sprintf("cat > /tmp/%s << 'EOF'\n%s\nEOF", baseFileName, string(content)),
	)

	var stderr bytes.Buffer

	echoCmd.Stderr = &stderr

	if err := echoCmd.Run(); err != nil {
		return fmt.Errorf("failed to create file in container: %v, stderr: %s", err, stderr.String())
	}

	return nil
}

func ExecuteInContainer(containerID string, language string, filename string) (string, string, error) {
	execCmd := getExecutionCommand(language, filename)

	cmd := exec.Command("docker", "exec", containerID, "sh", "-c", execCmd)

	fmt.Printf("Executing command: %v\n", cmd)

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", "", err
	}

	return stdout.String(), stderr.String(), nil
}

func getExecutionCommand(language string, filename string) string {
	commands := map[string]string{
		"javascript": "node /tmp/" + filename,
	}

	command, exists := commands[strings.ToLower(language)]

	if !exists {
		return "echo 'Unsupported language'"
	}

	return "timeout 10s " + command
}
