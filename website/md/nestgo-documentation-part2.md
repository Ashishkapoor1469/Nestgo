# NestGo Framework - Complete Documentation (Part 2)

---

# Advanced Features

## Guards

### What are Guards?

Guards determine whether a request should be handled by the route handler. They're used for:
- Authentication
- Authorization (RBAC)
- Feature flags
- Request validation

### Guard Interface

```go
type Guard interface {
    CanActivate(ctx *Context) (bool, error)
}
```

### Creating a Guard

```go
package guards

import "github.com/nestgo/nestgo/common"

type AuthGuard struct {
    tokenService *TokenService
}

func NewAuthGuard(tokenService *TokenService) *AuthGuard {
    return &AuthGuard{tokenService: tokenService}
}

func (g *AuthGuard) CanActivate(ctx *common.Context) (bool, error) {
    token := ctx.Header("Authorization")
    if token == "" {
        return false, common.NewUnauthorizedError("Missing authentication token")
    }
    
    // Validate token
    user, err := g.tokenService.Validate(token)
    if err != nil {
        return false, common.NewUnauthorizedError("Invalid token")
    }
    
    // Store user in context
    ctx.Set("user", user)
    return true, nil
}
```

### Applying Guards

**Global Guard:**
```go
app.UseGuard(guards.NewAuthGuard(tokenService))
```

**Controller-level:**
```go
func (c *UsersController) Guards() []common.Guard {
    return []common.Guard{
        guards.NewAuthGuard(),
    }
}
```

**Route-level:**
```go
{
    Method: "DELETE",
    Path: "/users/{id}",
    Handler: c.Delete,
    Guards: []common.Guard{
        guards.NewAuthGuard(),
        guards.NewAdminGuard(),
    },
}
```

### RBAC Guard

```go
type Role string

const (
    RoleUser  Role = "user"
    RoleAdmin Role = "admin"
    RoleSuper Role = "superadmin"
)

type RolesGuard struct {
    allowedRoles []Role
}

func NewRolesGuard(roles ...Role) *RolesGuard {
    return &RolesGuard{allowedRoles: roles}
}

func (g *RolesGuard) CanActivate(ctx *common.Context) (bool, error) {
    user := ctx.Get("user").(*User)
    if user == nil {
        return false, common.NewUnauthorizedError("User not authenticated")
    }
    
    for _, allowed := range g.allowedRoles {
        if user.Role == allowed {
            return true, nil
        }
    }
    
    return false, common.NewForbiddenError("Insufficient permissions")
}

// Usage
{
    Method: "POST",
    Path: "/admin/users",
    Handler: c.CreateUser,
    Guards: []common.Guard{
        guards.NewRolesGuard(RoleAdmin, RoleSuper),
    },
}
```

### Feature Flag Guard

```go
type FeatureFlagGuard struct {
    flagService *FeatureFlagService
    featureName string
}

func NewFeatureFlagGuard(service *FeatureFlagService, feature string) *FeatureFlagGuard {
    return &FeatureFlagGuard{
        flagService: service,
        featureName: feature,
    }
}

func (g *FeatureFlagGuard) CanActivate(ctx *common.Context) (bool, error) {
    enabled := g.flagService.IsEnabled(g.featureName)
    if !enabled {
        return false, common.NewNotFoundError("Feature not available")
    }
    return true, nil
}
```

### Rate Limiting Guard

```go
type RateLimitGuard struct {
    limiter *RateLimiter
    max     int
    window  time.Duration
}

func (g *RateLimitGuard) CanActivate(ctx *common.Context) (bool, error) {
    ip := ctx.Request.RemoteAddr
    
    allowed := g.limiter.Allow(ip, g.max, g.window)
    if !allowed {
        return false, common.NewTooManyRequestsError("Rate limit exceeded")
    }
    
    return true, nil
}
```

---

## Interceptors

### What are Interceptors?

Interceptors are AOP-style components that can:
- Transform request/response data
- Bind extra logic before/after method execution
- Extend basic functionality
- Cache responses
- Log method execution time

### Interceptor Interface

```go
type Interceptor interface {
    Intercept(ctx *Context, next HandlerFunc) error
}
```

### Creating an Interceptor

