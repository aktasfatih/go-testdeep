// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestZero(t *testing.T) {
	checkOK(t, 0, Zero())
	checkOK(t, int64(0), Zero())
	checkOK(t, float64(0), Zero())
	checkOK(t, nil, Zero())
	checkOK(t, (map[string]int)(nil), Zero())
	checkOK(t, ([]int)(nil), Zero())
	checkOK(t, [3]int{}, Zero())
	checkOK(t, MyStruct{}, Zero())
	checkOK(t, (*MyStruct)(nil), Zero())
	checkOK(t, &MyStruct{}, Ptr(Zero()))
	checkOK(t, (chan int)(nil), Zero())
	checkOK(t, (func())(nil), Zero())
	checkOK(t, false, Zero())

	checkError(t, 12, Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int) 12"),
		Expected: mustBe("(int) 0"),
	})
	checkError(t, int64(12), Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int64) 12"),
		Expected: mustBe("(int64) 0"),
	})
	checkError(t, float64(12), Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(float64) 12"),
		Expected: mustBe("(float64) 0"),
	})
	checkError(t, map[string]int{}, Zero(), expectedError{
		Message:  mustBe("nil map"),
		Path:     mustBe("DATA"),
		Got:      mustBe("not nil"),
		Expected: mustBe("nil"),
	})
	checkError(t, []int{}, Zero(), expectedError{
		Message:  mustBe("nil slice"),
		Path:     mustBe("DATA"),
		Got:      mustBe("not nil"),
		Expected: mustBe("nil"),
	})
	checkError(t, [3]int{0, 12}, Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA[1]"),
		Got:      mustBe("(int) 12"),
		Expected: mustBe("(int) 0"),
	})
	checkError(t, MyStruct{ValInt: 12}, Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA.ValInt"),
		Got:      mustBe("(int) 12"),
		Expected: mustBe("(int) 0"),
	})
	checkError(t, &MyStruct{}, Zero(), expectedError{
		Message: mustBe("values differ"),
		Path:    mustBe("*DATA"),
		// in fact, pointer on 0'ed struct contents
		Got:      mustContain(`ValInt: (int) 0`),
		Expected: mustBe("nil"),
	})
	checkError(t, true, Zero(), expectedError{
		Message:  mustBe("values differ"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(bool) true"),
		Expected: mustBe("(bool) false"),
	})

	//
	// String
	equalStr(t, Zero().String(), "Zero()")
}

func TestNotZero(t *testing.T) {
	checkOK(t, 12, NotZero())
	checkOK(t, int64(12), NotZero())
	checkOK(t, float64(12), NotZero())
	checkOK(t, map[string]int{}, NotZero())
	checkOK(t, []int{}, NotZero())
	checkOK(t, [3]int{1}, NotZero())
	checkOK(t, MyStruct{ValInt: 1}, NotZero())
	checkOK(t, &MyStruct{}, NotZero())
	checkOK(t, make(chan int), NotZero())
	checkOK(t, func() {}, NotZero())
	checkOK(t, true, NotZero())

	checkError(t, nil, NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("nil"),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, 0, NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int) 0"),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, int64(0), NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(int64) 0"),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, float64(0), NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(float64) 0"),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, (map[string]int)(nil), NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(map[string]int) <nil>"),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, ([]int)(nil), NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("([]int) <nil>"),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, [3]int{}, NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustContain("(int) 0"),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, MyStruct{}, NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustContain(`ValInt: (int) 0`),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, &MyStruct{}, Ptr(NotZero()), expectedError{
		Message: mustBe("zero value"),
		Path:    mustBe("*DATA"),
		// in fact, pointer on 0'ed struct contents
		Got:      mustContain(`ValInt: (int) 0`),
		Expected: mustBe("NotZero()"),
	})
	checkError(t, false, NotZero(), expectedError{
		Message:  mustBe("zero value"),
		Path:     mustBe("DATA"),
		Got:      mustBe("(bool) false"),
		Expected: mustBe("NotZero()"),
	})

	//
	// String
	equalStr(t, NotZero().String(), "NotZero()")
}

func TestZeroTypeBehind(t *testing.T) {
	equalTypes(t, Zero(), nil)
	equalTypes(t, NotZero(), nil)
}