package proxyprotocol

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/mholt/caddy"
	"time"
	"net"
)

func TestProxyprotocolParserOneCDIR(t *testing.T) {
	c := caddy.NewTestController(
		"http",
		`proxyprotocol 127.0.0.1/32`)
	cfg, err := parse(c)
	assert.Nil(t, err, "Should not return any erros, receive %v", err)
	if assert.NotNil(t, cfg, "Config should be set") {
		_, n, _ := net.ParseCIDR("127.0.0.1/32")
		assert.Equal(
			t,
			time.Duration(0),
			cfg.Timeout,
			"Timeout shouldn't be set")
		assert.Equal(
			t,
			[]*net.IPNet{n},
			cfg.Subnets,
			"Should get only one subnet")
	}
}
//
func TestProxyprotocolParserWithTimout(t *testing.T) {
	c := caddy.NewTestController(
		"http",
		`proxyprotocol 0.0.0.0/0 {
			timeout 2s
		}`)
	cfg, err := parse(c)
	assert.Nil(t, err, "Should not return any erros, receive %v", err)
	if assert.NotNil(t, cfg, "Config should be set") {
		_, n, _ := net.ParseCIDR("0.0.0.0/0")
		assert.Equal(
			t,
			time.Duration(2) * time.Second,
			cfg.Timeout,
			"Timeout must be 2s")
		assert.Equal(
			t,
			[]*net.IPNet{n},
			cfg.Subnets,
			"Should get only one subnet")
	}
}

func TestProxyprotocolMultipleCDIR(t *testing.T) {
	c := caddy.NewTestController(
		"http",
		`proxyprotocol 0.0.0.0/0 ::/0 127.0.0.1/32`)
	cfg, err := parse(c)
	assert.Nil(t, err, "Should not return any erros, receive %v", err)
	if assert.NotNil(t, cfg, "Config should be set") {
		_, n1, _ := net.ParseCIDR("0.0.0.0/0")
		_, n2, _ := net.ParseCIDR("::/0")
		_, n3, _ := net.ParseCIDR("127.0.0.1/32")
		assert.Equal(
			t,
			time.Duration(0),
			cfg.Timeout,
			"Timeout shouldn't be set")
		assert.Equal(
			t,
			[]*net.IPNet{n1, n2, n3},
			cfg.Subnets,
			"Should get 3 subnet")
	}
}

func TestProxyprotocolParserMultipleCDIRWithTimout(t *testing.T) {
	c := caddy.NewTestController(
		"http",
		`proxyprotocol 0.0.0.0/0 1234:321::1/24 {
			timeout 25m
		}`)
	cfg, err := parse(c)
	assert.Nil(t, err, "Should not return any erros, receive %v", err)
	if assert.NotNil(t, cfg, "Config should be set") {
		_, n1, _ := net.ParseCIDR("0.0.0.0/0")
		_, n2, _ := net.ParseCIDR("1234:321::1/24")
		assert.Equal(
			t,
			time.Duration(25) * time.Minute,
			cfg.Timeout,
			"Timeout must be 25 minutes")
		assert.Equal(
			t,
			[]*net.IPNet{n1, n2},
			cfg.Subnets,
			"Should get 2 subnet")
	}
}
