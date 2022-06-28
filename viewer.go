package conditional

import (
	"github.com/coredns/coredns/request"
)

func (c *Conditional) Filter(state request.Request) bool {
	// TODO: Context need to be passed to Filter, otherwise metadata eval cannot work.  Requires interface change.
	params := Parameters{state: &state, extractors: c.extractors}
	// return true if all expressions evaluate to true
	for _, expr := range c.viewRules {
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
