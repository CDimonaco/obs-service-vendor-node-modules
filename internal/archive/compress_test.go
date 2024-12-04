package archive_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"testing"

	"github.com/cdimonaco/obs-service-vendor_node_modules/internal/archive"
	"github.com/stretchr/testify/assert"
)

func TestGzFolderCompressionSuccess(t *testing.T) {
	ctx := context.TODO()
	fixturesFolder := "../../testfixtures/tocompress"
	outDir := t.TempDir()
	archiveName := "test.tar.gz"

	err := archive.CompressFolder(
		ctx,
		fixturesFolder,
		outDir,
		archiveName,
		archive.GZ,
	)
	assert.NoError(t, err)
	assert.NoError(t, checkIsTarGz(path.Join(outDir, archiveName)))
}

func TestZstFolderCompressionSuccess(t *testing.T) {
	ctx := context.TODO()
	fixturesFolder := "../../testfixtures/tocompress"
	outDir := t.TempDir()
	archiveName := "test.tar.zst"

	err := archive.CompressFolder(
		ctx,
		fixturesFolder,
		outDir,
		archiveName,
		archive.ZST,
	)
	assert.NoError(t, err)
	assert.NoError(t, checkIsTarZst(path.Join(outDir, archiveName)))
}

func TestCompressError(t *testing.T) {
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
	}

	for name, test := range tc {
		t.Run(name, func(t *testing.T) {
			err := archive.CompressFolder(context.TODO(), test.inputFolder, test.outDir, "test.tar.gz", archive.GZ)
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
