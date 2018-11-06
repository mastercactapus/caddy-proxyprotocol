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
		t := 5 * time.Second
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
				if c.Val() == "none" {
					t = 0
					break
				}
				t, err = time.ParseDuration(c.Val())
				if err != nil {
					return nil, err
				}
				if t < 0 {
					return nil, c.ArgErr()
				}

			default:
				return nil, c.ArgErr()
			}
		}

		if len(subnets) == 0 {
			subnets = append(subnets,
				&net.IPNet{Mask: make([]byte, 4), IP: make([]byte, 4)},
				&net.IPNet{Mask: make([]byte, 16), IP: make([]byte, 16)},
			)
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
