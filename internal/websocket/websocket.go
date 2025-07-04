// Package websocket
// Author: Jon Brown
// Date: Mar 30, 2024
// URL: https://github.com/brojonat/websocket

package websocket

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"slices"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// DefaultSetupClient is an example implementation of a function that sets up a
// websocket connection.
func DefaultSetupConn(c *websocket.Conn) {
	pw := 60 * time.Second
	c.SetReadLimit(512)
	c.SetReadDeadline(time.Now().Add(pw))
	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(pw))
		return nil
	})
}

func DefaultUpgrader(origins []string) websocket.Upgrader {
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return slices.Contains(origins, r.Header.Get("Origin"))
	}
	return upgrader
}

// Client is an interface for reading from and writing to a websocket
// connection. It is designed to be used as a middleman between a service and a
// client websocket connection.
type Client interface {
	io.Writer
	io.Closer

	// WriteForever is responsible for writing messages to the client (including
	// the regularly spaced ping messages)
	WriteForever(context.Context, func(Client), time.Duration)

	// ReadForever is responsible for reading messages from the client, and passing
	// them to the message handlers
	ReadForever(context.Context, func(Client), ...MessageHandler)

	// SetLogger allows consumers to inject their own logging dependencies
	SetLogger(any) error

	// Log allows implementors to use their own logging dependencies
	Log(int, string, ...any)

	Conn() *websocket.Conn

	// Wait blocks until the client is done processing messages
	Wait()
}

type MessageHandler func(Client, []byte)

// ServeWS upgrades HTTP connections to WebSocket, creates the Client, calls the
// onCreate callback, and starts goroutines that handle reading (writing)
// from (to) the client.
func ServeWS(
	// upgrader upgrades the connection
	upgrader websocket.Upgrader,
	// connSetup is called on the upgraded WebSocket connection to configure
	// the connection
	connSetup func(*websocket.Conn),
	// clientFactory is a function that takes a connection and returns a new Client
	clientFactory func(*websocket.Conn) Client,
	// onCreate is a function to call once the Client is created (e.g.,
	// store it in a some collection on the service for later reference)
	onCreate func(context.Context, context.CancelFunc, Client),
	// onDestroy is a function to call after the WebSocket connection is closed
	// (e.g., remove it from the collection on the service)
	onDestroy func(Client),
	// ping is the interval at which ping messages are aren't sent
	ping time.Duration,
	// msgHandlers are callbacks that handle messages received from the client
	msgHandlers []MessageHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			// if Upgrade fails it closes the connection, so just return
			return
		}
		connSetup(conn)
		client := clientFactory(conn)
		ctx := context.Background()
		ctx, cf := context.WithCancel(ctx)
		onCreate(ctx, cf, client)

		// all writes will happen in this goroutine, ensuring only one write on
		// the connection at a time
		go client.WriteForever(ctx, onDestroy, ping)

		// all reads will happen in this goroutine, ensuring only one reader on
		// the connection at a time
		go client.ReadForever(ctx, onDestroy, msgHandlers...)
	}
}

type client struct {
	lock   *sync.RWMutex
	wg     *sync.WaitGroup
	conn   *websocket.Conn
	egress chan []byte
	logger *slog.Logger
}

// Conn implements Client.
func (c *client) Conn() *websocket.Conn {
	return c.conn
}

// NewClient returns a new Client from a *websocket.Conn. This can be passed to
// ServeWS as the client factory arg.
func NewClient(c *websocket.Conn) Client {
	// add 2 to the wait group for the read/write goroutines
	wg := &sync.WaitGroup{}
	wg.Add(2)
	return &client{
		lock:   &sync.RWMutex{},
		wg:     wg,
		conn:   c,
		egress: make(chan []byte, 32),
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

// Write implements the Writer interface.
func (c *client) Write(p []byte) (int, error) {
	c.egress <- p
	return len(p), nil
}

// Close implements the Closer interface. Note the behavior of calling Close()
// multiple times is undefined; this implementation swallows all errors.
func (c *client) Close() error {
	c.conn.WriteControl(websocket.CloseMessage, []byte{}, time.Time{})
	c.conn.Close()
	return nil
}

// WriteForever serially processes messages from the egress channel and writes them
// to the client, ensuring that all writes to the underlying connection are
// performed here.
func (c *client) WriteForever(ctx context.Context, onDestroy func(Client), ping time.Duration) {
	pingTicker := time.NewTicker(ping)
	defer func() {
		c.wg.Done()
		pingTicker.Stop()
		onDestroy(c)
	}()

	for {
		select {
		case <-ctx.Done():
			c.conn.WriteMessage(websocket.CloseMessage, nil)
			return
		case msgBytes, ok := <-c.egress:
			// ok will be false in case the egress channel is closed
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			// write a message to the connection
			if err := c.conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
				c.Log(int(slog.LevelError), fmt.Sprintf("error writing message: %v", err))
				return
			}
		case <-pingTicker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.Log(int(slog.LevelError), fmt.Sprintf("error writing ping: %v", err))
				return
			}
		}
	}
}

