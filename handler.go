package serve

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

// Name implements the Handler interface
func (c *Serve) Name() string { return "serve" }

// ServeDNS implements the Handler interface.
func (c *Serve) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return plugin.NextOrFailure(c.Name(), c.Next, ctx, w, r)
}
