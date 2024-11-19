package nodemodules

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

const archiveName = "node_vendor.tar.gz"

func Compress(nodeModulesInput string, packDest string, workingDir string) (string, error) {
	err := folderPathValid(nodeModulesInput)
	if err != nil {
		return "", fmt.Errorf("error during the opening of node_modules directory: %w", err)

	}

	err = folderPathValid(packDest)
	if err != nil {
		return "", fmt.Errorf("destination path is not valid: %w", err)
	}

	fullDestPath := path.Join(packDest, archiveName)

	cmd := exec.Command("tar", "-zcvf", fullDestPath, nodeModulesInput)
	cmd.Dir = workingDir
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error during the compression of node_modules directory: %w", err)
	}

	return fullDestPath, nil
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
