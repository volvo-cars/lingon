package bla

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/nats-io/nats.go"
)

var _ (Actor) = (*ContainerActor)(nil)

type ContainerActor struct {
	ctx  context.Context
	Conn *nats.Conn

	cli *client.Client

	subs []*nats.Subscription
}

// Subscribe implements Actor.
func (ca *ContainerActor) Subscribe() error {
	ca.ctx = context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("new docker client: %w", err)
	}
	ca.cli = cli

	{
		sub, err := SubscribeToSubjectWithAccount[ContainerRunMsg, ContainerRunReply](
			ca.Conn,
			fmt.Sprintf(ActionSubjectSubscribeForAccount, "container", "run"),
			ca.run,
		)
		if err != nil {
			return fmt.Errorf("subscribing: %w", err)
		}
		ca.subs = append(ca.subs, sub)
	}
	return nil
}

// Unsubscribe implements Actor.
func (ca *ContainerActor) Unsubscribe() error {
	for _, sub := range ca.subs {
		_ = sub.Unsubscribe()
	}
	_ = ca.cli.Close()
	return nil
}

type ContainerRunMsg struct {
	Host  string
	Image string
	Env   []string
	Cmd   []string
}

type ContainerRunReply struct {
	ID string
}

func (ca *ContainerActor) run(accountID string, msg *ContainerRunMsg) (*ContainerRunReply, error) {
	// TODO: need to pull the image also.
	// ca.cli.ImagePull(ca.ctx, msg.Image, types.ImagePullOptions{})
	resp, err := ca.cli.ContainerCreate(ca.ctx, &container.Config{
		Image: msg.Image,
		Env:   msg.Env,
		Cmd:   msg.Cmd,
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		return nil, fmt.Errorf("create container: %w", err)
	}
	return &ContainerRunReply{
		ID: resp.ID,
	}, nil
}
