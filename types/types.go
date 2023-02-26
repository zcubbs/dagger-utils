package types

import (
	"context"
	"dagger.io/dagger"
)

type Options struct {
	Src          *dagger.Directory
	Ctx          context.Context
	DaggerClient *dagger.Client
}

type GoOptions struct {
	Options
	BuildImg  string
	BinDir    string
	BinName   string
	BuildArgs []string
	TestArgs  []string
}

type NpmOptions struct {
	Options
	BuildImg string
}
