package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"flowforge-api/common/realtime"
	"flowforge-api/middleware"
)

type RealtimeHandler struct {
	hub *realtime.Hub
}

func NewRealtimeHandler(hub *realtime.Hub) *RealtimeHandler {
	return &RealtimeHandler{hub: hub}
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func (h *RealtimeHandler) WebSocket(c *gin.Context) {
	ch, unsub := h.hub.Subscribe(c.GetString(middleware.TenantIDKey))
	defer unsub()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for ev := range ch {
		if err := conn.WriteJSON(ev); err != nil {
			return
		}
	}
}

func (h *RealtimeHandler) Events(c *gin.Context) {
	ch, unsub := h.hub.Subscribe(c.GetString(middleware.TenantIDKey))
	defer unsub()

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Stream(func(w io.Writer) bool {
		ev, ok := <-ch
		if !ok {
			return false
		}
		_, _ = fmt.Fprintf(w, "data: %s\n\n", ev.JSON())
		return true
	})
}
