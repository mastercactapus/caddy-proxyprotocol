package proxyprotocol

import (
	"net"
	"time"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type Config struct {
	Timeout time.Duration
	Subnets []*net.IPNet
}

func init() {
	caddy.RegisterPlugin("proxyprotocol", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	configs, err := parseConfig(c)
	if err != nil {
		return err
	}

	if len(configs) > 0 {
		httpserver.GetConfig(c).AddListenerMiddleware(Configs(configs).NewListener)
	}

	return nil
}

func parseConfig(c *caddy.Controller) (cfgs []Config, err error) {
	for c.Next() {
		if c.Val() != "proxyprotocol" {
			continue
		}

		var cfg Config
		for _, arg := range c.RemainingArgs() {
			_, n, err := net.ParseCIDR(arg)
			if err != nil {
				return nil, err
			}
			cfg.Subnets = append(cfg.Subnets, n)
		}

		if c.NextBlock() {
			switch c.Val() {
			case "timeout":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				cfg.Timeout, err = time.ParseDuration(c.Val())
				if err != nil {
					return nil, err
				}
			default:
				return nil, c.ArgErr()
			}
		}
		if len(cfg.Subnets) == 0 {
			continue
		}
		if c.NextBlock() {
			return nil, c.ArgErr()
		}
		cfgs = append(cfgs, cfg)
	}
	return cfgs, nil
}
