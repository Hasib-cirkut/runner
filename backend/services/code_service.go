package services

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"runner-api/containers"
	"strings"
	"time"
)

func RunCodeInContainer(code string, lang string) (string, []string, error) {
	containerID, err := containers.CreateContainer(lang)
	var stderrLines []string

	fmt.Printf("container created: %v\n", containerID)

	if err != nil {
		return "", stderrLines, fmt.Errorf("failed to create container, %v", err)
	}

	defer func(containerID string) {
		err := containers.RemoveContainer(containerID)
		if err != nil {
			fmt.Printf("failed to remove container: %v\n", err)
		}
	}(containerID)

	time.Sleep(1000 * time.Millisecond)

	checkCmd := exec.Command("docker", "ps", "--filter", "id="+containerID, "--format", "{{.ID}}")

	output, err := checkCmd.Output()

	fmt.Printf("%v\n", output)

	if err != nil || len(output) == 0 {
		return "", stderrLines, fmt.Errorf("container %s is not running", containerID)
	}

	filename, err := saveCodeToFile(code, lang)

	fmt.Printf("filename generated: %v\n", filename)

	if err != nil {
		return "", stderrLines, fmt.Errorf("failed to save code: %v", err)
	}

	err = containers.CreateFileInContainer(containerID, filename)

	if err != nil {
		return "", stderrLines, fmt.Errorf("failed to copy code to container: %v", err)
	}

	stdout, stderr, err := containers.ExecuteInContainer(containerID, lang, filepath.Base(filename))

	if stderr != "" {
		stderrLines = strings.Split(stderr, "\n")
	}

	return stdout, stderrLines, err
}

func saveCodeToFile(code string, lang string) (string, error) {
	ext := getFileExtension(lang)

	if ext == "" {
		return "", errors.New("unsupported language")
	}

	filename := fmt.Sprintf("/tmp/code_%d.%s", time.Now().UnixNano(), ext)

	cmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' > %s", strings.ReplaceAll(code, "'", "'\\'"), filename))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to save file: %v, stderr: %s", err, stderr.String())
	}

	return filename, nil
}

func getFileExtension(lang string) string {
	extensions := map[string]string{
		"javascript": "js",
	}

	return extensions[strings.ToLower(lang)]
}
