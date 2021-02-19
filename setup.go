package serve

import (
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("serve", setup) }

func setup(c *caddy.Controller) error {
	cond, err := parseOnce(c)
	if err != nil {
		return plugin.Error("serve", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		cond.Next = next
		return cond
	})

	return nil
}

func parseOnce(c *caddy.Controller) (*Serve, error) {
	var (
		cond *Serve
		err  error
		i    int
	)
	for c.Next() {
		if i > 0 {
			return nil, plugin.ErrOnce
		}
		i++
		cond, err = parse(c)
		if err != nil {
			return nil, err
		}
	}
	return cond, nil
}

func parse(c *caddy.Controller) (*Serve, error) {
	cond := &Serve{}

	cond.extractors = makeExtractors()
	funcs := makeFunctions()

	opts := c.RemainingArgs()
	if len(opts) != 0 {
		return cond, c.ArgErr()
	}

	for c.NextBlock() {
			args := c.RemainingArgs()
			expr, err := govaluate.NewEvaluableExpressionWithFunctions(strings.Join(args, " "), funcs)
			if err != nil {
				return cond, err
			}
			cond.rules = append(cond.rules, expr)
	}

	return cond, nil
}

func parseBlock(c *caddy.Controller, cond *Serve) error {

	return nil
}
