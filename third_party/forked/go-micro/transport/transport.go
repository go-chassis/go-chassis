package transport

// Socket wrapp codec and net.conn
type Socket interface {
	Send(Header []byte, Body []byte, Metadata map[string]string, ID int) error
	Recv() ([]byte, []byte, map[string]string, int, error)
	Close() error
}

// Client client interface to implement send, receive, and close methods
type Client interface {
	Socket
}

// Listener interface for to listen the requests and accept the connections
type Listener interface {
	Addr() string
	Close() error
	Accept(func(Socket)) error
}

// Transport is an interface for communication between services.
// It uses socket send/recv semantics and could be implemented with various
// protocol: HTTP, RabbitMQ, NATS, ...
type Transport interface {
	//return a new client
	Dial(addr string, opts ...DialOption) (Client, error)
	Listen(addr string, opts ...ListenOption) (Listener, error)
	String() string
}
