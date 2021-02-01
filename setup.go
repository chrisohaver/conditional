package conditional

import (
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("conditional", setup) }

func setup(c *caddy.Controller) error {
	cond, err := parseOnce(c)
	if err != nil {
		return plugin.Error("forward", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		cond.Next = next
		return cond
	})

	return nil
}

func parseOnce(c *caddy.Controller) (*Conditional, error) {
	var (
		cond *Conditional
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

func parse(c *caddy.Controller) (*Conditional, error) {
	cond := &Conditional{}

	cond.extractors = makeExtractors()
	funcs := makeFunctions()

	opts := c.RemainingArgs()
	if len(opts) != 0 {
		return cond, c.ArgErr()
	}

	groups := make(map[string][]int)

	for c.NextBlock() {
		switch c.Val() {
		case "group":
			args := c.RemainingArgs()
			group := make([]int, len(args[1:]))
			for i, up := range args[1:] {
				u, err := strconv.Atoi(up)
				if err != nil {
					return cond, err
				}
				group[i] = u
			}
			groups[args[0]] = group
		case "use":
			args := c.RemainingArgs()
			if len(args) == 0 {
				return cond, c.ArgErr()
			}
			var r rule
			r.group = args[0]

			if len(args) > 2 {
				if args[1] != "if" {
					return cond, c.Errf("expected 'if' got '%s'", args[1])
				}
				// get expression args[2:]
				expr, err := govaluate.NewEvaluableExpressionWithFunctions(strings.Join(args[2:], " "), funcs)
				if err != nil {
					return cond, err
				}
				r.expr = expr
			}
			cond.rules = append(cond.rules, r)
		default:
			return cond, c.Errf("unknown property '%s'", c.Val())
		}
	}

	for i := range cond.rules {
		if ups, ok := groups[cond.rules[i].group]; ok {
			cond.rules[i].upstreams = ups
			continue
		}
		return cond, c.Errf("unknown group '%s'", cond.rules[i].group)
	}

	return cond, nil
}

func parseBlock(c *caddy.Controller, cond *Conditional) error {

	return nil
}
