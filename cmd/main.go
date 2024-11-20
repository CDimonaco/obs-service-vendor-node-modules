package main

import (
	"log/slog"
	"os"
	"path"

	"github.com/cdimonaco/obs-service-vendor-node-modules/internal/nodemodules"
	"github.com/jessevdk/go-flags"
)

type opts struct {
	ArchiveName string `long:"archive-name" description:"node_modules archive name" default:"node_vendor.tar.gz"`
	OutputDir   string `long:"outdir" description:"Archive output directory"`
	SubDir      string `long:"subdir" description:"Service working directory"`
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

	outputDir := cwd
	if opts.OutputDir != "" {
		outputDir = opts.OutputDir
	}

	logger.Info(
		"starting obs-service-vendor-node-modules",
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

	err = nodemodules.Compress(outputDir, cwd, opts.ArchiveName)
	if err != nil {
		logger.Error("could not compress node_modules archive", "error", err)
		os.Exit(1)
	}

	logger.Info("node_modules archive created", "archive", path.Join(outputDir, opts.ArchiveName))
}
