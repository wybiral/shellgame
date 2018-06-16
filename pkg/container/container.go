package container

import (
	"context"
	"github.com/docker/docker/api/types"
	"net"
)

// Instance of a Linux container.
type Container struct {
	cl *Client
	ID string
}

// Attach to container IO and return as a net.Conn instance.
func (c *Container) Attach(ctx context.Context) (net.Conn, error) {
	options := types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	}
	hijack, err := c.cl.docker.ContainerAttach(ctx, c.ID, options)
	if err != nil {
		return nil, err
	}
	return hijack.Conn, nil
}

// Start the container.
func (c *Container) Start(ctx context.Context) error {
	options := types.ContainerStartOptions{}
	return c.cl.docker.ContainerStart(ctx, c.ID, options)
}

// Send kill signal to container.
func (c *Container) Kill(ctx context.Context) error {
	return c.cl.docker.ContainerKill(ctx, c.ID, "SIGKILL")
}

// Remove container (after killing).
func (c *Container) Remove(ctx context.Context) error {
	options := types.ContainerRemoveOptions{
		RemoveVolumes: true,
	}
	return c.cl.docker.ContainerRemove(ctx, c.ID, options)
}
