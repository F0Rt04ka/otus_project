//go:build !race

package test

import (
	"bytes"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"
)

func Test_runDaemonAndClient(t *testing.T) {
	o, err := exec.Command("pwd").Output()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	cmdDir := strings.Replace(strings.Trim(string(o), "\n "), "/test", "", 1)

	daemonCmd := exec.Command("./daemon.out", "-server-port", "44040")
	daemonCmd.Dir = cmdDir
	err = daemonCmd.Start()
	if err != nil {
		t.Fatalf("Failed to run daemon: %v", err)
	}
	defer daemonCmd.Process.Kill()

	// Запускаем клиент и проверяем его работу
	clientCmd := exec.Command("./client.out", "-N", "3", "-server-port", "44040")
	clientCmd.Dir = cmdDir
	stdout := bytes.Buffer{}
	clientCmd.Stdout = &stdout

	err = clientCmd.Start()
	if err != nil {
		t.Fatalf("Failed to run client: %v", err)
	}

	time.Sleep(5 * time.Second) // Ждем 5 секунд, чтобы клиент успел выполнить запросы

	err = clientCmd.Process.Kill()
	if err != nil {
		t.Fatalf("Failed to kill client process: %v", err)
	}
	stringOutput := stdout.String()
	t.Logf("Client output:\n%s", stringOutput)

	if !validateOutput(stringOutput) {
		t.Fatalf("Invalid output from client:\n%s", stringOutput)
	}
}

func validateOutput(output string) bool {
	re := regexp.MustCompile(
		`(?m)^CPU Usage: (\d+\.\d+%) (\d+\.\d+%) (\d+\.\d+%)\s*$` +
			`\nLoad Average: (\d+\.\d+) (\d+\.\d+) (\d+\.\d+)\s*$` +
			`\nDisk Load: (\d+\.\d+) TPS (\d+\.\d+ KB/s) (\d+\.\d+ KB/s)\s*$` +
			`\nFilesystem Usage:\s*$` +
			`(?:\n\s+[a-zA-Z0-9\-\/]+: (\d+\.\d+) used MB, (\d+\.\d+%); (\d+\.\d+) used inodes, (\d+\.\d+%)\s*)+`)

	return re.MatchString(output)
}
