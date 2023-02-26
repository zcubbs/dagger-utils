package golang

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

func (b Builder) GoTest(options types.GoOptions) error {
	setupDefaults(&options)

	testCmd := []string{"go"}
	testCmd = append(testCmd, options.TestArgs...)

	_, err := b.Options.DaggerClient.Container().
		From(options.BuildImg).
		WithMountedDirectory("/src", b.Options.Src).
		WithWorkdir("/src").
		WithExec(testCmd).
		ExitCode(b.Options.Ctx)
	if err != nil {
		return err
	}

	return nil
}

func (b Builder) GoBuild(options types.GoOptions) error {
	setupDefaults(&options)

	buildCmd := []string{"go"}
	buildCmd = append(buildCmd, options.BuildArgs...)

	builder := b.Options.DaggerClient.
		Container().
		From(options.BuildImg).
		WithMountedDirectory("/src", b.Options.Src).
		WithWorkdir("/src").
		WithExec(buildCmd)

	_, err := builder.File(options.BinName).
		Export(b.Options.Ctx, fmt.Sprintf("%s/%s", options.BinDir, options.BinName))
	if err != nil {
		return err
	}

	return nil
}

func setupDefaults(options *types.GoOptions) {
	if options.BuildImg == "" {
		options.BuildImg = "golang:1.20"
	}

	if options.BinDir == "" {
		options.BinDir = "bin"
	}

	if options.BinName == "" {
		options.BinName = "app"
	}

	if options.BuildArgs == nil {
		options.BuildArgs = []string{"build", "-o", options.BinName}
	}

	if options.TestArgs == nil {
		options.TestArgs = []string{"test", "-v"}
	}
}
