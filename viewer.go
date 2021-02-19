package serve

import (
	"github.com/coredns/coredns/request"
)

func (c *Serve) Filter(state request.Request) bool {
	params := Parameters{state: &state, extractors: c.extractors}
	// return true if all expressions evaluate to true
	for _, expr := range c.rules {
		result, err := expr.Eval(params)
		if err != nil {
			return false
		}
		if b, ok := result.(bool); ok && b {
			continue
		}
		return false
	}
	return true
}
