package container

import (
	"encoding/json"
	"fmt"
	"github.com/zcubbs/dagger-utils/types"
	"os"
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

	_, err = scanner.
		WithWorkdir("/src").
		WithExec([]string{"sbom:/sbom.json", "-v", "-o", "json", "--file", "vuln.json"}).
		File("vuln.json").
		Export(ctx, fmt.Sprintf("%s/vuln.json", s.ScanOptions.BinDir))
	if err != nil {
		return err
	}

	//fmt.Printf("Vulnerability report saved to %s/vuln.json", dir)
	//// Export the vulnerability report
	//_, err = scanner.
	//	WithWorkdir("/src").
	//	File(fmt.Sprintf("%s/vuln.json", dir)).
	//	Export(ctx, fmt.Sprintf("%s/vuln.json", s.ScanOptions.BinDir))
	//if err != nil {
	//	return err
	//}

	return nil
}

// ParseVulnForSeverityLevels parses the vuln report for a specific severity level and returns count of each level
func (s *Scanner) ParseVulnForSeverityLevels(vulns []Vulnerability) (map[string]int, int) {
	levels := make(map[string]int, 0)
	fixes := 0
	for _, vuln := range vulns {
		if levels[vuln.Severity] == 0 {
			levels[vuln.Severity] = 1
		} else {
			levels[vuln.Severity]++
		}
		if len(vuln.Fix.Versions) > 0 {
			fixes++
		}
	}

	return levels, fixes
}

// ScanVuln scans the vuln report for vulnerabilities
func (s *Scanner) ScanVuln() ([]Vulnerability, error) {
	vulnJSON, err := os.ReadFile("./bin/vuln.json")
	if err != nil {
		return nil, err
	}
	vulns := make([]Vulnerability, 0)
	doc := &Document{}
	err = json.Unmarshal(vulnJSON, &doc)
	if err != nil {
		return nil, err
	}

	for _, match := range doc.Matches {
		if match.Vulnerability.ID != "" {
			vulns = append(vulns, match.Vulnerability)
		}
	}

	return vulns, nil
}

type Vulnerability struct {
	VulnerabilityMetadata
	Fix        Fix        `json:"fix"`
	Advisories []Advisory `json:"advisories"`
}

type Fix struct {
	Versions []string `json:"versions"`
	State    string   `json:"state"`
}

type Advisory struct {
	ID   string `json:"id"`
	Link string `json:"link"`
}

type Coordinates struct {
	RealPath     string `json:"path" cyclonedx:"path"`                 // The path where all path ancestors have no hardlinks / symlinks
	FileSystemID string `json:"layerID,omitempty" cyclonedx:"layerID"` // An ID representing the filesystem. For container images, this is a layer digest. For directories or a root filesystem, this is blank.
}

type VulnerabilityMetadata struct {
	ID          string   `json:"id"`
	DataSource  string   `json:"dataSource"`
	Namespace   string   `json:"namespace,omitempty"`
	Severity    string   `json:"severity,omitempty"`
	URLs        []string `json:"urls"`
	Description string   `json:"description,omitempty"`
	Cvss        []Cvss   `json:"cvss"`
}

type Cvss struct {
	Version        string      `json:"version"`
	Vector         string      `json:"vector"`
	Metrics        CvssMetrics `json:"metrics"`
	VendorMetadata interface{} `json:"vendorMetadata"`
}

type CvssMetrics struct {
	BaseScore           float64  `json:"baseScore"`
	ExploitabilityScore *float64 `json:"exploitabilityScore,omitempty"`
	ImpactScore         *float64 `json:"impactScore,omitempty"`
}

type Document struct {
	Matches        []Match        `json:"matches"`
	IgnoredMatches []IgnoredMatch `json:"ignoredMatches,omitempty"`
	Source         *source        `json:"source"`
	Distro         distribution   `json:"distro"`
	Descriptor     descriptor     `json:"descriptor"`
}

// Match is a single item for the JSON array reported
type Match struct {
	Vulnerability          Vulnerability           `json:"vulnerability"`
	RelatedVulnerabilities []VulnerabilityMetadata `json:"relatedVulnerabilities"`
	MatchDetails           []MatchDetails          `json:"matchDetails"`
	Artifact               Package                 `json:"artifact"`
}

// MatchDetails contains all data that indicates how the result match was found
type MatchDetails struct {
	Type       string      `json:"type"`
	Matcher    string      `json:"matcher"`
	SearchedBy interface{} `json:"searchedBy"`
	Found      interface{} `json:"found"`
}

type Package struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Type         Type              `json:"type"`
	Locations    []Coordinates     `json:"locations"`
	Language     Language          `json:"language"`
	Licenses     []string          `json:"licenses"`
	CPEs         []string          `json:"cpes"`
	PURL         string            `json:"purl"`
	Upstreams    []UpstreamPackage `json:"upstreams"`
	MetadataType MetadataType      `json:"metadataType,omitempty"`
	Metadata     interface{}       `json:"metadata,omitempty"`
}

type UpstreamPackage struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type IgnoredMatch struct {
	Match
	AppliedIgnoreRules []IgnoreRule `json:"appliedIgnoreRules"`
}

type IgnoreRule struct {
	Vulnerability string             `json:"vulnerability,omitempty"`
	FixState      string             `json:"fix-state,omitempty"`
	Package       *IgnoreRulePackage `json:"package,omitempty"`
}

type IgnoreRulePackage struct {
	Name     string `json:"name,omitempty"`
	Version  string `json:"version,omitempty"`
	Type     string `json:"type,omitempty"`
	Location string `json:"location,omitempty"`
}

type source struct {
	Type   string      `json:"type"`
	Target interface{} `json:"target"`
}

// distribution provides information about a detected Linux distribution.
type distribution struct {
	Name    string   `json:"name"`    // Name of the Linux distribution
	Version string   `json:"version"` // Version of the Linux distribution (major or major.minor version)
	IDLike  []string `json:"idLike"`  // the ID_LIKE field found within the /etc/os-release file
}

type descriptor struct {
	Name                  string      `json:"name"`
	Version               string      `json:"version"`
	Configuration         interface{} `json:"configuration,omitempty"`
	VulnerabilityDBStatus interface{} `json:"db,omitempty"`
}

type Type string

type Language string

type MetadataType string
