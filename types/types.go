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
