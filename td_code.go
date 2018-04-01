package testdeep

import (
	"fmt"
	"reflect"
)

type tdCode struct {
	TestDeepBase
	function reflect.Value
	argType  reflect.Type
}

var _ TestDeep = &tdCode{}

func Code(fn interface{}) TestDeep {
	vfn := reflect.ValueOf(fn)

	if vfn.Kind() != reflect.Func {
		panic("usage: Code(FUNC)")
	}

	fnType := vfn.Type()
	if fnType.NumIn() != 1 {
		panic("Code(FUNC): FUNC must take only one argument")
	}

	switch fnType.NumOut() {
	case 2:
		if fnType.Out(1).Kind() != reflect.String {
			break
		}
		fallthrough

	case 1:
		if fnType.Out(0).Kind() == reflect.Bool {
			return &tdCode{
				TestDeepBase: NewTestDeepBase(3),
				function:     vfn,
				argType:      fnType.In(0),
			}
		}
	}

	panic("Code(FUNC): FUNC must return bool or (bool, string)")
}

func (c *tdCode) Match(ctx Context, got reflect.Value) *Error {
	if !got.Type().AssignableTo(c.argType) {
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "incompatible parameter type",
			Got:      rawString(got.Type().String()),
			Expected: rawString(c.argType.String()),
			Location: c.GetLocation(),
		}
	}

	ret := c.function.Call([]reflect.Value{got})
	if ret[0].Bool() {
		return nil
	}

	if ctx.booleanError {
		return booleanError
	}

	err := Error{
		Context:  ctx,
		Message:  "ran code with %% as argument",
		Location: c.GetLocation(),
	}

	if len(ret) > 1 {
		err.Summary = tdCodeResult{
			Value:  got,
			Reason: ret[1].String(),
		}
	} else {
		err.Summary = tdCodeResult{
			Value: got,
		}
	}

	return &err
}

func (c *tdCode) String() string {
	return "Code(" + c.function.Type().String() + ")"
}

type tdCodeResult struct {
	Value  reflect.Value
	Reason string
}

var _ testDeepStringer = tdCodeResult{}

func (r tdCodeResult) ___testDeep___() {}

func (r tdCodeResult) String() string {
	if r.Reason == "" {
		return fmt.Sprintf("  value: %s\nit failed but didn't say why",
			toString(r.Value))
	}
	return fmt.Sprintf("        value: %s\nit failed coz: %s",
		toString(r.Value), r.Reason)
}