package golang

import (
	"fmt"
	"github.com/zcubbs/dagger-utils/types"
)

type Builder struct {
	types.Options
	types.GoOptions
}

func (b Builder) GoLint() error {
	types.SetDefaults(&b.Options)
	setupDefaults(&b.GoOptions)

	lintCmd := []string{"golangci-lint"}
	lintCmd = append(lintCmd, b.LintArgs...)
	_, err := b.Options.DaggerClient.Container().
		From(b.LintImg).
		WithMountedDirectory("/src", b.Options.Src).
		WithWorkdir("/src").
		WithExec(lintCmd).
		ExitCode(b.Options.Ctx)
	if err != nil {
		return err
	}

	return nil
}

func (b Builder) GoTest() error {
	types.SetDefaults(&b.Options)
	setupDefaults(&b.GoOptions)

	testCmd := []string{"go"}
	testCmd = append(testCmd, b.TestArgs...)

	_, err := b.Options.DaggerClient.Container().
		From(b.BuildImg).
		WithMountedDirectory("/src", b.Options.Src).
		WithWorkdir("/src").
		WithExec(testCmd).
		ExitCode(b.Options.Ctx)
	if err != nil {
		return err
	}

	return nil
}

func (b Builder) GoBuild() error {
	types.SetDefaults(&b.Options)
	setupDefaults(&b.GoOptions)

	buildCmd := []string{"go"}
	buildCmd = append(buildCmd, b.BuildArgs...)

	builder := b.Options.DaggerClient.
		Container().
		From(b.BuildImg).
		WithMountedDirectory("/src", b.Options.Src).
		WithWorkdir("/src").
		WithExec(buildCmd)

	_, err := builder.File(b.BinName).
		Export(b.Options.Ctx, fmt.Sprintf("%s/%s", b.BinDir, b.BinName))
	if err != nil {
		return err
	}

	return nil
}

func setupDefaults(options *types.GoOptions) {
	types.SetDefaults(&options.Options)

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

	if options.LintImg == "" {
		options.LintImg = "golangci/golangci-lint:v1.51.2"
	}

	if options.LintTimeout == "" {
		options.LintTimeout = "5m"
	}

	if options.LintArgs == nil {
		options.LintArgs = []string{
			"run",
			fmt.Sprintf("--timeout=%s", options.LintTimeout),
			"-v",
		}
	}
}
