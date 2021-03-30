package rule

// based on an example here: https://github.com/mgechev/revive/blob/master/rule/argument-limit.go

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// ByteArrayLimitRule lints given else constructs.
type ByteArrayLimitRule struct{}

// Apply applies the rule to given file.
func (*ByteArrayLimitRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	if len(arguments) != 1 {
		panic(`invalid configuration for "byte-array-limit"`)
	}

	total, ok := arguments[0].(int64) // Alt. non panicking version
	if !ok {
		panic(`invalid value passed as argument number to the "byte-array-limit" rule`)
	}

	var failures []lint.Failure

	walker := lintByteArrayLimit{
		total: int(total),
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, file.AST)

	return failures
}

// Name returns the rule name.
func (*ByteArrayLimitRule) Name() string {
	return "byte-array-limit"
}

type lintByteArrayLimit struct {
	total     int
	onFailure func(lint.Failure)
}

func (w lintByteArrayLimit) Visit(n ast.Node) ast.Visitor {
	node, ok := n.(*ast.CompositeLit)
	if ok {
		num := len(node.Elts)
		if at, ok := node.Type.(*ast.ArrayType); ok {
			if id, ok := at.Elt.(*ast.Ident); ok {
				if id.String() == "byte" && num > w.total {
					w.onFailure(lint.Failure{
						Confidence: 1,
						Failure:    fmt.Sprintf("for byte arrays longer than %d please use byte(\"abc\\x00\") format", w.total),
						Node:       node.Type,
					})
					return w
				}
			}
		}
	}
	return w
}
