package i18n

import (
	"encoding/json"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func Test_New(t *testing.T) {
	t.Parallel()

	New(&Config{})
}

func Test_I18nPrinter(t *testing.T) {
	testCases := []struct {
		name string
		i18n func() *I18n
		in   string
		out  string
	}{
		{
			name: "translate doesn't exists",
			i18n: func() *I18n {
				i18n := New(&Config{Fallback: language.English})

				i18n.printer[language.English] = message.NewPrinter(language.English)
				i18n.printer[language.Russian] = message.NewPrinter(language.Russian)

				err := message.SetString(language.English, "apple", "Apple")
				assert.NoError(t, err)

				return i18n
			},
			in:  "apple",
			out: "apple",
		},
		{
			name: "native language text",
			i18n: func() *I18n {
				i18n := New(&Config{Fallback: language.English})

				i18n.printer[language.English] = message.NewPrinter(language.English)
				i18n.printer[language.Russian] = message.NewPrinter(language.Russian)

				err := message.SetString(language.English, "apple", "Apple")
				assert.NoError(t, err)

				err = message.SetString(language.Russian, "apple", "Яблоко")
				assert.NoError(t, err)

				return i18n
			},
			in:  "apple",
			out: "Яблоко",
		},
		{
			name: "fallback language",
			i18n: func() *I18n {
				i18n := New(&Config{Fallback: language.English})

				i18n.printer[language.English] = message.NewPrinter(language.English)

				err := message.SetString(language.English, "apple", "Apple")
				assert.NoError(t, err)

				return i18n
			},
			in:  "apple",
			out: "Apple",
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				i18n := tc.i18n()
				p := i18n.Printer(language.Russian)

				assert.Equal(t, tc.out, p.Sprintf(tc.in))
			},
		)
	}
}

