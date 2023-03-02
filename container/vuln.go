package container

import (
	"fmt"
	"github.com/zcubbs/dagger-utils/types"
)

type Scanner struct {
	types.Options
	types.ScanOptions
}

// GenerateVulnReport scans the SBOM for vulnerabilities
func (s *Scanner) GenerateVulnReport(targetImage string) error {
	types.SetDefaults(&s.ScanOptions.Options)
	types.SetupScanDefaults(&s.ScanOptions)

	ctx := s.ScanOptions.Options.Ctx

	sbomFileID, err := s.GenerateSBOM(targetImage)
	if err != nil {
		return err
	}

	file := s.ScanOptions.DaggerClient.File(*sbomFileID)

	scanner := s.ScanOptions.DaggerClient.
		Container().
		From(s.ScanOptions.GrypeImg).
		WithMountedFile("/sbom.json", file)

	_, err = scanner.WithExec([]string{"sbom:/sbom.json", "-v", "-o", "json", "--file", "vuln.json"}).
		ExitCode(ctx)
	if err != nil {
		return err
	}

	// Export the vulnerability report
	_, err = scanner.File("./vuln.json").
		Export(ctx, fmt.Sprintf("%s/vuln.json", s.ScanOptions.BinDir))
	if err != nil {
		return err
	}

	return nil
}
