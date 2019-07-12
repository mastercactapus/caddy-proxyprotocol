package proxyprotocol

import (
	"testing"
	"time"

	"github.com/caddyserver/caddy"
	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {

	type exp struct {
		subnet  string
		timeout time.Duration
	}

	check := func(name, cfg string, expected ...exp) {
		t.Run(name, func(t *testing.T) {
			cfgs, err := parseConfig(caddy.NewTestController("http", cfg))
			assert.NoError(t, err)
			assert.Len(t, cfgs, len(expected))
			for i, exp := range expected {
				if len(cfgs) <= i {
					break
				}
				cfg := cfgs[i]
				assert.Equal(t, exp.subnet, cfg.Subnet.String(), "Subnet")
				assert.Equal(t, exp.timeout.String(), cfg.Timeout.String(), "Timeout")
			}
		})
	}

	check(
		"empty",
		``,
	)
	check(
		"default",
		`proxyprotocol`,
		exp{subnet: "0.0.0.0/0", timeout: 5 * time.Second},
		exp{subnet: "::/0", timeout: 5 * time.Second},
	)
	check(
		"default-options",
		`proxyprotocol {
			timeout 1s
		}`,
		exp{subnet: "0.0.0.0/0", timeout: time.Second},
		exp{subnet: "::/0", timeout: time.Second},
	)
	check(
		"single-subnet",
		`proxyprotocol 127.0.0.1/32`,
		exp{subnet: "127.0.0.1/32", timeout: 5 * time.Second},
	)
	check(
		"multi-subnet",
		`proxyprotocol 0.0.0.0/0 ::/0`,
		exp{subnet: "0.0.0.0/0", timeout: 5 * time.Second},
		exp{subnet: "::/0", timeout: 5 * time.Second},
	)
	check(
		"duplicate",
		`proxyprotocol 0.0.0.0/0
		proxyprotocol ::/0`,
		exp{subnet: "0.0.0.0/0", timeout: 5 * time.Second},
		exp{subnet: "::/0", timeout: 5 * time.Second},
	)
	check(
		"no-timeout",
		`proxyprotocol 0.0.0.0/0 {
			timeout 0
		}`,
		exp{subnet: "0.0.0.0/0"},
	)
	check(
		"no-timeout",
		`proxyprotocol 0.0.0.0/0 {
			timeout none
		}`,
		exp{subnet: "0.0.0.0/0"},
	)

	check(
		"block-single",
		`proxyprotocol 0.0.0.0/0 {
			timeout 2s
		}`,
		exp{subnet: "0.0.0.0/0", timeout: 2 * time.Second},
	)
	check(
		"block-multi",
		`proxyprotocol 0.0.0.0/0 1234:321::1/24 {
			timeout 25m
		}`,
		exp{subnet: "0.0.0.0/0", timeout: 25 * time.Minute},
		// normalized subnet str
		exp{subnet: "1234:300::/24", timeout: 25 * time.Minute},
	)
	check(
		"block-duplicate",
		`proxyprotocol 0.0.0.0/0 {
			timeout 25m
		}
		proxyprotocol 1234:321::1/24 {
			timeout 30m
		}`,
		exp{subnet: "0.0.0.0/0", timeout: 25 * time.Minute},
		exp{subnet: "1234:300::/24", timeout: 30 * time.Minute},
	)
	check(
		"multi-site",
		`example.com {
			proxyprotocol 0.0.0.0/0
		}
		foo.com {
			proxyprotocol 1234:321::1/24 {
				timeout 30m
			}
		}`,
		exp{subnet: "0.0.0.0/0", timeout: 5 * time.Second},
		exp{subnet: "1234:300::/24", timeout: 30 * time.Minute},
	)
}
