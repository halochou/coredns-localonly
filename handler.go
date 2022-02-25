package localonly

import (
	"context"
	"regexp"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

type LocalOnly struct {
	NetMask int
	Zones   []*regexp.Regexp
	Next    plugin.Handler
}

// ServeDNS implements the plugin.Handler interface.
func (rr LocalOnly) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	match := false
	for _, zone := range rr.Zones {
		if zone.MatchString(state.Name()) {
			match = true
			break
		}
	}
	if !match {
		return plugin.NextOrFailure(rr.Name(), rr.Next, ctx, w, r)
	}
	wrr := &LocalOnlyResponseWriter{w, rr.NetMask}
	return plugin.NextOrFailure(rr.Name(), rr.Next, ctx, wrr, r)
}

// Name implements the Handler interface.
func (rr LocalOnly) Name() string { return "localonly" }
