package i18n

import (
	"fmt"

	"golang.org/x/text/language"
)

// ErrorLanguageTagAlreadyExists reports a duplicate language tag.
type ErrorLanguageTagAlreadyExists struct {
	Tag language.Tag
}

// Error message.
func (i *ErrorLanguageTagAlreadyExists) Error() string {
	return fmt.Sprintf("language tag `%v` already exists", i.Tag)
}

// ErrorFallbackTagNotExists reports that translation files for fallback
// language do not exist.
type ErrorFallbackTagNotExists struct {
	Tag language.Tag
}

// Error message.
func (i *ErrorFallbackTagNotExists) Error() string {
	return fmt.Sprintf("fallback language tag `%v` doesn't exists", i.Tag)
}

// ErrorMessageValidate reports a validation error.
type ErrorMessageValidate struct {
	Tag       language.Tag
	FilePath  string
	MessageID string
	Message   string
}

// Error message.
func (i *ErrorMessageValidate) Error() string {
	return fmt.Sprintf(
		"validation error in file %v for language %v (id: %v):\n%v\n",
		i.FilePath,
		i.Tag.String(),
		i.MessageID,
		i.Message,
	)
}
