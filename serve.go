package serve

import (
	"github.com/Knetic/govaluate"
	"github.com/coredns/coredns/plugin"
)

type Serve struct {
	rules      []*govaluate.EvaluableExpression
	extractors extractorMap
	Next       plugin.Handler
}

