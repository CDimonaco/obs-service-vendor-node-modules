package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type SupportedArchive string

const (
	GZ  SupportedArchive = "gz"
	ZST SupportedArchive = "zst"
)

func CompressFolder(
	inputPath string,
	outputPath string,
	archiveName string,
) error {
	err := folderPathValid(inputPath)
	if err != nil {
		return fmt.Errorf("input folder to compress is invalid: %w", err)
	}

	err = folderPathValid(outputPath)
	if err != nil {
		return fmt.Errorf("output folder for compressed archive is invalid: %w", err)
	}

	absInputPath, err := filepath.Abs(stripTrailingSlashes(inputPath))
	if err != nil {
		return err
	}

	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		return err
	}

	out, err := os.Create(path.Join(absOutputPath, archiveName))
	if err != nil {
		return fmt.Errorf("error writing archive %s, %w", archiveName, err)
	}
	defer out.Close()

	gw := gzip.NewWriter(out)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	return filepath.Walk(absInputPath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(strings.Replace(file, absInputPath, "", -1), string(filepath.Separator))

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		f.Close()

		return nil
	})

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

func stripTrailingSlashes(path string) string {
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}

	return path
}
