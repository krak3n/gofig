package gofig

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	type Config struct {
		String          string                                 `gofig:"string"`
		Int             int                                    `gofig:"int"`
		Slice           []int                                  `gofig:"slice"`
		Map             map[string]string                      `gofig:"map"`
		NestedMap       map[string]map[string][]int            `gofig:"nestedmap"`
		DeeplyNestedMap map[string]map[string]map[string][]int `gofig:"deeplynestedmap"`
	}

	cases := map[string]struct {
		parser *InMemoryParser
		opts   []Option
		want   Config
	}{
		"String": {
			parser: func() *InMemoryParser {
				p := NewInMemoryParser()
				p.Add("string", "bar")

				return p
			}(),
			want: Config{
				String: "bar",
			},
		},
		"Int": {
			parser: func() *InMemoryParser {
				p := NewInMemoryParser()
				p.Add("int", int(1))

				return p
			}(),
			want: Config{
				Int: int(1),
			},
		},
		"Slice": {
			parser: func() *InMemoryParser {
				p := NewInMemoryParser()
				p.Add("slice", []int{1, 2, 3})

				return p
			}(),
			want: Config{
				Slice: []int{1, 2, 3},
			},
		},
		"Map": {
			parser: func() *InMemoryParser {
				p := NewInMemoryParser()
				p.Add("map.key", "value")

				return p
			}(),
			want: Config{
				Map: map[string]string{
					"key": "value",
				},
			},
		},
		"NestedMap": {
			parser: func() *InMemoryParser {
				p := NewInMemoryParser()
				p.Add("nestedmap.foo.bar", []int{1, 2, 3})
				p.Add("nestedmap.foo.baz", []int{4, 5, 6})

				return p
			}(),
			want: Config{
				NestedMap: map[string]map[string][]int{
					"foo": {
						"bar": {1, 2, 3},
						"baz": {4, 5, 6},
					},
				},
			},
		},
		"DeeplyNestedMap": {
			parser: func() *InMemoryParser {
				p := NewInMemoryParser()
				p.Add("deeplynestedmap.foo.bar.fizz", []int{1, 2, 3})
				p.Add("deeplynestedmap.foo.bar.buzz", []int{4, 5, 6})
				p.Add("deeplynestedmap.foo.baz.fizz", []int{1, 2, 3})
				p.Add("deeplynestedmap.foo.baz.buzz", []int{4, 5, 6})

				return p
			}(),
			want: Config{
				DeeplyNestedMap: map[string]map[string]map[string][]int{
					"foo": {
						"bar": {
							"fizz": {1, 2, 3},
							"buzz": {4, 5, 6},
						},
						"baz": {
							"fizz": {1, 2, 3},
							"buzz": {4, 5, 6},
						},
					},
				},
			},
		},
	}

	for name, testCase := range cases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var cfg Config

			opts := append(
				tc.opts,
				WithDebug(),
				SetLogger(LoggerFunc(func(v ...interface{}) {
					t.Log(v...)
				})))

			g, err := New(&cfg, opts...)

			if err != nil {
				t.Fatal("want nil error, got:", err)
			}

			if err := g.Parse(tc.parser); err != nil {
				t.Fatal("want nil error, got:", err)
			}

			if !cmp.Equal(tc.want, cfg) {
				t.Errorf("want %+v, got %+v", tc.want, cfg)
			}
		})
	}
}

func TestKeyFormatting(t *testing.T) {
	type Fizz struct {
		Buzz string `gofig:"buzz"`
	}

	type Config struct {
		Foo  string `gofig:"foo"`
		Fizz Fizz   `gofig:"FIZZ"`
	}

	cases := map[string]struct {
		parser *InMemoryParser
		opts   []Option
		want   Config
	}{
		"ParseUpperToLower": {
			parser: func() *InMemoryParser {
				p := NewInMemoryParser()
				p.Add("FOO", "bar")
				p.Add("FIZZ.BUZZ", "baz")

				return p
			}(),
			want: Config{
				Foo: "bar",
				Fizz: Fizz{
					Buzz: "baz",
				},
			},
		},
	}

	for name, testCase := range cases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var cfg Config

			opts := append(
				tc.opts,
				WithDebug(),
				SetLogger(LoggerFunc(func(v ...interface{}) {
					t.Log(v...)
				})))

			g, err := New(&cfg, opts...)

			if err != nil {
				t.Fatal("want nil error, got:", err)
			}

			if err := g.Parse(tc.parser); err != nil {
				t.Fatal("want nil error, got:", err)
			}

			if !cmp.Equal(tc.want, cfg) {
				t.Errorf("want %+v, got %+v", tc.want, cfg)
			}
		})
	}
}
