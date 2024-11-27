package source

import (
	"fmt"
	"os"
	"os/exec"
)

// OpenArchive open the source archive using tar command
// returns the path of the opened archive
func OpenArchive(archiveName string) (string, error) {
	cmd := exec.Command("tar", "-zcf", archiveName, "node_modules")
	cmd.Dir = workingDir
	cmd.Stderr = os.Stderr
	_, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("error during the compression of node_modules directory: %w", err)
	}
}
