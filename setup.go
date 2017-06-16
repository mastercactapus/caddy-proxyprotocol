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
	var cfg Config
	var configs []Config
	var err error

	cfg, err = parse(c)
	if err != nil {
		return err
	}
	configs = append(configs, cfg)

	if configs != nil {
		httpserver.GetConfig(c).AddListenerMiddleware(Configs(configs).NewListener)
	}

	return nil
}

func parse(c *caddy.Controller) (Config, error) {
	var cfg Config
	var err error
	for c.Next() {
		for _, arg := range c.RemainingArgs(){
			_, n, err := net.ParseCIDR(arg)
			if err != nil {
				return cfg, err
			}
			cfg.Subnets = append(cfg.Subnets, n)
		}

		if c.NextBlock() {
			switch c.Val() {
			case "timeout":
				if !c.NextArg() {
					return cfg, c.ArgErr()
				}
				cfg.Timeout, err = time.ParseDuration(c.Val())
				if err != nil {
					return cfg, err
				}
			default:
				return cfg, c.ArgErr()
			}
		}
	}
	return cfg, nil
}
