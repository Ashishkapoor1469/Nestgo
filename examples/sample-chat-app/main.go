package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/Ashishkapoor1469/Nestgo/cache"
	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/core"
	"github.com/Ashishkapoor1469/Nestgo/di"
	nestgrpc "github.com/Ashishkapoor1469/Nestgo/grpc"
	"github.com/Ashishkapoor1469/Nestgo/middleware"
	"github.com/Ashishkapoor1469/Nestgo/ws"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ChatMessage represents a chat message.
type ChatMessage struct {
	User      string    `json:"user"`
	Message   string    `json:"message"`
	Room      string    `json:"room"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	messagesMu   sync.Mutex
	messagesList = []ChatMessage{
		{User: "System", Message: "Welcome to the NestGo Chat App!", Room: "general", Timestamp: time.Now()},
	}
)

func main() {
	logger := slog.Default()
	logger.Info("Starting NestGo chat application...")

	// 1. Redis Connection Setup with Fallback
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	var rateStore middleware.RateLimitStore
	var cacheStore middleware.CacheStore

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err == nil {
		logger.Info("✅ Connected to Redis successfully! Using Redis for Caching and Rate Limiting.")
		rateStore = middleware.NewRedisStore(rdb, "chat_ratelimit")
		cacheStore = cache.NewRedisCacheStore(rdb, "chat_cache")
	} else {
		logger.Info("⚠️ Redis not running (or unreachable). Falling back to In-Memory Caching and Rate Limiting.")
		rateStore = middleware.NewMemoryStore()
		cacheStore = middleware.NewMemoryCacheStore()
	}

	// 2. Middleware instances
	messageRateLimiter := middleware.NewRateLimiter(5, 10*time.Second).WithStore(rateStore)
	historyCache := middleware.NewCache(30*time.Second).WithStore(cacheStore).WithTags("history")

	// 3. Create NestGo App
	app := core.New(
		core.WithAddress(":3000"),
		core.WithGRPC(":50051"),
	)

	// 4. Setup WebSocket Module
	wsModule := ws.NewWebSocketModule("/ws", logger)
	app.RegisterModule(wsModule)
	setupWebSocket(wsModule.Gateway(), cacheStore)

	// 5. Register Controllers
	apiController := &ApiController{
		cacheStore:             cacheStore,
		grpcAddr:               "localhost:50051",
		historyCacheMiddleware: historyCache.Middleware(),
		messageRateLimiter:     messageRateLimiter,
	}
	app.RegisterModule(&ChatModule{
		apiController:          apiController,
	})

	// 6. Register gRPC Service Registrar
	pingServer := &MyPingServer{}
	_ = app.Container().ProvideValue(pingServer)

	// Start the server
	logger.Info("Starting NestGo App (HTTP on :3000, gRPC on :50051)...")
	if err := app.Start(); err != nil {
		logger.Error("Application shutdown error", "error", err)
	}
}

// Setup WebSocket gateway events
func setupWebSocket(gw *ws.Gateway, cacheStore middleware.CacheStore) {
	gw.OnConnect(func(conn *ws.Connection) {
		gw.JoinRoom(conn, "general")
		welcome := ws.Message{
			Event: "user-joined",
			Room:  "general",
			Data:  json.RawMessage(fmt.Sprintf(`{"user":"System","message":"%s has connected to general"}`, conn.ID)),
		}
		gw.BroadcastToRoom("general", welcome)
	})

	gw.OnDisconnect(func(conn *ws.Connection) {
		for r := range conn.Rooms {
			bye := ws.Message{
				Event: "user-left",
				Room:  r,
				Data:  json.RawMessage(fmt.Sprintf(`{"user":"System","message":"%s left %s"}`, conn.ID, r)),
			}
			gw.BroadcastToRoom(r, bye)
		}
	})

	gw.On("send-message", func(conn *ws.Connection, data json.RawMessage) error {
		var payload struct {
			User    string `json:"user"`
			Message string `json:"message"`
			Room    string `json:"room"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return err
		}

		msg := ChatMessage{
			User:      payload.User,
			Message:   payload.Message,
			Room:      payload.Room,
			Timestamp: time.Now(),
		}

		messagesMu.Lock()
		messagesList = append(messagesList, msg)
		messagesMu.Unlock()

		// Invalidate cache
		cacheStore.InvalidateTags(context.Background(), "history")

		rawMsg, _ := json.Marshal(msg)
		gw.BroadcastToRoom(payload.Room, ws.Message{
			Event: "new-message",
			Room:  payload.Room,
			Data:  rawMsg,
		})
		return nil
	})

	gw.On("join-room", func(conn *ws.Connection, data json.RawMessage) error {
		var payload struct {
			User string `json:"user"`
			Room string `json:"room"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return err
		}

		// Leave all current rooms
		for r := range conn.Rooms {
			gw.LeaveRoom(conn, r)
			gw.BroadcastToRoom(r, ws.Message{
				Event: "user-left",
				Room:  r,
				Data:  json.RawMessage(fmt.Sprintf(`{"user":"System","message":"%s left %s"}`, payload.User, r)),
			})
		}

		// Join new room
		gw.JoinRoom(conn, payload.Room)
		gw.BroadcastToRoom(payload.Room, ws.Message{
			Event: "user-joined",
			Room:  payload.Room,
			Data:  json.RawMessage(fmt.Sprintf(`{"user":"System","message":"%s joined %s"}`, payload.User, payload.Room)),
		})
		return nil
	})
}

// ChatModule holds routing setup
type ChatModule struct {
	apiController          *ApiController
}

func (m *ChatModule) Module() common.ModuleConfig {
	return common.ModuleConfig{
		Name: "ChatModule",
		Controllers: []common.Controller{
			&HomeController{},
			m.apiController,
		},
		Providers: []di.Provider{},
	}
}

// HomeController serves the frontend SPA
type HomeController struct{}

func (c *HomeController) Prefix() string { return "/" }

func (c *HomeController) Routes() []common.Route {
	return []common.Route{
		{
			Method: "GET",
			Path:   "/",
			Handler: func(ctx *common.Context) error {
				ctx.Writer.Header().Set("Content-Type", "text/html")
				_, err := ctx.Writer.Write([]byte(indexHTML))
				return err
			},
		},
	}
}

// ApiController handles REST API endpoints
type ApiController struct {
	cacheStore             middleware.CacheStore
	grpcAddr               string
	historyCacheMiddleware common.Middleware
	messageRateLimiter     *middleware.RateLimiter
}

func (c *ApiController) Prefix() string { return "/api" }

func (c *ApiController) Routes() []common.Route {
	return []common.Route{
		{
			Method:      "GET",
			Path:        "/history",
			Handler:     c.GetHistory,
			Middlewares: []common.Middleware{c.historyCacheMiddleware},
			Summary:     "Get message history",
			Description: "Retrieves chat history. Cached for 30s.",
		},
		{
			Method:      "POST",
			Path:        "/message",
			Handler:     c.PostMessage,
			Middlewares: []common.Middleware{common.Middleware(c.messageRateLimiter.Middleware())},
			Summary:     "Post message via REST",
			Description: "Allows posting a message via REST. Rate limited to 5 requests per 10 seconds.",
		},
		{
			Method:      "GET",
			Path:        "/ping-grpc",
			Handler:     c.PingGRPC,
			Summary:     "Ping gRPC server",
			Description: "Pings the parallel gRPC server and returns response.",
		},
	}
}

func (c *ApiController) GetHistory(ctx *common.Context) error {
	messagesMu.Lock()
	defer messagesMu.Unlock()
	return ctx.OK(messagesList)
}

func (c *ApiController) PostMessage(ctx *common.Context) error {
	var body struct {
		User    string `json:"user"`
		Message string `json:"message"`
		Room    string `json:"room"`
	}
	if err := ctx.Bind(&body); err != nil {
		return ctx.ValidationErrorResponse(err)
	}

	msg := ChatMessage{
		User:      body.User,
		Message:   body.Message,
		Room:      body.Room,
		Timestamp: time.Now(),
	}

	messagesMu.Lock()
	messagesList = append(messagesList, msg)
	messagesMu.Unlock()

	// Invalidate history cache
	c.cacheStore.InvalidateTags(ctx.Request.Context(), "history")

	return ctx.Created(msg)
}

func (c *ApiController) PingGRPC(ctx *common.Context) error {
	msg := ctx.QueryDefault("msg", "Hello NestGo")
	res, err := callGRPCPing(c.grpcAddr, msg)
	if err != nil {
		return ctx.Error(http.StatusInternalServerError, err.Error())
	}
	return ctx.OK(map[string]string{"reply": res})
}

// --- gRPC custom manual reflection ---
type PingRequest struct {
	Message string `json:"message"`
}

type PingResponse struct {
	Message string `json:"message"`
}

type MyPingServer struct{}

var _ nestgrpc.ServiceRegistrar = (*MyPingServer)(nil)

func (s *MyPingServer) Ping(ctx context.Context, req *PingRequest) (*PingResponse, error) {
	return &PingResponse{Message: "Pong from gRPC: " + req.Message}, nil
}

func (s *MyPingServer) RegisterGRPC(server *grpc.Server) {
	server.RegisterService(&PingService_ServiceDesc, s)
}

var PingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "PingService",
	HandlerType: (*PingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _PingService_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ping.proto",
}

type PingServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
}

func _PingService_Ping_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PingServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PingService/Ping",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(PingServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func callGRPCPing(addr string, msg string) (string, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	req := &PingRequest{Message: msg}
	res := new(PingResponse)
	err = conn.Invoke(context.Background(), "/PingService/Ping", req, res)
	if err != nil {
		return "", err
	}
	return res.Message, nil
}

// Premium SPA HTML template
const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>NestGo v0.6.0 Enterprise Showcase</title>
  <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;500;600;700&display=swap" rel="stylesheet">
  <style>
    * {
      box-sizing: border-box;
      margin: 0;
      padding: 0;
      font-family: 'Outfit', sans-serif;
    }
    body {
      background: #090d16;
      color: #f3f4f6;
      height: 100vh;
      display: flex;
      flex-direction: column;
      overflow: hidden;
    }
    header {
      background: rgba(15, 23, 42, 0.6);
      backdrop-filter: blur(16px);
      border-bottom: 1px solid #1e293b;
      padding: 16px 28px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      z-index: 10;
    }
    .brand {
      display: flex;
      align-items: center;
      gap: 12px;
    }
    .logo {
      background: linear-gradient(135deg, #6366f1, #a855f7);
      width: 38px;
      height: 38px;
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 700;
      color: white;
      box-shadow: 0 4px 14px rgba(99, 102, 241, 0.4);
      font-size: 18px;
    }
    .title-area h1 {
      font-size: 18px;
      font-weight: 600;
      letter-spacing: -0.01em;
    }
    .title-area p {
      font-size: 12px;
      color: #64748b;
    }
    .status {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 13px;
      color: #94a3b8;
      background: rgba(30, 41, 59, 0.5);
      padding: 6px 12px;
      border-radius: 20px;
      border: 1px solid #1e293b;
    }
    .dot {
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background: #ef4444;
      transition: all 0.3s;
    }
    .dot.online {
      background: #10b981;
      box-shadow: 0 0 10px #10b981;
    }
    .container {
      display: flex;
      flex: 1;
      overflow: hidden;
    }
    .sidebar {
      width: 250px;
      background: #0c1322;
      border-right: 1px solid #1e293b;
      padding: 24px 16px;
      display: flex;
      flex-direction: column;
      gap: 28px;
    }
    .section-title {
      font-size: 11px;
      text-transform: uppercase;
      letter-spacing: 0.08em;
      color: #475569;
      font-weight: 600;
      margin-bottom: 10px;
    }
    .room-list {
      display: flex;
      flex-direction: column;
      gap: 6px;
    }
    .room-btn {
      background: transparent;
      border: none;
      color: #94a3b8;
      padding: 10px 14px;
      border-radius: 8px;
      text-align: left;
      cursor: pointer;
      font-weight: 500;
      font-size: 14px;
      transition: all 0.2s;
      display: flex;
      align-items: center;
      gap: 10px;
    }
    .room-btn:hover {
      background: rgba(30, 41, 59, 0.4);
      color: #f3f4f6;
    }
    .room-btn.active {
      background: rgba(99, 102, 241, 0.12);
      color: #818cf8;
      border: 1px solid rgba(99, 102, 241, 0.25);
    }
    .chat-area {
      flex: 1;
      display: flex;
      flex-direction: column;
      background: #090d16;
    }
    .messages-container {
      flex: 1;
      padding: 24px;
      overflow-y: auto;
      display: flex;
      flex-direction: column;
      gap: 16px;
    }
    .msg-wrap {
      display: flex;
      flex-direction: column;
      max-width: 65%;
      padding: 12px 18px;
      border-radius: 14px;
      background: #111a2e;
      border: 1px solid #1e293b;
      align-self: flex-start;
      animation: msgPop 0.25s cubic-bezier(0.16, 1, 0.3, 1);
    }
    .msg-wrap.me {
      align-self: flex-end;
      background: linear-gradient(135deg, #4f46e5, #6366f1);
      border: none;
    }
    .msg-header {
      display: flex;
      justify-content: space-between;
      gap: 32px;
      font-size: 11px;
      font-weight: 600;
      color: #64748b;
      margin-bottom: 6px;
    }
    .msg-wrap.me .msg-header {
      color: #a5b4fc;
    }
    .msg-body {
      font-size: 14px;
      line-height: 1.5;
    }
    @keyframes msgPop {
      from { transform: translateY(8px); opacity: 0; }
      to { transform: translateY(0); opacity: 1; }
    }
    .chat-input-bar {
      padding: 20px 24px;
      border-top: 1px solid #1e293b;
      background: rgba(15, 23, 42, 0.3);
      display: flex;
      gap: 12px;
      align-items: center;
    }
    .chat-input-bar input {
      flex: 1;
      background: #0d1527;
      border: 1px solid #1e293b;
      border-radius: 8px;
      padding: 12px 16px;
      color: #f3f4f6;
      font-size: 14px;
      outline: none;
      transition: all 0.2s;
    }
    .chat-input-bar input:focus {
      border-color: #6366f1;
      box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.2);
    }
    .chat-input-bar button {
      background: #6366f1;
      color: white;
      border: none;
      border-radius: 8px;
      padding: 12px 24px;
      font-weight: 600;
      font-size: 14px;
      cursor: pointer;
      transition: background-color 0.2s;
    }
    .chat-input-bar button:hover {
      background: #4f46e5;
    }
    .control-panel {
      width: 320px;
      background: #0c1322;
      border-left: 1px solid #1e293b;
      display: flex;
      flex-direction: column;
      overflow-y: auto;
    }
    .panel-card {
      padding: 24px;
      border-bottom: 1px solid #1e293b;
    }
    .panel-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 14px;
    }
    .panel-header h3 {
      font-size: 14px;
      font-weight: 600;
    }
    .badge {
      display: inline-block;
      padding: 4px 8px;
      border-radius: 6px;
      font-size: 10px;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }
    .badge.hit {
      background: rgba(16, 185, 129, 0.12);
      color: #10b981;
      border: 1px solid rgba(16, 185, 129, 0.25);
    }
    .badge.miss {
      background: rgba(239, 68, 68, 0.12);
      color: #ef4444;
      border: 1px solid rgba(239, 68, 68, 0.25);
    }
    .card-body {
      font-size: 12px;
      color: #64748b;
      line-height: 1.6;
    }
    .panel-btn {
      width: 100%;
      background: #111a2e;
      border: 1px solid #1e293b;
      color: #f3f4f6;
      padding: 10px;
      border-radius: 6px;
      cursor: pointer;
      font-weight: 600;
      font-size: 12px;
      transition: all 0.2s;
      margin-top: 12px;
    }
    .panel-btn:hover {
      background: #1e293b;
      border-color: #334155;
    }
    .panel-btn.accent {
      background: #a855f7;
      border: none;
    }
    .panel-btn.accent:hover {
      background: #9333ea;
    }
    .frame-box {
      background: #060911;
      border: 1px solid #111a2e;
      border-radius: 6px;
      padding: 12px;
      height: 100px;
      overflow-y: auto;
      font-family: monospace;
      font-size: 11px;
      color: #38bdf8;
      margin-top: 10px;
    }
    .alert-box {
      padding: 8px 12px;
      border-radius: 6px;
      font-size: 11px;
      margin-top: 10px;
      display: none;
      animation: alertSlide 0.2s ease-out;
    }
    .alert-box.error {
      background: rgba(239, 68, 68, 0.12);
      color: #fca5a5;
      border: 1px solid rgba(239, 68, 68, 0.25);
      display: block;
    }
    @keyframes alertSlide {
      from { transform: translateY(-4px); opacity: 0; }
      to { transform: translateY(0); opacity: 1; }
    }
  </style>
</head>
<body>

  <header>
    <div class="brand">
      <div class="logo">N</div>
      <div class="title-area">
        <h1>NestGo Chat Showcase</h1>
        <p>Enterprise stack v0.6.0 running live</p>
      </div>
    </div>
    <div class="status">
      <div class="dot" id="status-dot"></div>
      <span id="status-text">Disconnected</span>
    </div>
  </header>

  <div class="container">
    
    <div class="sidebar">
      <div>
        <div class="section-title">Rooms</div>
        <div class="room-list">
          <button class="room-btn active" onclick="switchRoom('general')"># general</button>
          <button class="room-btn" onclick="switchRoom('tech')"># tech-talk</button>
          <button class="room-btn" onclick="switchRoom('random')"># random</button>
        </div>
      </div>
      
      <div style="margin-top: auto;">
        <div class="section-title">Session User</div>
        <input type="text" id="username" class="chat-input-bar input" style="width:100%; padding:8px 12px; border-radius:6px;" value="User" placeholder="Enter name...">
      </div>
    </div>

    <div class="chat-area">
      <div class="messages-container" id="message-container">
        <!-- Messages stream -->
      </div>
      <div class="chat-input-bar">
        <input type="text" id="chat-message" placeholder="Type message..." onkeydown="if(event.key==='Enter') sendMessage()">
        <button onclick="sendMessage()">Send</button>
      </div>
    </div>

    <div class="control-panel">
      
      <!-- Caching Showcase -->
      <div class="panel-card">
        <div class="panel-header">
          <h3>Redis Caching</h3>
          <span class="badge" id="cache-badge">NO REQ</span>
        </div>
        <div class="card-body">
          The message list history endpoint (<code>/api/history</code>) is cached for 30s. Sending any new message invalidates this tag automatically.
        </div>
        <button class="panel-btn" onclick="fetchHistory()">Fetch /api/history</button>
      </div>

      <!-- Rate Limiting Showcase -->
      <div class="panel-card">
        <div class="panel-header">
          <h3>Redis Rate Limiter</h3>
        </div>
        <div class="card-body">
          Posting via the REST endpoint (<code>/api/message</code>) is restricted to 5 requests per 10 seconds. Check headers & block alerts.
        </div>
        <button class="panel-btn accent" onclick="triggerRESTPost()">Post via REST API</button>
        <div class="alert-box" id="rate-limit-alert">Rate limit exceeded! Please wait.</div>
      </div>

      <!-- gRPC Showcase -->
      <div class="panel-card">
        <div class="panel-header">
          <h3>gRPC Transport</h3>
        </div>
        <div class="card-body">
          Fires a gateway HTTP query that calls a backend gRPC service (on port :50051) and prints response.
        </div>
        <button class="panel-btn" onclick="pingGRPC()">Ping gRPC Backend</button>
        <div class="frame-box" id="grpc-logs">No RPC events logged yet.</div>
      </div>

      <!-- Frames Logging -->
      <div class="panel-card" style="flex:1; border-bottom:none;">
        <div class="section-title">WebSocket Frames</div>
        <div class="frame-box" id="ws-logs" style="height: calc(100% - 20px);">Ready for events...</div>
      </div>

    </div>

  </div>

  <script>
    let wsConn;
    let currentRoom = 'general';
    const usernameInput = document.getElementById('username');
    const wsLog = document.getElementById('ws-logs');
    const grpcLog = document.getElementById('grpc-logs');
    const msgContainer = document.getElementById('message-container');

    // Create custom random username
    usernameInput.value = 'User_' + Math.floor(Math.random() * 900 + 100);

    function connectWS() {
      const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = proto + '//' + window.location.host + '/ws';
      
      logWSFrame('Connecting to ' + wsUrl + '...');
      wsConn = new WebSocket(wsUrl);

      wsConn.onopen = () => {
        document.getElementById('status-dot').className = 'dot online';
        document.getElementById('status-text').innerText = 'Connected';
        logWSFrame('WebSocket open');
        
        // Initial history fetch
        fetchHistory();
      };

      wsConn.onclose = () => {
        document.getElementById('status-dot').className = 'dot';
        document.getElementById('status-text').innerText = 'Disconnected';
        logWSFrame('WebSocket closed. Retrying in 3s...');
        setTimeout(connectWS, 3000);
      };

      wsConn.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        logWSFrame('IN: ' + event.data);
        
        if (msg.event === 'new-message') {
          const payload = JSON.parse(msg.data);
          appendMessage(payload);
        } else if (msg.event === 'user-joined' || msg.event === 'user-left') {
          const payload = JSON.parse(msg.data);
          appendSystemMessage(payload.message);
        }
      };
    }

    function switchRoom(room) {
      if (room === currentRoom) return;
      
      // Update UI active state
      document.querySelectorAll('.room-btn').forEach(btn => {
        if (btn.innerText.includes(room)) btn.classList.add('active');
        else btn.classList.remove('active');
      });

      // Notify server
      sendWSFrame('join-room', { user: usernameInput.value, room: room });
      currentRoom = room;
      msgContainer.innerHTML = '';
      
      // Fetch history for new room
      fetchHistory();
    }

    function sendMessage() {
      const msgInput = document.getElementById('chat-message');
      const text = msgInput.value.trim();
      if (!text) return;

      sendWSFrame('send-message', {
        user: usernameInput.value,
        message: text,
        room: currentRoom
      });
      
      msgInput.value = '';
    }

    function sendWSFrame(event, data) {
      if (!wsConn || wsConn.readyState !== WebSocket.OPEN) return;
      const frame = { event: event, data: data };
      const raw = JSON.stringify(frame);
      wsConn.send(raw);
      logWSFrame('OUT: ' + raw);
    }

    function logWSFrame(text) {
      const time = new Date().toLocaleTimeString();
      wsLog.innerHTML = '[' + time + '] ' + text + '\n' + wsLog.innerHTML;
    }

    function appendMessage(msg) {
      if (msg.room !== currentRoom) return;
      const isMe = msg.user === usernameInput.value;
      const wrap = document.createElement('div');
      wrap.className = 'msg-wrap' + (isMe ? ' me' : '');
      wrap.innerHTML = '<div class="msg-header"><span>' + msg.user + '</span><span>' + new Date(msg.timestamp).toLocaleTimeString() + '</span></div><div class="msg-body">' + msg.message + '</div>';
      msgContainer.appendChild(wrap);
      msgContainer.scrollTop = msgContainer.scrollHeight;
    }

    function appendSystemMessage(text) {
      const wrap = document.createElement('div');
      wrap.style.cssText = 'text-align: center; color: #475569; font-size: 12px; margin: 8px 0; font-style: italic;';
      wrap.innerText = text;
      msgContainer.appendChild(wrap);
      msgContainer.scrollTop = msgContainer.scrollHeight;
    }

    // --- REST API Actions ---
    function fetchHistory() {
      fetch('/api/history')
        .then(res => {
          const cacheHeader = res.headers.get('X-Cache');
          const cacheBadge = document.getElementById('cache-badge');
          cacheBadge.className = 'badge ' + (cacheHeader === 'HIT' ? 'hit' : 'miss');
          cacheBadge.innerText = cacheHeader || 'MISS';

          return res.json();
        })
        .then(data => {
          msgContainer.innerHTML = '';
          if (data && data.data) {
            data.data.forEach(msg => {
              appendMessage(msg);
            });
          }
        });
    }

    function triggerRESTPost() {
      const alertBox = document.getElementById('rate-limit-alert');
      alertBox.className = 'alert-box';

      fetch('/api/message', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user: usernameInput.value + ' (REST)',
          message: 'REST API message trigger',
          room: currentRoom
        })
      })
      .then(res => {
        if (res.status === 429) {
          alertBox.className = 'alert-box error';
          alertBox.innerText = 'Rate limit exceeded! Wait ' + (res.headers.get('Retry-After') || '1') + 's.';
        } else {
          fetchHistory();
        }
      });
    }

    function pingGRPC() {
      const time = new Date().toLocaleTimeString();
      grpcLog.innerHTML = '[' + time + '] Calling PingRequest...';

      fetch('/api/ping-grpc?msg=Greetings%20from%20NestGo')
        .then(res => res.json())
        .then(data => {
          if (data && data.data && data.data.reply) {
            grpcLog.innerHTML = '[' + time + '] SUCCESS: ' + data.data.reply;
          } else {
            grpcLog.innerHTML = '[' + time + '] ERROR: ' + JSON.stringify(data);
          }
        })
        .catch(err => {
          grpcLog.innerHTML = '[' + time + '] FAILED: ' + err.message;
        });
    }

    // Connect immediately
    connectWS();
  </script>
</body>
</html>
`
