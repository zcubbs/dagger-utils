package types

import (
	"context"
	"dagger.io/dagger"
	"os"
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
