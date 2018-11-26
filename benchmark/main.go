package main

import (
	"github.com/go-chassis/go-chassis/benchmark/util"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-chassis benchmark tool"
	app.Description = "example: ./benchmark -c 10 -d 30s -u http://service/hello"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "c",
			Usage:       "concurrency number",
			Destination: &bench.Configs.Concurrency,
		},
		cli.StringFlag{
			Name:        "d",
			Usage:       "how long request should send, example -d 30s",
			Destination: &bench.Configs.Duration,
		},
		cli.StringFlag{
			Name:        "u",
			Usage:       "URI, for example: -u http://SomeHTTPService/rest/api, -u grpc://SomeGRPCService",
			Destination: &bench.Configs.Target,
		},
		cli.StringFlag{
			Name:        "-m",
			Usage:       "method, only http support this option, for example: -m POST",
			Destination: &bench.Configs.Method,
		},
		cli.StringFlag{
			Name:        "D",
			Usage:       "only http support this option, it is the path of body, for example: -D /path/to/body.json",
			Destination: &bench.Configs.BodyPath,
		},
	}
	app.Action = func(c *cli.Context) error {
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
	bench.Benchmark()
}
