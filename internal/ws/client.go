package ws

import (
	"context"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type WsServiceType string

const (
	WsCustomerDisplayType WsServiceType = "customer_display"
)

type WsClient struct {
	conn        *websocket.Conn
	userID      uuid.UUID
	branchID    string
	serviceType WsServiceType
	sendC       chan string
}

func NewWsClient(
	conn *websocket.Conn,
	userID uuid.UUID,
	branchID string,
	serviceType WsServiceType,
) *WsClient {
	return &WsClient{
		conn:        conn,
		userID:      userID,
		branchID:    branchID,
		serviceType: serviceType,
		sendC:       make(chan string, 5),
	}
}

func (c *WsClient) Read(ctx context.Context) error {
	for {
		_, _, err := c.conn.Read(ctx)
		if err != nil {
			return err
		}
	}
}

func (c *WsClient) Write(ctx context.Context) error {
	for {
		select {
		case message := <-c.sendC:
			if err := c.conn.Write(
				ctx,
				websocket.MessageText,
				[]byte(message),
			); err != nil {
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *WsClient) Start(parentCtx context.Context) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	errC := make(chan error, 2)

	go func() {
		errC <- c.Read(ctx)
	}()

	go func() {
		errC <- c.Write(ctx)
	}()

	err := <-errC

	return err
}
