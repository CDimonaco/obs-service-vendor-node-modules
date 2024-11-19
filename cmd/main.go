package main

import (
	"log/slog"
	"os"

	"github.com/cdimonaco/obs-service-vendor-node-modules/internal/nodemodules"
	"github.com/jessevdk/go-flags"
)

type opts struct {
	NodeModulesDir string `long:"node-modules-folder" description:"Input node_modules directory" required:"true"`
	OutputDir      string `long:"output-dir" description:"Archive output directory" required:"true"`
	WorkingDir     string `long:"workdir" description:"Service working directory" required:"false"`
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

	if opts.WorkingDir != "" {
		cwd = opts.WorkingDir
	}

	logger.Info(
		"starting obs-service-vendor-node-modules",
		"node-modules-folder",
		opts.NodeModulesDir,
		"output-dir",
		opts.OutputDir,
		"workdir",
		cwd,
	)

	archivePath, err := nodemodules.Compress(opts.NodeModulesDir, opts.OutputDir, cwd)
	if err != nil {
		logger.Error("could not compress node_modules archive", "error", err)
		os.Exit(1)
	}

	logger.Info("node_modules archive created", "path", archivePath)
}
