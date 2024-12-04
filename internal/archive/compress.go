package archive

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mholt/archives"
)

type SupportedArchive string

const (
	GZ  SupportedArchive = "gz"
	ZST SupportedArchive = "zst"
)

func compressWithZst(
	ctx context.Context,
	absInputPath string,
	archiveFile *os.File,
) error {
	archiveFiles, err := buildArchiveFileStructure(absInputPath)
	if err != nil {
		return err
	}

	files, err := archives.FilesFromDisk(ctx, nil, archiveFiles)
	if err != nil {
		return err
	}

	format := archives.CompressedArchive{
		Compression: archives.Zstd{},
		Archival:    archives.Tar{},
	}

	err = format.Archive(ctx, archiveFile, files)
	if err != nil {
		return err
	}

	return nil
}

func compressWithGz(
	ctx context.Context,
	absInputPath string,
	archiveFile *os.File,
) error {
	archiveFiles, err := buildArchiveFileStructure(absInputPath)
	if err != nil {
		return err
	}

	files, err := archives.FilesFromDisk(ctx, nil, archiveFiles)
	if err != nil {
		return err
	}

	format := archives.CompressedArchive{
		Compression: archives.Gz{},
		Archival:    archives.Tar{},
	}

	err = format.Archive(ctx, archiveFile, files)
	if err != nil {
		return err
	}

	return nil
}

func CompressFolder(
	ctx context.Context,
	inputPath string,
	outputPath string,
	archiveName string,
	archiveType SupportedArchive,
) error {
	if archiveType != GZ && archiveType != ZST {
		return fmt.Errorf("archive type %s, not supported, use gz or zst", archiveType)
	}

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

	if archiveType == GZ {
		return compressWithGz(ctx, absInputPath, out)
	}

	return compressWithZst(ctx, absInputPath, out)
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

func buildArchiveFileStructure(
	rootPath string,
) (map[string]string, error) {
	archiveFiles := make(map[string]string)
	archiveBaseFolder := path.Base(rootPath)

	err := filepath.Walk(rootPath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		fileName := strings.TrimPrefix(strings.Replace(file, rootPath, "", -1), string(filepath.Separator))

		archiveFiles[file] = path.Join(archiveBaseFolder, fileName)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return archiveFiles, nil
}
