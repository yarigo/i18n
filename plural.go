package i18n

import (
	"reflect"
	"strings"

	_plural "golang.org/x/text/feature/plural"
	"golang.org/x/text/message/catalog"
)

type plural struct{}

// Set a plural rule.
func (i plural) set(arg int, key string, rules reflect.Value) catalog.Message {
	var match string
	cases := strings.Split(key, ":")
	if len(cases) == 2 {
		match = cases[1]
	} else {
		match = cases[0]
	}

	// values := reflect.ValueOf(rules)
	if rules.Kind() == reflect.Interface {
		return _plural.Selectf(arg, "", match, rules.Interface())
	}
	return nil
}