```go
package interceptors

type LoggingInterceptor struct {
    logger *slog.Logger
}

func NewLoggingInterceptor(logger *slog.Logger) *LoggingInterceptor {
    return &LoggingInterceptor{logger: logger}
}

func (i *LoggingInterceptor) Intercept(ctx *common.Context, next common.HandlerFunc) error {
    start := time.Now()
    
    i.logger.Info("Request started",
        "method", ctx.Request.Method,
        "path", ctx.Request.URL.Path,
    )
    
    // Execute handler
    err := next(ctx)
    
    duration := time.Since(start)
    i.logger.Info("Request completed",
        "duration", duration,
        "status", ctx.Response.Status,
    )
    
    return err
}
```

### Response Transformation Interceptor

```go
type TransformInterceptor struct{}

func (i *TransformInterceptor) Intercept(ctx *common.Context, next common.HandlerFunc) error {
    // Execute handler
    err := next(ctx)
    if err != nil {
        return err
    }
    
    // Transform response
    data := ctx.Get("response")
    wrapped := map[string]interface{}{
        "success": true,
        "data":    data,
        "timestamp": time.Now(),
    }
    
    ctx.Set("response", wrapped)
    return nil
}
```

### Cache Interceptor

```go
type CacheInterceptor struct {
    cache *redis.Client
    ttl   time.Duration
}

func (i *CacheInterceptor) Intercept(ctx *common.Context, next common.HandlerFunc) error {
    // Only cache GET requests
    if ctx.Request.Method != "GET" {
        return next(ctx)
    }
    
    cacheKey := generateCacheKey(ctx)
    
    // Check cache
    cached, err := i.cache.Get(ctx.Request.Context(), cacheKey).Result()
    if err == nil {
        ctx.Set("response", cached)
        ctx.SetHeader("X-Cache", "HIT")
        return nil
    }
    
    // Execute handler
    err = next(ctx)
    if err != nil {
        return err
    }
    
    // Cache response
    response := ctx.Get("response")
    i.cache.Set(ctx.Request.Context(), cacheKey, response, i.ttl)
    ctx.SetHeader("X-Cache", "MISS")
    
    return nil
}
```

### Timeout Interceptor

```go
type TimeoutInterceptor struct {
    timeout time.Duration
}

func (i *TimeoutInterceptor) Intercept(ctx *common.Context, next common.HandlerFunc) error {
    done := make(chan error, 1)
    
    go func() {
        done <- next(ctx)
    }()
    
    select {
    case err := <-done:
        return err
    case <-time.After(i.timeout):
        return common.NewRequestTimeoutError("Request timeout")
    }
}
```

### Applying Interceptors

```go
// Global
app.UseInterceptor(interceptors.NewLoggingInterceptor(logger))

// Controller-level
func (c *UsersController) Interceptors() []common.Interceptor {
    return []common.Interceptor{
        interceptors.NewCacheInterceptor(cache, 5*time.Minute),
    }
}

// Route-level
{
    Method: "GET",
    Path: "/expensive-operation",
    Handler: c.ExpensiveOp,
    Interceptors: []common.Interceptor{
        interceptors.NewTimeoutInterceptor(30 * time.Second),
    },
}
```

---

## Pipes

### What are Pipes?

Pipes transform or validate input data before it reaches the route handler.

### Built-in Pipes

```go
// Parse and validate UUID
pipes.UUID()

// Parse integer with validation
pipes.Int(1, 100) // min: 1, max: 100

// Parse boolean
pipes.Bool()

// Trim whitespace
pipes.Trim()

// Parse JSON
pipes.JSON(&dto)
```

### Custom Pipe

```go
type ParseDatePipe struct {
    format string
}

func (p *ParseDatePipe) Transform(value string) (interface{}, error) {
    date, err := time.Parse(p.format, value)
    if err != nil {
        return nil, common.NewBadRequestError("Invalid date format")
    }
    return date, nil
}

// Usage
func (c *Controller) GetOrders(ctx *common.Context) error {
    startDateStr := ctx.Query("startDate")
    
    pipe := &ParseDatePipe{format: "2006-01-02"}
    startDate, err := pipe.Transform(startDateStr)
    if err != nil {
        return err
    }
    
    orders, err := c.service.GetOrders(startDate.(time.Time))
    return ctx.JSON(200, orders)
}
```

---

# Data Management

## Database Integration

### Supported Databases

- PostgreSQL
- MySQL
- SQLite
- MongoDB
- Redis

### Database Module

