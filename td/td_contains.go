// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package td

import (
	"reflect"
	"strings"

	"github.com/maxatome/go-testdeep/helpers/tdutil"
	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/types"
	"github.com/maxatome/go-testdeep/internal/util"
)

type tdContains struct {
	tdSmugglerBase
}

var _ TestDeep = &tdContains{}

// summary(Contains): checks that a string, error or fmt.Stringer
// interfaces contain a sub-string; or an array, slice or map contain
// a value
// input(Contains): str,array,slice,map

// Contains is a smuggler operator with a little convenient exception
// for strings. Contains has to be applied on arrays, slices, maps or
// strings. It compares each item of data array/slice/map/string (rune
// for strings) against "expectedValue".
//
//   list := []int{12, 34, 28}
//   td.Cmp(t, list, td.Contains(34))                 // succeeds
//   td.Cmp(t, list, td.Contains(td.Between(30, 35))) // succeeds too
//   td.Cmp(t, list, td.Contains(35))                 // fails
//
//   hash := map[string]int{"foo": 12, "bar": 34, "zip": 28}
//   td.Cmp(t, hash, td.Contains(34))                 // succeeds
//   td.Cmp(t, hash, td.Contains(td.Between(30, 35))) // succeeds too
//   td.Cmp(t, hash, td.Contains(35))                 // fails
//
//   got := "foo bar"
//   td.Cmp(t, got, td.Contains('o'))                  // succeeds
//   td.Cmp(t, got, td.Contains(rune('o')))            // succeeds
//   td.Cmp(t, got, td.Contains(td.Between('n', 'p'))) // succeeds
//
// When Contains(nil) is used, nil is automatically converted to a
// typed nil on the fly to avoid confusion (if the array/slice/map
// item type allows it of course.) So all following Cmp calls
// are equivalent (except the (*byte)(nil) one):
//
//   num := 123
//   list := []*int{&num, nil}
//   td.Cmp(t, list, td.Contains(nil))         // succeeds → (*int)(nil)
//   td.Cmp(t, list, td.Contains((*int)(nil))) // succeeds
//   td.Cmp(t, list, td.Contains(td.Nil()))    // succeeds
//   // But...
//   td.Cmp(t, list, td.Contains((*byte)(nil))) // fails: (*byte)(nil) ≠ (*int)(nil)
//
// As well as these ones:
//
//   hash := map[string]*int{"foo": nil, "bar": &num}
//   td.Cmp(t, hash, td.Contains(nil))         // succeeds → (*int)(nil)
//   td.Cmp(t, hash, td.Contains((*int)(nil))) // succeeds
//   td.Cmp(t, hash, td.Contains(td.Nil()))    // succeeds
//
// As a special case for string (or convertible), error or
// fmt.Stringer interface (error interface is tested before
// fmt.Stringer), "expectedValue" can be a string, a rune or a
// byte. In this case, it tests if the got string contains this
// expected string, rune or byte.
//
//   type Foobar string
//   td.Cmp(t, Foobar("foobar"), td.Contains("ooba")) // succeeds
//
//   err := errors.New("error!")
//   td.Cmp(t, err, td.Contains("ror")) // succeeds
//
//   bstr := bytes.NewBufferString("fmt.Stringer!")
//   td.Cmp(t, bstr, td.Contains("String")) // succeeds
func Contains(expectedValue interface{}) TestDeep {
	c := tdContains{
		tdSmugglerBase: newSmugglerBase(expectedValue),
	}

	if !c.isTestDeeper {
		c.expectedValue = reflect.ValueOf(expectedValue)
	}
	return &c
}

func (c *tdContains) doesNotContain(ctx ctxerr.Context, got interface{}) *ctxerr.Error {
	if ctx.BooleanError {
		return ctxerr.BooleanError
	}
	return ctx.CollectError(&ctxerr.Error{
		Message:  "does not contain",
		Got:      got,
		Expected: c,
	})
}

// getExpectedValue returns the expected value handling the
// Contains(nil) case: in this case it returns a typed nil (same type
// as the items of got).
// got is an array, a slice or a map (it's the caller responsibility to check)
func (c *tdContains) getExpectedValue(got reflect.Value) reflect.Value {
	// If the expectValue is non-typed nil
	if !c.expectedValue.IsValid() {
		// AND the kind of items in got is...
		switch got.Type().Elem().Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface,
			reflect.Map, reflect.Ptr, reflect.Slice:
			// returns a typed nil
			return reflect.Zero(got.Type().Elem())
		}
	}
	return c.expectedValue
}

func (c *tdContains) Match(ctx ctxerr.Context, got reflect.Value) *ctxerr.Error {
	switch got.Kind() {
	case reflect.Array, reflect.Slice:
		expectedValue := c.getExpectedValue(got)
		for index := got.Len() - 1; index >= 0; index-- {
			if deepValueEqualOK(got.Index(index), expectedValue) {
				return nil
			}
		}
		return c.doesNotContain(ctx, got)

	case reflect.Map:
		expectedValue := c.getExpectedValue(got)
		if !tdutil.MapEachValue(got, func(v reflect.Value) bool {
			return !deepValueEqualOK(v, expectedValue)
		}) {
			return nil
		}
		return c.doesNotContain(ctx, got)

		// For String kind *AND* TestDeep operator, applies this operator on
		// each character of the string
	case reflect.String:
		if c.isTestDeeper {
			for _, chr := range got.String() {
				if deepValueEqualOK(reflect.ValueOf(chr), c.expectedValue) {
					return nil
				}
			}
			return c.doesNotContain(ctx, got)
		}
	}

	// A TestDeep operator can only be applied on arrays, slices, map
	// and as a special feature on strings (all handled in switch
	// above). For all other cases, it is an error.
	if !c.isTestDeeper {
		// If expectedValue is a string, a rune or a byte, we try to get a
		// string from got and check whether expectedValue is contained in
		// this string or not
		switch expectedKind := c.expectedValue.Kind(); expectedKind {
		case reflect.String, reflect.Int32, reflect.Uint8: // string, rune & byte
			str, err := getString(ctx, got)
			if err != nil {
				return err
			}

			switch expectedKind {
			case reflect.String:
				if strings.Contains(str, c.expectedValue.String()) {
					return nil
				}
			case reflect.Int32:
				if strings.ContainsRune(str, rune(c.expectedValue.Int())) {
					return nil
				}
			default: // = case reflect.Uint8:
				if strings.IndexByte(str, byte(c.expectedValue.Uint())) >= 0 {
					return nil
				}
			}
			return c.doesNotContain(ctx, str)
		}
	}

	if ctx.BooleanError {
		return ctxerr.BooleanError
	}
	var expectedType interface{}
	if c.expectedValue.IsValid() {
		expectedType = types.RawString(c.expectedValue.Type().String())
	} else {
		expectedType = c
	}

	return ctx.CollectError(&ctxerr.Error{
		Message:  "cannot check contains",
		Got:      types.RawString(got.Type().String()),
		Expected: expectedType,
	})
}

func (c *tdContains) String() string {
	return "Contains(" + util.ToString(c.expectedValue) + ")"
}