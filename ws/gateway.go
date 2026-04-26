// Package ws provides WebSocket gateway support for the NestGo framework.
package ws

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Message is a WebSocket message.
type Message struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
	Room  string          `json:"room,omitempty"`
}

// Connection wraps a WebSocket connection with metadata.
type Connection struct {
	ID     string
	Conn   *websocket.Conn
	Rooms  map[string]bool
	mu     sync.Mutex
	values map[string]any
}

// Send sends a message to this connection.
func (c *Connection) Send(msg Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Conn.WriteJSON(msg)
}

// Set stores a value on the connection.
func (c *Connection) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[key] = value
}

// Get retrieves a value from the connection.
func (c *Connection) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.values[key]
	return v, ok
}

// MessageHandler handles a WebSocket event.
type MessageHandler func(conn *Connection, data json.RawMessage) error

// Gateway manages WebSocket connections and message routing.
type Gateway struct {
	mu          sync.RWMutex
	connections map[string]*Connection
	rooms       map[string]map[string]*Connection
	handlers    map[string]MessageHandler
	upgrader    websocket.Upgrader
	logger      *slog.Logger
	connCounter int

	// Lifecycle callbacks.
	onConnect    func(conn *Connection)
	onDisconnect func(conn *Connection)
}

// NewGateway creates a new WebSocket gateway.
func NewGateway(logger *slog.Logger) *Gateway {
	return &Gateway{
		connections: make(map[string]*Connection),
		rooms:       make(map[string]map[string]*Connection),
		handlers:    make(map[string]MessageHandler),
		logger:      logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// On registers a handler for an event.
func (gw *Gateway) On(event string, handler MessageHandler) {
	gw.mu.Lock()
	defer gw.mu.Unlock()
	gw.handlers[event] = handler
}

// OnConnect sets the connection callback.
func (gw *Gateway) OnConnect(fn func(conn *Connection)) {
	gw.onConnect = fn
}

// OnDisconnect sets the disconnection callback.
func (gw *Gateway) OnDisconnect(fn func(conn *Connection)) {
	gw.onDisconnect = fn
}

// HandleHTTP upgrades HTTP to WebSocket and manages the connection.
func (gw *Gateway) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	wsConn, err := gw.upgrader.Upgrade(w, r, nil)
	if err != nil {
		gw.logger.Error("websocket upgrade failed", "error", err)
		return
	}

	gw.mu.Lock()
	gw.connCounter++
	connID := fmt.Sprintf("conn_%d", gw.connCounter)
	conn := &Connection{
		ID:     connID,
		Conn:   wsConn,
		Rooms:  make(map[string]bool),
		values: make(map[string]any),
	}
	gw.connections[connID] = conn
	gw.mu.Unlock()

	gw.logger.Info("websocket connected", "connection_id", connID)

	if gw.onConnect != nil {
		gw.onConnect(conn)
	}

	// Read loop.
	go gw.readLoop(conn)
}

// JoinRoom adds a connection to a room.
func (gw *Gateway) JoinRoom(conn *Connection, room string) {
	gw.mu.Lock()
	defer gw.mu.Unlock()

	if gw.rooms[room] == nil {
		gw.rooms[room] = make(map[string]*Connection)
	}
	gw.rooms[room][conn.ID] = conn
	conn.Rooms[room] = true
}

// LeaveRoom removes a connection from a room.
func (gw *Gateway) LeaveRoom(conn *Connection, room string) {
	gw.mu.Lock()
	defer gw.mu.Unlock()

	if gw.rooms[room] != nil {
		delete(gw.rooms[room], conn.ID)
		if len(gw.rooms[room]) == 0 {
			delete(gw.rooms, room)
		}
	}
	delete(conn.Rooms, room)
}

// BroadcastToRoom sends a message to all connections in a room.
func (gw *Gateway) BroadcastToRoom(room string, msg Message) {
	gw.mu.RLock()
	conns := gw.rooms[room]
	gw.mu.RUnlock()

	for _, conn := range conns {
		if err := conn.Send(msg); err != nil {
			gw.logger.Error("broadcast failed", "room", room, "connection", conn.ID, "error", err)
		}
	}
}

// Broadcast sends a message to all connected clients.
func (gw *Gateway) Broadcast(msg Message) {
	gw.mu.RLock()
	conns := make([]*Connection, 0, len(gw.connections))
	for _, c := range gw.connections {
		conns = append(conns, c)
	}
	gw.mu.RUnlock()

	for _, conn := range conns {
		if err := conn.Send(msg); err != nil {
			gw.logger.Error("broadcast failed", "connection", conn.ID, "error", err)
		}
	}
}

// ConnectionCount returns the number of active connections.
func (gw *Gateway) ConnectionCount() int {
	gw.mu.RLock()
	defer gw.mu.RUnlock()
	return len(gw.connections)
}

// Handler returns an http.HandlerFunc for the WebSocket endpoint.
func (gw *Gateway) Handler() http.HandlerFunc {
	return gw.HandleHTTP
}

func (gw *Gateway) readLoop(conn *Connection) {
	defer func() {
		gw.removeConnection(conn)
		_ = conn.Conn.Close()
	}()

	for {
		var msg Message
		err := conn.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				gw.logger.Error("websocket read error", "connection", conn.ID, "error", err)
			}
			break
		}

		gw.mu.RLock()
		handler, exists := gw.handlers[msg.Event]
		gw.mu.RUnlock()

		if !exists {
			gw.logger.Warn("unhandled websocket event", "event", msg.Event)
			continue
		}

		if err := handler(conn, msg.Data); err != nil {
			gw.logger.Error("websocket handler error", "event", msg.Event, "error", err)
		}
	}
}

func (gw *Gateway) removeConnection(conn *Connection) {
	gw.mu.Lock()
	defer gw.mu.Unlock()

	// Remove from all rooms.
	for room := range conn.Rooms {
		if gw.rooms[room] != nil {
			delete(gw.rooms[room], conn.ID)
			if len(gw.rooms[room]) == 0 {
				delete(gw.rooms, room)
			}
		}
	}

	delete(gw.connections, conn.ID)
	gw.logger.Info("websocket disconnected", "connection_id", conn.ID)

	if gw.onDisconnect != nil {
		gw.onDisconnect(conn)
	}
}
