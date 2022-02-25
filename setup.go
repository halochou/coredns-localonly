package localonly

import (
	"regexp"
	"strconv"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

var log = clog.NewWithPlugin("localonly")

func init() { plugin.Register("localonly", setup) }

func setup(c *caddy.Controller) error {
	mask, zones, err := parse(c)
	if err != nil {
		return plugin.Error("localonly", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return LocalOnly{
			NetMask: mask,
			Zones:   zones,
			Next:    next,
		}
	})

	return nil
}

func parse(c *caddy.Controller) (int, []*regexp.Regexp, error) {
	mask := 24
	var zones []*regexp.Regexp
	for c.Next() {
		args := c.RemainingArgs()
		for i, arg := range args {
			if i == 0 {
				v, err := strconv.Atoi(arg)
				if err != nil {
					return 0, nil, c.ArgErr()
				}
				mask = v
			} else {
				pat, err := regexp.Compile(arg)
				if err != nil {
					return 0, nil, c.ArgErr()
				}
				zones = append(zones, pat)
			}
		}
	}
	return mask, zones, nil
}
