package dagger_utils

import (
	"fmt"
	"github.com/zcubbs/dagger-utils/types"
)

type Builder struct {
	types.Options
}

func (b Builder) GoLint(lintImg string, lintTimeout string) error {
	_, err := b.Options.DaggerClient.Container().
		From(lintImg).
		WithMountedDirectory("/src", b.Options.Src).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run",
			fmt.Sprintf("--timeout=%s", lintTimeout), "-v"}).
		ExitCode(b.Options.Ctx)
	if err != nil {
		return err
	}

	return nil
}

func (b Builder) GoBuild(buildImg string, binDir string, binName string) error {
	builder := b.Options.DaggerClient.
		Container().
		From(buildImg).
		WithMountedDirectory("/src", b.Options.Src).
		WithWorkdir("/src").
		WithExec([]string{"go", "test"}).
		WithExec([]string{"go", "build", "-o", binName})

	_, err := builder.File(binName).Export(b.Options.Ctx, fmt.Sprintf("%s/%s", binDir, binName))
	if err != nil {
		return err
	}

	return nil
}