```go
package database

import (
    "database/sql"
    "github.com/nestgo/nestgo/common"
    _ "github.com/lib/pq" // PostgreSQL driver
)

type DatabaseModule struct{}

func (m *DatabaseModule) Module() common.ModuleConfig {
    db, err := sql.Open("postgres", "postgres://user:pass@localhost/dbname")
    if err != nil {
        panic(err)
    }
    
    if err := db.Ping(); err != nil {
        panic(err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    return common.ModuleConfig{
        Name:   "database",
        Global: true,
        Providers: []di.Provider{
            {Instance: db, Scope: di.Singleton},
        },
        Exports: []di.Provider{
            {Instance: db},
        },
    }
}
```

### Using Database in Services

```go
type UsersRepository struct {
    db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
    return &UsersRepository{db: db}
}

func (r *UsersRepository) FindAll(ctx context.Context) ([]*User, error) {
    rows, err := r.db.QueryContext(ctx, "SELECT id, email, name FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []*User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Email, &user.Name); err != nil {
            return nil, err
        }
        users = append(users, &user)
    }
    
    return users, nil
}

func (r *UsersRepository) Create(ctx context.Context, user *User) error {
    query := `INSERT INTO users (email, name, password_hash) VALUES ($1, $2, $3) RETURNING id`
    return r.db.QueryRowContext(ctx, query, user.Email, user.Name, user.PasswordHash).Scan(&user.ID)
}
```

---

## ORMs & Query Builders

### GORM Integration

```go
package database

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func NewGormDB(config *Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
        config.DB.Host, config.DB.User, config.DB.Password, config.DB.Name, config.DB.Port)
    
    return gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
}

// Usage in repository
type UsersRepository struct {
    db *gorm.DB
}

func (r *UsersRepository) FindAll(ctx context.Context) ([]*User, error) {
    var users []*User
    result := r.db.WithContext(ctx).Find(&users)
    return users, result.Error
}

func (r *UsersRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}
```

### SQLx Integration

```go
import "github.com/jmoiron/sqlx"

type UsersRepository struct {
    db *sqlx.DB
}

func (r *UsersRepository) FindAll(ctx context.Context) ([]*User, error) {
    var users []*User
    err := r.db.SelectContext(ctx, &users, "SELECT * FROM users")
    return users, err
}

func (r *UsersRepository) FindByID(ctx context.Context, id string) (*User, error) {
    var user User
    err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
    return &user, err
}
```

---

## Migrations

### Using golang-migrate

```go
package migrations

import (
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseURL string) error {
    m, err := migrate.New(
        "file://migrations",
        databaseURL,
    )
    if err != nil {
        return err
    }
    
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }
    
    return nil
}
```

**Migration Files:**

`migrations/001_create_users_table.up.sql`:
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

`migrations/001_create_users_table.down.sql`:
```sql
DROP TABLE IF EXISTS users;
```

### CLI Integration

```bash
# Create migration
nestgo migration:create create_users_table

# Run migrations
nestgo migration:run

# Rollback
nestgo migration:rollback

# Reset database
nestgo migration:reset
```

---

## Transactions

### Manual Transaction Handling

```go
func (s *OrdersService) PlaceOrder(ctx context.Context, dto *CreateOrderDTO) error {
    // Begin transaction
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // Create order
    orderID, err := s.createOrder(ctx, tx, dto)
    if err != nil {
        return err
    }
    
    // Create order items
    for _, item := range dto.Items {
        if err := s.createOrderItem(ctx, tx, orderID, item); err != nil {
            return err
        }
        
        // Update product stock
        if err := s.updateStock(ctx, tx, item.ProductID, -item.Quantity); err != nil {
            return err
        }
    }
    
    // Commit transaction
    return tx.Commit()
}
```

### Transaction Wrapper

```go
func (r *Repository) WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    if err := fn(tx); err != nil {
        return err
    }
    
    return tx.Commit()
}

// Usage
err := r.WithTransaction(ctx, func(tx *sql.Tx) error {
    if err := createUser(tx, user); err != nil {
        return err
    }
    if err := createProfile(tx, profile); err != nil {
        return err
    }
    return nil
})
```

### GORM Transactions

```go
func (s *Service) TransferMoney(ctx context.Context, from, to string, amount float64) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Deduct from sender
        if err := tx.Model(&Account{}).Where("id = ?", from).
            Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
            return err
        }
        
        // Add to receiver
        if err := tx.Model(&Account{}).Where("id = ?", to).
            Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
            return err
        }
        
        // Create transaction record
        return tx.Create(&Transaction{
            FromAccount: from,
            ToAccount:   to,
            Amount:      amount,
        }).Error
    })
}
```

