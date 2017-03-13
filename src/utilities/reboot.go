package utilities

import (
	"os/exec"
)

func Reboot(localId string) error {
	cmdName := "gnome-terminal"
	cmdArgs := []string{"-x", "sh", "-c", "go run main.go -id=" + localId}
	return exec.Command(cmdName, cmdArgs...).Run()
}
