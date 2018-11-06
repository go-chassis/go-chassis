package grpc

import (
	"context"
	"github.com/go-chassis/go-chassis/core/client"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/invocation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func init() {
	client.InstallPlugin("grpc", New)
}

//Client is grpc client holder
type Client struct {
	c    *grpc.ClientConn
	opts client.Options
}

//New create new grpc client
func New(opts client.Options) (client.ProtocolClient, error) {
	var err error
	var conn *grpc.ClientConn
	if opts.TLSConfig == nil {
		conn, err = grpc.Dial(opts.Endpoint, grpc.WithInsecure())
	} else {
		conn, err = grpc.Dial(opts.Endpoint,
			grpc.WithTransportCredentials(credentials.NewTLS(opts.TLSConfig)))
	}
	if err != nil {
		return nil, err
	}
	return &Client{
		c:    conn,
		opts: opts,
	}, nil
}

//TransformContext will deliver header in chassis context key to grpc context key
func TransformContext(ctx context.Context) context.Context {
	m := common.FromContext(ctx)
	md := metadata.New(m)
	return metadata.NewOutgoingContext(ctx, md)
}

//Call remote server
func (c *Client) Call(ctx context.Context, addr string, inv *invocation.Invocation, rsp interface{}) error {
	ctx = TransformContext(ctx)
	return c.c.Invoke(ctx, "/"+inv.SchemaID+"/"+inv.OperationID, inv.Args, rsp)
}

//String return name
func (c *Client) String() string {
	return "grpc"
}

// Close close conn
func (c *Client) Close() error {
	return c.c.Close()
}
