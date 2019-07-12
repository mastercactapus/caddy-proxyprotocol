package proxyprotocol

import (
	"os"

	"github.com/caddyserver/caddy"
	pp "github.com/mastercactapus/proxyprotocol"
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
		f := ppL.Filter()
		f = append(f, r...)
		ppL.SetFilter(f)
		return l
	}

	ppL := pp.NewListener(l, 0)
	ppL.SetFilter(r)

	return &Listener{
		Listener: ppL,
		cl:       l,
	}
}
