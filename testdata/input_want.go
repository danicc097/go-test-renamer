package testdata

import "testing"

func TestSomething(t *testing.T) {
	someConstant := " something"
	someOtherConstant := " something else"
	t.Run("test_with_spaces", func(t *testing.T) {
		t.Parallel()

		t.Run("deeply nested"+someConstant, func(t *testing.T) {
			t.Parallel()
		})
	})

	t.Run("another__t.Run()", func(t *testing.T) {
		t.Parallel()
	})

	t.Run(someOtherConstant, func(t *testing.T) {
		t.Parallel()
	})

	type testCase struct {
		name string
	}

	tests := []testCase{
		{
			name: "a_b_c",
		},
		{
			name: "1_2_3",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
		})
	}
}
