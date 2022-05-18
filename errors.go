package i18n

import (
	"fmt"

	"golang.org/x/text/language"
)

// Language tag already exists.
type ErrorLanguageTagAlreadyExists struct {
	Tag language.Tag
}

// Error message.
func (i *ErrorLanguageTagAlreadyExists) Error() string {
	return fmt.Sprintf("language tag `%v` already exists", i.Tag)
}

// Fallback language tag doesn't exists.
type ErrorFallbackTagNotExists struct {
	Tag language.Tag
}

// Error message.
func (i *ErrorFallbackTagNotExists) Error() string {
	return fmt.Sprintf("fallback language tag `%v` doesn't exists", i.Tag)
}

// Message structure validation error.
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