// ReadForever serially processes messages from the client and passes them to the
// supplied message handlers in their own goroutine. Each message will be processed
// serially, but the handlers are executed concurrently.
func (c *client) ReadForever(ctx context.Context, onDestroy func(Client), handlers ...MessageHandler) {
	defer func() {
		c.wg.Done()
		onDestroy(c)
	}()

	ingress := make(chan []byte)
	errCancel := make(chan error)
	loop := true

	// read forever and push into ingress
	go func() {
		for {
			_, payload, err := c.conn.ReadMessage()
			if err != nil {
				errCancel <- err
				return
			}
			ingress <- payload
		}
	}()

	// read while waiting for shutdown signals
	for loop {
		select {
		case <-ctx.Done():
			c.Log(0, "read loop cancelled, shutting down")
			loop = false
		case err := <-errCancel:
			c.Log(0, "client connection closed in read loop")
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				c.Log(0, "read loop encountered error, shutting down", "error", err.Error())
			}
			loop = false
		case payload := <-ingress:
			// handle the message so each message is processed serially (though note
			// that handlers are called concurrently)
			var wg sync.WaitGroup
			wg.Add(len(handlers))
			for _, h := range handlers {
				go func(h func(Client, []byte)) {
					h(c, payload)
					wg.Done()
				}(h)
			}
			wg.Wait()
		}
	}
}

func (c *client) SetLogger(v any) error {
	l, ok := v.(*slog.Logger)
	if !ok {
		return fmt.Errorf("bad logger value supplied")
	}
	c.logger = l
	return nil
}

func (c *client) Log(level int, s string, args ...any) {
	_, f, l, ok := runtime.Caller(1)
	if ok {
		args = append(args, "caller_source", fmt.Sprintf("%s %d", f, l))
	}
	switch level {
	case int(slog.LevelDebug):
		c.logger.Debug(s, args...)
	case int(slog.LevelInfo):
		c.logger.Info(s, args...)
	case int(slog.LevelWarn):
		c.logger.Warn(s, args...)
	case int(slog.LevelError):
		c.logger.Error(s, args...)
	}
}

// Done blocks until the read/write goroutines have completed
func (c *client) Wait() {
	c.wg.Wait()
}

// Manager maintains a set of Clients.
type Manager interface {
	Clients() []Client
	RegisterClient(context.Context, context.CancelFunc, Client)
	UnregisterClient(Client)
	Run(context.Context)
}

type manager struct {
	mu         *sync.RWMutex
	clients    map[Client]context.CancelFunc
	register   chan regreq
	unregister chan regreq
}

type regreq struct {
	context context.Context
	cancel  context.CancelFunc
	client  Client
	done    chan struct{}
}

func NewManager() Manager {
	return &manager{
		mu:         &sync.RWMutex{},
		clients:    make(map[Client]context.CancelFunc),
		register:   make(chan regreq),
		unregister: make(chan regreq),
	}
}

// Clients returns the currently managed Clients.
func (m *manager) Clients() []Client {
	res := []Client{}
	m.mu.RLock()
	defer m.mu.RUnlock()

	for c := range m.clients {
		res = append(res, c)
	}
	return res
}

// RegisterClient adds the Client to the Manager's store.
func (m *manager) RegisterClient(ctx context.Context, cf context.CancelFunc, c Client) {
	done := make(chan struct{})
	rr := regreq{
		context: ctx,
		cancel:  cf,
		client:  c,
		done:    done,
	}
	m.register <- rr
	<-done
}

// UnregisterClient removes the Client from the Manager's store.
func (m *manager) UnregisterClient(c Client) {
	done := make(chan struct{})
	rr := regreq{
		client: c,
		done:   done,
	}
	m.unregister <- rr
	<-done
}

// Run runs in its own goroutine processing (un)registration requests.
func (m *manager) Run(ctx context.Context) {
	// helper fn for cleaning up client
	cleanupClient := func(c Client) {
		cancel, ok := m.clients[c]
		if ok {
			cancel()
		}
		delete(m.clients, c)
		c.Close()
	}

	for {
		select {
		case <-ctx.Done():
			m.mu.Lock()
			for client := range m.clients {
				cleanupClient(client)
			}
			m.mu.Unlock()
		case rr := <-m.register:
			m.mu.Lock()
			m.clients[rr.client] = rr.cancel
			m.mu.Unlock()
			rr.done <- struct{}{}

		case rr := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[rr.client]; ok {
				cleanupClient(rr.client)
			}
			m.mu.Unlock()
			rr.done <- struct{}{}
		}
	}
}

// Broadcaster is an example implementation of Manager that has a
// Broadcast method that writes the supplied message to all clients.
type broadcaster struct {
	*manager
}

func NewBroadcaster() Broadcaster {
	m := manager{
		mu:         &sync.RWMutex{},
		clients:    make(map[Client]context.CancelFunc),
		register:   make(chan regreq),
		unregister: make(chan regreq),
	}
	return &broadcaster{
		manager: &m,
	}
}

func (bb *broadcaster) Broadcast(b []byte) error {
	bb.mu.Lock()
	defer bb.mu.Unlock()
	var errs []error
	for w := range bb.clients {
		_, err := w.Write(b)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

type Broadcaster interface {
	Manager
	Broadcast([]byte) error
}
