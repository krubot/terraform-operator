package opa

import (
	"context"

	"github.com/open-policy-agent/opa/loader"
	"github.com/open-policy-agent/opa/rego"
)

func Validation() bool {
	ctx := context.Background()

	result, err := loader.All([]string{"/opt/."})
	if err != nil {
		panic(err)
	}

	compiler, err := result.Compiler()
	if err != nil {
		panic(err)
	}

	store, err := result.Store()
	if err != nil {
		panic(err)
	}

	rs, err := rego.New(rego.Compiler(compiler), rego.Store(store), rego.Query(`data.terraform.authz`)).Eval(ctx)
	if err != nil {
		panic(err)
	}

	for _, r := range rs {
		if r.Expressions[0].Value == true {
			return true
		}
	}

	return false
}
