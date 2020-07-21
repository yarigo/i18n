package i18n

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

// Rule type.
type rule interface {
	set(int, string, reflect.Value) catalog.Message
}

// Plural.
type pluraleRule struct{}

// I18n structure data.
type I18n struct {
	languages map[language.Tag]*message.Printer
	keys      map[language.Tag][]string
}

// Structure of the translation phrase.
type phrase struct {
	ID    string      `json:"id"`
	Text  string      `json:"text"`
	Rules interface{} `json:"rules"`
}

// Structure of the translation file.
type translation struct {
	Phrase []phrase `json:"messages"`
}

// New instance.
func New() *I18n {
	return &I18n{
		languages: make(map[language.Tag]*message.Printer),
		keys:      make(map[language.Tag][]string),
	}
}

// Load translation data.
func (i *I18n) Load(lang language.Tag, data *[]byte) error {
	var t translation
	err := json.Unmarshal(*data, &t)
	if err != nil {
		return err
	}

	err = i.rules(lang, &t)
	if err != nil {
		return err
	}

	i.languages[lang] = message.NewPrinter(lang)

	return nil
}

// L10n return current localization printer.
func (i *I18n) L10n(lang language.Tag) *message.Printer {
	if p, ok := i.languages[lang]; ok {
		return p
	}
	return message.NewPrinter(language.Und)
}

// Set rules for language.
func (i *I18n) rules(lang language.Tag, data *translation) error {
	for _, val := range data.Phrase {
		err := i.contain(lang, val.ID)
		if err != nil {
			return err
		}
		if val.Rules != nil {
			if reflect.TypeOf(val.Rules).Kind() == reflect.Map {
				values := reflect.ValueOf(val.Rules)
				keys := values.MapKeys()
				for _, key := range keys {
					rule := i.getType(key.String())
					if rule != nil {
						err = message.Set(
							lang,
							val.ID,
							rule.set(1, key.String(), values.MapIndex(key)),
						)
						if err != nil {
							return err
						}
					}
				}
			}
			// TODO: Make multiple rules.
			// else if reflect.TypeOf(val.Rules).Kind() == reflect.Slice {
			// 	values := reflect.ValueOf(val.Rules)
			// 	for k := 0; k < values.Len(); k++ {
			// 		rule := i.getType(key.String())
			// 		fmt.Println(values.Index(k))
			// 	}
			// }
		} else {
			err := message.SetString(lang, val.ID, val.Text)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Get rule type.
func (i *I18n) getType(key string) rule {
	rule := strings.Split(key, ":")
	if len(rule) == 1 || rule[0] == "plural" {
		return plural{}
	}

	return nil
}

// Check duplicate key per language.
func (i *I18n) contain(lang language.Tag, key string) error {
	if _, ok := i.keys[lang]; ok {
		for _, val := range i.keys[lang] {
			if val == key {
				return fmt.Errorf("duplicate id `%v`", key)
			}
		}
	} else {
		i.keys[lang] = make([]string, 0)
	}

	i.keys[lang] = append(i.keys[lang], key)

	return nil
}
