// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"reflect"
)

type tdExpectedType struct {
	Base
	expectedType reflect.Type
	isPtr        bool
}

func (t *tdExpectedType) errorTypeMismatch(ctx Context, gotType rawString) *Error {
	return ctx.CollectError(&Error{
		Message:  "type mismatch",
		Got:      gotType,
		Expected: rawString(t.expectedTypeStr()),
	})
}

func (t *tdExpectedType) checkPtr(ctx Context, pGot *reflect.Value) *Error {
	if t.isPtr {
		got := *pGot
		if got.Kind() != reflect.Ptr {
			if ctx.booleanError {
				return booleanError
			}
			return t.errorTypeMismatch(ctx, rawString(got.Type().String()))
		}
		*pGot = got.Elem()
	}
	return nil
}

func (t *tdExpectedType) checkType(ctx Context, got reflect.Value) *Error {
	if got.Type() != t.expectedType {
		if ctx.booleanError {
			return booleanError
		}
		var gotType rawString
		if t.isPtr {
			gotType = "*"
		}
		gotType += rawString(got.Type().String())
		return t.errorTypeMismatch(ctx, gotType)
	}
	return nil
}

func (t *tdExpectedType) TypeBehind() reflect.Type {
	if t.isPtr {
		return reflect.New(t.expectedType).Type()
	}
	return t.expectedType
}

func (t *tdExpectedType) expectedTypeStr() string {
	if t.isPtr {
		return "*" + t.expectedType.String()
	}
	return t.expectedType.String()
}