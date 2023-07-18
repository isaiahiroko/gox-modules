package container

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

var (
	engine *client.Client
	ctx    context.Context
)

func init() {
	ctx = context.Background()
	var err error
	engine, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
}

type container struct {
	logger io.Writer
}

func (c *container) Build(directory, name string) error {
	source, err := archive.TarWithOptions(directory, &archive.TarOptions{})
	if err != nil {
		return err
	}

	res, err := engine.ImageBuild(
		ctx,
		source,
		types.ImageBuildOptions{
			SuppressOutput: false,
			Remove:         true,
			ForceRemove:    true,
			PullParent:     true,
			Tags:           []string{name},
			Dockerfile:     "Dockerfile",
		},
	)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(c.logger, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *container) Push(username, password, address, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	credential, err := json.Marshal(types.AuthConfig{
		Username:      username,
		Password:      password,
		ServerAddress: address,
	})
	if err != nil {
		return err
	}

	auth := base64.URLEncoding.EncodeToString(credential)

	res, err := engine.ImagePush(ctx, name, types.ImagePushOptions{
		RegistryAuth: auth,
	})
	if err != nil {
		return err
	}
	defer res.Close()

	_, err = io.Copy(c.logger, res)
	if err != nil {
		return err
	}

	return nil
}

func Close() error {
	return engine.Close()
}

func New(logger io.Writer) *container {
	return &container{logger}
}