---

# Real-time & Events

## WebSockets

### WebSocket Gateway

```go
package ws

import (
    "github.com/gorilla/websocket"
    "github.com/nestgo/nestgo/common"
)

type ChatGateway struct {
    clients   map[*websocket.Conn]bool
    broadcast chan Message
    upgrader  websocket.Upgrader
}

func NewChatGateway() *ChatGateway {
    return &ChatGateway{
        clients:   make(map[*websocket.Conn]bool),
        broadcast: make(chan Message),
        upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool {
                return true // Configure CORS properly
            },
        },
    }
}

func (g *ChatGateway) HandleConnection(ctx *common.Context) error {
    conn, err := g.upgrader.Upgrade(ctx.Response, ctx.Request, nil)
    if err != nil {
        return err
    }
    defer conn.Close()
    
    g.clients[conn] = true
    defer delete(g.clients, conn)
    
    for {
        var msg Message
        err := conn.ReadJSON(&msg)
        if err != nil {
            break
        }
        
        g.broadcast <- msg
    }
    
    return nil
}

func (g *ChatGateway) BroadcastMessages() {
    for {
        msg := <-g.broadcast
        for client := range g.clients {
            err := client.WriteJSON(msg)
            if err != nil {
                client.Close()
                delete(g.clients, client)
            }
        }
    }
}
```

### WebSocket Module

```go
type WebSocketModule struct{}

func (m *WebSocketModule) Module() common.ModuleConfig {
    gateway := NewChatGateway()
    
    // Start broadcast goroutine
    go gateway.BroadcastMessages()
    
    return common.ModuleConfig{
        Name: "websocket",
        Controllers: []common.Controller{
            NewWSController(gateway),
        },
    }
}
```

### WebSocket Controller

```go
type WSController struct {
    gateway *ChatGateway
}

func (c *WSController) Prefix() string {
    return "/ws"
}

func (c *WSController) Routes() []common.Route {
    return []common.Route{
        {Method: "GET", Path: "/chat", Handler: c.gateway.HandleConnection},
    }
}
```

### Client Example

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/chat');

ws.onopen = () => {
    console.log('Connected');
    ws.send(JSON.stringify({ type: 'join', user: 'John' }));
};

ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    console.log('Received:', message);
};

ws.send(JSON.stringify({
    type: 'message',
    content: 'Hello World!',
    user: 'John'
}));
```

---

## Event Emitters

### Event Bus

```go
package events

type EventBus struct {
    subscribers map[string][]EventHandler
    mu          sync.RWMutex
}

type EventHandler func(event Event) error

type Event struct {
    Type    string
    Payload interface{}
}

func NewEventBus() *EventBus {
    return &EventBus{
        subscribers: make(map[string][]EventHandler),
    }
}

func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    
    eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

func (eb *EventBus) Emit(event Event) error {
    eb.mu.RLock()
    handlers := eb.subscribers[event.Type]
    eb.mu.RUnlock()
    
    for _, handler := range handlers {
        if err := handler(event); err != nil {
            return err
        }
    }
    
    return nil
}

func (eb *EventBus) EmitAsync(event Event) {
    go eb.Emit(event)
}
```

### Using Events

```go
// Register event handlers
eventBus.Subscribe("user.created", func(event events.Event) error {
    user := event.Payload.(*User)
    log.Printf("User created: %s", user.Email)
    
    // Send welcome email
    return emailService.SendWelcome(user.Email)
})

eventBus.Subscribe("user.created", func(event events.Event) error {
    user := event.Payload.(*User)
    // Create user profile
    return profileService.CreateDefault(user.ID)
})

// Emit events
func (s *UsersService) Create(ctx context.Context, dto *CreateUserDTO) (*User, error) {
    user, err := s.repository.Create(ctx, dto)
    if err != nil {
        return nil, err
    }
    
    // Emit event asynchronously
    s.eventBus.EmitAsync(events.Event{
        Type:    "user.created",
        Payload: user,
    })
    
    return user, nil
}
```

---

## Message Queues

### RabbitMQ Integration

```go
package queue

import "github.com/streadway/amqp"

type RabbitMQ struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }
    
    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }
    
    return &RabbitMQ{conn: conn, channel: ch}, nil
}

