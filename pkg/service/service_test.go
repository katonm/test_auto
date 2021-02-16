package service

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
)

type testFix struct {
	value  string
	result string
}

type testValidate struct {
	value  string
	result string
}

func TestFix(t *testing.T) {
	var ctx context.Context
	var testsFix = []testFix{
		{"[()]{}{[()()]()}", "[()]{}{[()()]()}"},
		{"[(])", "[()]"},
		{"[({", "[({})]"},
		{")}]", ""},
		{"[()]{}{[()()]()}}}}}}}}", "[()]{}{[()()]()}"},
	}
	svc := &service{logger: log.NewNopLogger()}

	for _, v := range testsFix {
		a := svc.Fix(ctx, v.value)
		if a != v.result {
			t.Error(
				"For", v.value,
				"expected", v.result,
				"got", v,
			)
		}
	}
}

func TestValidate(t *testing.T) {
	var ctx context.Context
	var testsCases = []testValidate{
		{"[()]{}{[()()]()}", "Balanced"},
		{"[(])", "Not Balanced"},
		{"[({", "Not Balanced"},
		{")}]", "Not Balanced"},
		{"[()]{}{[()()]()}}}}}}}}", "Not Balanced"},
	}
	svc := &service{logger: log.NewNopLogger()}

	for _, v := range testsCases {
		a := svc.Validate(ctx, v.value)
		if a != v.result {
			t.Error(
				"For", v.value,
				"expected", v.result,
				"got", v,
			)
		}
	}
}
