package bench

import (
	"context"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/core/client"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	"time"
)

func callHTTP(ctx context.Context, restInvoker *core.RestInvoker, method string, url string, body []byte) {
	for {
		start := time.Now()
		req, err := rest.NewRequest(method, url, body)
		if err != nil {
			panic(err)
		}
		resp, err := restInvoker.ContextDo(ctx, req)
		if err != nil {
			if err != client.ErrCanceled {
				errRequest++
			}
			continue
		}
		_ = httputil.ReadBody(resp)
		resp.Body.Close()
		totalRequest++
		latency.UpdateSince(start)

		select {
		case <-ctx.Done():
			break
		default:
			continue
		}
	}
}