func Test_I18nTag(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		i18n *I18n
		dir  func() []fs.DirEntry
		err  bool
	}{
		{
			name: "successfully",
			i18n: &I18n{
				languages: make([]lang, 0),
				printer:   make(map[language.Tag]*message.Printer),
				config:    &Config{Fallback: language.English},
			},
			dir: func() []fs.DirEntry {
				var result []fs.DirEntry

				lang := []string{"ru", "en"}

				for _, l := range lang {
					result = append(result, newFakeDirEntry(fakeDirEntry{
						name: l,
					}))
				}

				return result
			},
		},
		{
			name: "wrong file name",
			i18n: &I18n{
				languages: make([]lang, 0),
				printer:   make(map[language.Tag]*message.Printer),
				config:    &Config{Fallback: language.English},
			},
			dir: func() []fs.DirEntry {
				var result []fs.DirEntry

				lang := []string{"-"}

				for _, l := range lang {
					result = append(result, newFakeDirEntry(fakeDirEntry{
						name: l,
					}))
				}

				return result
			},
			err: true,
		},
		{
			name: "duplicate language tag",
			i18n: &I18n{
				languages: make([]lang, 0),
				printer:   make(map[language.Tag]*message.Printer),
				config:    &Config{Fallback: language.English},
			},
			dir: func() []fs.DirEntry {
				var result []fs.DirEntry

				lang := []string{"ru", "en", "ru"}

				for _, l := range lang {
					result = append(result, newFakeDirEntry(fakeDirEntry{
						name: l,
					}))
				}

				return result
			},
			err: true,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				err := tc.i18n.tag(tc.dir())

				if tc.err {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}

func Test_I18nFallback(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		i18n *I18n
		err  bool
	}{
		{
			name: "successfully",
			i18n: &I18n{
				languages: []lang{{tag: language.English}, {tag: language.Russian}},
				config:    &Config{Fallback: language.English},
			},
		},
		{
			name: "fallback language is not defined",
			i18n: &I18n{
				languages: []lang{{tag: language.English}, {tag: language.Russian}},
				config:    &Config{Fallback: language.Und},
			},
		},
		{
			name: "language doesn't exists",
			i18n: &I18n{
				languages: []lang{{tag: language.English}, {tag: language.Russian}},
				config:    &Config{Fallback: language.German},
			},
			err: true,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				err := tc.i18n.fallback()

				if tc.err {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}

func Test_TranslationAppend(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		data string
		err  bool
	}{
		{
			name: "successfully",
			data: "[{\"id\": \"message id\", \"message\": \"message text\"}]",
		},
		{
			name: "wrong format",
			err:  true,
		},
	}

	tr := &translation{tag: language.English}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				err := tr.append([]byte(tc.data))

				if tc.err {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}

func Test_TranslationLoadMessages(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		messages string
		err      bool
	}{
		{
			name: "successfully",
			messages: `[
        {
          "id": "id #1",
          "message": "message #1"
        },
        {
          "id": "id #2",
          "message": "message #2"
        }
      ]`,
		},
		{
			name: "validation error",
			messages: `[
        {
          "id": "id #1",
          "message": "message #1"
        },
        {
          "message": "message #2"
        }
      ]`,
			err: true,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				tr := &translation{tag: language.English}

				var data []translationMessage

				err := json.Unmarshal([]byte(tc.messages), &data)
				assert.NoError(t, err)

				err = tr.loadMessages(data)

				if tc.err {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}

func Test_TranslationValidateMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		message string
		err     bool
	}{
		{
			name: "successfully",
			message: `{
        "id": "id",
        "message": "message"
      }`,
		},
		{
			name: "id is not set",
			message: `{
        "message": "message"
      }`,
			err: true,
		},
		{
			name: "message and rules is not set",
			message: `{
        "id": "id"
      }`,
			err: true,
		},
		{
			name: "message and rules is set",
			message: `{
        "id": "id",
        "message": "message",
        "rules": {
          "1": {
            "=0": "яблок нет",
            "one": "есть один ящик с яблоками",
            "few": "есть %d ящика с яблоками",
            "many": "есть %d ящиков с яблоками"
          }
        }
      }`,
			err: true,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				tr := &translation{tag: language.English}

				var data translationMessage

				err := json.Unmarshal([]byte(tc.message), &data)
				assert.NoError(t, err)

				err = tr.validateMessage(data)

				if tc.err {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}

func Test_TranslationLoadMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		message string
		argv    []interface{}
		out     string
	}{
		{
			name: "successfully",
			message: `{
        "id": "apple",
        "message": "яблоки"
      }`,
			out: "яблоки",
		},
		{
			name: "successfully with rules",
			message: `{
        "id": "apple",
        "rules": {
          "=0":    "яблок нет",
          "one":   "есть одно яблоко",
          "few":   "есть %d яблока",
          "many":  "есть %d яблок"
        }
      }`,
			argv: []interface{}{1},
			out:  "есть одно яблоко",
		},
		{
			name: "successfully with rules",
			message: `{
        "id": "apple",
        "rules": {
          "2": {
            "=0": "яблок нет",
            "one": "есть %d яблок в %d ящике",
            "few": "есть %d яблок в %d ящиках",
            "many": "есть %d яблок в %d ящиках"
          }
        }
      }`,
			argv: []interface{}{100, 1},
			out:  "есть 100 яблок в 1 ящике",
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				tr := &translation{tag: language.Russian}
				p := message.NewPrinter(language.Russian)

				var data translationMessage

				err := json.Unmarshal([]byte(tc.message), &data)
				assert.NoError(t, err)

				assert.NoError(t, tr.loadMessage(data))

				assert.Equal(t, tc.out, p.Sprintf("apple", tc.argv...))
			},
		)
	}
}

func Test_TranslationLoadRules(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		message string
		err     bool
	}{
		{
			name: "successfully",
			message: `{
        "id" : "apple",
        "rules": {
          "=0":    "яблок нет",
          "one":   "есть одно яблоко",
          "few":   "есть %d яблока",
          "many":  "есть %d яблок",
          "other": "есть яблоки"
        }
      }`,
		},
		{
			name: "successfully",
			message: `{
        "id" : "apple",
        "rules": {
          "1": {
            "=0":    "яблок нет",
            "one":   "есть одно яблоко",
            "few":   "есть %d яблока",
            "many":  "есть %d яблок",
            "other": "есть яблоки"
          }
        }
      }`,
		},
		{
			name: "wrong rules",
			message: `{
        "id" : "apple",
        "rules": ""
      }`,
			err: true,
		},
		{
			name: "wrong argument number",
			message: `{
        "id" : "apple",
        "rules": {
          "-": {
            "=0":    "яблок нет",
            "one":   "есть одно яблоко",
            "few":   "есть %d яблока",
            "many":  "есть %d яблок",
            "other": "есть яблоки"
          }
        }
      }`,
			err: true,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				tr := &translation{tag: language.Russian}

				var data translationMessage

				err := json.Unmarshal([]byte(tc.message), &data)
				assert.NoError(t, err)

				err = tr.loadRules(data)

				if tc.err {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}

// Implement fs.DirEntry interface.
type fakeDirEntry struct {
	name     string
	isDir    bool
	fileMode fs.FileMode
	fileInfo fs.FileInfo
	err      error
}

// Name returns the name of the file (or subdirectory) described by the entry.
// This name is only the final element of the path (the base name), not the entire path.
// For example, Name would return "hello.go" not "home/gopher/hello.go".
func (i *fakeDirEntry) Name() string {
	return i.name
}

// IsDir reports whether the entry describes a directory.
func (i *fakeDirEntry) IsDir() bool {
	return i.isDir
}

// Type returns the type bits for the entry.
// The type bits are a subset of the usual FileMode bits, those returned by the FileMode.Type method.
func (i *fakeDirEntry) Type() fs.FileMode {
	return i.fileMode
}

// Info returns the FileInfo for the file or subdirectory described by the entry.
// The returned FileInfo may be from the time of the original directory read
// or from the time of the call to Info. If the file has been removed or renamed
// since the directory read, Info may return an error satisfying errors.Is(err, ErrNotExist).
// If the entry denotes a symbolic link, Info reports the information about the link itself,
// not the link's target.
func (i *fakeDirEntry) Info() (fs.FileInfo, error) {
	return i.fileInfo, i.err
}

// Create a new fake dir entry.
func newFakeDirEntry(f fakeDirEntry) *fakeDirEntry {
	return &f
}
