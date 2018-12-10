// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/dark"
	"github.com/maxatome/go-testdeep/internal/test"
)

type MyStructBase struct {
	ValBool bool
}

type MyStructMid struct {
	MyStructBase
	ValStr string
}

type MyStruct struct {
	MyStructMid
	ValInt int
	Ptr    *int
}

func (s *MyStruct) MyString() string {
	return "!"
}

type MyInterface interface {
	MyString() string
}

type MyStringer struct{}

func (s MyStringer) String() string { return "pipo bingo" }

type expectedError struct {
	Path     expectedErrorMatch
	Message  expectedErrorMatch
	Got      expectedErrorMatch
	Expected expectedErrorMatch
	Summary  expectedErrorMatch
	Located  bool
	Origin   *expectedError
}

type expectedErrorMatch struct {
	Exact   string
	Match   *regexp.Regexp
	Contain string
}

func mustBe(str string) expectedErrorMatch {
	return expectedErrorMatch{Exact: str}
}

func mustMatch(str string) expectedErrorMatch {
	return expectedErrorMatch{Match: regexp.MustCompile(str)}
}

func mustContain(str string) expectedErrorMatch {
	return expectedErrorMatch{Contain: str}
}

func indent(str string, numSpc int) string {
	return strings.Replace(str, "\n", "\n\t"+strings.Repeat(" ", numSpc), -1)
}

func cmpErrorStr(t *testing.T, err *ctxerr.Error,
	got string, expected expectedErrorMatch, fieldName string,
	args ...interface{}) bool {
	t.Helper()

	if expected.Exact != "" && got != expected.Exact {
		t.Errorf(`%sError.%s mismatch
	     got: %s
	expected: %s
	Full error:
	> %s`,
			test.BuildTestName(args),
			fieldName, indent(got, 10), indent(expected.Exact, 10),
			strings.Replace(err.Error(), "\n\t", "\n\t> ", -1))
		return false
	}

	if expected.Contain != "" && !strings.Contains(got, expected.Contain) {
		t.Errorf(`%sError.%s mismatch
	           got: %s
	should contain: %s
	Full error:
	> %s`,
			test.BuildTestName(args),
			fieldName,
			indent(got, 16), indent(expected.Contain, 16),
			strings.Replace(err.Error(), "\n\t", "\n\t> ", -1))
		return false
	}

	if expected.Match != nil && !expected.Match.MatchString(got) {
		t.Errorf(`%sError.%s mismatch
	         got: %s
	should match: %s
	Full error:
	> %s`,
			test.BuildTestName(args),
			fieldName,
			indent(got, 14), indent(expected.Match.String(), 14),
			strings.Replace(err.Error(), "\n\t", "\n\t> ", -1))
		return false
	}

	return true
}

func matchError(t *testing.T, err *ctxerr.Error, expectedError expectedError,
	expectedIsTestDeep bool, args ...interface{}) bool {
	t.Helper()

	if !cmpErrorStr(t, err, err.Message, expectedError.Message,
		"Message", args...) {
		return false
	}

	if !cmpErrorStr(t, err, err.Context.Path, expectedError.Path,
		"Context.Path", args...) {
		return false
	}

	if !cmpErrorStr(t, err, err.GotString(), expectedError.Got, "Got", args...) {
		return false
	}

	if !cmpErrorStr(t, err,
		err.ExpectedString(), expectedError.Expected, "Expected", args...) {
		return false
	}

	if !cmpErrorStr(t, err,
		err.SummaryString(), expectedError.Summary, "Summary", args...) {
		return false
	}

	// If expected is a TestDeep, the Location should be set
	if expectedIsTestDeep {
		expectedError.Located = true
	}
	if expectedError.Located != err.Location.IsInitialized() {
		t.Errorf(`%sLocation of the origin of the error
	     got: %v
	expected: %v`,
			test.BuildTestName(args), err.Location.IsInitialized(), expectedError.Located)
		return false
	}

	if expectedError.Located &&
		!strings.HasSuffix(err.Location.File, "_test.go") {
		t.Errorf(`%sFile of the origin of the error
	     got: line %d of %s
	expected: *_test.go`,
			test.BuildTestName(args), err.Location.Line, err.Location.File)
		return false
	}

	if expectedError.Origin != nil {
		if err.Origin == nil {
			t.Errorf(`%sError should originate from another Error`,
				test.BuildTestName(args))
			return false
		}

		return matchError(t, err.Origin, *expectedError.Origin,
			expectedIsTestDeep, args...)
	}
	if err.Origin != nil {
		t.Errorf(`%sError should NOT originate from another Error`,
			test.BuildTestName(args))
		return false
	}

	return true
}

