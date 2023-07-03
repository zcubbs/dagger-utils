package types

import (
	"context"
	"dagger.io/dagger"
	"os"
)

var (
	syftImage      = "anchore/syft:latest"
	grypeImage     = "anchore/grype:latest"
	binDir         = "bin"
	buildPackImage = "paketobuildpacks/builder:base"
)

type Options struct {
	Src          *dagger.Directory
	Ctx          context.Context
	DaggerClient *dagger.Client
}

type GoOptions struct {
	Options
	BuildImg    string
	LintImg     string
	LintArgs    []string
	LintTimeout string
	BinDir      string
	BinName     string
	BuildArgs   []string
	TestArgs    []string
}

type MavenOptions struct {
	Options
	BuildImg    string
	MvnArgs     []string
	CacheKey    string
	EnableCache bool
}

type NpmOptions struct {
	Options
	BuildImg string
}

type ScanOptions struct {
	Options
	GrypeImg string
	SyftImg  string
	BinDir   string
}

type Scanner struct {
	ScanOptions ScanOptions
}

type DockerConfig struct {
	Auths map[string]AuthConfig `json:"auths"`
}

type AuthConfig struct {
	Auth  string `json:"auth"`
	Email string `json:"email"`
}

type RegistryInfo struct {
	RegistryServer   string
	RegistryUsername string
	RegistryPassword string
	RegistryEmail    string
}

type ImageBuilderOptions struct {
	BuildImg string
}

func SetDefaults(options *Options) {
	if options == nil {
		options = &Options{}
	}

	if options.Ctx == nil {
		options.Ctx = context.Background()
	}

	if options.DaggerClient == nil {
		client, err := dagger.Connect(options.Ctx, dagger.WithLogOutput(os.Stdout))
		if err != nil {
			panic(err)
		}
		options.DaggerClient = client
	}

	if options.Src == nil {
		options.Src = options.DaggerClient.Host().Directory(".")
	}
}

func SetupScanDefaults(options *ScanOptions) {
	if options == nil {
		options = &ScanOptions{}
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

func SetupImageBuilderDefaults(options *ImageBuilderOptions) {
	if options == nil {
		options = &ImageBuilderOptions{}
	}

	if options.BuildImg == "" {
		options.BuildImg = buildPackImage
	}
}
