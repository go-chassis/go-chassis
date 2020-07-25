package client

import (
	"github.com/go-mesh/openlogging"
	"net/url"
)

func getProtocolMap(eps []string) map[string]string {
	m := make(map[string]string)
	for _, ep := range eps {
		u, err := url.Parse(ep)
		if err != nil {
			openlogging.Error("url err: " + err.Error())
			continue
		}
		m[u.Scheme] = u.Host
	}
	return m
}
