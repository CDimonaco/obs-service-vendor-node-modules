package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/cdimonaco/obs-service-vendor_node_modules/internal/archive"
	"github.com/cdimonaco/obs-service-vendor_node_modules/internal/nodemodules"
	"github.com/jessevdk/go-flags"
)

type opts struct {
	SrcArchive  string `long:"srcarchive" description:"The source archive name" required:"true"`
	SubDir      string `long:"subdir" description:"Subdirectory relative to the source code root"`
	Compression string `long:"compression" description:"Compression method used by source archive" default:"gz" choice:"gz" choice:"zst"`
	ArchiveName string `long:"vendor-archive-name" description:"node_modules archive name, will be compress with the same method as source archive" default:"node_vendor"`
	OutputDir   string `long:"outdir" description:"Archive output directory" required:"true"`
}

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var opts opts
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	vendorArchiveName := fmt.Sprintf("%s.tar.%s", opts.ArchiveName, opts.Compression)

	logger.Info(
		"starting obs-service-vendor_node_modules",
		"srcarchive",
		opts.SrcArchive,
		"compression",
		opts.Compression,
		"subdir",
		opts.SubDir,
		"outdir",
		opts.OutputDir,
		"archive-name",
		vendorArchiveName,
	)

	logger.Info("unpacking source archive", "name", opts.SrcArchive)
	sourceUnpackDest := path.Join(cwd, "source_dest")

	err = os.MkdirAll(sourceUnpackDest, 0755)
	if err != nil {
		logger.Error("error during source decompress", "error", err)
		os.Exit(1)
	}

	err = archive.DecompressArchive(
		ctx,
		opts.SrcArchive,
		sourceUnpackDest,
		archive.SupportedArchive(opts.Compression),
	)
	if err != nil {
		logger.Error("error during source decompress", "error", err)
		os.Exit(1)
	}

	npmCwd := sourceUnpackDest
	if opts.SubDir != "" {
		npmCwd = path.Join(npmCwd, opts.SubDir)
	}

	logger.Info("installing node dependencies", "subdir", opts.SubDir)

	err = nodemodules.Install(npmCwd)
	if err != nil {
		logger.Error("could not install the node dependencies", "error", err)
		os.Exit(1)
	}

	logger.Info("node dependencies installed", "subdir", opts.SubDir)

	err = archive.CompressFolder(
		ctx,
		path.Join(npmCwd, "node_modules"),
		opts.OutputDir,
		vendorArchiveName,
		archive.SupportedArchive(opts.Compression),
	)
	if err != nil {
		logger.Error("could not compress node_modules archive", "error", err)
		os.Exit(1)
	}

	logger.Info("node_modules archive created", "archive", path.Join(opts.OutputDir, vendorArchiveName))
}
