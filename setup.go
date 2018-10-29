package proxyprotocol

import (
	"net"
	"time"

	pp "github.com/mastercactapus/proxyprotocol"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type ppRules []pp.Rule

func init() {
	caddy.RegisterPlugin("proxyprotocol", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	rules, err := parseConfig(c)
	if err != nil {
		return err
	}

	if len(rules) > 0 {
		httpserver.GetConfig(c).AddListenerMiddleware(ppRules(rules).NewListener)
	}

	return nil
}

func parseConfig(c *caddy.Controller) (cfgs []pp.Rule, err error) {
	for c.Next() {
		if c.Val() != "proxyprotocol" {
			continue
		}

		var subnets []*net.IPNet
		var t time.Duration
		for _, arg := range c.RemainingArgs() {
			_, n, err := net.ParseCIDR(arg)
			if err != nil {
				return nil, err
			}
			subnets = append(subnets, n)
		}

		if c.NextBlock() {
			switch c.Val() {
			case "timeout":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				t, err = time.ParseDuration(c.Val())
				if err != nil {
					return nil, err
				}
			default:
				return nil, c.ArgErr()
			}
		}
		if len(subnets) == 0 {
			continue
		}
		if c.NextBlock() {
			return nil, c.ArgErr()
		}
		for _, n := range subnets {
			cfgs = append(cfgs, pp.Rule{Subnet: n, Timeout: t})
		}
	}
	return cfgs, nil
}
