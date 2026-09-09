package ws

import (
	"context"
	"fmt"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	ctx    context.Context
	client *WsClient
	errC   chan error
}

type BroadcastRequest struct {
	message     string
	userID      uuid.UUID
	serviceType WsServiceType
}

type WebsocketHub struct {
	clients map[string]*WsClient

	registerC   chan RegisterRequest
	unregisterC chan *WsClient
	broadcastC  chan BroadcastRequest
}

func NewWebsocketHub() *WebsocketHub {
	return &WebsocketHub{
		clients:     make(map[string]*WsClient),
		registerC:   make(chan RegisterRequest, 5),
		unregisterC: make(chan *WsClient, 5),
		broadcastC:  make(chan BroadcastRequest, 5),
	}
}

func (h *WebsocketHub) Register(
	ctx context.Context,
	conn *websocket.Conn,
	userID uuid.UUID,
	branchID string,
	serviceType WsServiceType,
) (*WsClient, error) {
	switch serviceType {
	case WsCustomerDisplayType:
		// supported
	default:
		return nil, fmt.Errorf("unsupported service type: %s", serviceType)
	}

	client := NewWsClient(conn, userID, branchID, serviceType)

	req := RegisterRequest{
		ctx:    ctx,
		client: client,
		errC:   make(chan error),
	}

	select {
	case h.registerC <- req:
		// Request đã được gửi vào Hub
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	select {
	case err := <-req.errC:
		return client, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (h *WebsocketHub) Unregister(client *WsClient) {
	h.unregisterC <- client
}

func (h *WebsocketHub) Broadcast(
	ctx context.Context,
	message string,
	userID uuid.UUID,
	serviceType WsServiceType,
) error {
	req := BroadcastRequest{
		message:     message,
		userID:      userID,
		serviceType: serviceType,
	}

	select {
	case h.broadcastC <- req:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (h *WebsocketHub) Run(ctx context.Context) {
	for {
		select {
		case registerReq := <-h.registerC:
			select {
			case <-registerReq.ctx.Done():
				// do nothing
			default:
				key := registerReq.client.userID.String() + "_" + string(registerReq.client.serviceType)
				h.clients[key] = registerReq.client
				registerReq.errC <- nil
			}

		case client := <-h.unregisterC:
			key := client.userID.String() + "_" + string(client.serviceType)
			if _, ok := h.clients[key]; ok {
				delete(h.clients, key)
			}
		case req := <-h.broadcastC:
			key := req.userID.String() + "_" + string(req.serviceType)
			client, ok := h.clients[key]
			if !ok {
				continue
			}

			client.sendC <- req.message
		case <-ctx.Done():
			return
		}
	}
}
