# i18n ![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)

i18n is a Go package that helps you translate Go programs into multiple languages.

## Features

- plural

## Installation

```shell
go get -u github.com/yarigo/i18n
```

## Rules

For set a current localization, use a language mather function of all supported languages and match from it by language string:

```go
matcher := language.NewMatcher([]language.Tag{
	language.English, // The first language is used as fallback.
	language.MustParse("en-AU"),
	language.Danish,
	language.Chinese,
	language.Russian,
})

lang, _ := language.MatchStrings(matcher, "ru")
```

Package support multiple definition by its number identification. As identification use a position into array rules. See tests for understanding this.

### plural

If you want to use plural function, you can set it as "plural:one" or "one" (without plural prefix).

A selector matches an argument if:

- it is "other" or Other
- it matches the plural form of the argument: "zero", "one", "two", "few",
  or "many", or the equivalent Form
- it is of the form "=x" where x is an integer that matches the value of
  the argument.
- it is of the form "<x" where x is an integer that is larger than the
  argument.

For use a format, set it as suffix. For example:

```json
"one:%d": "..."
```

or

```json
"plural:one:%d": "..."
```

Examples of format strings are:

- %.2f decimal with scale 2
- %.2e scientific notation with precision 3 (scale + 1)
- %d integer

## Locale structure

```json
{
	"messages": [
		{
			"id": "id of message",
			"text": "message",
			"rules": ["rules"]
		}
	]
}
```
