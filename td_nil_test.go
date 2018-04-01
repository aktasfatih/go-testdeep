package testdeep_test

import (
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestNil(t *testing.T) {
	checkOK(t, (func())(nil), Nil())
	checkOK(t, ([]int)(nil), Nil())
	checkOK(t, (map[bool]bool)(nil), Nil())
	checkOK(t, (*int)(nil), Nil())
	checkOK(t, (chan int)(nil), Nil())
	checkOK(t, nil, Nil())

	checkError(t, 42, Nil(),
		expectedError{
			Message:  mustBe("non-nil"),
			Path:     mustBe("DATA"),
			Got:      mustBe("(int) 42"),
			Expected: mustBe("nil"),
		})

	num := 42
	checkError(t, &num, Nil(),
		expectedError{
			Message:  mustBe("non-nil"),
			Path:     mustBe("DATA"),
			Got:      mustMatch(`\(\*int\).*42`),
			Expected: mustBe("nil"),
		})

	//
	// String
	equalStr(t, Nil().String(), "nil")
}

func TestNotNil(t *testing.T) {
	num := 42
	checkOK(t, func() {}, NotNil())
	checkOK(t, []int{}, NotNil())
	checkOK(t, map[bool]bool{}, NotNil())
	checkOK(t, &num, NotNil())
	checkOK(t, 42, NotNil())

	checkError(t, (func())(nil), NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, ([]int)(nil), NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, (map[bool]bool)(nil), NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, (*int)(nil), NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, (chan int)(nil), NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustContain("nil"),
			Expected: mustBe("not nil"),
		})
	checkError(t, nil, NotNil(),
		expectedError{
			Message:  mustBe("nil value"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("not nil"),
		})

	//
	// String
	equalStr(t, NotNil().String(), "not nil")
}