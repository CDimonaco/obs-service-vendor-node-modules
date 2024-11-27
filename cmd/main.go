package main

import (
	"log/slog"
	"os"
	"path"

	"github.com/cdimonaco/obs-service-vendor_node_modules/internal/nodemodules"
	"github.com/jessevdk/go-flags"
)

type opts struct {
	SrcArchive  string `long:"srcarchive" description:"The source archive name" required:"true"`
	SubDir      string `long:"subdir" description:"Subdirectory relative to the source code root"`
	Compression string `long:"compression" description:"Compression method used by source archive" default:"gz" choice:"gz" choice:"zst"`
	ArchiveName string `long:"vendor-archive-name" description:"node_modules archive name, will be compress with the same method as source archive" default:"node_vendor.tar.gz"`
	OutputDir   string `long:"outdir" description:"Archive output directory" required:"true"`
}

func main() {
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

	if opts.SubDir != "" {
		cwd = path.Join(cwd, opts.SubDir)
	}

	logger.Info(
		"starting obs-service-vendor_node_modules",
		"archive-name",
		opts.ArchiveName,
		"output-dir",
		opts.OutputDir,
		"workdir",
		cwd,
	)

	err = nodemodules.Install(cwd)
	if err != nil {
		logger.Error("could not install the node dependencies", "error", err)
		os.Exit(1)
	}

	logger.Info("node dependencies installed")

	err = nodemodules.Compress(opts.OutputDir, cwd, opts.ArchiveName)
	if err != nil {
		logger.Error("could not compress node_modules archive", "error", err)
		os.Exit(1)
	}

	logger.Info("node_modules archive created", "archive", path.Join(opts.OutputDir, opts.ArchiveName))
}
