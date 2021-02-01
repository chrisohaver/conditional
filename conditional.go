package conditional

import (
	"github.com/Knetic/govaluate"
	"github.com/coredns/coredns/plugin"
)

type Conditional struct {
	rules      []rule
	extractors extractorMap
	Next       plugin.Handler
}

type rule struct {
	expr      *govaluate.EvaluableExpression
	upstreams []int
	group     string
}
