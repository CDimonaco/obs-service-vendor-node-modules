package nodemodules

import (
	"fmt"
	"os/exec"
)

func Install(workingDir string) error {
	cmd := exec.Command("npm", "install")
	cmd.Dir = workingDir
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error during npm install: %w", err)
	}

	return nil
}
