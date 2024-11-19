package nodemodules

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func Compress(
	packDest string,
	workingDir string,
	archiveName string,
) error {
	err := folderPathValid(path.Join(workingDir, "node_modules"))
	if err != nil {
		return fmt.Errorf("error during the opening of node_modules directory: %w", err)
	}

	err = folderPathValid(packDest)
	if err != nil {
		return fmt.Errorf("destination path is not valid: %w", err)
	}

	cmd := exec.Command("tar", "-zcf", archiveName, "node_modules")
	cmd.Dir = workingDir
	cmd.Stderr = os.Stderr
	_, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("error during the compression of node_modules directory: %w", err)
	}

	return nil
}

func folderPathValid(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	return nil
}
