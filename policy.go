package conditional

import (
	"context"

	"github.com/coredns/coredns/plugin/forward"
	"github.com/coredns/coredns/plugin/metadata"
	"github.com/coredns/coredns/request"
)

func (c *Conditional) String() string { return "conditional" }

func (c *Conditional) List(ctx context.Context, p []*forward.Proxy, state *request.Request) []*forward.Proxy {
	params := Parameters{state: state, extractors: c.extractors}
	for _, r := range c.fwdRules {
		result, err := r.expr.Eval(params)
		if err != nil {
			return nil
		}
		if b, ok := result.(bool); ok && b {
			ups := make([]*forward.Proxy, len(r.upstreams))
			for i, n := range r.upstreams {
				ups[i] = p[n]
			}
			metadata.SetValueFunc(ctx, "forward/group", func() string {
				return r.group
			})
			return ups
		}
	}
	return nil
}
