package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func Test_Contains(t *testing.T) {
	t.Parallel()

	type result struct {
		key int
		ok  bool
	}

	testCases := []struct {
		name   string
		slice  []lang
		value  language.Tag
		result result
	}{
		{
			name: "value is exists",
			slice: []lang{
				{tag: language.English},
				{tag: language.Russian},
			},
			value:  language.Russian,
			result: result{key: 1, ok: true},
		},
		{
			name: "value doesn't exists",
			slice: []lang{
				{tag: language.English},
				{tag: language.Russian},
			},
			value:  language.German,
			result: result{key: -1, ok: false},
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				key, ok := contains(tc.slice, tc.value)

				assert.Equal(t, tc.result.key, key)
				assert.Equal(t, tc.result.ok, ok)
			},
		)
	}
}
