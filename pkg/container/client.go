package container

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Instance of client for managing containers.
type Client struct {
	docker *client.Client
}

// Return new client instance.
func NewClient() (*Client, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cl := &Client{docker: docker}
	return cl, nil
}

// Create a new container from an image.
func (cl *Client) Create(ctx context.Context, img string) (*Container, error) {
	config := &container.Config{
		Image:           img,
		Hostname:        img,
		Tty:             true,
		OpenStdin:       true,
		AttachStdin:     true,
		AttachStdout:    true,
		AttachStderr:    true,
		NetworkDisabled: true,
	}
	res, err := cl.docker.ContainerCreate(ctx, config, nil, nil, "")
	if err != nil {
		return nil, err
	}
	return &Container{cl: cl, ID: res.ID}, nil
}
