package archive_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cdimonaco/obs-service-vendor_node_modules/internal/archive"
	"github.com/stretchr/testify/assert"
)

func TestGzFolderCompressionSuccess(t *testing.T) {
	ctx := context.TODO()
	fixturesFolder := "../../testfixtures/tocompress"
	outDir := t.TempDir()
	archiveName := "test.tar.gz"
	fullArchivePath := path.Join(outDir, archiveName)

	err := archive.CompressFolder(
		ctx,
		fixturesFolder,
		outDir,
		archiveName,
		archive.GZ,
	)
	assert.NoError(t, err)
	assert.NoError(t, checkIsTarGz(fullArchivePath))
	assert.NoError(t, uncompressArchiveWithCli(fullArchivePath, outDir))
	assert.NoError(t, checkCompressedArchiveStructure(fixturesFolder, path.Join(outDir, "tocompress")))
}

func TestZstFolderCompressionSuccess(t *testing.T) {
	ctx := context.TODO()
	fixturesFolder := "../../testfixtures/tocompress"
	outDir := t.TempDir()
	archiveName := "test.tar.zst"
	fullArchivePath := path.Join(outDir, archiveName)

	err := archive.CompressFolder(
		ctx,
		fixturesFolder,
		outDir,
		archiveName,
		archive.ZST,
	)
	assert.NoError(t, err)
	assert.NoError(t, checkIsTarZst(path.Join(outDir, archiveName)))
	assert.NoError(t, uncompressArchiveWithCli(fullArchivePath, outDir))
	assert.NoError(t, checkCompressedArchiveStructure(fixturesFolder, path.Join(outDir, "tocompress")))
}

func TestCompressErrors(t *testing.T) {
	tc := map[string]struct {
		inputFolder   string
		outDir        string
		errorContains string
		archiveType   archive.SupportedArchive
	}{
		"should return error when input folder is not valid and the archive is compressed with gzip": {
			inputFolder:   "../invalid",
			outDir:        t.TempDir(),
			errorContains: "input folder to compress is invalid",
			archiveType:   archive.GZ,
		},
		"should return error when input folder is not valid and the archive is compressed with zst": {
			inputFolder:   "../invalid",
			outDir:        t.TempDir(),
			errorContains: "input folder to compress is invalid",
			archiveType:   archive.ZST,
		},
		"should return error when output folder is not valid and the archive is compressed with gzip": {
			inputFolder:   "../../testfixtures/tocompress",
			outDir:        "../invalid",
			errorContains: "output folder for compressed archive is invalid",
			archiveType:   archive.GZ,
		},
		"should return error when output folder is not valid and the archive is compressed with zst": {
			inputFolder:   "../../testfixtures/tocompress",
			outDir:        "../invalid",
			errorContains: "output folder for compressed archive is invalid",
			archiveType:   archive.ZST,
		},
		"should return error when archive type is not supported": {
			inputFolder:   "../../testfixtures/tocompress",
			outDir:        t.TempDir(),
			errorContains: "archive type invalid, not supported, use gz or zst",
			archiveType:   archive.SupportedArchive("invalid"),
		},
	}

	for name, test := range tc {
		t.Run(name, func(t *testing.T) {
			err := archive.CompressFolder(context.TODO(), test.inputFolder, test.outDir, "test.tar.gz", test.archiveType)
			assert.ErrorContains(t, err, test.errorContains)
		})
	}
}

func checkIsTarGz(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	gzHeader := []byte{0x1f, 0x8b}
	buf := make([]byte, len(gzHeader))

	_, err = io.ReadFull(f, buf)
	if err != nil {
		return err
	}

	if bytes.Equal(buf, gzHeader) {
		return nil
	}

	return fmt.Errorf("could not identify tar archive as compressed with gzip")
}

func checkIsTarZst(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	zstHeader := []byte{0x28, 0xB5, 0x2F, 0xFD}
	buf := make([]byte, len(zstHeader))

	_, err = io.ReadFull(f, buf)
	if err != nil {
		return err
	}

	if bytes.Equal(buf, zstHeader) {
		return nil
	}

	return fmt.Errorf("could not identify tar archive as compressed with zst")
}

func uncompressArchiveWithCli(tarPath string, workDir string) error {
	cmd := exec.Command("tar", "-xvf", tarPath)
	cmd.Stderr = os.Stderr
	cmd.Dir = workDir
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func checkCompressedArchiveStructure(expectedFolderRoot string, uncompressedArchiveRoot string) error {
	return filepath.Walk(expectedFolderRoot, func(file string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		fileName := strings.TrimPrefix(strings.Replace(file, expectedFolderRoot, "", -1), string(filepath.Separator))
		// open expected file
		expectedFile, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		dpath := path.Join(uncompressedArchiveRoot, fileName)
		// decompressed file
		decompressedFile, err := os.ReadFile(dpath)
		if err != nil {
			return err
		}

		if !bytes.Equal(expectedFile, decompressedFile) {
			return fmt.Errorf("file %s not equal to decompressed file", file)
		}

		return nil

	})
}
