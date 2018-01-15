package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ServiceComb/go-chassis/third_party/forked/valyala/fasthttp"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

func main() {
	flag.Parse()

	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world!\n\n")

	fmt.Fprintf(ctx, "Request method is %q", ctx.Method())
	fmt.Fprintf(ctx, "RequestURI is %q", ctx.RequestURI())
	fmt.Fprintf(ctx, "Requested path is %q", ctx.Path())
	fmt.Fprintf(ctx, "Host is %q", ctx.Host())
	fmt.Fprintf(ctx, "Query string is %q", ctx.QueryArgs())
	fmt.Fprintf(ctx, "User-Agent is %q", ctx.UserAgent())
	fmt.Fprintf(ctx, "Connection has been established at %s", ctx.ConnTime())
	fmt.Fprintf(ctx, "Request has been started at %s", ctx.Time())
	fmt.Fprintf(ctx, "Serial request number for the current connection is %d", ctx.ConnRequestNum())
	fmt.Fprintf(ctx, "Your ip is %q\n", ctx.RemoteIP())

	fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)

	ctx.SetContentType("text/plain; charset=utf8")

	// Set arbitrary headers
	ctx.Response.Header.Set("X-My-Header", "my-header-value")

	// Set cookies
	var c fasthttp.Cookie
	c.SetKey("cookie-name")
	c.SetValue("cookie-value")
	ctx.Response.Header.SetCookie(&c)
}
