package proxyprotocol

import (
	"os"

	pp "github.com/mastercactapus/proxyprotocol"
	"github.com/mholt/caddy"
)

// Listener adds PROXY protocol support to a caddy.Listener.
type Listener struct {
	*pp.Listener
	cl caddy.Listener
}

// File implements the caddy.Listener interface.
func (l *Listener) File() (*os.File, error) { return l.cl.File() }

func (r ppRules) NewListener(l caddy.Listener) caddy.Listener {
	if ppL, ok := l.(*Listener); ok {
		// merge existing
		ppL.AddRules(r)
		return l
	}
	return &Listener{
		Listener: pp.NewListener(l, r),
		cl:       l,
	}
}
