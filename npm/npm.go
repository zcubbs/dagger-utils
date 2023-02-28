package npm

import "github.com/zcubbs/dagger-utils/types"

type Builder struct {
	types.Options
}

func (b Builder) NpmTest(options types.NpmOptions) error {
	setupOptions(&options)
	_, err := b.Options.DaggerClient.Container().
		From(options.BuildImg).
		WithMountedDirectory("/src", b.Options.Src).
		WithWorkdir("/src").
		WithExec([]string{"npm", "test"}).
		ExitCode(b.Options.Ctx)
	if err != nil {
		return err
	}

	return nil
}

func (b Builder) NpmBuild(options types.NpmOptions) error {
	setupOptions(&options)
	_, err := b.Options.DaggerClient.Container().
		From(options.BuildImg).
		WithMountedDirectory("/src", b.Options.Src).
		WithWorkdir("/src").
		WithExec([]string{"npm", "run", "build"}).
		ExitCode(b.Options.Ctx)
	if err != nil {
		return err
	}

	return nil
}

func (b Builder) NpmLint(options types.NpmOptions) error {
	setupOptions(&options)
	_, err := b.Options.DaggerClient.Container().
		From(options.BuildImg).
		WithMountedDirectory("/src", b.Options.Src).
		WithWorkdir("/src").
		WithExec([]string{"npm", "run", "lint"}).
		ExitCode(b.Options.Ctx)
	if err != nil {
		return err
	}

	return nil
}

func (b Builder) NpmInstall(options types.NpmOptions) error {
	setupOptions(&options)
	_, err := b.Options.DaggerClient.Container().
		From(options.BuildImg).
		WithMountedDirectory("/src", b.Options.Src).
		WithWorkdir("/src").
		WithExec([]string{"npm", "install"}).
		ExitCode(b.Options.Ctx)
	if err != nil {
		return err
	}

	return nil
}

func (b Builder) NpmBuildImage(options types.NpmOptions) error {
	setupOptions(&options)
	// TODO: add support for custom Dockerfile
	return nil
}

func setupOptions(options *types.NpmOptions) {
	types.SetDefaults(&options.Options)

	if options.BuildImg == "" {
		options.BuildImg = "node:16"
	}
}
