package proxyprotocol

import (
	"testing"
	"time"

	"github.com/mholt/caddy"
	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {

	type expectedCfg struct {
		subnets []string
		timeout time.Duration
	}

	check := func(name, cfg string, expected ...expectedCfg) {
		t.Run(name, func(t *testing.T) {
			cfgs, err := parseConfig(caddy.NewTestController("http", cfg))
			assert.NoError(t, err)
			assert.Len(t, cfgs, len(expected))
			for i, exp := range expected {
				if len(cfgs) <= i {
					break
				}
				cfg := cfgs[i]
				assert.Len(t, cfg.Subnets, len(exp.subnets))
				for i, sub := range exp.subnets {
					assert.Equal(t, sub, cfg.Subnets[i].String(), "Subnet[%d]", i)
				}
				assert.Equal(t, exp.timeout.String(), cfg.Timeout.String(), "Timeout")
			}
		})
	}

	check(
		"empty",
		"",
	)
	check(
		"single-subnet",
		`proxyprotocol 127.0.0.1/32`,
		expectedCfg{subnets: []string{"127.0.0.1/32"}},
	)
	check(
		"multi-subnet",
		`proxyprotocol 0.0.0.0/0 ::/0`,
		expectedCfg{subnets: []string{"0.0.0.0/0", "::/0"}},
	)
	check(
		"duplicate",
		`proxyprotocol 0.0.0.0/0
		proxyprotocol ::/0`,
		expectedCfg{subnets: []string{"0.0.0.0/0"}},
		expectedCfg{subnets: []string{"::/0"}},
	)
	check(
		"block-single",
		`proxyprotocol 0.0.0.0/0 {
			timeout 2s
		}`,
		expectedCfg{subnets: []string{"0.0.0.0/0"}, timeout: 2 * time.Second},
	)
	check(
		"block-multi",
		`proxyprotocol 0.0.0.0/0 1234:321::1/24 {
			timeout 25m
		}`,
		expectedCfg{subnets: []string{"0.0.0.0/0", "1234:300::/24"}, timeout: 25 * time.Minute},
	)
	check(
		"block-duplicate",
		`proxyprotocol 0.0.0.0/0 {
			timeout 25m
		}
		proxyprotocol 1234:321::1/24 {
			timeout 30m
		}`,
		expectedCfg{subnets: []string{"0.0.0.0/0"}, timeout: 25 * time.Minute},
		expectedCfg{subnets: []string{"1234:300::/24"}, timeout: 30 * time.Minute},
	)

}
