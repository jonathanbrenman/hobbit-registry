package clients

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"log"
	"os"
)

type DockerClient interface {
	Pull(image string) error
	Tag(image string) error
	Push(image string) error
	Delete(image string) error
}

type dockerClient struct{
	PrivateRegistry string
	Docker *client.Client
	Context context.Context
	Credentials *string
}

func NewDockerClient(user, pass, privateRegistry string) DockerClient {
	var authStr string
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        log.Fatal("is docker running ?")
    }
    if user != "" && pass != "" {
    	authConfig := types.AuthConfig{
			Username: user,
			Password: pass,
		}
		encodedJSON, _ := json.Marshal(authConfig)
		authStr = base64.URLEncoding.EncodeToString(encodedJSON)
	}
	return &dockerClient{
		PrivateRegistry: privateRegistry,
		Context: context.Background(),
		Docker: cli,
		Credentials: &authStr,
	}
}

func (d *dockerClient) Pull(image string) error {
	reader, err := d.Docker.ImagePull(d.Context, image, types.ImagePullOptions{
		RegistryAuth: *d.Credentials,
	})
    if err != nil {
        return err
    }
    defer reader.Close()
    io.Copy(os.Stdout, reader)
	return nil
}

func (d *dockerClient) Tag(image string) error {
	err := d.Docker.ImageTag(d.Context, image, fmt.Sprintf("%s/%s", d.PrivateRegistry, image))
    if err != nil {
        return err
    }
	return nil
}

func (d *dockerClient) Push(image string) error {
	opts := types.ImagePushOptions{
		RegistryAuth: "123", // Workarround
	}
	if *d.Credentials != "" {
		opts = types.ImagePushOptions{
			RegistryAuth: *d.Credentials,
		}
	}
	reader, err := d.Docker.ImagePush(d.Context,
		fmt.Sprintf("%s/%s", d.PrivateRegistry, image),
		opts,
	)
    if err != nil {
        return err
    }
    defer reader.Close()
    io.Copy(os.Stdout, reader)
	return nil
}

func (d *dockerClient) Delete(image string) error {
	_, err := d.Docker.ImageRemove(d.Context, image, types.ImageRemoveOptions{})
    if err != nil {
        return err
    }
    _, err = d.Docker.ImageRemove(d.Context,
    	fmt.Sprintf("%s/%s", d.PrivateRegistry, image,
    	), types.ImageRemoveOptions{})
    if err != nil {
        return err
    }
	return nil
}