func _checkError(t *testing.T, got, expected interface{},
	expectedError expectedError, args ...interface{}) bool {
	t.Helper()

	err := testdeep.EqDeeplyError(got, expected)
	if err == nil {
		t.Errorf("%sAn Error should have occurred", test.BuildTestName(args))
		return false
	}

	_, expectedIsTestDeep := expected.(testdeep.TestDeep)
	if !matchError(t, err.(*ctxerr.Error), expectedError, expectedIsTestDeep, args...) {
		return false
	}

	if testdeep.EqDeeply(got, expected) {
		t.Errorf(`%sBoolean context failed
	     got: true
	expected: false`, test.BuildTestName(args))
		return false
	}

	return true
}

func ifaceExpectedError(t *testing.T, expectedError expectedError) expectedError {
	t.Helper()

	if !strings.Contains(expectedError.Path.Exact, "DATA") {
		return expectedError
	}

	newExpectedError := expectedError
	newExpectedError.Path.Exact = strings.Replace(expectedError.Path.Exact,
		"DATA", "DATA.Iface", 1)

	if newExpectedError.Origin != nil {
		newOrigin := ifaceExpectedError(t, *newExpectedError.Origin)
		newExpectedError.Origin = &newOrigin
	}

	return newExpectedError
}

// checkError calls _checkError twice. The first time with the same
// parameters, the second time in an interface{} context.
func checkError(t *testing.T, got, expected interface{},
	expectedError expectedError, args ...interface{}) bool {
	t.Helper()

	if ok := _checkError(t, got, expected, expectedError, args...); !ok {
		return false
	}

	type tmpStruct struct {
		Iface interface{}
	}

	return _checkError(t, tmpStruct{Iface: got},
		testdeep.Struct(
			tmpStruct{},
			testdeep.StructFields{
				"Iface": expected,
			}),
		ifaceExpectedError(t, expectedError),
		args...)
}

func checkErrorForEach(t *testing.T,
	gotList []interface{}, expected interface{},
	expectedError expectedError, args ...interface{}) (ret bool) {
	t.Helper()

	globalTestName := test.BuildTestName(args)
	ret = true

	for idx, got := range gotList {
		testName := fmt.Sprintf("Got #%d", idx)
		if globalTestName != "" {
			testName += ", " + globalTestName
		}
		ret = checkError(t, got, expected, expectedError, testName) && ret
	}
	return
}

func _checkOK(t *testing.T, got, expected interface{},
	args ...interface{}) bool {
	t.Helper()

	if !testdeep.CmpDeeply(t, got, expected, args...) {
		return false
	}

	if !testdeep.EqDeeply(got, expected) {
		t.Errorf(`%sBoolean context failed
	     got: false
	expected: true`, test.BuildTestName(args))
		return false
	}

	if err := testdeep.EqDeeplyError(got, expected); err != nil {
		t.Errorf(`%sEqDeeplyError returned an error: %s`,
			test.BuildTestName(args), err)
		return false
	}

	return true
}

// checkOK calls _checkOK twice. The first time with the same
// parameters, the second time in an interface{} context.
func checkOK(t *testing.T, got, expected interface{},
	args ...interface{}) bool {
	t.Helper()

	if ok := _checkOK(t, got, expected, args...); !ok {
		return false
	}

	type tmpStruct struct {
		Iface interface{}
	}

	// Dirty hack to force got be passed as an interface kind
	return _checkOK(t, tmpStruct{Iface: got},
		testdeep.Struct(
			tmpStruct{},
			testdeep.StructFields{
				"Iface": expected,
			}),
		args...)
}

func checkOKOrPanicIfUnsafeDisabled(t *testing.T, got, expected interface{},
	args ...interface{}) bool {
	t.Helper()

	var ret bool
	cmp := func() {
		t.Helper()
		ret = _checkOK(t, got, expected, args...)
	}

	// Should panic if unsafe package is not available
	if dark.UnsafeDisabled {
		return test.CheckPanic(t, cmp,
			"dark.GetInterface() does not handle private ")
	}

	cmp()
	return ret
}

func checkOKForEach(t *testing.T, gotList []interface{}, expected interface{},
	args ...interface{}) (ret bool) {
	t.Helper()

	globalTestName := test.BuildTestName(args)
	ret = true

	for idx, got := range gotList {
		testName := fmt.Sprintf("Got #%d", idx)
		if globalTestName != "" {
			testName += ", " + globalTestName
		}
		ret = checkOK(t, got, expected, testName) && ret
	}
	return
}

func equalTypes(t *testing.T, got testdeep.TestDeep, expected interface{},
	args ...interface{}) bool {
	gotType := got.TypeBehind()
	expectedType := reflect.TypeOf(expected)

	if gotType == expectedType {
		return true
	}

	var gotStr, expectedStr string

	if gotType == nil {
		gotStr = "nil"
	} else {
		gotStr = gotType.String()
	}

	if expected == nil {
		expectedStr = "nil"
	} else {
		expectedStr = expectedType.String()
	}

	t.Helper()
	t.Errorf(`%sFailed test
	     got: %s
	expected: %s`,
		test.BuildTestName(args), gotStr, expectedStr)
	return false
}
