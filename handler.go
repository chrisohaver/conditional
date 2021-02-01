package conditional

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

// Name implements the Handler interface
func (c *Conditional) Name() string { return "conditional" }

// ServeDNS implements the Handler interface.
func (c *Conditional) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return plugin.NextOrFailure(c.Name(), c.Next, ctx, w, r)
}
