// Package tcp provides a TCP transport
package tcp

import (
	"bufio"
	"crypto/tls"
	"encoding/gob"
	"net"
	"time"

	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/transport"
	microTransport "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport"
)

type tcpTransport struct {
	opts microTransport.Options
}

type tcpTransportListener struct {
	listener net.Listener
	timeout  time.Duration
}

func (t *tcpTransportListener) Addr() string {
	return t.listener.Addr().String()
}

func (t *tcpTransportListener) Close() error {
	return t.listener.Close()
}

var baseConnStatic baseConn

func (t *tcpTransportListener) Accept(fn func(microTransport.Socket)) error {
	var tmpDelay time.Duration

	for {
		c, err := t.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tmpDelay == 0 {
					tmpDelay = 5 * time.Millisecond
				} else {
					tmpDelay *= 2
				}
				if max := 1 * time.Second; tmpDelay > max {
					tmpDelay = max
				}
				lager.Logger.Errorf(err, "http: Accept error: %v; retrying in %v", err, tmpDelay)
				time.Sleep(tmpDelay)
				continue
			}
			return err
		}

		encBuf := bufio.NewWriter(c)

		baseConnStatic = baseConn{
			r: bufio.NewReader(c),
			w: bufio.NewWriter(c),
			c: c,
		}

		sock := &tcpTransportSocket{
			baseConn: baseConnStatic,
			timeout:  t.timeout,
			conn:     c,
			encBuf:   encBuf,
			enc:      gob.NewEncoder(encBuf),
			dec:      gob.NewDecoder(c),
		}

		go func() {
			// TODO: design a better error response strategy
			defer func() {
				if r := recover(); r != nil {
					sock.Close()
				}
			}()

			fn(sock)
		}()
	}
}

func (t *tcpTransport) Dial(addr string, opts ...microTransport.DialOption) (microTransport.Client, error) {
	dopts := microTransport.DialOptions{
		Timeout: microTransport.DefaultDialTimeout,
	}

	if opts != nil {
		for _, o := range opts {
			o(&dopts)
		}
	}

	var c net.Conn
	var err error

	// TODO: support dial options
	if t.opts.Secure || t.opts.TLSConfig != nil {
		config := t.opts.TLSConfig
		if config == nil {
			config = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		c, err = tls.DialWithDialer(&net.Dialer{Timeout: dopts.Timeout}, "tcp", addr, config)
	} else {
		c, err = net.DialTimeout("tcp", addr, dopts.Timeout)
	}

	if err != nil {
		return nil, err
	}

	encBuf := bufio.NewWriter(c)

	return &tcpTransportClient{
		baseConn: baseConn{
			r: bufio.NewReader(c),
			w: bufio.NewWriter(c),
			c: c,
		},
		dialOpts: dopts,
		conn:     c,
		encBuf:   encBuf,
		enc:      gob.NewEncoder(encBuf),
		dec:      gob.NewDecoder(c),
		timeout:  t.opts.Timeout,
	}, nil
}

func (t *tcpTransport) Listen(addr string, opts ...microTransport.ListenOption) (microTransport.Listener, error) {
	var lopts microTransport.ListenOptions
	for _, o := range opts {
		o(&lopts)
	}

	var l net.Listener
	var err error

	// TODO: use listen options
	if t.opts.Secure || t.opts.TLSConfig != nil {
		config := t.opts.TLSConfig

		fn := func(addr string) (net.Listener, error) {
			return tls.Listen("tcp", addr, config)
		}

		l, err = microTransport.Listen(addr, fn)
	} else {
		fn := func(addr string) (net.Listener, error) {
			return net.Listen("tcp", addr)
		}

		l, err = microTransport.Listen(addr, fn)
	}

	if err != nil {
		return nil, err
	}

	return &tcpTransportListener{
		timeout:  t.opts.Timeout,
		listener: l,
	}, nil
}

var protocol = "tcp"

func (t *tcpTransport) String() string {
	return protocol
}

//NewTransport is a function
func NewTransport(opts ...microTransport.Option) microTransport.Transport {
	var topts microTransport.Options
	for _, o := range opts {
		o(&topts)
	}
	return &tcpTransport{opts: topts}
}

func init() {
	transport.InstallPlugin(protocol, NewTransport)
}
