package conditional

import (
	"github.com/Knetic/govaluate"
	"github.com/coredns/coredns/plugin"
)

type Conditional struct {
	fwdRules   []fwdRule
	viewRules  []*govaluate.EvaluableExpression
	extractors extractorMap
	Next       plugin.Handler
}

type fwdRule struct {
	expr      *govaluate.EvaluableExpression
	upstreams []int
	group     string
}
