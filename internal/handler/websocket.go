package handler

import (
	"net/http"
	"supermarket-backend/internal/middleware"
	"supermarket-backend/internal/response"
	"supermarket-backend/internal/ws"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WebsocketHandler struct {
	hub            *ws.WebsocketHub
	allowedOrigins []string
}

func NewWebsocketHandler(hub *ws.WebsocketHub, allowedOrigins []string) *WebsocketHandler {
	return &WebsocketHandler{
		hub:            hub,
		allowedOrigins: allowedOrigins,
	}
}

func (h *WebsocketHandler) Handle(c *gin.Context) {
	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{OriginPatterns: h.allowedOrigins})

	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			"Failed to upgrade websocket connection",
		)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "connection closed")

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		conn.Close(
			websocket.StatusPolicyViolation,
			"Unauthorized",
		)
		return
	}

	branchID := middleware.GetUserBranchID(c)

	if branchID == "" {
		conn.Close(
			websocket.StatusPolicyViolation,
			"Branch not found",
		)
		return
	}

	rawServiceType := c.Query("service_type")

	if rawServiceType == "" {
		conn.Close(
			websocket.StatusPolicyViolation,
			"Service type is required",
		)
		return
	}

	serviceType := ws.WsServiceType(rawServiceType)

	client, err := h.hub.Register(c.Request.Context(), conn, userID, branchID, serviceType)
	if err != nil {
		conn.Close(
			websocket.StatusPolicyViolation,
			err.Error(),
		)
		return
	}
	defer h.hub.Unregister(client)

	err = client.Start(c.Request.Context())
	if err != nil {
		conn.Close(
			websocket.StatusPolicyViolation,
			err.Error(),
		)
		return
	}
}
