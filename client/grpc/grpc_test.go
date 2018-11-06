package grpc_test

import (
	"github.com/go-chassis/go-chassis/client/grpc"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestTransformContext(t *testing.T) {
	ctx := common.NewContext(map[string]string{
		"1": "2",
		"3": "4",
	})
	ctx = grpc.TransformContext(ctx)
	md, ok := metadata.FromOutgoingContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, "2", md["1"][0])
	assert.Equal(t, "4", md["3"][0])
}

func TestNew(t *testing.T) {

}
