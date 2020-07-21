package i18n

import (
	"reflect"
	"testing"

	"golang.org/x/text/language"
)

// Test settings.
var cfg struct {
	parallel bool // run in parallel mode
}

func init() {
	cfg.parallel = true
}

func TestRules(t *testing.T) {
	if cfg.parallel {
		t.Parallel()
	}

	type test struct {
		name     string
		loadLang language.Tag
		data     []byte
		in       string
		lang     language.Tag
		params   interface{}
		out      string
		err      bool
	}

	tests := []test{
		{
			name:     "unknown language",
			loadLang: language.Russian,
			data: []byte(`{
				"messages": [
					{
						"id": "test",
						"text": "яблоко"
					}
				]
			}`),
			in:   "banana",
			lang: language.English,
			out:  "banana",
		},
		{
			name:     "unknown id",
			loadLang: language.Russian,
			data: []byte(`{
				"messages": [
					{
						"id": "test",
						"text": "яблоко"
					}
				]
			}`),
			in:   "banana",
			lang: language.Russian,
			out:  "banana",
		},
		{
			name:     "string",
			loadLang: language.Russian,
			data: []byte(`{
				"messages": [
					{
						"id": "test",
						"text": "яблоко"
					}
				]
			}`),
			in:   "test",
			lang: language.Russian,
			out:  "яблоко",
		},
		{
			name:     "duplicate id",
			loadLang: language.Russian,
			data: []byte(`{
				"messages": [
					{
						"id": "test",
						"text": "яблоко"
					},
					{
						"id": "test",
						"text": "банан"
					}
				]
			}`),
			in:   "test",
			lang: language.Russian,
			out:  "яблоко",
			err:  true,
		},
		{
			name:     "multiple variables",
			loadLang: language.Russian,
			data: []byte(`{
				"messages": [
					{
						"id": "test",
						"text": "%s был в %s"
					}
				]
			}`),
			in:     "test",
			lang:   language.Russian,
			params: []interface{}{"Джо", "Париже"},
			out:    "Джо был в Париже",
		},
		{
			name:     "plural without rule prefix",
			loadLang: language.Russian,
			data: []byte(`{
				"messages": [
					{
						"id": "test",
						"text": "есть яблоки",
						"rules": {
							"one": "есть %d яблоко",
							"=6": "есть %d яблок"
						}
					}
				]
			}`),
			in:     "test",
			lang:   language.Russian,
			params: 6,
			out:    "есть 6 яблок",
		},
		{
			name:     "plural with rule prefix",
			loadLang: language.Russian,
			data: []byte(`{
				"messages": [
					{
						"id": "test",
						"text": "есть яблоки",
						"rules": {
							"plural:one": "есть %d яблоко",
							"plural:many": "есть %d яблок"
						}
					}
				]
			}`),
			in:     "test",
			lang:   language.Russian,
			params: 10,
			out:    "есть 10 яблок",
		},
		{
			name:     "variant selector",
			loadLang: language.Russian,
			data: []byte(`{
				"messages": [
					{
						"id": "test",
						"text": "%[1]s invite %[2]s and %[3]d other guests to their party.",
						"rules": [
							{
								"=0": "There is no party. Move on!",
								"=1": ""
							}
						]
					}
				]
			}`),
			in:     "test",
			lang:   language.Russian,
			params: 10,
			out:    "есть 10 яблок",
		},
	}

	for _, val := range tests {
		p := New()
		err := p.Load(val.loadLang, &val.data)
		if err != nil && !val.err {
			t.Fatalf("%v is failed: %v", val.name, err)
		} else if err == nil && val.err {
			t.Fatalf("%v is failed: we expected a error", val.name)
		} else if err != nil && val.err {
			continue
		}

		l := p.L10n(val.lang)
		var out string
		if reflect.ValueOf(val.params).Kind() == reflect.Slice {
			out = l.Sprintf(val.in, val.params.([]interface{})...)
		} else {
			out = l.Sprintf(val.in, val.params)
		}
		if val.out != out {
			t.Fatalf(
				"%v is failed:\n\twant(%v)\n\thave(%v)",
				val.name,
				val.out,
				out,
			)
		}
	}
}
