package testdeep

import (
	"reflect"
)

type tdAny struct {
	tdList
}

var _ TestDeep = &tdAny{}

func Any(items ...interface{}) TestDeep {
	return &tdAny{
		tdList: newList(items...),
	}
}

func (a *tdAny) Match(ctx Context, got reflect.Value) *Error {
	for _, item := range a.items {
		if deepValueEqualOK(got, item) {
			return nil
		}
	}

	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "comparing with Any",
		Got:      got,
		Expected: a,
		Location: a.GetLocation(),
	}
}