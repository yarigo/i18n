package i18n

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Config of i18n.
type Config struct {
	// Path to the languages folder.
	// File or directory name use as language tag.
	Path string
	// Fallback language.
	Fallback language.Tag
}

// I18n data.
type I18n struct {
	languages []lang
	printer   map[language.Tag]*message.Printer
	config    *Config
}

// Language properties.
type lang struct {
	tag   language.Tag
	entry fs.DirEntry
}

// New instance of i18n.
func New(cfg *Config) *I18n {
	return &I18n{
		languages: make([]lang, 0),
		printer:   make(map[language.Tag]*message.Printer),
		config:    cfg,
	}
}

// Load all locales.
func (i *I18n) Load() (err error) {
	var files []fs.DirEntry

	// Read locales directory.
	files, err = os.ReadDir(i.config.Path)
	if err != nil {
		return
	}

	// Create a languages tags.
	if err = i.tag(files); err != nil {
		return
	}

	return i.load()
}

// Printer implements language-specific formatted I/O analogous to the fmt
// package.
func (i *I18n) Printer(tag language.Tag) *message.Printer {
	if printer, ok := i.printer[tag]; ok {
		return printer
	}

	return message.NewPrinter(i.config.Fallback)
}

// Get a language tag.
func (i *I18n) tag(dir []fs.DirEntry) (err error) {
	var lang lang
	var ok bool

	for _, lang.entry = range dir {
		lang.tag, err = language.Parse(lang.entry.Name())
		if err != nil {
			return
		}

		if _, ok = contains(i.languages, lang.tag); ok {
			return &ErrorLanguageTagAlreadyExists{Tag: lang.tag}
		}

		i.languages = append(i.languages, lang)

		i.printer[lang.tag] = message.NewPrinter(lang.tag)
	}

	return i.fallback()
}

// Test fallback language.
func (i *I18n) fallback() error {
	if i.config.Fallback == language.Und {
		return nil
	}

	if _, ok := contains(i.languages, i.config.Fallback); !ok {
		return &ErrorFallbackTagNotExists{Tag: i.config.Fallback}
	}

	return nil
}

// Load all languages files.
func (i *I18n) load() (err error) {
	for _, lang := range i.languages {
		if err = loadLanguage(lang.tag, lang.entry, i.config.Path); err != nil {
			return
		}
	}

	// Clean languages properties.
	i.languages = nil

	return
}

// Translation properties.
type translation struct {
	tag      language.Tag
	filePath string
}

// Load language files.
func loadLanguage(
	tag language.Tag,
	file fs.DirEntry,
	rootPath string,
) error {
	currentPath := filepath.Join(rootPath, file.Name())

	if file.IsDir() {
		files, err := os.ReadDir(currentPath)
		if err != nil {
			return err
		}

		for _, entry := range files {
			return loadLanguage(tag, entry, currentPath)
		}
	}

	t := &translation{tag: tag, filePath: currentPath}

	return t.loadLanguageFile()
}

// Append language file.
func (i *translation) loadLanguageFile() error {
	b, err := ioutil.ReadFile(i.filePath)
	if err != nil {
		return err
	}

	return i.append(b)
}

// Translation message properties.
type translationMessage struct {
	ID      string      `json:"id"`
	Message *string     `json:"message,omitempty"`
	Rules   interface{} `json:"rules,omitempty"`
}

// Append language messages.
func (i *translation) append(b []byte) (err error) {
	var m []translationMessage
	if err = json.Unmarshal(b, &m); err != nil {
		return
	}

	return i.loadMessages(m)
}

// Load translation messages.
func (i *translation) loadMessages(m []translationMessage) (err error) {
	for _, message := range m {
		if err = i.validateMessage(message); err != nil {
			return
		}

		if err = i.loadMessage(message); err != nil {
			return
		}
	}

	return
}

// Validate message structure.
func (i *translation) validateMessage(message translationMessage) error {
	if len(message.ID) == 0 {
		return &ErrorMessageValidate{
			Tag:      i.tag,
			FilePath: i.filePath,
			Message:  "message id is not set",
		}
	}

	if message.Message == nil && message.Rules == nil {
		return &ErrorMessageValidate{
			Tag:       i.tag,
			FilePath:  i.filePath,
			MessageID: message.ID,
			Message:   "`message` or `rules` field should be set",
		}
	}

	if message.Message != nil && message.Rules != nil {
		return &ErrorMessageValidate{
			Tag:       i.tag,
			FilePath:  i.filePath,
			MessageID: message.ID,
			Message:   "only one of field `message` or `rules` should be set",
		}
	}

	return nil
}

// Load translation message.
func (i *translation) loadMessage(m translationMessage) error {
	if m.Rules != nil {
		return i.loadRules(m)
	}

	return message.SetString(i.tag, m.ID, *m.Message)
}

// Load translation rules.
func (i *translation) loadRules(m translationMessage) (err error) {
	if reflect.TypeOf(m.Rules).Kind() != reflect.Map {
		return &ErrorMessageValidate{
			Tag:       i.tag,
			FilePath:  i.filePath,
			MessageID: m.ID,
			Message:   "`rules` filed should be a map",
		}
	}

	value := reflect.ValueOf(m.Rules)
	iter := value.MapRange()

	// The arg-th substitution argument
	arg := 1
	msg := make([]interface{}, 0)

	for iter.Next() {
		if iter.Value().Elem().Type().Kind() == reflect.Map {
			arg, err = strconv.Atoi(iter.Key().String())
			if err != nil {
				return
			}

			subIter := iter.Value().Elem().MapRange()

			for subIter.Next() {
				msg = append(msg, subIter.Key().String(), subIter.Value().Interface())
			}

			err = message.Set(i.tag, m.ID, plural.Selectf(arg, "", msg...))
			if err != nil {
				return
			}

			// Reset messages.
			msg = nil
		} else {
			msg = append(msg, iter.Key().String(), iter.Value().Interface())
		}
	}

	if len(msg) > 0 {
		err = message.Set(i.tag, m.ID, plural.Selectf(arg, "", msg...))
		if err != nil {
			return
		}
	}

	return
}
