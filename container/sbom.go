package container

import (
	"dagger.io/dagger"
	"fmt"
	"github.com/zcubbs/dagger-utils/types"
)

// GenerateSBOM generates a software bill of materials for the container image
func (s *Scanner) GenerateSBOM(targetImage string) (*dagger.FileID, error) {
	types.SetDefaults(&s.ScanOptions.Options)
	types.SetupScanDefaults(&s.ScanOptions)

	client := s.ScanOptions.DaggerClient
	ctx := s.ScanOptions.Options.Ctx
	bom := client.Container().From(s.ScanOptions.SyftImg)

	buildCmd := []string{targetImage, "--scope", "all-layers", "-v", "-o", "spdx-json", "--file", "sbom.json"}
	bom = bom.WithWorkdir("/src").
		WithExec(buildCmd)

	fileID, err := bom.File("./sbom.json").ID(ctx)
	if err != nil {
		return nil, err
	}

	_, err = bom.File("./sbom.json").
		Export(ctx, fmt.Sprintf("%s/sbom.json", s.ScanOptions.BinDir))
	if err != nil {
		return nil, err
	}

	return &fileID, nil
}
