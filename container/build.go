package container

import (
	"dagger.io/dagger"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/zcubbs/dagger-utils/types"
	"os"
)

type ImageBuilder struct {
	types.Options
	types.ImageBuilderOptions
	types.RegistryInfo
}

func (b *ImageBuilder) Build(imgName, imgTag string) (*string, error) {
	types.SetDefaults(&b.Options)
	types.SetupImageBuilderDefaults(&b.ImageBuilderOptions)

	return b.PackBuild(imgName, imgTag)
}

func (b *ImageBuilder) PackBuild(imgName, imgTag string) (*string, error) {
	client := b.DaggerClient
	ctx := b.Options.Ctx
	packBuilder := client.Container().From(b.BuildImg)

	dcBytes, err := generateDC(&b.RegistryInfo)
	if err != nil {
		return nil, err
	}

	ddir := client.Directory().WithNewFile(
		"/tmp/config.json",
		string(dcBytes),
		dagger.DirectoryWithNewFileOpts{},
	)

	packBuilder = packBuilder.WithMountedDirectory("/mnt", ddir)
	packBuilder = packBuilder.
		WithMountedDirectory("/src", b.Options.Src).
		WithWorkdir("/src")

	imageName := fmt.Sprintf("%s/%s:%s", b.RegistryInfo.RegistryServer, imgName, imgTag)
	packBuilder = packBuilder.WithExec(
		[]string{"mkdir", "-p", "/home/cnb/.docker/"},
	)
	packBuilder = packBuilder.WithExec(
		[]string{"sh", "-c", "cp /mnt/tmp/config.json /home/cnb/.docker/config.json"},
	)
	build := packBuilder.WithExec(
		[]string{"bash", "-c", fmt.Sprintf("CNB_PLATFORM_API=0.8 /cnb/lifecycle/creator -app=. %s", imageName)},
	)

	_, err = build.Stdout(ctx)
	return &imageName, os.RemoveAll("./src")
}

func NewDockerConfig(username, password, email, server string) *types.DockerConfig {
	if username == "" || password == "" || email == "" || server == "" {
		return nil
	}
	authStr := fmt.Sprintf("%s:%s", username, password)
	authBytes := []byte(authStr)
	encodedAuth := base64.StdEncoding.EncodeToString(authBytes)

	return &types.DockerConfig{
		Auths: map[string]types.AuthConfig{
			fmt.Sprintf("https://%s", server): {
				Auth:  encodedAuth,
				Email: email,
			},
		},
	}
}

func generateDC(regInfo *types.RegistryInfo) ([]byte, error) {
	var dcBytes []byte
	if regInfo != nil {
		var err error
		dc := NewDockerConfig(regInfo.RegistryUsername, regInfo.RegistryPassword, regInfo.RegistryEmail, regInfo.RegistryServer)
		if dc != nil {
			dcBytes, err = json.Marshal(dc)
			if err != nil {
				return nil, err
			}
		}
	}
	return dcBytes, nil
}
