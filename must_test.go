package gofig

import (
	"errors"
	"testing"
)

func TestMust(t *testing.T) {
	panicerr := errors.New("boom")

	cases := map[string]struct {
		fn   func() error
		want error
	}{
		"PanicsOnError": {
			fn: func() error {
				return panicerr
			},
			want: panicerr,
		},
		"DoesNotPanicOnNilError": {
			fn: func() error {
				return nil
			},
			want: nil,
		},
	}

	for name, testCase := range cases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				r := recover()

				if tc.want == nil && r == nil {
					return
				}

				err, ok := r.(error)
				if !ok {
					t.Errorf("r is not an error: %T", r)
					t.FailNow()
				}

				t.Log(err)

				if !errors.Is(err, tc.want) {
					t.Errorf("want %v, got %v", tc.want, err)
				}
			}()

			Must(tc.fn())
		})
	}
}
