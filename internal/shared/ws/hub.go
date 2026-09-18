package ws

import (
	"encoding/json"

	"github.com/olahol/melody"
)

// Hub centraliza los WebSockets de la aplicacion. Los modulos publican eventos
// aqui y Melody los transmite a los clientes conectados.
type Hub struct {
	melody *melody.Melody
}

func NewHub() *Hub {
	return &Hub{melody: melody.New()}
}

// Melody expone la instancia para registrar la ruta HTTP del WebSocket.
func (h *Hub) Melody() *melody.Melody { return h.melody }

// Broadcast envia {"type":"...","payload":...} a todos los clientes.
func (h *Hub) Broadcast(eventType string, payload any) {
	if h == nil || h.melody == nil {
		return
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return
	}
	data := []byte(`{"type":"` + eventType + `","payload":` + string(raw) + `}`)
	_ = h.melody.Broadcast(data)
}
