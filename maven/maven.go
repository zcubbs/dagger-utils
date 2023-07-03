package maven

import (
	"github.com/zcubbs/dagger-utils/types"
)

type Builder struct {
	types.Options
	types.MavenOptions
}

func (b Builder) MavenTest() error {
	return nil
}

func (b Builder) MavenBuild() error {
	types.SetDefaults(&b.Options)
	setupDefaults(&b.MavenOptions)

	buildCmd := []string{"mvn"}
	buildCmd = append(buildCmd, b.MvnArgs...)

	builder := b.Options.DaggerClient.
		Container().
		From(b.BuildImg).
		WithMountedDirectory("/src", b.Options.Src).
		WithWorkdir("/src")

	if b.MavenOptions.EnableCache {
		cacheVolume := b.Options.DaggerClient.CacheVolume(b.CacheKey)
		builder = builder.WithMountedCache(b.CacheKey, cacheVolume)
	}

	_, err := builder.ExitCode(b.Options.Ctx)
	if err != nil {
		return err
	}

	return nil
}

func setupDefaults(options *types.MavenOptions) {
	types.SetDefaults(&options.Options)

	if options.BuildImg == "" {
		options.BuildImg = "maven:3.9.3-eclipse-temurin-17-focal"
	}

	if options.CacheKey == "" {
		options.CacheKey = "maven-cache"
	}

	if options.MvnArgs == nil {
		options.MvnArgs = []string{"install", "-DskipTests=true", "-Dmaven.javadoc.skip=true", "-B", "-V"}
	}
}
