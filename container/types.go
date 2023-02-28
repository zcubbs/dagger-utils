package container

import "github.com/zcubbs/dagger-utils/types"

var (
	syftImage  = "anchore/syft:latest"
	grypeImage = "anchore/grype:latest"
	binDir     = "bin"
)

type Scanner struct {
	ScanOptions types.ScanOptions
}

func setupDefaults(options *types.ScanOptions) {
	if options == nil {
		options = &types.ScanOptions{}
	}

	if options.BinDir == "" {
		options.BinDir = binDir
	}

	if options.GrypeImg == "" {
		options.GrypeImg = grypeImage
	}

	if options.SyftImg == "" {
		options.SyftImg = syftImage
	}
}
