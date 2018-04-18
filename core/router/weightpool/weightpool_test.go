package weightpool_test

import (
	"testing"

	"github.com/ServiceComb/go-chassis/core/config/model"
	wp "github.com/ServiceComb/go-chassis/core/router/weightpool"
	"github.com/stretchr/testify/assert"
)

var (
	tagsOff100 = []*model.RouteTag{
		{Weight: 25, Tags: map[string]string{"version": "A"}},
		{Weight: 30, Tags: map[string]string{"version": "B"}},
		{Weight: 40, Tags: map[string]string{"version": "C"}},
	}
	tags50 = []*model.RouteTag{
		{Weight: 50, Tags: map[string]string{"version": "A"}},
		{Weight: 50, Tags: map[string]string{"version": "B"}},
	}
)

func TestPoolPickOne(t *testing.T) {
	p, ok := wp.GetPool().Get("test")
	if !ok {
		p = wp.NewPool(tagsOff100...)
		wp.GetPool().Set("test", p)
	}

	var a, b, c, d int
	for i := 0; i < 100; i++ {
		t := p.PickOne()
		switch t.Tags["version"] {
		case "A":
			a++
		case "B":
			b++
		case "C":
			c++
		case "latest":
			d++
		}
	}
	assert.Equal(t, a, 25)
	assert.Equal(t, b, 30)
	assert.Equal(t, c, 40)
	assert.Equal(t, d, 5)

	wp.GetPool().Reset("test")
}

func BenchmarkPickOne(b *testing.B) {
	p := wp.NewPool(tagsOff100...)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.PickOne()
	}
	b.ReportAllocs()
}

func BenchmarkPickOneParallel(b *testing.B) {
	p := wp.NewPool(tags50...)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.PickOne()
		}
	})

	b.ReportAllocs()
}
