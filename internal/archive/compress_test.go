package archive_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"testing"

	"github.com/cdimonaco/obs-service-vendor_node_modules/internal/archive"
	"github.com/stretchr/testify/assert"
)

func TestCompressFolderSuccess(t *testing.T) {
	fixturesFolder := "../../testfixtures/tocompress"
	outDir := t.TempDir()
	archiveName := "test.tar.gz"

	err := archive.CompressFolder(fixturesFolder, outDir, archiveName)
	assert.NoError(t, err)
	assert.NoError(t, checkIsTar(path.Join(outDir, archiveName)))
}

func checkIsTar(path string) error {
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

	return fmt.Errorf("could not identify tar archive")
}
