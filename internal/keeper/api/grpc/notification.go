package grpc

import (
	"context"
	"log"
	"time"

	"gophkeeper/internal/keeper/tui"
	pb "gophkeeper/pkg/proto/keeper/grpcapi"

	tea "github.com/charmbracelet/bubbletea"
)

// Subscribes for notifications and sending signal to tea program to reload list
func (c *GRPCClient) Notifications(p *tea.Program) {
	var (
		stream pb.Notification_SubscribeV1Client
		err    error
	)

	for {
		// Subscribe to notifications
		if stream == nil {
			if stream, err = c.subscribe(); err != nil {
				log.Printf("failed to subscribe: %v\n", err)
				c.sleep()
				continue
			}
		}

		response, err := stream.Recv()
		if err != nil {
			log.Printf("failed to recv msg: %v\n", err)
			stream = nil
			c.sleep()

			// Retry
			continue
		}

		log.Println("received", response)

		// Trigger secret list reload
		if p != nil {
			p.Send(tui.ReloadSecretList{})
		}
	}
}

func (c *GRPCClient) sleep() {
	time.Sleep(time.Second * 2)
}

func (c *GRPCClient) subscribe() (pb.Notification_SubscribeV1Client, error) {
	return c.notifyClient.SubscribeV1(context.Background(), &pb.SubscribeV1Request{
		Id: c.clientID,
	})
}
