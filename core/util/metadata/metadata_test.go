package metadata_test

import (
	context17 "context"
	"github.com/ServiceComb/go-chassis/core/util/metadata"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	_ "golang.org/x/net/context"
	"testing"
)

func TestMetadata(t *testing.T) {
	var ctx context.Context
	//ctx:=context.TODO()
	m, b := metadata.FromContext(ctx)
	assert.Nil(t, m)
	assert.Equal(t, false, b)

	var ctx17 context17.Context
	m1, b1 := metadata.FromContext17(ctx17)
	assert.Nil(t, m1)
	assert.Equal(t, false, b1)

	//var ctx context.Context
	ctx1 := context.TODO()
	m, b = metadata.FromContext(ctx1)
	assert.Nil(t, m)
	assert.Equal(t, false, b)

	ctx171 := context.TODO()
	m1, b1 = metadata.FromContext17(ctx171)
	assert.Nil(t, m1)
	assert.Equal(t, false, b1)

	var mt metadata.Metadata = make(map[string]string)
	mt["abc"] = "abc"

	ct := metadata.NewContext(ctx1, mt)
	assert.NotNil(t, ct)

	ct17 := metadata.NewContext17(ctx171, mt)
	assert.NotNil(t, ct17)

}
