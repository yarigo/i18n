package i18n

import (
	"golang.org/x/text/language"
)

// Check if slice of languages tags contains the language tag name.
func contains(s []lang, val language.Tag) (int, bool) {
	for i, v := range s {
		if v.tag == val {
			return i, true
		}
	}

	return -1, false
}
