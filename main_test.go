package main

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"golang.org/x/text/unicode/rangetable"
)

func getTestParams() *gopter.TestParameters {
	tp := gopter.DefaultTestParameters()
	if !testing.Short() {
		tp.MinSuccessfulTests = 1000
	}
	return tp
}

func overrideTestingRun(p *gopter.Properties, t *testing.T) {
	b := &strings.Builder{}
	if _, err := b.WriteString("\n"); err != nil {
		t.Errorf("overrideTestingRun broke appending newline: %v", err)
	}
	reporter := gopter.NewFormatedReporter(true, 75, b)
	if !p.Run(reporter) {
		t.Error(b.String())
	} else {
		t.Log(b.String())
	}
}

func TestMain(t *testing.T) {
	testCases := []struct {
		name     string
		in       string
		expected string
	}{
		{"ascii", "hello world", "ｈｅｌｌｏ ｗｏｒｌｄ"},
		{"full-width", "ｈｅｌｌｏ ｗｏｒｌｄ", "ｈｅｌｌｏ ｗｏｒｌｄ"},
		{"japanese", "日本語", "日本語"},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			out := &strings.Builder{}
			if err := convert(strings.NewReader(tc.in), out); err != nil {
				t.Error(err)
			}
			if got := out.String(); got != tc.expected {
				t.Logf("in = %x; got = %x; expected = %x", tc.in, got, tc.expected)
				t.Errorf("convert(%q) = %q; want %q", tc.in, got, tc.expected)
			} else {
				t.Logf("convert(%q) = %q", tc.in, got)
			}
		})
	}

	t.Run("gopter", func(t *testing.T) {
		properties := gopter.NewProperties(getTestParams())

		properties.Property("strange unchanged", prop.ForAll(
			func(s string) bool {
				out := &strings.Builder{}
				if err := convert(strings.NewReader(s), out); err != nil {
					t.Log(err)
					return false
				}
				return out.String() == s
			},
			// gen.UnicodeString(&unicode.RangeTable{
			// 	R16: []unicode.Range16{unicode.Range16{0xFF00, 0xFFEF, 1}},
			// }),
			gen.UnicodeString(rangetable.Merge([]*unicode.RangeTable{
				unicode.Braille,
				unicode.Cyrillic,
				unicode.Greek,
				unicode.Han,
				unicode.Hangul,
				unicode.Hebrew,
				unicode.Hiragana,
				unicode.Katakana,
				//unicode.Latin,
			}...)),
		))

		properties.Property("basic-latin changed", prop.ForAll(
			func(s string) bool {
				out := &strings.Builder{}
				if err := convert(strings.NewReader(s), out); err != nil {
					t.Log(err)
					return false
				}
				return out.String() != s
			},
			gen.RegexMatch(`[!-~]+`),
		))

		properties.Property("length invariant", prop.ForAll(
			func(s string) bool {
				out := &strings.Builder{}
				if err := convert(strings.NewReader(s), out); err != nil {
					t.Log(err)
					return false
				}
				return utf8.RuneCountInString(out.String()) == utf8.RuneCountInString(s)
			},
			gen.AnyString(),
		))

		overrideTestingRun(properties, t)
	})
}
