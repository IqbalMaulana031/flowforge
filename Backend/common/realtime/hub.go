package realtime

import (
	"encoding/json"
	"sync"
)

type Event struct {
	TenantID string `json:"tenant_id"`
	RunID    string `json:"run_id"`
	StepID   string `json:"step_id,omitempty"`
	Status   string `json:"status"`
	Type     string `json:"type"`
	Payload  any    `json:"payload,omitempty"`
}
type Hub struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan Event]struct{}
}

func NewHub() *Hub { return &Hub{subscribers: map[string]map[chan Event]struct{}{}} }
func (h *Hub) Subscribe(tenantID string) (chan Event, func()) {
	ch := make(chan Event, 16)
	h.mu.Lock()
	if h.subscribers[tenantID] == nil {
		h.subscribers[tenantID] = map[chan Event]struct{}{}
	}
	h.subscribers[tenantID][ch] = struct{}{}
	h.mu.Unlock()
	return ch, func() { h.mu.Lock(); delete(h.subscribers[tenantID], ch); close(ch); h.mu.Unlock() }
}
func (h *Hub) Publish(ev Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.subscribers[ev.TenantID] {
		select {
		case ch <- ev:
		default:
		}
	}
}
func (e Event) JSON() []byte { b, _ := json.Marshal(e); return b }
