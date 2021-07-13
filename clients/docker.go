package clients

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"log"
	"os"
)

type DockerClient interface {
	Pull(image string) error
}

type dockerClient struct{
	Docker *client.Client
	Context context.Context
}

func NewDockerClient() DockerClient {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        log.Fatal("is docker running ?")
    }
	return &dockerClient{
		Context: context.Background(),
		Docker: cli,
	}
}

func (d *dockerClient) Pull(image string) error {
	reader, err := d.Docker.ImagePull(d.Context, "docker.io/library/alpine", types.ImagePullOptions{})
    if err != nil {
        return err
    }
    defer reader.Close()
    io.Copy(os.Stdout, reader)
	return nil
}