func (mq *RabbitMQ) Publish(queue string, message []byte) error {
    _, err := mq.channel.QueueDeclare(queue, true, false, false, false, nil)
    if err != nil {
        return err
    }
    
    return mq.channel.Publish("", queue, false, false, amqp.Publishing{
        ContentType: "application/json",
        Body:        message,
    })
}

func (mq *RabbitMQ) Consume(queue string, handler func([]byte) error) error {
    msgs, err := mq.channel.Consume(queue, "", false, false, false, false, nil)
    if err != nil {
        return err
    }
    
    for msg := range msgs {
        if err := handler(msg.Body); err != nil {
            msg.Nack(false, true) // Requeue on error
        } else {
            msg.Ack(false)
        }
    }
    
    return nil
}
```

### Queue Consumer

```go
type EmailQueueConsumer struct {
    emailService *EmailService
    mq           *RabbitMQ
}

func (c *EmailQueueConsumer) Start() {
    c.mq.Consume("emails", func(body []byte) error {
        var email EmailMessage
        if err := json.Unmarshal(body, &email); err != nil {
            return err
        }
        
        return c.emailService.Send(email.To, email.Subject, email.Body)
    })
}
```

---

# Background Processing

## Task Scheduling

### Cron Module

```go
package jobs

import "github.com/robfig/cron/v3"

type CronModule struct {
    cron *cron.Cron
}

func NewCronModule() *CronModule {
    return &CronModule{
        cron: cron.New(),
    }
}

func (m *CronModule) AddJob(spec string, cmd func()) error {
    _, err := m.cron.AddFunc(spec, cmd)
    return err
}

func (m *CronModule) Start() {
    m.cron.Start()
}

func (m *CronModule) Stop() {
    m.cron.Stop()
}
```

### Scheduled Tasks

```go
// Register cron jobs
cronModule.AddJob("0 0 * * *", func() {
    // Run daily at midnight
    log.Println("Running daily cleanup")
    cleanupService.RemoveExpiredData()
})

cronModule.AddJob("*/5 * * * *", func() {
    // Run every 5 minutes
    log.Println("Syncing data")
    syncService.Sync()
})

cronModule.AddJob("0 9 * * MON", func() {
    // Run every Monday at 9 AM
    reportService.GenerateWeeklyReport()
})

cronModule.Start()
```

---

## Background Jobs

### Job Queue

```go
package jobs

type Job struct {
    ID      string
    Type    string
    Payload interface{}
    Retries int
}

type JobQueue struct {
    jobs    chan Job
    workers int
}

func NewJobQueue(workers int) *JobQueue {
    return &JobQueue{
        jobs:    make(chan Job, 1000),
        workers: workers,
    }
}

func (jq *JobQueue) Start(handlers map[string]JobHandler) {
    for i := 0; i < jq.workers; i++ {
        go jq.worker(handlers)
    }
}

func (jq *JobQueue) worker(handlers map[string]JobHandler) {
    for job := range jq.jobs {
        handler, exists := handlers[job.Type]
        if !exists {
            log.Printf("No handler for job type: %s", job.Type)
            continue
        }
        
        if err := handler(job); err != nil {
            log.Printf("Job failed: %v", err)
            if job.Retries < 3 {
                job.Retries++
                jq.jobs <- job // Retry
            }
        }
    }
}

func (jq *JobQueue) Enqueue(job Job) {
    jq.jobs <- job
}
```

### Job Handlers

```go
type JobHandler func(Job) error

// Email job handler
func SendEmailHandler(job Job) error {
    email := job.Payload.(EmailData)
    return emailService.Send(email.To, email.Subject, email.Body)
}

// Image processing handler
func ProcessImageHandler(job Job) error {
    data := job.Payload.(ImageData)
    return imageService.Resize(data.Path, data.Width, data.Height)
}

// Register handlers
handlers := map[string]JobHandler{
    "send_email":    SendEmailHandler,
    "process_image": ProcessImageHandler,
}

jobQueue.Start(handlers)

// Enqueue jobs
jobQueue.Enqueue(Job{
    Type: "send_email",
    Payload: EmailData{
        To:      "user@example.com",
        Subject: "Welcome",
        Body:    "Welcome to our app!",
    },
})
```

---

This completes Part 2 of the documentation. Would you like me to continue with Part 3 covering Security, Testing, CLI, Deployment, and API Reference?
