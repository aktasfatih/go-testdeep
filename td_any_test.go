// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep_test

import (
	"testing"

	"github.com/maxatome/go-testdeep"
	"github.com/maxatome/go-testdeep/internal/test"
)

func TestAny(t *testing.T) {
	checkOK(t, 6, testdeep.Any(nil, 5, 6, 7))
	checkOK(t, nil, testdeep.Any(5, 6, 7, nil))

	checkError(t, 6, testdeep.Any(5),
		expectedError{
			Message:  mustBe("comparing with Any"),
			Path:     mustBe("DATA"),
			Got:      mustBe("6"),
			Expected: mustBe("Any(5)"),
		})

	checkError(t, 6, testdeep.Any(nil),
		expectedError{
			Message:  mustBe("comparing with Any"),
			Path:     mustBe("DATA"),
			Got:      mustBe("6"),
			Expected: mustBe("Any(nil)"),
		})

	checkError(t, nil, testdeep.Any(6),
		expectedError{
			Message:  mustBe("comparing with Any"),
			Path:     mustBe("DATA"),
			Got:      mustBe("nil"),
			Expected: mustBe("Any(6)"),
		})

	//
	// String
	test.EqualStr(t, testdeep.Any(6).String(), "Any(6)")
	test.EqualStr(t, testdeep.Any(6, 7).String(), "Any(6,\n    7)")
}

func TestAnyTypeBehind(t *testing.T) {
	equalTypes(t, testdeep.Any(6), nil)
}
