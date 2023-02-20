package builder

import (
	"dagger-utils/types"
	"fmt"
)

type Builder struct {
	types.Options
}

func (b Builder) goLint(lintImg string, lintTimeout string) {
	_, err := b.Options.DaggerClient.Container().
		From(lintImg).
		WithMountedDirectory("/src", b.Options.Src).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run",
			fmt.Sprintf("--timeout=%s", lintTimeout), "-v"}).
		ExitCode(b.Options.Ctx)
	if err != nil {
		panic(err)
	}
}

func (b Builder) goBuild(buildImg string, binDir string, binName string) {
	builder := b.Options.DaggerClient.
		Container().
		From(buildImg).
		WithMountedDirectory("/src", b.Options.Src).
		WithWorkdir("/src").
		WithExec([]string{"go", "test"}).
		WithExec([]string{"go", "build", "-o", binName})

	_, err := builder.File(binName).Export(b.Options.Ctx, fmt.Sprintf("%s/%s", binDir, binName))
	if err != nil {
		panic(err)
	}
}
