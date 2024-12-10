package archive

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archives"
)

func handleDecompressedFiles(outputPath string) func(ctx context.Context, info archives.FileInfo) error {
	return func(ctx context.Context, info archives.FileInfo) error {
		// Normalize extracted paths
		archiveDestPath := filepath.Clean("/" + info.NameInArchive)
		archiveDestPath = strings.TrimPrefix(archiveDestPath, string(os.PathSeparator))
		realPath := filepath.Join(outputPath, archiveDestPath)

		if info.IsDir() {
			// Create directory
			err := os.MkdirAll(realPath, 0755)
			if err != nil {
				return err
			}
			return nil
		}

		r, err := info.Open()
		if err != nil {
			return fmt.Errorf("error opening archive file: %w", err)
		}
		defer r.Close()

		f, err := os.OpenFile(realPath, os.O_CREATE|os.O_WRONLY, info.Mode())
		if err != nil {
			return fmt.Errorf("error creating decompressed file: %w", err)
		}
		defer f.Close()

		if _, err := io.Copy(f, r); err != nil {
			return fmt.Errorf("error copying archive file to destination: %w", err)
		}

		return nil
	}
}

func decompressGZ(
	ctx context.Context,
	archivePath string,
	outputPath string,
) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	var compression archives.Gz
	rc, err := compression.OpenReader(f)
	if err != nil {
		return err
	}
	defer rc.Close()

	var archive archives.Tar
	return archive.Extract(ctx, rc, handleDecompressedFiles(outputPath))
}

func decompressZST(
	ctx context.Context,
	archivePath string,
	outputPath string,
) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	var compression archives.Zstd
	rc, err := compression.OpenReader(f)
	if err != nil {
		return err
	}
	defer rc.Close()

	var archive archives.Tar
	return archive.Extract(ctx, rc, handleDecompressedFiles(outputPath))
}

func DecompressArchive(
	ctx context.Context,
	archivePath string,
	outputPath string,
	archiveType SupportedArchive,
) error {
	if archiveType != GZ && archiveType != ZST {
		return fmt.Errorf("archive type %s, not supported, use gz or zst", archiveType)
	}

	err := folderPathValid(outputPath)
	if err != nil {
		return fmt.Errorf("output folder to decompress archive is invalid: %w", err)
	}

	if archiveType == GZ {
		return decompressGZ(ctx, archivePath, outputPath)
	}

	return decompressZST(ctx, archivePath, outputPath)
}
