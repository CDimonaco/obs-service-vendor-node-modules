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
	sourceUnpackDest := path.Join(opts.OutputDir, "source_dest")

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

	// cwd into the extracted source directory
	// it's the first and only directory in the source_dest directory
	dirs, err := os.ReadDir(sourceUnpackDest)
	if err != nil {
		logger.Error("could not find read the source_dest directory after source extraction", "error", err)
		os.Exit(1)
	}

	if len(dirs) != 1 {
		logger.Error("more than one directory in source_dest folder, fatal error occured")
		os.Exit(1)
	}

	npmCwd := path.Join(sourceUnpackDest, dirs[0].Name())
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

	logger.Info("compressing node modules")

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

	err = os.RemoveAll(sourceUnpackDest)
	if err != nil {
		logger.Error("error during cleanup, please clean manually source_dest folder")
	}
}
