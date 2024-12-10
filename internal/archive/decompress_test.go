package archive_test

import (
	"context"
	"path"
	"testing"

	"github.com/cdimonaco/obs-service-vendor_node_modules/internal/archive"
	"github.com/stretchr/testify/assert"
)

func TestGzDecompressionSuccess(t *testing.T) {
	ctx := context.TODO()
	input := "../../testfixtures/compressed.tar.gz"
	fixturesFolder := "../../testfixtures/tocompress"
	outDir := t.TempDir()

	err := archive.DecompressArchive(
		ctx,
		input,
		outDir,
		archive.GZ,
	)

	assert.NoError(t, err)
	assert.NoError(t, checkCompressedArchiveStructure(fixturesFolder, path.Join(outDir, "tocompress")))
}

func TestZstDecompressionSuccess(t *testing.T) {
	ctx := context.TODO()
	input := "../../testfixtures/compressed.tar.zst"
	fixturesFolder := "../../testfixtures/tocompress"
	outDir := t.TempDir()

	err := archive.DecompressArchive(
		ctx,
		input,
		outDir,
		archive.ZST,
	)

	assert.NoError(t, err)
	assert.NoError(t, checkCompressedArchiveStructure(fixturesFolder, path.Join(outDir, "tocompress")))
}

func TestDecompressionErrors(t *testing.T) {
	tc := map[string]struct {
		inputArchive  string
		outDir        string
		errorContains string
		archiveType   archive.SupportedArchive
	}{
		"should return error when output folder is not valid and the archive is compressed with gzip": {
			inputArchive:  "../../testfixtures/compressed.tar.gz",
			outDir:        "../invalid",
			errorContains: "output folder to decompress archive is invalid",
			archiveType:   archive.GZ,
		},
		"should return error when output folder is not valid and the archive is compressed with zst": {
			inputArchive:  "../../testfixtures/compressed.tar.zst",
			outDir:        "../invalid",
			errorContains: "output folder to decompress archive is invalid",
			archiveType:   archive.ZST,
		},
		"should return error when archive type is not supported": {
			inputArchive:  "../../testfixtures/compressed.tar.zst",
			outDir:        t.TempDir(),
			errorContains: "archive type invalid, not supported, use gz or zst",
			archiveType:   archive.SupportedArchive("invalid"),
		},
	}

	for name, test := range tc {
		t.Run(name, func(t *testing.T) {
			err := archive.DecompressArchive(context.TODO(), test.inputArchive, test.outDir, test.archiveType)
			assert.ErrorContains(t, err, test.errorContains)
		})
	}
}
